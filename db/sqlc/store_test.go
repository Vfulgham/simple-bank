package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	// db object from main_test.go
	store := NewStore(testDB)

	// create 2 accounts
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">>before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions
	n := 2
	amount := int64(10)

	// create channels for sharing data amount go routines
	errs := make(chan error)
	results := make(chan TransferTxResult)

	// test TransferTx function and send objects to channels
	// test n transfers from account1 to account2
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

		// check each accounts' balance
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Println(">>tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance // money that goes out of account1
		diff2 := toAccount.Balance - account2.Balance // money that goes in to account2
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
	}

	// check final updated balances
	// using GetAccountForUpdate() for db transaction lock 
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">>after:", updatedAccount1.Balance, updatedAccount2.Balance)

}
