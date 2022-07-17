package repository

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/moeryomenko/highload-architect-otus/social/internal/domain"
)

const (
	insertProfileQueyr = `INSERT INTO profiles (id, first_name, last_name, age, gender, interests, city) VALUES (UNHEX(?), ?, ?, ?, ?, ?, ?)`
	paginatedListQuery = `SELECT first_name, last_name, age, gender, interests, city, created_at FROM profiles %s ORDER BY created_at DESC LIMIT %d`

	searchSubstring = ` WHERE %s `
	nextPage        = ` created_at < '%s' `
)

// Users incapsulates data access layer for user profiles.
type Users struct {
	writePool *client.Pool
	readPool  *client.Pool
}

// NewUsers returns new instance of user repository.
func NewUsers(writePool, readPool *client.Pool) *Users {
	return &Users{writePool: writePool, readPool: readPool}
}

// Save saves user profile to repository.
func (r *Users) Save(ctx context.Context, user *domain.User) error {
	return transaction(ctx, r.writePool, func(conn *client.Conn) (err error) {
		_, err = conn.Execute(
			insertProfileQueyr,
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

	listQuery := queryBuilder.getQuery()

	page := make([]domain.User, 0, queryBuilder.pageSize)
	var nextToken string
	err := query(ctx, r.readPool, func(conn *client.Conn) error {
		var result mysql.Result
		defer result.Close()
		err := conn.ExecuteSelectStreaming(listQuery, &result, func(row []mysql.FieldValue) error {
			page = append(page, domain.User{
				Info: &domain.Profile{
					FirstName: string(row[0].AsString()),
					LastName:  string(row[1].AsString()),
					Age:       int(row[2].AsInt64()),
					Gender:    domain.Gender(row[3].AsString()),
					Interests: strings.Split(string(row[4].AsString()), ","),
					City:      string(row[5].AsString()),
				},
			})
			nextToken = string(row[6].AsString())
			return nil
		}, nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}

	return page, b64.URLEncoding.EncodeToString([]byte(nextToken)), nil
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

// WithSearchByFirstName sets search by first name predicate.
func WithSearchByFirstName(name string) PageOption {
	return func(pq *pageQuery) {
		pq.searchPredicates = append(pq.searchPredicates, predicate{"first_name", name})
	}
}

// WithSearchByLastName sets search by last name predicate.
func WithSearchByLastName(name string) PageOption {
	return func(pq *pageQuery) {
		pq.searchPredicates = append(pq.searchPredicates, predicate{"last_name", name})
	}
}

// DecodeToken decodes token to time for WithPageAt.
func DecodeToken(token string) (time.Time, error) {
	dst, err := b64.URLEncoding.DecodeString(token)
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse("2006-01-02 15:04:05", string(dst))
}

// pageQuery is query builder for iterating in List.
type pageQuery struct {
	pageSize int
	at       *time.Time
	// searchPredicates used for search profiles by given attributes.
	searchPredicates []predicate
}

// getQuery returns query specified by options.
func (pq *pageQuery) getQuery() string {
	if len(pq.searchPredicates) == 0 && pq.at == nil {
		return fmt.Sprintf(paginatedListQuery, "", pq.pageSize)
	}

	predicates := make([]string, 0, 2)

	query := searchQuery(pq.searchPredicates)
	if query != "" {
		predicates = append(predicates, query)
	}
	if pq.at != nil {
		predicates = append(predicates, fmt.Sprintf(nextPage, pq.at.Format("2006-01-02 15:04:05")))
	}

	query = strings.Join(predicates, "AND")

	return fmt.Sprintf(paginatedListQuery, fmt.Sprintf(searchSubstring, query), pq.pageSize)
}

func searchQuery(searchPredicates []predicate) string {
	searchParams := []string{}
	for _, predict := range searchPredicates {
		searchParams = append(searchParams, fmt.Sprintf(" %s LIKE '%v' ", predict.column, predict.value))
	}
	return strings.Join(searchParams, "AND")
}

// predicate is search query builder.
type predicate struct {
	column string
	value  any
}

type PageOption func(*pageQuery)
