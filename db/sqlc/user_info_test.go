package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/randongz/save_plus/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) UsersInfo {
	hashedPassword, err := utils.HashPassword(utils.RandomStringWithSpecifiedLenth(8))
	require.NoError(t, err)

	randomUserName := utils.RandomStringWithSpecifiedLenth(10)
	randomFullName := utils.RandomStringWithSpecifiedLenth(10)
	randomPhoneNumber := utils.RandomPhoneNumber();
	randomAvatar := utils.RandomAvatar()


	arg := CreateNewUserParams{
		Username: randomUserName,
		HashedPassword: hashedPassword,
		FullName: randomFullName,
		Email: utils.RandomEmail(),
		Phone: pgtype.Text{
			String: randomPhoneNumber,
			Valid: true,
		},
		Gender: 1,
		Avatar: pgtype.Text{
			String: randomAvatar,
			Valid: true,
		},
	}

	user, err := testQueries.CreateNewUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.Time.IsZero())

	return user
} 