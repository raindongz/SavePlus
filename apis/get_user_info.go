package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nioliu/commons/log"
	db "github.com/randongz/save_plus/db/sqlc"
	"github.com/randongz/save_plus/token"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type getUserInfoReq struct {
	userId int64
}

type getUserInfoRsp struct {
	UserInfo *db.UsersInfo `json:"user_info,omitempty"`
	Result   int8          `json:"result,omitempty"` // 0: unknown; 1: success; 2: failed
	Details  string        `json:"details,omitempty"`
}

func (server *Server) getUserInfo(ctx *gin.Context) {
	req := new(getUserInfoReq)
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
		req.userId = int64(uid)
	}
	rsp := new(getUserInfoRsp)
	statusCode, err := server.doGetUserInfo(ctx, req, rsp)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "do get user info failed", zap.Error(err))
		rsp.Details = err.Error()
	}
	ctx.JSON(statusCode, rsp)
}

func (server *Server) doGetUserInfo(ctx context.Context, req *getUserInfoReq,
	rsp *getUserInfoRsp) (int, error) {
	// 1. get from db
	usersInfo, err := server.store.GetUserById(ctx, req.userId)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "get user id from db failed", zap.Error(err))
		return http.StatusInternalServerError, err
	}

	// 2. filter
	rsp.UserInfo = &db.UsersInfo{
		Username: usersInfo.Username,
		FullName: usersInfo.FullName,
		Email:    usersInfo.Email,
		Phone:    usersInfo.Phone,
		Gender:   usersInfo.Gender,
		Avatar:   usersInfo.Avatar,
	}

	return http.StatusOK, nil
}
