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
	"github.com/mileusna/useragent"
)

const (
	usersURL  = "/users"
	userURL   = "/users/:id"
	loginURL  = "/login"
	logoutURL = "/logout"
)

type Service interface {
	FindById(context.Context, user.UserID) (*user.User, error)
	FindByUsername(context.Context, string) (*user.User, error)
	FindByAccessToken(context.Context, string) (*user.User, error)
	FindByActiveAccessToken(context.Context, string) (*user.User, error)
	FindAll(context.Context) ([]*user.User, error)
	Create(context.Context, *user.CreateUserDTO) (user.UserID, error)
	Update(context.Context, *user.UpdateUserDTO) error
	Login(context.Context, *user.LoginUserDTO, *access.UserAgent) (*access.Token, error)
	Logout(ctx context.Context, token string) error
}

type handler struct {
	logger  loggerinterface.Logger
	service Service
	test    string
}

func NewHandler(logger loggerinterface.Logger, service Service) *handler {
	return &handler{
		logger:  logger,
		service: service,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	md := middleware.NewMiddleware(h.logger)

	// router.GET(usersURL, md.AdapterMiddleware(md.AuthMiddleware(md.DefaultMiddlewares(h.GetUserList), h.service)))
	router.GET(usersURL, md.AdapterMiddleware(md.DefaultMiddlewares(md.AuthMiddleware(h.GetUserList, h.service))))
	router.GET(userURL, md.AdapterMiddleware(md.DefaultMiddlewares(h.GetUserByID)))
	router.POST(usersURL, md.AdapterMiddleware(md.DefaultMiddlewares(h.CreateUser)))
	router.PATCH(userURL, md.AdapterMiddleware(md.DefaultMiddlewares(h.UpdateUser)))
	router.POST(loginURL, md.AdapterMiddleware(md.DefaultMiddlewares(h.LoginUser)))
	router.GET(logoutURL, md.AdapterMiddleware(md.DefaultMiddlewares(md.AuthMiddleware(h.LogoutUser, h.service))))
}

func (h *handler) GetUserList(hc *handlerContext.HandleContext) error {
	h.logger.Info("UserList USER", hc.GetUser())
	users, err := h.service.FindAll(context.Background())
	if err != nil {
		return apperror.NewHandlerError(err, http.StatusInternalServerError)
	}
	usersJSON, err := json.Marshal(users)
	if err != nil {
	}
	hc.W.Write(usersJSON)
	return nil
}
func (h *handler) GetUserByID(hc *handlerContext.HandleContext) error {
	h.logger.Info("ByID USER", hc.GetUser())
	id, err := strconv.Atoi(hc.P.ByName("id"))
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, "Invalid url param (id)", http.StatusBadRequest)
	}
	user, err := h.service.FindById(context.Background(), user.UserID(id))
	if err != nil {
		return apperror.NewHandlerError(err, http.StatusInternalServerError)
	}

	if user == nil {
		return apperror.ErrNotFound
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return apperror.NewHandlerError(err, http.StatusInternalServerError)
	}
	hc.W.Write(userJSON)
	return nil
}
func (h *handler) CreateUser(hc *handlerContext.HandleContext) error {
	var CreateUserDTO user.CreateUserDTO
	err := json.NewDecoder(hc.R.Body).Decode(&CreateUserDTO)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusBadRequest)
	}
	userID, err := h.service.Create(context.Background(), &CreateUserDTO)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusInternalServerError)
	}
	hc.W.Write([]byte(strconv.Itoa(int(userID))))
	return nil
}

func (h *handler) UpdateUser(hc *handlerContext.HandleContext) error {
	id, err := strconv.Atoi(hc.P.ByName("id"))
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, "Invalid url param (id)", http.StatusBadRequest)
	}
	userModel, err := h.service.FindById(context.Background(), user.UserID(id))
	if err != nil {
		return apperror.NewHandlerError(err, http.StatusInternalServerError)
	}

	if userModel == nil {
		return apperror.ErrNotFound
	}
	UpdateUserDTO := user.UpdateUserDTO{
		ID:       userModel.ID,
		Email:    userModel.Email,
		Username: userModel.Username,
	}
	err = json.NewDecoder(hc.R.Body).Decode(&UpdateUserDTO)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusBadRequest)
	}
	err = h.service.Update(context.Background(), &UpdateUserDTO)
	if err != nil {
		return apperror.NewHandlerError(err, http.StatusInternalServerError)
	}
	hc.W.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *handler) LoginUser(hc *handlerContext.HandleContext) error {
	var LoginUserDTO user.LoginUserDTO

	err := json.NewDecoder(hc.R.Body).Decode(&LoginUserDTO)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusBadRequest)
	}

	ua := useragent.Parse(hc.R.Header.Get("User-Agent"))

	token, err := h.service.Login(
		context.Background(),
		&LoginUserDTO,
		&access.UserAgent{
			Browser:        &ua.Name,
			BrowserVersion: &ua.Version,
			OS:             &ua.OS,
			OSVersion:      &ua.OSVersion,
			Device:         &ua.Device,
			IsMobile:       &ua.Mobile,
			IsTablet:       &ua.Tablet,
			IsDesktop:      &ua.Desktop,
			IsBot:          &ua.Bot,
			URL:            &ua.URL,
			FullUserAgent:  &ua.String,
		},
	)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusInternalServerError)
	}
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusBadRequest)
	}
	hc.W.WriteHeader(http.StatusOK)
	hc.W.Write(tokenBytes)
	return nil
}

func (h *handler) LogoutUser(hc *handlerContext.HandleContext) error {
	h.logger.Info("LOGOUT USER HANDLER")
	noAuthErr := apperror.NewAuthError("user not auth")
	if hc.IsGuest() {
		return apperror.NewHandlerErrorWithMessage(noAuthErr, noAuthErr.Error(), http.StatusUnauthorized)
	}
	u := hc.GetUser()
	if u == nil {
		return apperror.NewHandlerErrorWithMessage(noAuthErr, noAuthErr.Error(), http.StatusUnauthorized)
	}
	err := h.service.Logout(context.Background(), hc.GetToken())

	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusInternalServerError)
	}
	hc.W.WriteHeader(http.StatusNoContent)
	return nil
}
