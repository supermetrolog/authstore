package http

import (
	"authstore/internal/apperror"
	"authstore/internal/common/http/middleware"
	"authstore/internal/domain/user/entity/user"
	"authstore/pkg/logging"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mileusna/useragent"
)

const (
	usersURL = "/users"
	userURL  = "/users/:id"
	loginURL = "/login"
)

type handler struct {
	logger  *logging.Logger
	service user.Service
}

func NewHandler(logger *logging.Logger, service user.Service) *handler {
	return &handler{
		logger:  logger,
		service: service,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	md := middleware.NewMiddleware(*h.logger)

	router.GET(usersURL, md.DefaultMiddlewares(h.GetUserList))
	router.GET(userURL, md.DefaultMiddlewares(h.GetUserByID))
	router.POST(usersURL, md.DefaultMiddlewares(h.CreateUser))
	router.PATCH(userURL, md.DefaultMiddlewares(h.UpdateUser))
	router.POST(loginURL, md.DefaultMiddlewares(h.LoginUser))
}

func (h *handler) GetUserList(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	users, err := h.service.FindAll(context.Background())
	if err != nil {
		return apperror.NewHandlerError(err, http.StatusInternalServerError)
	}
	usersJSON, err := json.Marshal(users)
	if err != nil {
	}
	w.Write(usersJSON)
	return nil
}
func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	id, err := strconv.Atoi(params.ByName("id"))
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
	w.Write(userJSON)
	return nil
}
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	var CreateUserDTO user.CreateUserDTO
	err := json.NewDecoder(r.Body).Decode(&CreateUserDTO)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusBadRequest)
	}
	userID, err := h.service.Create(context.Background(), &CreateUserDTO)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(strconv.Itoa(int(userID))))
	return nil
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	id, err := strconv.Atoi(params.ByName("id"))
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
	err = json.NewDecoder(r.Body).Decode(&UpdateUserDTO)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusBadRequest)
	}
	err = h.service.Update(context.Background(), &UpdateUserDTO)
	if err != nil {
		return apperror.NewHandlerError(err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	var LoginUserDTO user.LoginUserDTO

	err := json.NewDecoder(r.Body).Decode(&LoginUserDTO)
	if err != nil {
		return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusBadRequest)
	}

	ua := useragent.Parse(r.Header.Get("User-Agent"))

	token, err := h.service.Login(
		context.Background(),
		&LoginUserDTO,
		&user.UserAgent{
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
	w.WriteHeader(http.StatusOK)
	w.Write(tokenBytes)
	return nil
}
