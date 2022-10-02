package service

import (
	"authstore/internal/common/loggerinterface"
	access "authstore/internal/domain/access/entity"
	"context"
)

type Repository interface {
	FindAll(context.Context) ([]*access.Access, error)
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
