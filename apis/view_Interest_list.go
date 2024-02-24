package apis

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nioliu/commons/log"
	db "github.com/randongz/save_plus/db/sqlc"
	"github.com/randongz/save_plus/token"
	"go.uber.org/zap"
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
	Price   string `json:"price"`
	Images  string `json:"images"`
}

func (server *Server) viewMyInterestList(ctx *gin.Context) {
	req := new(viewMyInterestReq)

	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "get user id from ctx failed", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}
	req.userID = userId

	rsp := new(viewMyInterestRsp)
	statusCode, err := server.doViewMyInterestList(ctx, req, rsp)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "do view my interest failed", zap.Error(err))
		rsp.Detail = err.Error()
	}
	ctx.JSON(statusCode, rsp)
}

func getUserIdFromCtx(ctx *gin.Context) (int64, error) {
	v, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		log.WarnWithCtxFields(ctx, "no id in token")

		return 0, errors.New("no id in token")
	}
	if payload, is := v.(*token.Payload); !is {
		log.WarnWithCtxFields(ctx, "unexpected payload in auth key")
		return 0, errors.New("unexpected payload in auth key")
	} else {
		uid, err := strconv.Atoi(payload.Uid)
		if err != nil {
			log.WarnWithCtxFields(ctx, "unexpected payload in auth key")
			return 0, errors.New("unexpected payload in auth key")
		}
		return int64(uid), nil
	}
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
func transInterestListOuterRsp(ctx context.Context, list []db.PostInfo) ([]InterestItem, error) {
	outerItems := make([]InterestItem, 0, len(list))
	for _, row := range list {
		i := InterestItem{
			PostId:  row.ID,
			Status:  row.PostStatus,
			Title:   row.Title,
			Content: row.Content,
			Images:  row.Images,
		}
		outerItems = append(outerItems, i)
	}
	return outerItems, nil
}
