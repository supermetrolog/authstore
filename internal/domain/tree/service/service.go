package service

import (
	"authstore/internal/apperror"
	"authstore/internal/common/loggerinterface"
	tree "authstore/internal/domain/tree/entity"
	user "authstore/internal/domain/user/entity"
	"authstore/pkg/validator"
	"context"
)

type Repository interface {
	CreateNode(context.Context, *tree.CreateNodeDTO) (tree.NodeID, error)
	FindTreeByUserID(context.Context, user.UserID) (*tree.Node, error)
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

func (s *Service) CreateNode(ctx context.Context, node *tree.CreateNodeDTO) (tree.NodeID, error) {
	errs := validator.New().Validate(node)
	if errs != nil {
		return 0, apperror.NewValidationError(errs)
	}

	return s.repository.CreateNode(ctx, node)
}

func (s Service) FindTreeByUserID(ctx context.Context, userID user.UserID) (*tree.Node, error) {
	return s.repository.FindTreeByUserID(ctx, userID)
}
