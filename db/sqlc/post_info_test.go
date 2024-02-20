package db

import (
	"context"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/randongz/save_plus/utils"
	"github.com/stretchr/testify/require"
)

func TestCreatNewPost(t *testing.T) {
	randomTitle := utils.RandomStringWithSpecifiedLenth(utils.RandomInt(0, 70))
	randomContent := utils.RandomStringWithSpecifiedLenth(utils.RandomInt(0, 2048))
	randomPrice := utils.RandomInt(0, 100)
	randomItemNum := utils.RandomInt(0, 100)
	randomImages := utils.RandomStringWithSpecifiedLenth(512)
	arg := CreateNewPostParams{
		Title:        randomTitle,
		Content:      randomContent,
		TotalPrice:   strconv.Itoa(int(randomPrice)),
		DeliveryType: 0,
		Area: pgtype.Text{
			String: "boston",
			Valid:  true,
		},
		ItemNum:    int32(randomItemNum),
		PostStatus: 0,
		Negotiable: 0,
		Images:     randomImages,
	}

	post, err := testQueries.CreateNewPost(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, post)
}
