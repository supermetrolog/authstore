package handler

import (
	user "authstore/internal/domain/user/entity"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handle func(*HandleContext) error

type HandleContext struct {
	HttpContext
	UserContext
}

func NewHandleContext(w http.ResponseWriter, r *http.Request, p httprouter.Params) *HandleContext {
	return &HandleContext{
		HttpContext: HttpContext{
			W: w,
			R: r,
			P: p,
		},
	}
}

type HttpContext struct {
	W http.ResponseWriter
	R *http.Request
	P httprouter.Params
}
type UserContext struct {
	user   *user.User
	isAuth bool
}

func (u *UserContext) GetUser() *user.User {
	return u.user
}

// once
func (u *UserContext) SetUser(model *user.User) {
	if u.user != nil || model == nil {
		return
	}
	u.user = model
	u.isAuth = true
}
func (u *UserContext) IsGuest() bool {
	return !u.isAuth
}

func (u *UserContext) IsAuth() bool {
	return u.isAuth
}
