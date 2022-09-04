package middleware

import (
	"authstore/internal/apperror"
	"authstore/pkg/httpheader"
	"authstore/pkg/logging"
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

type Handle func(http.ResponseWriter, *http.Request, httprouter.Params) error

type Middleware struct {
	logger logging.Logger
}

func NewMiddleware(logger logging.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m *Middleware) DefaultMiddlewares(handle Handle) httprouter.Handle {
	var result httprouter.Handle
	result = m.ErrorHandlingMiddleware(handle)
	result = m.ContentTypeJSONMiddleware(result)
	return result
}
func (m *Middleware) ContentTypeJSONMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Add(
			httpheader.ContentTypeKey,
			httpheader.ContentTypeJSON,
		)
		next(w, r, p)
	}

}

func (m *Middleware) ErrorHandlingMiddleware(next Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := next(w, r, p)
		if err == nil {
			return
		}

		switch e := err.(type) {
		case *apperror.HandlerError:
			switch eo := e.OriginError().(type) {
			case apperror.ValidationError:
				fmt.Println(eo)
				e.SetStatusCode(http.StatusBadRequest)
				e.SetName(apperror.ValidationErrorName)
			case apperror.LoginError:
				fmt.Println(eo)
				e.SetStatusCode(http.StatusUnauthorized)
				e.SetName(apperror.LoginErrorName)
			}

			w.WriteHeader(e.StatusCode())

			if e.StatusCode() == http.StatusInternalServerError {
				m.logger.Errorf("Internal server error %v", err)
			}

			responseErr := apperror.ResponseError{
				Error: err,
			}
			errBytes, jsonErr := json.Marshal(responseErr)

			if jsonErr != nil {
				m.logger.Errorf("error marshaling exited with error: %v", jsonErr)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Server error"))
				return
			}

			w.Write(errBytes)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unknown error"))
		}

	}
}
