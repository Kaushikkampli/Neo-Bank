package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kaushikkampli/neobank/token"
	"github.com/kaushikkampli/neobank/utils"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddleWare(t *testing.T) {
	testCases := []struct {
		name          string
		setAuth       func(t *testing.T, req *http.Request, token token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "HappyPath",
			setAuth: func(t *testing.T, req *http.Request, token token.Maker) {
				authValue, err := token.CreateToken("user", time.Minute)
				require.NoError(t, err)
				req.Header.Set(AuthorizationHeader, AuthorizationType+" "+authValue)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UnHappyPath - NoAuthorisation",
			setAuth: func(t *testing.T, req *http.Request, token token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnHappyPath - UnsupportedAuthorisationType",
			setAuth: func(t *testing.T, req *http.Request, token token.Maker) {
				authValue, err := token.CreateToken("user", time.Minute)
				require.NoError(t, err)
				req.Header.Set(AuthorizationHeader, utils.RandomString(6)+" "+authValue)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnHappyPath - MissingAuthorisationToken",
			setAuth: func(t *testing.T, req *http.Request, token token.Maker) {
				req.Header.Set(AuthorizationHeader, AuthorizationType)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {

		server, err := NewTestServer(t, nil)
		require.NoError(t, err)

		path := "/auth"
		server.router.GET(path,
			authMiddleware(server.token),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, nil)
			})

		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, path, nil)
		require.NoError(t, err)

		tc.setAuth(t, request, server.token)

		server.router.ServeHTTP(recorder, request)
		tc.checkResponse(t, recorder)
	}
}
