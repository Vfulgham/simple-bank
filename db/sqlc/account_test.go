package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Vfulgham/simple-bank/db/util"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	// expected
	arg := CreateAccountParams{
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	// this is testQueries object created in main_test.go
	// actual
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	// check expected vs. actual
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	// check actual values are returned and they are not 0
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

// you need a Queries object
// and connection to db to test, created in main_test.go
func TestCreateAccount(t *testing.T){
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T){

	// create account1 in db
	account1 := createRandomAccount(t)

	// get account from db
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	// test expections (account1) with actual (account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	// compare within seconds
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
		// create account1 in db
		account1 := createRandomAccount(t)

		arg := UpdateAccountParams{
			ID: account1.ID,
			Balance: util.RandomMoney(),
		}

		// get account from db
		account2, err := testQueries.UpdateAccount(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, account2)
	
		// test expections (account1) with actual (account2)
		require.Equal(t, account1.ID, account2.ID)
		require.Equal(t, account1.Owner, account2.Owner)
		require.Equal(t, arg.Balance, account2.Balance) // new balance
		require.Equal(t, account1.Currency, account2.Currency)
	
		// compare within seconds
		require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T){
			// create account1 in db
			account1 := createRandomAccount(t)
	
			// get account from db
			err := testQueries.DeleteAccount(context.Background(), account1.ID)
			require.NoError(t, err)

			account2, err := testQueries.GetAccount(context.Background(), account1.ID)
			require.Error(t, err)
			require.EqualError(t, err, sql.ErrNoRows.Error()) // check error is equal to no rows
			require.Empty(t, account2)
}

func TestListAccounts(t *testing.T){

	// create accounts in db
	for i :=0; i < 10; i++{
		createRandomAccount(t)
	}

	// skip first 5 records and return next 5
	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5) // slice length

	for _, account := range accounts{
		require.NotEmpty(t, account)
	}
}

