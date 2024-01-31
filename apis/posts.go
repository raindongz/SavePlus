package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/randongz/save_plus/db/sqlc"
)

type CreateNewPostRequest struct {
	Title        string `json:"title" binding:"required,min=6"`
	Content      string `json:"content" binding:"required,min=10,max=2048"`
	TotalPrice   string `json:"total_price" binding:"required,min=1"`
	PostUserID   *int64 `json:"post_user_id" binding:"required"`
	DeliveryType *int16 `json:"delivery_type" binding:"required,oneof=0 1"`
	Area         string `json:"area" binding:"required,min=1"`
	ItemNum      *int32 `json:"item_num" binding:"required,min=1"`
	PostStatus   *int16 `json:"post_status" binding:"required,oneof=0 1"`
	Negotiable   *int16 `json:"negotiable" binding:"required,oneof=0 1"`
	Images       string `json:"images" binding:"required"`
}

func (server *Server) createNewPost(ctx *gin.Context) {
	var req CreateNewPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//TODO below line will be used later for authentication
	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateNewPostParams{
		Title:      req.Title,
		Content:    req.Content,
		TotalPrice: req.TotalPrice,
		//TODO Authorization risk, need to compare username in payload with username related to userid here
		PostUserID:   *req.PostUserID,
		DeliveryType: *req.DeliveryType,
		Area: pgtype.Text{
			String: req.Area,
			Valid:  true,
		},
		ItemNum:    *req.ItemNum,
		PostStatus: *req.PostStatus,
		Negotiable: *req.Negotiable,
		Images:     req.Images,
	}

	post, err := server.store.CreateNewPost(ctx, arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, post)
}
