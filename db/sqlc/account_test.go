package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kaushikkampli/neobank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	CreateAccountDetails := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), CreateAccountDetails)
	require.NoError(t, err)

	require.NotEmpty(t, account)
	require.Equal(t, CreateAccountDetails.Owner, account.Owner)
	require.Equal(t, CreateAccountDetails.Balance, account.Balance)
	require.Equal(t, CreateAccountDetails.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)

	accountFromDB, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountFromDB)
	require.Equal(t, account.ID, accountFromDB.ID)
	require.Equal(t, account.Owner, accountFromDB.Owner)
	require.Equal(t, account.Balance, accountFromDB.Balance)
	require.Equal(t, account.Currency, accountFromDB.Currency)
	require.WithinDuration(t, account.CreatedAt, accountFromDB.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	updateAccountDetails := UpdateAccountParams{
		ID:      account.ID,
		Balance: utils.RandomMoney(),
	}

	account, err := testQueries.UpdateAccount(context.Background(), updateAccountDetails)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, updateAccountDetails.ID, account.ID)
	require.Equal(t, updateAccountDetails.Balance, account.Balance)
}

func TestDeleteAccount(t *testing.T) {
	accountToBeDeleted := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), accountToBeDeleted.ID)

	require.NoError(t, err)
	accountToBeDeleted, err = testQueries.GetAccount(context.Background(), accountToBeDeleted.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountToBeDeleted)
}
