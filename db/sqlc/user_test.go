package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rizkiromadoni/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username: util.RandomOwner(),
		Email:    util.RandomEmail(),
		FullName: util.RandomOwner(),
		Password: hashedPassword,
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{
			String: util.RandomOwner(),
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.NotEqual(t, oldUser.FullName, newUser.FullName)
	require.WithinDuration(t, oldUser.PasswordChangedAt, newUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, oldUser.CreatedAt, newUser.CreatedAt, time.Second)
}

func TestUpdateUserEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: pgtype.Text{
			String: util.RandomEmail(),
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.NotEqual(t, oldUser.Email, newUser.Email)
	require.WithinDuration(t, oldUser.PasswordChangedAt, newUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, oldUser.CreatedAt, newUser.CreatedAt, time.Second)
}

func TestUpdateUserPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Password: pgtype.Text{
			String: util.RandomString(6),
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.NotEqual(t, oldUser.Password, newUser.Password)
	require.WithinDuration(t, oldUser.PasswordChangedAt, newUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, oldUser.CreatedAt, newUser.CreatedAt, time.Second)
}
