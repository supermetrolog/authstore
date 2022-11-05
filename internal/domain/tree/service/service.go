package service

import (
	"authstore/internal/apperror"
	"authstore/internal/common/loggerinterface"
	tree "authstore/internal/domain/tree/entity"
	user "authstore/internal/domain/user/entity"
	"authstore/pkg/validator"
	"context"
	"fmt"
)

type Repository interface {
	CreateNode(context.Context, *tree.CreateNodeDTO) (tree.NodeID, error)
	FindTreeByUserID(context.Context, user.UserID) (*tree.Node, error)
	FindRootByUserID(context.Context, user.UserID) (*tree.Node, error)
	FindNodeByID(context.Context, tree.NodeID) (*tree.Node, error)
	FindNodeByName(context.Context, string) (*tree.Node, error)
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

	root, err := s.repository.FindRootByUserID(ctx, user.UserID(*node.UserID))
	if err != nil {
		return 0, err
	}

	if root == nil && node.ParentID != nil {
		return 0, fmt.Errorf("Parent node with ID (%d) do not exist", *node.ParentID)
	}
	if root != nil && node.ParentID == nil {
		return 0, fmt.Errorf("Root node for user with id (%d) already exist", *node.UserID)
	}
	parentNode, err := s.repository.FindNodeByID(ctx, *node.ParentID)
	if err != nil {
		return 0, err
	}
	if parentNode == nil {
		return 0, fmt.Errorf("Parent node with ID (%d) do not exist", *node.ParentID)
	}
	existWithNameNode, err := s.repository.FindNodeByName(ctx, *node.Name)
	if err != nil {
		return 0, err
	}
	if existWithNameNode != nil {
		return 0, fmt.Errorf("Node with name (%s) already exist", *node.Name)
	}
	return s.repository.CreateNode(ctx, node)
}

func (s Service) FindTreeByUserID(ctx context.Context, userID user.UserID) (*tree.Node, error) {
	return s.repository.FindRootByUserID(ctx, userID)
}
