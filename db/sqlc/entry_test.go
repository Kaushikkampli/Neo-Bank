package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kaushikkampli/neobank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	CreateEntryDetails := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), CreateEntryDetails)

	if err != nil {
		t.Error(err)
	}

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, CreateEntryDetails.AccountID, entry.AccountID)
	require.Equal(t, CreateEntryDetails.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)

	entryFromDB, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entryFromDB)
	require.Equal(t, entry.ID, entryFromDB.ID)
	require.Equal(t, entry.AccountID, entryFromDB.AccountID)
	require.Equal(t, entry.Amount, entryFromDB.Amount)
	require.WithinDuration(t, entry.CreatedAt, entryFromDB.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)

	updateEntryDetails := UpdateEntryParams{
		ID:     entry.ID,
		Amount: utils.RandomMoney(),
	}

	entry, err := testQueries.UpdateEntry(context.Background(), updateEntryDetails)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, updateEntryDetails.ID, entry.ID)
	require.Equal(t, updateEntryDetails.Amount, entry.Amount)

}

func TestDeleteEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	entry, err = testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry)
}
