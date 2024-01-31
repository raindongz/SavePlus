package apis

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nioliu/commons/log"
	"github.com/randongz/save_plus/token"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

type userLoginReq struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type userLoginRsp struct {
	AccessToken string `json:"accessToken,omitempty"`
	Detail      string `json:"detail,omitempty"`
	Status      int8   `json:"status,omitempty"` // 0: unknown; 1: successfully; 2: wrong password;
}

func (server *Server) userLogin(ctx *gin.Context) {
	var statusCode int
	var err error
	var rspStatus int8
	rsp := new(userLoginRsp)
	defer func() {
		if err != nil {
			log.ErrorWithCtxFields(ctx, "user login failed", zap.Error(err))
			rsp.Detail = err.Error()
		}
		rsp.Status = rspStatus
		ctx.JSON(statusCode, rsp)
	}()
	req := new(userLoginReq)
	if err = ctx.ShouldBindJSON(req); err != nil {
		log.ErrorWithCtxFields(ctx, "bind json failed", zap.Error(err))
		statusCode = http.StatusBadRequest
		return
	}

	// 1. check basic user info
	if err = checkLoginUserInfoParams(ctx, req); err != nil {
		log.ErrorWithCtxFields(ctx, "check basic user info failed", zap.Error(err))
		statusCode = http.StatusBadRequest
		return
	}

	// 2. get user info from database
	usersInfo, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "get user by email failed", zap.Error(err))
		statusCode = http.StatusBadRequest
		return
	}

	// 3. compare user info
	if !(strings.EqualFold(usersInfo.HashedPassword, req.Password)) {
		log.InfoWithCtxFields(ctx, "password is not equal")
		statusCode = http.StatusUnauthorized
		rspStatus = 2
		return
	}

	// 4. generate userToken
	userToken, err := server.tokenMaker.CreateToken(strconv.Itoa(int(usersInfo.ID)), token.DefaultTokenDuration)

	// 5. return
	rsp.AccessToken = userToken
	statusCode = http.StatusOK
	rspStatus = 1

}

func checkLoginUserInfoParams(ctx context.Context, userInfo *userLoginReq) error {
	if userInfo.Password == "" {
		log.ErrorWithCtxFields(ctx, "password is empty")
		return fmt.Errorf("password is empty")
	}
	if userInfo.Email == "" {
		log.ErrorWithCtxFields(ctx, "email is empty")
		return fmt.Errorf("email is empty")
	}
	return nil
}
