package service

import (
	"authstore/internal/apperror"
	"authstore/internal/domain/user/entity/user"
	"authstore/pkg/logging"
	"authstore/pkg/security"
	"authstore/pkg/validator"
	"context"
)

type Service struct {
	logger     *logging.Logger //fucking logger
	repository user.Repository
}

func NewService(logger *logging.Logger, repository user.Repository) *Service {
	return &Service{
		logger:     logger,
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, user *user.CreateUserDTO) (user.UserID, error) {

	errs := validator.New().Validate(user)
	if errs != nil {
		return 0, apperror.NewValidationError(errs)
	}
	hash, err := security.HashPassword(*user.Password, 14)
	if err != nil {
		return 0, err
	}
	*user.Password = hash
	return s.repository.Create(ctx, user)
}

func (s *Service) Update(ctx context.Context, user *user.UpdateUserDTO) error {
	errs := validator.New().Validate(user)
	if errs != nil {
		return apperror.NewValidationError(errs)
	}
	if user.Password == nil {
		return s.repository.Update(ctx, user)
	}
	hash, err := security.HashPassword(*user.Password, 14)
	if err != nil {
		return err
	}
	*user.Password = hash
	return s.repository.Update(ctx, user)
}

//Returned user list or error
func (s *Service) FindAll(ctx context.Context) ([]*user.User, error) {
	users, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (s *Service) FindById(ctx context.Context, id user.UserID) (*user.User, error) {
	user, err := s.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *Service) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	user, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, dto *user.LoginUserDTO, useragent *user.UserAgent) (*user.Token, error) {
	errs := validator.New().Validate(dto)
	if errs != nil {
		return nil, apperror.NewValidationError(errs)
	}
	model, err := s.FindByUsername(ctx, *dto.Username)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, apperror.NewLoginError("invalid username or password")
	}
	if !security.CheckPasswordHash(*dto.Password, *model.PasswordHash) {
		return nil, apperror.NewLoginError("invalid username or password")
	}

	token := security.GenerateRandomString(32)
	var tokenExpire uint64 = 3400 * 24
	Token := &user.Token{
		Token:  &token,
		Expire: &tokenExpire,
	}
	createAccessDTO := user.CreateAccessDTO{
		UserID:    model.ID,
		Token:     Token,
		UserAgent: useragent,
	}
	if errs := validator.New().Validate(&createAccessDTO); err != nil {
		return nil, apperror.NewValidationError(errs)
	}
	if _, err := s.repository.CreateAccess(ctx, &createAccessDTO); err != nil {
		return nil, err
	}
	return Token, nil
}
