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

type soldToReq struct {
	PostId   int64 `json:"postId"`
	SellerId int64 `json:"sellerId"`
	BuyerId  int64 `json:"buyerId"`
}

func (server *Server) soldTo(ctx *gin.Context) {
	req := new(soldToReq)
	if err := ctx.BindJSON(req); err != nil {
		log.ErrorWithCtxFields(ctx, "bind json failed", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}
	code, err := server.doSoldTo(ctx, req)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "do sold to failed", zap.Error(err))

	}
	ctx.Status(code)
}

func (server *Server) doSoldTo(ctx context.Context, req *soldToReq) (int, error) {
	if err := checkSoldReqParams(ctx, req); err != nil {
		log.ErrorWithCtxFields(ctx, "check sold req param failed", zap.Error(err))
		return http.StatusBadRequest, err
	}
	// 1. check the post item
	postInfo, err := server.store.GetPost(ctx, req.PostId)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "get post info failed", zap.Error(err))
		return http.StatusInternalServerError, err
	}
	if postInfo.PostStatus != 0 {
		log.ErrorWithCtxFields(ctx, "current post info is inactive",
			zap.Int64("post_id", postInfo.PostUserID))
		return http.StatusBadRequest, errors.New("current post info is inactive")
	}

	// 2. check buyer
	buyerInfo, err := server.store.GetUserById(ctx, req.BuyerId)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "get buyer info failed", zap.Error(err))
		return http.StatusBadRequest, errors.New("get buyer info failed")
	}
	if buyerInfo.DeletedFlag != 0 {
		log.ErrorWithCtxFields(ctx, "buyer is not active", zap.Int64("buyer id", req.BuyerId))
		return http.StatusBadRequest, errors.New("buyer is not active")
	}

	// 3. check seller
	sellerInfo, err := server.store.GetUserById(ctx, req.SellerId)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "get seller info failed", zap.Error(err))
		return http.StatusBadRequest, errors.New("get buyer info failed")
	}
	if sellerInfo.DeletedFlag != 0 {
		log.ErrorWithCtxFields(ctx, "seller is not active", zap.Int64("buyer id", req.SellerId))
		return http.StatusBadRequest, errors.New("buyer is not active")
	}

	// 4. sold
	_, err = server.store.CreateTradingRecord(ctx, db.CreateTradingRecordParams{
		PostID:       req.PostId,
		SoldToUserID: req.BuyerId,
		SellerID:     req.SellerId,
		Price:        postInfo.TotalPrice,
	})
	if err != nil {
		log.ErrorWithCtxFields(ctx, "add trading record failed", zap.Error(err))
		return http.StatusInternalServerError, errors.New("add trading record failed")
	}
	return http.StatusOK, nil
}

func checkSoldReqParams(ctx context.Context, soldToReq *soldToReq) error {
	if soldToReq.PostId == 0 {
		log.ErrorWithCtxFields(ctx, "post_id is empty")
		return fmt.Errorf("post_id is empty")
	}
	if soldToReq.SellerId == 0 {
		log.ErrorWithCtxFields(ctx, "seller_id is empty")
		return fmt.Errorf("seller_id is empty")
	}
	if soldToReq.BuyerId == 0 {
		log.ErrorWithCtxFields(ctx, "buyer_id is empty")
		return fmt.Errorf("buyer_id is empty")
	}
	return nil
}
