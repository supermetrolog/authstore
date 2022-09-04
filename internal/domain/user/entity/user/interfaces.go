package user

import (
	"context"
)

type Service interface {
	FindById(context.Context, UserID) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
	FindAll(context.Context) ([]*User, error)
	Create(context.Context, *CreateUserDTO) (UserID, error)
	Update(context.Context, *UpdateUserDTO) error
	Login(context.Context, *LoginUserDTO, *UserAgent) (*Token, error)
}
type Repository interface {
	FindById(context.Context, UserID) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
	FindAll(context.Context) ([]*User, error)
	Create(context.Context, *CreateUserDTO) (UserID, error)
	Update(context.Context, *UpdateUserDTO) error
	CreateAccess(context.Context, *CreateAccessDTO) (TokenID, error)
}
