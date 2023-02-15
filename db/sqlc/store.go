package db

import (
	"fmt"
	"context"
	"database/sql"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries // to access its functions
	db *sql.DB // this db object is required to create a new db transaction
}

// intakes a db object and returns a Store object
func (s *Store) NewStore (db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: NewQueries(db),
	}
}

// execTx executes a function within a database transaction
// generic so it can be reused in multiple transactions
func (store *Store) execTx(ctx context.Context, funcName func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := NewQueries(tx)
	err = funcName(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// contains the input parameters for the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update 
// accounts' balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error){
	var result TransferTxResult
	
	err := store.execTx(ctx, func(q *Queries) error{
		return nil
	})

	return result, err
}