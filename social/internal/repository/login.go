package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/moeryomenko/highload-architect-otus/social/internal/domain"
)

const (
	insertLoginQuery = `INSERT INTO users (id, nickname, password) VALUES (UNHEX(?), ?, ?)`
	selectLoginQuery = `SELECT HEX(id), password FROM users WHERE nickname = ?`
)

var ErrNotFound = errors.New("not found")

// Login incapsulates login/signup repository logic.
type Login struct {
	conn *sql.DB
}

// NewLogin returns new instance of login repository.
func NewLogin(conn *sql.DB) *Login {
	return &Login{conn: conn}
}

// Save saves signup credentials for login.
func (r *Login) Save(ctx context.Context, login *domain.Login) error {
	return transaction(ctx, r.conn, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, insertLoginQuery, uuidToBinary(login.UserID), login.Nickname, login.Password)
		return err
	})
}

// Get returns login credentials for check login.
func (r *Login) Get(ctx context.Context, nickname string) (*domain.Login, error) {
	login := &domain.Login{Nickname: nickname}

	row := r.conn.QueryRowContext(ctx, selectLoginQuery, login.Nickname)
	err := row.Err()
	switch err {
	case nil:
		var id string
		err = row.Scan(&id, &login.Password)
		if err != nil {
			return nil, err
		}
		login.UserID = binaryToUUID(id)
		return login, err
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
