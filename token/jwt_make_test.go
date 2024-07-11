package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kaushikkampli/neobank/utils"
	"github.com/stretchr/testify/require"
)

func TestJWTMakerHappyPath(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	jwt, err := maker.CreateToken(username, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, jwt)

	payload, err := maker.ValidateToken(jwt)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
}

func TestJWTMakerErrorInvalidSize(t *testing.T) {
	_, err := NewJWTMaker(utils.RandomString(2))
	require.Error(t, err, fmt.Errorf("invalid key size: min length %v", minKeySize))
}

func TestJWTMakerExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	jwt, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, jwt)

	_, err = maker.ValidateToken(jwt)
	require.Error(t, err, ErrorExpiredToken)
}

func TestJWTMakerHeadderTamperedAttack(t *testing.T) {
	secret := utils.RandomString(32)
	maker, err := NewJWTMaker(secret)
	require.NoError(t, err)

	username := utils.RandomOwner()

	payload := NewPayload(username, time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	jwtToken, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = maker.ValidateToken(jwtToken)
	require.Error(t, err, ErrorInvalidToken)
}
