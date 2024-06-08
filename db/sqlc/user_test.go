package db

import (
	"context"
	"testing"
	"time"

	"github.com/kaushikkampli/neobank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:     utils.RandomOwner(),
		HashedPasswd: "passwd",
		FullName:     utils.RandomOwner(),
		EmailID:      utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	if err != nil {
		t.Error(err)
	}

	require.NotEmpty(t, user)
	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.EmailID, arg.EmailID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	userFromDB, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, userFromDB)
	require.Equal(t, user.Username, userFromDB.Username)
	require.Equal(t, user.FullName, userFromDB.FullName)
	require.Equal(t, user.EmailID, userFromDB.EmailID)
	require.Equal(t, user.HashedPasswd, userFromDB.HashedPasswd)
	require.WithinDuration(t, user.CreatedAt, userFromDB.CreatedAt, time.Second)
}
