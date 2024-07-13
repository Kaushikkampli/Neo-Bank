package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/kaushikkampli/neobank/db/sqlc"
	"github.com/kaushikkampli/neobank/utils"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store) (Server, error) {

	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		TokenExpirationTime: time.Minute,
	}

	testServer, err := NewServer(config, store)
	require.NoError(t, err)

	return *testServer, nil
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
