package apis

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nioliu/commons/log"
	"github.com/randongz/save_plus/token"
	"go.uber.org/zap"
)

type PostItem struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Images     string `json:"images"`
	TotalPrice string `json:"total_price"`
	Area       string `json:"area"`
}
type MyPostHistoryResponse struct {
	PostItems []PostItem `json:"post_list,omitempty"`
}

func authenticationAndGetUserId(ctx *gin.Context) (int64, error) {
	// 1. check if token exist in context
	variable, exist := ctx.Get(authorizationPayloadKey)
	if !exist {
		log.WarnWithCtxFields(ctx, "paylaod doesn't exist in context")
		return 0, errors.New("paylaod doesn't exist in context")
	}
	// 2. check if token is right type
	if payload, is := variable.(*token.Payload); !is {
		log.WarnWithCtxFields(ctx, "unexpected format payload")
		return 0, errors.New("unexpected format payload")
	} else {
		// 3. check if userId can be extracted from payload
		userId, err := strconv.Atoi(payload.Uid)
		if err != nil {
			log.WarnWithCtxFields(ctx, "convert payload uid failed")
			return 0, errors.New("convertPayload uid failed")
		}
		return int64(userId), nil
	}
}

func (server *Server) getMyPostHistory(ctx *gin.Context) {
	// 1. authentication
	userId, err := authenticationAndGetUserId(ctx)
	if err != nil || userId == 0 {
		log.ErrorWithCtxFields(ctx, "authentication failed", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}
	// 2. get my post history from db
	postList, err := server.store.GetMyPostList(ctx, userId)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "get postList from db failed", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// 3. construct response and return
	rsp := new(MyPostHistoryResponse)
	for i := range postList {
		postItem := PostItem{
			ID:         postList[i].ID,
			Title:      postList[i].Title,
			Content:    postList[i].Content,
			Images:     postList[i].Images,
			TotalPrice: postList[i].TotalPrice,
			Area:       postList[i].Area.String,
		}
		rsp.PostItems = append(rsp.PostItems, postItem)
	}
	ctx.JSON(http.StatusOK, rsp)
}
