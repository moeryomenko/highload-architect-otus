package repository

import (
	"context"
	"database/sql"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func transaction(ctx context.Context, conn *sql.DB, trx func(context.Context, *sql.Tx) error) (err error) {
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	err = trx(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
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
