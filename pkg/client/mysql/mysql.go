package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewClient() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@/auth")
	if err != nil {
		return nil, fmt.Errorf("mysql open failed: %v", err)
	}
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("mysql connection failed: %v", err)
	}
	return db, nil
}
