package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	mockdb "github.com/vlad19930514/webApp/db/mock"
	db "github.com/vlad19930514/webApp/db/sqlc"
	"github.com/vlad19930514/webApp/util"
	"go.uber.org/mock/gomock"
)

func TestGetUserAPI(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name string

		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusAccepted, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "NotFound",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			id, _ := user.ID.Value()
			url := fmt.Sprintf("/user/%v", id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

type eqCreateUserParamsMatcher struct {
	arg     db.CreateUserParams
	id      pgtype.UUID
	created pgtype.Timestamptz
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	// TODO добавить сравнения
	//пример user.Created.Time.Equal(gotUser.Created.Time)
	e.arg.Created = arg.Created
	e.arg.ID = arg.ID
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v, id %v,created %v", e.arg, e.id, e.created)
}

func EqCreateUserParams(arg db.CreateUserParams, id pgtype.UUID, created pgtype.Timestamptz) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, id, created}
}

func TestCreateUserAPI(t *testing.T) {
	user := randomUser()
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"firstname": user.Firstname,
				"lastname":  user.Lastname,
				"email":     user.Email,
				"age":       user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Firstname: user.Firstname,
					Lastname:  user.Lastname,
					Email:     user.Email,
					Age:       user.Age,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, user.ID, user.Created)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusAccepted, recorder.Code)
				//fmt.Println(recorder.Body)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"firstname": user.Firstname,
				"lastname":  user.Lastname,
				"email":     "invalid-email",
				"age":       user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Email", "Это не email - invalid-email?")
			},
		},
		{
			name: "InvalidFirstname",
			body: gin.H{
				"firstname": "1234", // First name must be alpha
				"lastname":  user.Lastname,
				"email":     user.Email,
				"age":       user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Firstname", "Передаем только буквы - 1234")
			},
		},
		{
			name: "InvalidLastname",
			body: gin.H{
				"firstname": user.Firstname,
				"lastname":  "1234", // Last name must be alpha
				"email":     user.Email,
				"age":       user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Lastname", "Передаем только буквы - 1234")
			},
		},
		{
			name: "InvalidAge",
			body: gin.H{
				"firstname": user.Firstname,
				"lastname":  user.Lastname,
				"email":     user.Email,
				"age":       200, // Age must be between 1 and 130
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Age", "Многовато будет - 200")
			},
		}, {
			name: "MissingFirstname",
			body: gin.H{
				"lastname": user.Lastname,
				"email":    user.Email,
				"age":      user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Firstname", "Это обязательное поле")
			},
		},
		{
			name: "MissingLastname",
			body: gin.H{
				"firstname": user.Firstname,
				"email":     user.Email,
				"age":       user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Lastname", "Это обязательное поле")
			},
		},
		{
			name: "MissingEmail",
			body: gin.H{
				"firstname": user.Firstname,
				"lastname":  user.Lastname,
				"age":       user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Email", "Это обязательное поле")
			},
		},
		{
			name: "MissingAge",
			body: gin.H{
				"firstname": user.Firstname,
				"lastname":  user.Lastname,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Age", "Это обязательное поле")
			},
		},
		/* 		{// TODO ругается на парсинг json
			name: "EmptyBody",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Firstname", "Это обязательное поле")
				requireBodyContainsErrorMessage(t, recorder.Body, "Lastname", "Это обязательное поле")
				requireBodyContainsErrorMessage(t, recorder.Body, "Email", "Это обязательное поле")
				requireBodyContainsErrorMessage(t, recorder.Body, "Age", "Это обязательное поле")
			},
		}, */
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/user"
			// Marshall body data to JSON
			data, err := json.Marshal(tc.body)
			fmt.Println(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateUserAPI(t *testing.T) {
	user := randomUser()
	useFirstName := randomUser()
	userLastNameCase := randomUser() //TODO хардкод юзеров, иначе кейсы перезаписывают user
	randomName := util.RandomName()
	randomLastname := util.RandomName()
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ChangeFirstName",
			body: gin.H{
				"id":        useFirstName.ID,
				"firstname": randomName,
				"lastname":  useFirstName.Lastname,
				"email":     useFirstName.Email,
				"age":       useFirstName.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.UpdateUserParams{
					ID:        useFirstName.ID,
					Firstname: randomName,
					Lastname:  useFirstName.Lastname,
					Email:     useFirstName.Email,
					Age:       useFirstName.Age,
				}
				useFirstName.Firstname = randomName
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(useFirstName, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				requireBodyMatchUser(t, recorder.Body, db.User(useFirstName))
			},
		},
		{
			name: "ChangeLastName",
			body: gin.H{
				"id":        userLastNameCase.ID,
				"firstname": userLastNameCase.Firstname,
				"lastname":  randomLastname,
				"email":     userLastNameCase.Email,
				"age":       userLastNameCase.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.UpdateUserParams{
					ID:        userLastNameCase.ID,
					Firstname: userLastNameCase.Firstname,
					Lastname:  randomLastname,
					Email:     userLastNameCase.Email,
					Age:       userLastNameCase.Age,
				}
				userLastNameCase.Lastname = randomLastname
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(userLastNameCase, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				requireBodyMatchUser(t, recorder.Body, userLastNameCase)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"firstname": user.Firstname,
				"lastname":  user.Lastname,
				"email":     "invalid-email",
				"age":       user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Email", "Это не email - invalid-email?")
			},
		},
		{
			name: "InvalidFirstname",
			body: gin.H{
				"firstname": "1234", // First name must be alpha
				"lastname":  user.Lastname,
				"email":     user.Email,
				"age":       user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Firstname", "Передаем только буквы - 1234")
			},
		},
		{
			name: "InvalidLastname",
			body: gin.H{
				"firstname": user.Firstname,
				"lastname":  "1234", // Last name must be alpha
				"email":     user.Email,
				"age":       user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Lastname", "Передаем только буквы - 1234")
			},
		},
		{
			name: "InvalidAge",
			body: gin.H{
				"firstname": user.Firstname,
				"lastname":  user.Lastname,
				"email":     user.Email,
				"age":       200, // Age must be between 1 and 130
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Age", "Многовато будет - 200")
			},
		}, {
			name: "MissingFirstname",
			body: gin.H{
				"lastname": user.Lastname,
				"email":    user.Email,
				"age":      user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Firstname", "Это обязательное поле")
			},
		},
		{
			name: "MissingLastname",
			body: gin.H{
				"firstname": user.Firstname,
				"email":     user.Email,
				"age":       user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Lastname", "Это обязательное поле")
			},
		},
		{
			name: "MissingEmail",
			body: gin.H{
				"firstname": user.Firstname,
				"lastname":  user.Lastname,
				"age":       user.Age,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Email", "Это обязательное поле")
			},
		},
		{
			name: "MissingAge",
			body: gin.H{
				"firstname": user.Firstname,
				"lastname":  user.Lastname,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyContainsErrorMessage(t, recorder.Body, "Age", "Это обязательное поле")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/user"
			// Marshall body data to JSON
			data, err := json.Marshal(tc.body)
			fmt.Println(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
func randomUser() db.User {
	generatedUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("Ошибка при создании UUID: %v", err)
	}
	//TODO sqlc генерирует pgtype
	// Приведение сгенерированного UUID к типу pgtype.UUID
	dbUUID := pgtype.UUID{}
	copy(dbUUID.Bytes[:], generatedUUID[:]) // Устанавливаем bytes field
	dbUUID.Valid = true                     // Устанавливаем valid field

	id, _ := dbUUID.UUIDValue()
	return db.User{
		ID:        id,
		Firstname: util.RandomName(),
		Lastname:  util.RandomName(),
		Email:     util.RandomEmail(),
		Age:       util.RandomAge(),
		Created:   util.PgtypeCurrentTime(),
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Firstname, gotUser.Firstname)
	require.Equal(t, user.Lastname, gotUser.Lastname)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.Age, gotUser.Age)
	require.Equal(t, user.ID, gotUser.ID)
	require.True(t, user.Created.Time.Equal(gotUser.Created.Time))

}
func requireBodyContainsErrorMessage(t *testing.T, body *bytes.Buffer, field, expectedMessage string) {
	fmt.Println(body)
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var response struct {
		Errors []util.ErrorMsg `json:"errors"`
	}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)
	for _, e := range response.Errors {
		if e.Field == field {
			require.Equal(t, expectedMessage, e.Message)
			return
		}
	}
	require.Fail(t, fmt.Sprintf("no error message found for field %s", field))
}
