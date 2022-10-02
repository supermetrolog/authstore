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

	"github.com/julienschmidt/httprouter"
)

const (
	accessesURL = "/accesses"
)

type Service interface {
	FindAll(context.Context) ([]*access.Access, error)
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
