package apis

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nioliu/commons/log"
	db "github.com/randongz/save_plus/db/sqlc"
	"go.uber.org/zap"
	"net/http"
)

type createUserReq struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Gender   int    `json:"gender,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

type createUserRsp struct {
	Detail string `json:"detail,omitempty"`
	Result int8   `json:"result,omitempty"` // 0:unknown;1:success;2:duplicated
}

func (server *Server) createUser(ctx *gin.Context) {
	rsp := new(createUserRsp)
	req := new(createUserReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.ErrorWithCtxFields(ctx, "bind json failed", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}
	statusCode, err := server.doCreateUser(ctx, req, rsp)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "do create user failed", zap.Error(err))
		rsp.Detail = err.Error()
	}
	ctx.JSON(statusCode, rsp)

}

func (server *Server) doCreateUser(ctx context.Context, req *createUserReq,
	rsp *createUserRsp) (statusCode int, err error) {
	if err = checkBasicUserInfoParams(ctx, req); err != nil {
		log.ErrorWithCtxFields(ctx, "check basic user info failed", zap.Error(err))
		rsp.Detail = "check basic user info failed"
		return http.StatusUnauthorized, err
	}
	// 1. todo 增加邮箱验证功能

	// 2. insert directly
	_, err = server.store.CreateNewUser(ctx, db.CreateNewUserParams{
		Username:       req.Username,
		HashedPassword: req.Password,
		Email:          req.Email,
		//Phone:          pgtype.Text{},
		//Gender:         0,
		//Avatar:         pgtype.Text{},
	})
	if err != nil {
		log.ErrorWithCtxFields(ctx, "create new user failed", zap.Error(err))
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			rsp.Detail = "duplicated user"
			rsp.Result = 2
			return http.StatusBadRequest, errors.New("duplicated user")
		}
	}

	return http.StatusOK, nil
}

func checkBasicUserInfoParams(ctx context.Context, createUserReq *createUserReq) error {
	if createUserReq.Username == "" {
		log.ErrorWithCtxFields(ctx, "username is empty")
		return fmt.Errorf("username is empty")
	}
	if createUserReq.Password == "" {
		log.ErrorWithCtxFields(ctx, "password is empty")
		return fmt.Errorf("password is empty")
	}
	if createUserReq.Email == "" {
		log.ErrorWithCtxFields(ctx, "email is empty")
		return fmt.Errorf("email is empty")
	}
	return nil
}
