package repository

import (
	"context"
	"errors"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/moeryomenko/healing/decorators/mysql"
	"github.com/moeryomenko/highload-architect-otus/social/internal/domain"
)

const (
	insertLoginQuery = `INSERT INTO users (id, nickname, password) VALUES (UNHEX(?), ?, ?)`
	selectLoginQuery = `SELECT HEX(id), password FROM users WHERE nickname = ?`
)

var ErrNotFound = errors.New("not found")

// Login incapsulates login/signup repository logic.
type Login struct {
	writePool *mysql.Pool
	readPool  *mysql.Pool
}

// NewLogin returns new instance of login repository.
func NewLogin(writeConn, readConn *mysql.Pool) *Login {
	return &Login{writePool: writeConn, readPool: readConn}
}

// Save saves signup credentials for login.
func (r *Login) Save(ctx context.Context, login *domain.Login) error {
	return transaction(ctx, r.writePool, func(conn *client.Conn) error {
		_, err := conn.Execute(insertLoginQuery, uuidToBinary(login.UserID), login.Nickname, login.Password)
		return err
	})
}

// Get returns login credentials for check login.
func (r *Login) Get(ctx context.Context, nickname string) (*domain.Login, error) {
	login := &domain.Login{Nickname: nickname}

	err := query(ctx, r.readPool, func(conn *client.Conn) error {
		result, err := conn.Execute(selectLoginQuery, login.Nickname)
		if err != nil {
			return err
		}
		defer result.Close()

		if len(result.Values) == 0 {
			return ErrNotFound
		}

		id, err := result.GetString(0, 0)
		if err != nil {
			return err
		}
		login.UserID = binaryToUUID(id)

		login.Password, err = result.GetString(0, 1)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return login, nil
}
