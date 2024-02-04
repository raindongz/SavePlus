package apis

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nioliu/commons/log"
	db "github.com/randongz/save_plus/db/sqlc"
	"github.com/randongz/save_plus/token"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type viewMyInterestReq struct {
	userID int64
}

type viewMyInterestRsp struct {
	Items  []InterestItem `json:"items,omitempty"`
	Detail string         `json:"detail,omitempty"`
}
type InterestItem struct {
	PostId  int64  `json:"postId"`
	Status  int16  `json:"status"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (server *Server) viewMyInterestList(ctx *gin.Context) {
	req := new(viewMyInterestReq)

	v, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		log.WarnWithCtxFields(ctx, "no id in token")
		ctx.Status(http.StatusBadRequest)
		return
	}
	if payload, is := v.(*token.Payload); !is {
		log.WarnWithCtxFields(ctx, "unexpected payload in auth key")
		ctx.Status(http.StatusBadRequest)
		return
	} else {
		uid, err := strconv.Atoi(payload.Uid)
		if err != nil {
			log.WarnWithCtxFields(ctx, "unexpected payload in auth key")
			ctx.Status(http.StatusBadRequest)
			return
		}
		req.userID = int64(uid)
	}

	rsp := new(viewMyInterestRsp)
	statusCode, err := server.doViewMyInterestList(ctx, req, rsp)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "do view my interest failed", zap.Error(err))
		rsp.Detail = err.Error()
	}
	ctx.JSON(statusCode, rsp)
}

func (server *Server) doViewMyInterestList(ctx context.Context,
	req *viewMyInterestReq, rsp *viewMyInterestRsp) (int, error) {
	// 1. get information todo no limit?
	interestListByUserID, err := server.store.GetInterestListByUserID(ctx, req.userID)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "get interest list by user id failed", zap.Error(err))
		return http.StatusInternalServerError, errors.New("operation failed")
	}
	items, err := transInterestListOuterRsp(ctx, interestListByUserID)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "")
	}
	rsp.Items = items
	return http.StatusOK, nil
}

// trans interest list to outer rsp
func transInterestListOuterRsp(ctx context.Context, list []db.GetInterestListByUserIDRow) ([]InterestItem, error) {
	outerItems := make([]InterestItem, 0, len(list))
	for _, row := range list {
		i := InterestItem{
			PostId:  row.PostID,
			Status:  row.PostStatus,
			Title:   row.Title,
			Content: row.Content,
		}
		outerItems = append(outerItems, i)
	}
	return outerItems, nil
}
