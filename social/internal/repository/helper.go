package repository

import (
	"context"
	"strings"

	"github.com/go-mysql-org/go-mysql/client"
	uuid "github.com/satori/go.uuid"
)

func transaction(ctx context.Context, pool *client.Pool, tx func(conn *client.Conn) error) (err error) {
	return query(ctx, pool, func(conn *client.Conn) (err error) {
		err = conn.Begin()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				_ = conn.Rollback()
			}
		}()

		err = tx(conn)
		if err != nil {
			return err
		}
		return conn.Commit()
	})
}

func query(ctx context.Context, pool *client.Pool, fn func(conn *client.Conn) error) error {
	conn, err := pool.GetConn(ctx)
	if err != nil {
		return err
	}
	defer pool.PutConn(conn)

	return fn(conn)
}

func uuidToBinary(id uuid.UUID) string {
	return strings.ReplaceAll(id.String(), "-", "")
}

func binaryToUUID(id string) uuid.UUID {
	userID, _ := uuid.FromString(mapToUUID(id))
	return userID
}

func mapToUUID(id string) string {
	return id[:8] + "-" + id[8:12] + "-" + id[12:16] + "-" + id[16:20] + "-" + id[20:]
}
