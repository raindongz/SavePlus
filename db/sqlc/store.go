package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute all functions and db queries
type Store struct{
	*Queries
	conn *pgxpool.Pool
}

//create a new store
func NewStore(connPool *pgxpool.Pool) *Store{
	return &Store{
		conn: connPool,
		Queries: New(connPool),
	}
}

// execFunctionWithTransction execute a function within a database transaction
func (store *Store) execFunctionWithTransction (ctx context.Context, fc func (*Queries) error) error{
	tx, err := store.conn.Begin(ctx)
	if err != nil{
		return err
	}

	tran := New(tx)
	transactionErr := fc(tran)
	if transactionErr != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil{
			return fmt.Errorf("transaction error: %v, rollback error: %v", transactionErr, rollbackErr)
		}
		return err
	}
	return tx.Commit(ctx)
}