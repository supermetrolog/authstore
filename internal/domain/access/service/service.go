package service

import (
	"authstore/internal/apperror"
	"authstore/internal/common/loggerinterface"
	access "authstore/internal/domain/access/entity"
	"context"
	"errors"
)

type Repository interface {
	FindAll(context.Context) ([]*access.Access, error)
	FindByUserID(ctx context.Context, userID int64) ([]*access.Access, error)
	FindByID(context.Context, access.AccessID) (*access.Access, error)
	DisableAccess(context.Context, access.AccessID) error
	DisableAccesses(context.Context, ...access.AccessID) error
}

type Service struct {
	logger     loggerinterface.Logger
	repository Repository
}

func NewService(logger loggerinterface.Logger, repository Repository) *Service {
	return &Service{
		logger:     logger,
		repository: repository,
	}
}

func (s *Service) FindAll(ctx context.Context) ([]*access.Access, error) {
	accesses, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return accesses, nil
}

func (s *Service) DisableAccess(ctx context.Context, id access.AccessID) error {
	a, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if a == nil {
		return apperror.ErrNotFound
	}
	if a.IsInactive() {
		return errors.New("this access already disabled")
	}

	return s.repository.DisableAccess(ctx, id)
}

func (s *Service) DisableAccesses(ctx context.Context, userID int64) error {
	accesses, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if len(accesses) == 0 {
		return nil
	}
	var ids []access.AccessID
	for _, v := range accesses {
		ids = append(ids, *v.ID)
	}
	return s.repository.DisableAccesses(ctx, ids...)
}
