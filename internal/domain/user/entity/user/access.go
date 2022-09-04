package user

import (
	"authstore/pkg/validator"
)

type TokenID int64

type Token struct {
	Token  *string `json:"token"`
	Expire *uint64 `json:"expire"`
}

type UserAgent struct {
	Browser        *string `json:"browser"`
	BrowserVersion *string `json:"browser_version"`
	OS             *string `json:"os"`
	OSVersion      *string `json:"os_version"`
	Device         *string `json:"device"`
	IsMobile       *bool   `json:"is_mobile"`
	IsTablet       *bool   `json:"is_tablet"`
	IsDesktop      *bool   `json:"is_desktop"`
	IsBot          *bool   `json:"is_bot"`
	URL            *string `json:"url"`
	FullUserAgent  *string `json:"full_user_agent"`
}
type AccessID int64

type Access struct {
	ID        *AccessID  `json:"id"`
	UserID    *UserID    `json:"user_id"`
	Token     *Token     `json:"token"`
	UserAgent *UserAgent `json:"user_agent"`
	CreatedAt *string    `json:"created_at"`
}

type CreateAccessDTO struct {
	UserID    *UserID    `json:"user_id"`
	Token     *Token     `json:"token"`
	UserAgent *UserAgent `json:"user_agent"`
}

func (u *CreateAccessDTO) Validations() map[string][]validator.ValidatorHandler {
	return map[string][]validator.ValidatorHandler{
		"token": {
			validator.Required(u.Token.Token),
		},
		"expire": {
			validator.Required(u.Token.Expire),
		},
		"user_id": {
			validator.Required(u.UserID),
		},
	}

}
