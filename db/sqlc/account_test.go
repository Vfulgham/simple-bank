package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// you need a Queries object
// and connection to db to test, created in main_test.go
func TestCreateAccount(t *testing.T){

	// expected
	arg := CreateAccountParams{
		Owner: "tom",
		Balance: 100,
		Currency: "USD",
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
}