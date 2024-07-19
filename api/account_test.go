package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/kaushikkampli/neobank/db/mock"
	db "github.com/kaushikkampli/neobank/db/sqlc"
	"github.com/kaushikkampli/neobank/utils"
	"github.com/stretchr/testify/require"
)

var (
	account = randomAccount()
)

func TestGetAccount(t *testing.T) {

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "Happy Path",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatch(t, recorder.Body, account)
			},
		},
		{
			name:      "UnHappy Path",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal Error",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Invalid Input",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		ctrl := gomock.NewController(t)

		store := mockdb.NewMockStore(ctrl)

		tc.buildStubs(store)

		server, err := NewTestServer(t, store)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()

		url := fmt.Sprintf("/accounts/%d", tc.accountID)
		request, err := http.NewRequest(http.MethodGet, url, nil)

		SetAuthHeaderForTest(t, account.Owner, time.Minute, AuthorizationType, server.token, request)

		require.NoError(t, err)
		server.router.ServeHTTP(recorder, request)

		tc.checkResponse(t, recorder)
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

func requireBodyMatch(t *testing.T, body *bytes.Buffer, createdAccount db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var result db.Account
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	require.Equal(t, createdAccount, result)
}
