package apis

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nioliu/commons/log"
	db "github.com/randongz/save_plus/db/sqlc"
	"github.com/randongz/save_plus/token"
	"go.uber.org/zap"
)

type UpdateUserInfoReq struct {
	UserId   int64  `json:"id" binding:"required,min=1"`
	Username string `json:"username" binding:"required,min=3"`
	// Password string `json:"hashed_password" binding:"required,min=10"`
	FullName string `json:"full_name" binding:"required,min=4"`
	Email    string `json:"email" binding:"required,min=4"`
	Phone    string `json:"phone" binding:"required,min=4"`
	Avatar   string `json:"avatar" binding:"required,min=4"`
}

type UpdateUserInfoRsp struct {
	UserInfo *db.UsersInfo `json:"user_info"`
}

func (server *Server) updateUserInfo(ctx *gin.Context) {
	req := new(UpdateUserInfoReq)
	// 1. params varification
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.ErrorWithCtxFields(ctx, "params not valid", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 2. authentication
	v, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		log.WarnWithCtxFields(ctx, "no payload key exists")
		ctx.Status(http.StatusBadRequest)
		return
	}

	if payload, is := v.(*token.Payload); !is {
		log.WarnWithCtxFields(ctx, "unexpected payload type")
		ctx.Status(http.StatusBadRequest)
		return
	} else {
		uid, err := strconv.Atoi(payload.Uid)
		if err != nil {
			log.WarnWithCtxFields(ctx, "convert uid in payload from string to int err")
			ctx.Status(http.StatusBadRequest)
			return
		}
		// authorization
		if uid != int(req.UserId) {
			log.WarnWithCtxFields(ctx, fmt.Sprintf("unauthorized operation %d", uid))
			ctx.Status(http.StatusUnauthorized)
			return
		}
	}

	rsp := new(UpdateUserInfoRsp)
	// 3. check if user exist, if yes update user info
	if err := server.doUpdateUserInfo(ctx, req, rsp); err != nil {
		log.ErrorWithCtxFields(ctx, "doUpdateUserInfo failed: ", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// 4. return response
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) doUpdateUserInfo(ctx *gin.Context, req *UpdateUserInfoReq, rsp *UpdateUserInfoRsp) error {
	// 1. check if user exists in db
	userInfo, err := server.store.GetUserById(ctx, req.UserId)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "get user info from db failed: ", zap.Error(err))
		return err
	}

	// 2. update user info
	arg := db.UpdateUserInfoParams{
		ID:       userInfo.ID,
		Username: req.Username,
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    pgtype.Text{String: req.Phone, Valid: true},
		Avatar:   pgtype.Text{String: req.Avatar, Valid: true},
	}
	updatedUserInfo, err := server.store.UpdateUserInfo(ctx, arg)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "update user info to db failed: ", zap.Error(err))
		return err
	}

	// 3. setresponse
	rsp.UserInfo = &db.UsersInfo{
		ID:       updatedUserInfo.ID,
		Username: updatedUserInfo.Username,
		FullName: updatedUserInfo.FullName,
		Email:    updatedUserInfo.Email,
		Phone:    updatedUserInfo.Phone,
		Avatar:   updatedUserInfo.Avatar,
	}
	return nil
}
