package mysql

import (
	"authstore/internal/common/loggerinterface"
	tree "authstore/internal/domain/tree/entity"
	user "authstore/internal/domain/user/entity"
	"authstore/pkg/querybuilder"
	"context"
	"database/sql"
)

const Table = "user"

type repository struct {
	querybuilder.Builder
	logger loggerinterface.Logger
	client *sql.DB
}

func NewRepository(logger loggerinterface.Logger, client *sql.DB) *repository {
	return &repository{
		logger: logger,
		client: client,
	}
}
func (r *repository) Close() error {
	return r.client.Close()
}

func (r *repository) fetch(ctx context.Context, query string, args ...any) ([]*tree.Node, error) {
	response, err := r.client.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer response.Close()
	var nodes []*tree.Node

	for response.Next() {
		var item tree.Node
		response.Scan(&item.ID, &item.ParentID, &item.UserID, &item.Name, &item.Type, &item.Status, &item.CreatedAt, &item.UpdatedAt)
		nodes = append(nodes, &item)
	}
	return nodes, nil
}
func (r *repository) FindTreeByUserID(ctx context.Context, userID user.UserID) (*tree.Node, error) {
	sql := "SELECT id, parent_id, user_id, name, type, status, created_at, updated_at FROM tree WHERE user_id = ?"

	nodes, err := r.fetch(ctx, sql, userID)

	if err != nil {
		return nil, err
	}
	if len(nodes) != 0 {
		return nodes[0], nil
	}
	return nil, nil
}
func (r *repository) FindNodeByID(ctx context.Context, nodeID tree.NodeID) (*tree.Node, error) {
	sql := "SELECT id, parent_id, user_id, name, type, status, created_at, updated_at FROM tree WHERE id = ?"

	nodes, err := r.fetch(ctx, sql, nodeID)

	if err != nil {
		return nil, err
	}
	if len(nodes) != 0 {
		return nodes[0], nil
	}
	return nil, nil
}

func (r *repository) FindNodeByName(ctx context.Context, name string) (*tree.Node, error) {
	sql := "SELECT id, parent_id, user_id, name, type, status, created_at, updated_at FROM tree WHERE name = ?"

	nodes, err := r.fetch(ctx, sql, name)

	if err != nil {
		return nil, err
	}
	if len(nodes) != 0 {
		return nodes[0], nil
	}
	return nil, nil
}

func (r *repository) FindRootByUserID(ctx context.Context, userID user.UserID) (*tree.Node, error) {
	sql := "SELECT id, parent_id, user_id, name, type, status, created_at, updated_at FROM tree WHERE user_id = ? AND parent_id IS NULL"

	nodes, err := r.fetch(ctx, sql, userID)

	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, nil
	}
	return r.findRecursive(ctx, nodes[0])
}

func (r *repository) findRecursive(ctx context.Context, node *tree.Node) (*tree.Node, error) {
	sql := "SELECT id, parent_id, user_id, name, type, status, created_at, updated_at FROM tree WHERE parent_id = ?"

	nodes, err := r.fetch(ctx, sql, node.ID)

	if err != nil {
		return nil, err
	}

	node.Childrens = nodes

	for _, n := range node.Childrens {
		_, err := r.findRecursive(ctx, n)
		if err != nil {
			return nil, err
		}
	}

	return node, nil
}

func (r *repository) CreateNode(ctx context.Context, node *tree.CreateNodeDTO) (tree.NodeID, error) {
	sql := "INSERT INTO tree (parent_id, user_id, name, type, status) VALUES (?, ?, ?, ?, ?)"
	stmt, err := r.client.PrepareContext(ctx, sql)

	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, node.ParentID, node.UserID, node.Name, node.Type, node.Status)

	if err != nil {
		return 0, err
	}
	nodeID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return tree.NodeID(nodeID), nil
}
