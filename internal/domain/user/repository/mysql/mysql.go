package mysql

import (
	"authstore/internal/domain/user/entity/user"
	"authstore/pkg/logging"
	"authstore/pkg/querybuilder"
	"context"
	"database/sql"
)

type repository struct {
	querybuilder.Builder
	logger *logging.Logger
	client *sql.DB
}

func NewRepository(logger *logging.Logger, client *sql.DB) *repository {
	return &repository{
		logger: logger,
		client: client,
	}
}
func (r *repository) Close() error {
	return r.client.Close()
}
func (r *repository) fetchAccesses(ctx context.Context, query string, args ...any) ([]*user.Access, error) {
	response, err := r.client.QueryContext(ctx, query, args...)
	defer response.Close()

	if err != nil {
		return nil, err
	}
	var accesses []*user.Access

	for response.Next() {
		var item user.Access
		var Token user.Token
		var UserAgent user.UserAgent
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
		)

		item.Token = &Token
		item.UserAgent = &UserAgent
		accesses = append(accesses, &item)
	}
	return accesses, nil
}
func (r *repository) fetch(ctx context.Context, query string, args ...any) ([]*user.User, error) {
	response, err := r.client.QueryContext(ctx, query, args...)
	defer response.Close()

	if err != nil {
		return nil, err
	}
	var users []*user.User

	for response.Next() {
		var item user.User
		response.Scan(&item.ID, &item.Username, &item.Email, &item.PasswordHash, &item.RefreshToken)
		accessesSql := `SELECT 
		id, user_id, created_at, token, expire, browser, browser_version, os, os_version, device, is_mobile, is_tablet, is_desktop, is_bot, url, full_user_agent 
		FROM access 
		WHERE user_id = ?`
		accesses, err := r.fetchAccesses(ctx, accessesSql, item.ID)
		if err != nil {
			return nil, err
		}
		item.Accesses = accesses
		users = append(users, &item)
	}
	return users, nil
}

func (r *repository) FindById(ctx context.Context, id user.UserID) (*user.User, error) {
	query := "SELECT id, username, email, password_hash, refresh_token FROM user WHERE id = ? LIMIT 1"
	users, err := r.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		return users[0], nil
	}

	return nil, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	query := "SELECT id, username, email, password_hash, refresh_token FROM user WHERE username = ? LIMIT 1"
	users, err := r.fetch(ctx, query, username)
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		return users[0], nil
	}

	return nil, nil
}
func (r *repository) FindAll(ctx context.Context) ([]*user.User, error) {
	sql := "SELECT id, username, email, password_hash, refresh_token FROM user"

	users, err := r.fetch(ctx, sql)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) Create(ctx context.Context, dto *user.CreateUserDTO) (user.UserID, error) {
	sql := "INSERT INTO user (email, username, password_hash) VALUES (?, ?, ?)"
	stmt, err := r.client.PrepareContext(ctx, sql)

	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, dto.Email, dto.Username, dto.Password)

	if err != nil {
		return 0, err
	}
	userId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return user.UserID(userId), nil
}

func (r *repository) Update(ctx context.Context, dto *user.UpdateUserDTO) error {
	sql := "UPDATE user SET email = ?, username = ?, password_hash = ? WHERE id = ?"
	stmt, err := r.client.PrepareContext(ctx, sql)

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, dto.Email, dto.Username, dto.Password, dto.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CreateAccess(ctx context.Context, dto *user.CreateAccessDTO) (user.TokenID, error) {
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
			full_user_agent
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
	)

	if err != nil {
		return 0, err
	}
	tokenId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return user.TokenID(tokenId), nil
}
