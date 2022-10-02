package mysql

import (
	"authstore/internal/common/loggerinterface"
	"database/sql"
)

type Repository struct {
	logger loggerinterface.Logger
	client *sql.DB
}

func (r *Repository) Close() error {
	return r.client.Close()
}
