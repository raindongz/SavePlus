package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nioliu/commons/log"
	"go.uber.org/zap"
	"net/http"
)

type viewMyPurchaseRsp struct {
	Items  []InterestItem `json:"items,omitempty"`
	Detail string         `json:"detail,omitempty"`
}

func (server *Server) viewMyPurchaseHistory(ctx *gin.Context) {
	userIdFromCtx, err := getUserIdFromCtx(ctx)
	if err != nil || userIdFromCtx == 0 {
		log.ErrorWithCtxFields(ctx, "get user id failed", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}
	rsp := new(viewMyInterestRsp)
	statusCode, err := server.doViewPurchaseHistory(ctx, userIdFromCtx, rsp)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "do view purchase failed", zap.Error(err))
		rsp.Detail = err.Error()
	}
	ctx.JSON(statusCode, rsp)
}

func (server *Server) doViewPurchaseHistory(ctx context.Context, userId int64,
	rsp *viewMyInterestRsp) (statusCode int, err error) {
	tradingHistories, err := server.store.GetPurchaseByUserId(ctx, userId)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "get purchase by user id failed", zap.Error(err))
		return http.StatusInternalServerError, err
	}
	items := make([]InterestItem, 0, len(tradingHistories))
	for _, history := range tradingHistories {
		items = append(items,
			InterestItem{
				PostId:  history.PostID,
				Status:  history.PostStatus.Int16,
				Title:   history.Title.String,
				Content: history.Content.String,
				Price:   history.Price,
			},
		)
	}
	rsp.Items = items
	return http.StatusOK, nil
}
