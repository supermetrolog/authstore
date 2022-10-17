package http

import (
	"authstore/internal/apperror"
	handlerContext "authstore/internal/common/http/handler"
	"authstore/internal/common/http/middleware"
	"authstore/internal/common/loggerinterface"
	access "authstore/internal/domain/access/entity"
	user "authstore/internal/domain/user/entity"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const (
	accessesURL            = "/accesses"
	accessURL              = "/access/:id"
	userAccessesDisableURL = "/accesses/disable-user/:user_id"
)

type Service interface {
	FindAll(context.Context) ([]*access.Access, error)
	DisableAccess(context.Context, access.AccessID) error
	DisableAccesses(ctx context.Context, userID int64) error
}
type UserService interface {
	FindByActiveAccessToken(ctx context.Context, token string) (*user.User, error)
}
type handler struct {
	logger      loggerinterface.Logger
	service     Service
	userService UserService
}

func NewHandler(logger loggerinterface.Logger, service Service, userService UserService) *handler {
	return &handler{
		logger:      logger,
		service:     service,
		userService: userService,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	md := middleware.NewMiddleware(h.logger)
	router.GET(accessesURL, md.AdapterMiddleware(md.DefaultMiddlewares(md.AuthMiddleware(h.GetAccessList, h.userService))))
	router.DELETE(accessURL, md.AdapterMiddleware(md.DefaultMiddlewares(md.AuthMiddleware(h.DisableAccess, h.userService))))
	router.DELETE(userAccessesDisableURL, md.AdapterMiddleware(md.DefaultMiddlewares(md.AuthMiddleware(h.DisableAccesses, h.userService))))
}

func (h *handler) GetAccessList(hc *handlerContext.HandleContext) error {
	accesses, err := h.service.FindAll(context.Background())
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusInternalServerError)
	}
	bytes, err := json.Marshal(accesses)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusInternalServerError)
	}
	hc.W.WriteHeader(http.StatusOK)
	hc.W.Write(bytes)
	return nil
}

func (h *handler) DisableAccess(hc *handlerContext.HandleContext) error {
	id, err := strconv.Atoi(hc.P.ByName("id"))
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, "Invalid url param (id)", http.StatusBadRequest)
	}
	err = h.service.DisableAccess(context.Background(), access.AccessID(id))
	if err != nil {
		return apperror.NewHandlerError(err, http.StatusInternalServerError)
	}
	hc.W.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *handler) DisableAccesses(hc *handlerContext.HandleContext) error {
	id, err := strconv.Atoi(hc.P.ByName("user_id"))
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, "Invalid url param (id)", http.StatusBadRequest)
	}
	err = h.service.DisableAccesses(context.Background(), int64(id))
	if err != nil {
		return apperror.NewHandlerError(err, http.StatusInternalServerError)
	}
	hc.W.WriteHeader(http.StatusNoContent)
	return nil
}
