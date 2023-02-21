package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	// db object from main_test.go
	store := NewStore(testDB)

	// create 2 accounts
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	// create channels for sharing data amount go routines
	errs := make(chan error)
	results := make(chan TransferTxResult)

	// test TransferTx function and send objects to channels
	// test 2 transfers from account1 to account2
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			// send objects to channels
			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		// receiving object from channel
		err := <-errs
		require.NoError(t, err)

		// receiving object from channel
		result := <-results
		require.NotEmpty(t, result)

		// check expected vs actual
		// Transfer object
		transfer := result.Transfer
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// getting Transfer object from db
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// FromEntry object
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// ToEntry object
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO: check accounts' balance
	}

}
