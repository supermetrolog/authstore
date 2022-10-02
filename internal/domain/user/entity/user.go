package user

import (
	access "authstore/internal/domain/access/entity"
	"authstore/pkg/validator"
)

type UserID int64

const (
	PASSWORD_MIN_LENGTH = 5
	PASSWORD_MAX_LENGTH = 16
	USERNAME_MIN_LENGTH = 5
	USERNAME_MAX_LENGTH = 16
	EMAIL_MIN_LENGTH    = 8
	EMAIL_MAX_LENGTH    = 32
)

type User struct {
	ID           *UserID          `json:"id" bson:"_id,omitempty" db:"id"`
	Email        *string          `json:"email" bson:"email" db:"email"`
	Username     *string          `json:"username" bson:"username" db:"username"`
	PasswordHash *string          `json:"-" bson:"password_hash" db:"password_hash"`
	RefreshToken *string          `json:"refresh_token" bson:"refresh_token" db:"refresh_token"`
	Accesses     []*access.Access `json:"-"`
}

func (u *User) Validations() map[string][]validator.ValidatorHandler {
	return map[string][]validator.ValidatorHandler{
		"username": {
			validator.Required(u.Username),
			validator.MinLength(u.Username, USERNAME_MIN_LENGTH),
			validator.MaxLength(u.Username, USERNAME_MAX_LENGTH),
		},
		"password_hash": {
			validator.Required(u.PasswordHash),
		},
		"email": {
			validator.Required(u.Email),
			validator.MinLength(u.Email, EMAIL_MIN_LENGTH),
			validator.MaxLength(u.Email, EMAIL_MAX_LENGTH),
		},
		"id": {
			validator.Required(u.ID),
		},
	}

}

type CreateUserDTO struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func (dto *CreateUserDTO) Validations() map[string][]validator.ValidatorHandler {
	return map[string][]validator.ValidatorHandler{
		"username": {
			validator.Required(dto.Username),
			validator.MinLength(dto.Username, USERNAME_MIN_LENGTH),
			validator.MaxLength(dto.Username, USERNAME_MAX_LENGTH),
		},
		"password": {
			validator.Required(dto.Password),
			validator.MinLength(dto.Password, PASSWORD_MIN_LENGTH),
			validator.MaxLength(dto.Password, PASSWORD_MAX_LENGTH),
			validator.WithoutSymbols(dto.Password, []rune(validator.LETTERS_ONLY_TEMPLATE)...),
		},
		"email": {
			validator.Required(dto.Email),
			validator.MinLength(dto.Email, EMAIL_MIN_LENGTH),
			validator.MaxLength(dto.Email, EMAIL_MAX_LENGTH),
		},
	}

}

type UpdateUserDTO struct {
	ID       *UserID `json:"id"`
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func (dto *UpdateUserDTO) Validations() map[string][]validator.ValidatorHandler {
	return map[string][]validator.ValidatorHandler{
		"username": {
			validator.Required(dto.Username),
			validator.MinLength(dto.Username, USERNAME_MIN_LENGTH),
			validator.MaxLength(dto.Username, USERNAME_MAX_LENGTH),
		},
		"password": {
			validator.IfNotNil(dto.Password, validator.Required(dto.Password)),
			validator.IfNotNil(dto.Password, validator.MinLength(dto.Password, PASSWORD_MIN_LENGTH)),
			validator.IfNotNil(dto.Password, validator.MaxLength(dto.Password, PASSWORD_MAX_LENGTH)),
		},
		"email": {
			validator.Required(dto.Email),
			validator.MinLength(dto.Email, EMAIL_MIN_LENGTH),
			validator.MaxLength(dto.Email, EMAIL_MAX_LENGTH),
		},
		"id": {
			validator.Required(dto.ID),
		},
	}
}

type LoginUserDTO struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func (dto LoginUserDTO) Validations() map[string][]validator.ValidatorHandler {
	return map[string][]validator.ValidatorHandler{
		"username": {
			validator.Required(dto.Username),
		},
		"password": {
			validator.Required(dto.Password),
		},
	}
}
