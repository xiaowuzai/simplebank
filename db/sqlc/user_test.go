package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xiaowuzai/simplebank/util"
)

// 向数据库中随机添加用户
func createRandomUser(t *testing.T) User {
	password, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: password,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	fullName := util.RandomString(6)
	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: fullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateUser)
	require.Equal(t, fullName, updateUser.FullName)
	require.Equal(t, oldUser.Email, updateUser.Email)
	require.Equal(t, oldUser.HashedPassword, updateUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	email := util.RandomEmail()
	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateUser)
	require.Equal(t, email, updateUser.Email)
	require.Equal(t, oldUser.HashedPassword, updateUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updateUser.FullName)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	password := util.RandomString(6)
	hashPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: hashPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateUser)
	require.Equal(t, hashPassword, updateUser.HashedPassword)
	require.Equal(t, oldUser.Email, updateUser.Email)
	require.Equal(t, oldUser.FullName, updateUser.FullName)
}

func TestUpdateUserAllField(t *testing.T) {
	oldUser := createRandomUser(t)

	fullName := util.RandomString(6)
	email := util.RandomEmail()
	hashPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: hashPassword,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: fullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateUser)
	require.Equal(t, hashPassword, updateUser.HashedPassword)
	require.Equal(t, email, updateUser.Email)
	require.Equal(t, fullName, updateUser.FullName)
}
