package http

import (
	"authstore/internal/apperror"
	handlerContext "authstore/internal/common/http/handler"
	"authstore/internal/common/http/middleware"
	"authstore/internal/common/loggerinterface"
	tree "authstore/internal/domain/tree/entity"
	user "authstore/internal/domain/user/entity"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const (
	treeURL = "/tree"
	nodeURL = "/tree/:id"
)

type Service interface {
	FindTreeByUserID(ctx context.Context, id user.UserID) (*tree.Node, error)
	CreateNode(context.Context, *tree.CreateNodeDTO) (tree.NodeID, error)
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
	router.GET(treeURL, md.AdapterMiddleware(md.DefaultMiddlewares(md.AuthMiddleware(h.GetTree, h.userService))))
	router.POST(treeURL, md.AdapterMiddleware(md.DefaultMiddlewares(md.AuthMiddleware(h.CreateNode, h.userService))))
}

func (h *handler) GetTree(hc *handlerContext.HandleContext) error {
	tree, err := h.service.FindTreeByUserID(context.Background(), *hc.UserContext.GetUser().ID)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusInternalServerError)
	}

	bytes, err := json.Marshal(tree)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusInternalServerError)
	}

	hc.W.WriteHeader(http.StatusOK)
	hc.W.Write(bytes)
	return nil
}

func (h *handler) CreateNode(hc *handlerContext.HandleContext) error {

	var CreateNodeDTO tree.CreateNodeDTO
	err := json.NewDecoder(hc.R.Body).Decode(&CreateNodeDTO)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusBadRequest)
	}

	nodeID, err := h.service.CreateNode(context.Background(), &CreateNodeDTO)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusInternalServerError)
	}
	hc.W.WriteHeader(http.StatusOK)
	hc.W.Write([]byte(strconv.Itoa(int(nodeID))))
	return nil
}
