package entity

import (
	"authstore/pkg/validator"
)

const (
	TypeNode     = 1
	TypeLeaf     = 2
	StatusActive = 1
)

type NodeID int

type Node struct {
	ID        *NodeID
	ParentID  *NodeID
	UserID    *int
	Name      *string
	CreatedAt *string
	UpdatedAt *string
	Type      *int
	Status    *int
	Childrens []*Node
	Parent    *Node
}

func (n Node) IsNode() bool {
	return *n.Type == TypeNode
}

func (n Node) IsLeaf() bool {
	return *n.Type == TypeLeaf
}

type CreateNodeDTO struct {
	ParentID *NodeID `json:"parent_id"`
	UserID   *NodeID `json:"user_id"`
	Name     *string `json:"name"`
	Type     *int    `json:"type"`
	Status   *int    `json:"status"`
}

func (dto *CreateNodeDTO) Validations() map[string][]validator.ValidatorHandler {
	return map[string][]validator.ValidatorHandler{
		"name": {
			validator.Required(dto.Name),
		},
		"user_id": {
			validator.Required(dto.UserID),
		},
		"type": {
			validator.Required(dto.Type),
			validator.MinValue(dto.Status, 1),
			validator.MaxValue(dto.Status, 2),
		},
		"status": {
			validator.Required(dto.Status),
			validator.MinValue(dto.Status, 1),
			validator.MaxValue(dto.Status, 2),
		},
	}

}
