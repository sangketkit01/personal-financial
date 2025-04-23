package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) Store{
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error{
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil{
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil{
		if rbErr := tx.Rollback(ctx) ; rbErr != nil{
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit(ctx)
}