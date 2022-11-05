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
	ID        *NodeID `json:"id"`
	ParentID  *NodeID `json:"parent_id"`
	UserID    *int    `json:"user_id"`
	Name      *string `json:"name"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
	Type      *int    `json:"type"`
	Status    *int    `json:"status"`
	Childrens []*Node `json:"childrens"`
	Parent    *Node   `json:"parent"`
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
