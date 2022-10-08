package access

import (
	"authstore/pkg/validator"
	"strconv"
)

const (
	StatusActive   int8 = 1
	StatusInactive int8 = -1
)

type Token struct {
	Token  *string `json:"-"`
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

func (a AccessID) String() string {
	return strconv.Itoa(int(a))
}

type Access struct {
	ID        *AccessID  `json:"id"`
	UserID    *int64     `json:"user_id"`
	Token     *Token     `json:"token"`
	UserAgent *UserAgent `json:"user_agent"`
	CreatedAt *string    `json:"created_at"`
	Status    *int8      `json:"status"`
}

func (a *Access) IsActive() bool {
	return *a.Status == StatusActive
}
func (a *Access) IsInactive() bool {
	return *a.Status == StatusInactive
}

type CreateAccessDTO struct {
	UserID    *int64     `json:"user_id"`
	Token     *Token     `json:"token"`
	UserAgent *UserAgent `json:"user_agent"`
}

func (dto *CreateAccessDTO) Validations() map[string][]validator.ValidatorHandler {
	return map[string][]validator.ValidatorHandler{
		"token": {
			validator.Required(dto.Token.Token),
		},
		"expire": {
			validator.Required(dto.Token.Expire),
		},
		"user_id": {
			validator.Required(dto.UserID),
		},
	}

}
