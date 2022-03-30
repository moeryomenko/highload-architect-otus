package repository

import (
	"context"
	"database/sql"
	b64 "encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/moeryomenko/highload-architect-otus/social/internal/domain"
)

const (
	insertProfileQueyr = `INSERT INTO profiles (id, first_name, last_name, age, gender, interests, city) VALUES (UNHEX(?), ?, ?, ?, ?, ?, ?)`
	paginatedListQuery = `SELECT first_name, last_name, age, gender, interests, city, created_at FROM profiles %s ORDER BY created_at DESC LIMIT ?`

	nextPage = `WHERE created_at < ?`
)

// Users incapsulates data access layer for user profiles.
type Users struct {
	conn *sql.DB
}

// NewUsers returns new instance of user repository.
func NewUsers(conn *sql.DB) *Users {
	return &Users{conn: conn}
}

// Save saves user profile to repository.
func (r *Users) Save(ctx context.Context, user *domain.User) error {
	return transaction(ctx, r.conn, func(ctx context.Context, tx *sql.Tx) (err error) {
		_, err = tx.ExecContext(
			ctx, insertProfileQueyr,
			uuidToBinary(user.ID),
			user.Info.FirstName, user.Info.LastName,
			user.Info.Age, user.Info.Gender,
			strings.Join(user.Info.Interests, ","),
			user.Info.City,
		)
		return err
	})
}

// List returns list of users and token for future iteration through profiles.
func (r *Users) List(ctx context.Context, opts ...PageOption) ([]domain.User, string, error) {
	queryBuilder := &pageQuery{pageSize: 10}
	for _, opt := range opts {
		opt(queryBuilder)
	}

	query, params := queryBuilder.getQuery()

	rows, err := r.conn.QueryContext(ctx, query, params...)
	defer rows.Close()

	switch err {
	case nil:
		var (
			page      = make([]domain.User, 0, queryBuilder.pageSize)
			nextToken time.Time
		)
		for rows.Next() {
			var (
				interests string
				user      = domain.User{Info: &domain.Profile{}}
			)
			err := rows.Scan(
				&user.Info.FirstName, &user.Info.LastName,
				&user.Info.Age, &user.Info.Gender,
				&interests, &user.Info.City,
				&nextToken,
			)
			if err != nil {
				continue
			}
			user.Info.Interests = strings.Split(interests, ",")

			page = append(page, user)
		}
		return page, b64.StdEncoding.EncodeToString([]byte(nextToken.String())), nil
	case sql.ErrNoRows:
		return nil, "", ErrNotFound
	default:
		return nil, "", err
	}
}

// WithPageSize sets size of page for list.
func WithPageSize(size int) PageOption {
	return func(pq *pageQuery) {
		pq.pageSize = size
	}
}

// WithPageAt sets from what time take profiles for List.
func WithPageAt(at time.Time) PageOption {
	return func(pq *pageQuery) {
		pq.at = &at
	}
}

// DecodeToken decodes token to time for WithPageAt.
func DecodeToken(token string) (time.Time, error) {
	dst, err := b64.StdEncoding.DecodeString(token)
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(string(dst), "2006-01-02 15:04:05")
}

// pageQuery is query builder for iterating in List.
type pageQuery struct {
	pageSize int
	at       *time.Time
}

// getQuery returns query specified by options.
func (pq *pageQuery) getQuery() (string, []any) {
	if pq.at == nil {
		return fmt.Sprintf(paginatedListQuery, ""), []any{pq.pageSize}
	}
	return fmt.Sprintf(paginatedListQuery, nextPage), []any{pq.pageSize, *pq.at}
}

type PageOption func(*pageQuery)
