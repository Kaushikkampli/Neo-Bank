package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/kaushikkampli/neobank/db/mock"
	db "github.com/kaushikkampli/neobank/db/sqlc"
	"github.com/kaushikkampli/neobank/utils"
	"github.com/stretchr/testify/require"
)

type eqCreateUserRequestMatcher struct {
	arg    db.CreateUserParams
	passwd string
}

func (r eqCreateUserRequestMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.ComparePassword(r.passwd, arg.HashedPasswd)
	return err == nil
}

func (r eqCreateUserRequestMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", r.arg, r.passwd)
}

func EqCreateUserParams(arg db.CreateUserParams, passwd string) gomock.Matcher {
	return eqCreateUserRequestMatcher{arg, passwd}
}

func TestCreateUser(t *testing.T) {

	randomUser, passwd := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Happy Path",
			body: gin.H{
				"username": randomUser.Username,
				"password": passwd,
				"fullname": randomUser.FullName,
				"email":    randomUser.EmailID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:     randomUser.Username,
					FullName:     randomUser.FullName,
					HashedPasswd: randomUser.HashedPasswd,
					EmailID:      randomUser.EmailID,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, passwd)).
					Times(1).
					Return(randomUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, randomUser)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)
			server, err := NewTestServer(t, store)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			url := "/users"

			jsonData, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, passwd string) {
	passwd = utils.RandomString(6)

	hash, err := utils.HashPassword(passwd)
	require.NoError(t, err)
	user = db.User{
		Username:     utils.RandomOwner(),
		FullName:     utils.RandomOwner(),
		EmailID:      utils.RandomEmail(),
		HashedPasswd: hash,
	}

	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var result db.User
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	require.Equal(t, user.Username, result.Username)
}
