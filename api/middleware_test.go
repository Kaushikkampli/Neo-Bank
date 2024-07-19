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

func SetAuthHeaderForTest(t *testing.T, username string, duration time.Duration, authType string, token token.Maker, req *http.Request) {
	authValue, err := token.CreateToken(username, duration)
	require.NoError(t, err)
	req.Header.Set(AuthorizationHeader, authType+" "+authValue)
}

func TestAuthMiddleWare(t *testing.T) {
	testCases := []struct {
		name          string
		setAuth       func(t *testing.T, req *http.Request, token token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "HappyPath",
			setAuth: func(t *testing.T, req *http.Request, token token.Maker) {
				SetAuthHeaderForTest(t, "user", time.Minute, AuthorizationType, token, req)
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
				SetAuthHeaderForTest(t, "user", time.Minute, utils.RandomString(6), token, req)
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
