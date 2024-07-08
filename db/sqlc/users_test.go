package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/vlad19930514/webApp/util"
)

func createRandomUser() (User, CreateUserParams, error) {
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

	arg := CreateUserParams{
		ID:        id,
		Firstname: util.RandomName(),
		Lastname:  util.RandomName(),
		Email:     util.RandomEmail(),
		Age:       util.RandomAge(),
		Created:   util.PgtypeCurrentTime(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)

	return user, arg, err
}
func TestCreateUser(t *testing.T) {
	user, arg, err := createRandomUser()
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Firstname, user.Firstname)
	require.Equal(t, arg.Lastname, user.Lastname)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Age, user.Age)
	require.Equal(t, arg.ID, user.ID)
	require.NotZero(t, user.Created)
}

func TestGetUser(t *testing.T) {
	user1, _, err := createRandomUser()
	require.NoError(t, err)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Firstname, user2.Firstname)
	require.Equal(t, user1.Lastname, user2.Lastname)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Age, user2.Age)
	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Created, user2.Created)

	require.WithinDuration(t, user1.Created.Time, user1.Created.Time, time.Second)
}

func TestUpdateUser(t *testing.T) {
	createdUser, _, _ := createRandomUser()
	arg := UpdateUserParams{
		ID:        createdUser.ID,
		Firstname: util.RandomName(),
		Lastname:  util.RandomName(),
		Email:     util.RandomEmail(),
		Age:       util.RandomAge(),
	}
	updatedUser, _ := testQueries.UpdateUser(context.Background(), arg)
	require.NotEqual(t, createdUser.Firstname, updatedUser.Firstname)
	currentUser, _ := testQueries.GetUser(context.Background(), createdUser.ID)
	require.Equal(t, arg.Firstname, currentUser.Firstname)
	require.Equal(t, arg.Lastname, currentUser.Lastname)
	require.Equal(t, arg.Email, currentUser.Email)
	require.Equal(t, arg.Age, currentUser.Age)
}
