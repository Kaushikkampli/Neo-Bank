package db

import (
	"context"
	"testing"

	"github.com/kaushikkampli/neobank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)
	require.NotEmpty(t, transfer)
}

func TestGetTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	transferFromDB, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transferFromDB)
	require.Equal(t, transfer.ID, transferFromDB.ID)
	require.Equal(t, transfer.FromAccountID, transferFromDB.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transferFromDB.ToAccountID)
}

func TestUpdateTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	arg := UpdateTransferParams{
		ID:     transfer.ID,
		Amount: utils.RandomMoney(),
	}

	updatedTransfer, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedTransfer)
	require.Equal(t, transfer.ID, updatedTransfer.ID)
	require.Equal(t, arg.Amount, updatedTransfer.Amount)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)

	transferFromDB, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.Error(t, err)
	require.Empty(t, transferFromDB)
}
