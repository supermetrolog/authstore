package mysql

import (
	"authstore/internal/common/loggerinterface"
	access "authstore/internal/domain/access/entity"
	"context"
	"database/sql"
)

type repository struct {
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

func (r *repository) Create(ctx context.Context, dto *access.CreateAccessDTO) (access.AccessID, error) {
	sql := `INSERT INTO access (
			token,
			expire, 
			user_id, 
			browser, 
			browser_version, 
			os, 
			os_version, 
			device, 
			is_mobile, 
			is_tablet, 
			is_desktop, 
			is_bot, 
			url, 
			full_user_agent,
			status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := r.client.PrepareContext(ctx, sql)

	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(
		ctx,
		dto.Token.Token,
		dto.Token.Expire,
		dto.UserID,
		dto.UserAgent.Browser,
		dto.UserAgent.BrowserVersion,
		dto.UserAgent.OS,
		dto.UserAgent.OSVersion,
		dto.UserAgent.Device,
		dto.UserAgent.IsMobile,
		dto.UserAgent.IsTablet,
		dto.UserAgent.IsDesktop,
		dto.UserAgent.IsBot,
		dto.UserAgent.URL,
		dto.UserAgent.FullUserAgent,
		access.StatusActive,
	)

	if err != nil {
		return 0, err
	}
	accessId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return access.AccessID(accessId), nil
}
func (r *repository) fetch(ctx context.Context, query string, args ...any) ([]*access.Access, error) {
	response, err := r.client.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer response.Close()
	var accesses []*access.Access

	for response.Next() {
		var item access.Access
		var Token access.Token
		var UserAgent access.UserAgent
		response.Scan(
			&item.ID,
			&item.UserID,
			&item.CreatedAt,
			&Token.Token,
			&Token.Expire,
			&UserAgent.Browser,
			&UserAgent.BrowserVersion,
			&UserAgent.OS,
			&UserAgent.OSVersion,
			&UserAgent.Device,
			&UserAgent.IsMobile,
			&UserAgent.IsTablet,
			&UserAgent.IsDesktop,
			&UserAgent.IsBot,
			&UserAgent.URL,
			&UserAgent.FullUserAgent,
			&item.Status,
		)

		item.Token = &Token
		item.UserAgent = &UserAgent
		accesses = append(accesses, &item)
	}
	return accesses, nil
}
func (r repository) FindByAccessToken(ctx context.Context, token string) (*access.Access, error) {
	sql := `SELECT
		id, user_id, created_at, token, expire, browser, browser_version, os, os_version, device, is_mobile, is_tablet, is_desktop, is_bot, url, full_user_agent, status
		FROM access
		WHERE token = ? LIMIT 1`
	accesses, err := r.fetch(ctx, sql, token)
	if err != nil {
		return nil, err
	}
	if len(accesses) > 0 {
		return accesses[0], nil
	}
	return nil, nil
}

func (r repository) Delete(ctx context.Context, id access.AccessID) error {
	sql := `DELETE FROM access WHERE id = ?`

	stmt, err := r.client.PrepareContext(ctx, sql)

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, id)
	return err
}

func (r repository) DisableAccess(ctx context.Context, id access.AccessID) error {
	sql := "UPDATE access SET status = ? WHERE id = ?"

	stmt, err := r.client.PrepareContext(ctx, sql)

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, access.StatusInactive, id)
	return err
}
