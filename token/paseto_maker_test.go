package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/kaushikkampli/neobank/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoMakerHappyPath(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	jwt, err := maker.CreateToken(username, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, jwt)

	payload, err := maker.ValidateToken(jwt)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
}

func TestPasetoMakerErrorInvalidSize(t *testing.T) {
	_, err := NewPasetoMaker(utils.RandomString(2))
	require.Error(t, err, fmt.Errorf("invalid key size: min length %v", minKeySize))
}

func TestPasetoMakerExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	jwt, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, jwt)

	_, err = maker.ValidateToken(jwt)
	require.Error(t, err, ErrorExpiredToken)
}
