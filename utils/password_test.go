package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	input := RandomString(9)

	hash, err := HashPassword(input)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	err = ComparePassword(input, hash)
	require.NoError(t, err)

	err = ComparePassword(RandomString(9), hash)
	require.Error(t, err, bcrypt.ErrMismatchedHashAndPassword)

	hash2, err := HashPassword(hash)
	require.NoError(t, err)
	require.NotEqual(t, hash, hash2)
}
