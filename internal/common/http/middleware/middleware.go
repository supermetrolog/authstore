package middleware

import (
	"authstore/internal/apperror"
	"authstore/internal/common/http/handler"
	"authstore/internal/common/loggerinterface"
	user "authstore/internal/domain/user/entity"
	"authstore/pkg/httpheader"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HandlerError interface {
	Error() string
	StatusCode() int
	SetStatusCode(int)
	OriginError() error
}

type Middleware struct {
	logger loggerinterface.Logger
}

func NewMiddleware(logger loggerinterface.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m *Middleware) DefaultMiddlewares(handle handler.Handle) handler.Handle {
	var result handler.Handle
	result = m.ErrorHandlingMiddleware(handle)   // ^
	result = m.ContentTypeJSONMiddleware(result) // |
	return result
}
func (m *Middleware) AdapterMiddleware(next handler.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		m.logger.Info("ADAPTER")
		httpCtx := handler.NewHandleContext(w, r, p)
		next(httpCtx)
	}
}

type UserService interface {
	FindByActiveAccessToken(ctx context.Context, token string) (*user.User, error)
}

func (m *Middleware) AuthMiddleware(next handler.Handle, userService UserService) handler.Handle {
	return func(httpCtx *handler.HandleContext) error {
		m.logger.Info("AUTH")
		token := httpCtx.R.Header.Get("Authorization")
		m.logger.Info(token)
		if token == "" {
			err := apperror.NewAuthError("The request does not contain an authorization token")
			return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusUnauthorized)
		}
		u, err := userService.FindByActiveAccessToken(context.Background(), token)
		if err != nil {
			err := apperror.NewAuthError("Server error")
			return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusInternalServerError)
		}
		if u == nil {
			err := apperror.NewAuthError("Invalid access token")
			return apperror.NewHandlerErrorWithMessage(err, err.Error(), http.StatusUnauthorized)
		}
		httpCtx.SetUser(u, token)
		return next(httpCtx)
	}

}
func (m *Middleware) ContentTypeJSONMiddleware(next handler.Handle) handler.Handle {
	return func(httpCtx *handler.HandleContext) error {
		m.logger.Info("Content")
		httpCtx.W.Header().Add(
			httpheader.ContentTypeKey,
			httpheader.ContentTypeJSON,
		)
		return next(httpCtx)
	}

}

func (m *Middleware) ErrorHandlingMiddleware(next handler.Handle) handler.Handle {
	return func(httpCtx *handler.HandleContext) error {
		m.logger.Info("Error")
		err := next(httpCtx)
		if err == nil {
			return nil
		}

		switch e := err.(type) {
		case *apperror.HandlerError:
			fmt.Println(e.StatusCode())
			fmt.Println(e.Name)
			fmt.Println(e.OriginalError.Error())
			switch eo := e.OriginError().(type) {
			case apperror.ValidationError:
				fmt.Println(eo)
				e.SetStatusCode(http.StatusBadRequest)
				e.SetName(apperror.ValidationErrorName)
			case apperror.LoginError:
				e.SetStatusCode(http.StatusUnauthorized)
				e.SetName(apperror.LoginErrorName)
			}
			httpCtx.W.WriteHeader(e.StatusCode())

			if e.StatusCode() == http.StatusInternalServerError {
				m.logger.Errorf("Internal server error %v", err)
			}

			responseErr := apperror.ResponseError{
				Error: err,
			}
			errBytes, jsonErr := json.Marshal(responseErr)

			if jsonErr != nil {
				m.logger.Errorf("error marshaling exited with error: %v", jsonErr)
				httpCtx.W.WriteHeader(http.StatusInternalServerError)
				httpCtx.W.Write([]byte("Server error"))
				return jsonErr
			}

			httpCtx.W.Write(errBytes)
		default:
			httpCtx.W.WriteHeader(http.StatusInternalServerError)
			httpCtx.W.Write([]byte("Unknown error"))
		}
		return err
	}
}
