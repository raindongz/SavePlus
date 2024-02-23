package apis

import (
	"net/http"
	"strconv"

	"github.com/nioliu/commons/log"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/randongz/save_plus/token"
)

type SoldItemRequest struct {
	PostId int64 `json:"postid" binding:"required,min=1"`
}

func (server *Server) soldItem(ctx *gin.Context) {
	//check request params
	req := new(SoldItemRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.ErrorWithCtxFields(ctx, "bind params failed", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}

	// 1. authentication
	value, exist := ctx.Get(authorizationPayloadKey)
	if !exist {
		log.WarnWithCtxFields(ctx, "payload key doesn't exist")
		ctx.Status(http.StatusBadRequest)
		return
	}
	payload, is := value.(*token.Payload)
	if !is {
		log.WarnWithCtxFields(ctx, "unexpected payload key")
		ctx.Status(http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(payload.Uid)
	if err != nil {
		log.WarnWithCtxFields(ctx, "convert payload user id error")
		ctx.Status(http.StatusBadRequest)
		return
	}

	// 2. authorization
	post, err := server.store.GetPost(ctx, req.PostId)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "get post err", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if post.PostUserID != int64(userId) {
		log.ErrorWithCtxFields(ctx, "unauthorized operation", zap.Error(err))
		ctx.Status(http.StatusUnauthorized)
		return
	}

	// 3. mark post status as "sold"
	err = server.store.SoldPost(ctx, post.ID)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "complete sold error", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}
