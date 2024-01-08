package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/util"

	mockdb "github.com/xiaowuzai/simplebank/db/mock"

	"go.uber.org/mock/gomock"
)

var _ gomock.Matcher = (*eqCreateUserMatcher)(nil)

type eqCreateUserMatcher struct {
	arg      db.CreateUserParams //
	password string              // 未加密的
}

func (e eqCreateUserMatcher) Matches(x any) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserMatcher{
		arg:      arg,
		password: password,
	}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser()

	// hashPassword, err := util.HashPassword(password)
	// require.NoError(t, err)

	cases := []struct {
		name          string // test name
		body          gin.H  // 方便 Http 调用
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"fullName": user.FullName,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		// {
		// 	name: "InternalServerError",
		// 	body: gin.H{
		// 		"owner":    account.Owner,
		// 		"currency": account.Currency,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateAccount(gomock.Any(), gomock.Any()).
		// 			Times(1).
		// 			Return(db.Account{}, sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "BadRequest",
		// 	body: gin.H{
		// 		"owner":    account.Owner,
		// 		"currency": account.Currency,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		arg := db.CreateAccountParams{
		// 			Owner:    account.Owner,
		// 			Currency: account.Currency,
		// 			Balance:  0,
		// 		}
		// 		store.EXPECT().
		// 			CreateAccount(gomock.Any(), arg).
		// 			Times(1).
		// 			Return(db.Account{}, &pq.Error{Code: pq.ErrorCode("23505")})
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "InvalidId",
		// 	body: gin.H{
		// 		"owner": account.Owner,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateAccount(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
	}

	for i := range cases {
		tc := cases[i]

		t.Run(tc.name, func(t *testing.T) {
			// gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// mock db store
			store := mockdb.NewMockStore(ctrl)
			// run case
			tc.buildStubs(store)

			// new http server
			testServer := NewServer(store)
			recorder := httptest.NewRecorder()

			// marshal body
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			// run api test
			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			testServer.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(recorder)

		})
	}

}

func randomUser() (db.User, string) {
	user := db.User{
		Username: util.RandomOwner(),
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
	}

	return user, util.RandomString(6)
}
func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)

}
