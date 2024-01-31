package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/randongz/save_plus/db/sqlc"
)

type CreateNewPostRequest struct{
	Title string `json:"title" binding:"required,min=6"`
	Content string `json:"content" binding:"required,min=10,max=2048"`
	TotalPrice string `json:"total_price" binding:"required,min=1"`
	DeliveryType *int16 `json:"delivery_type" binding:"required,oneof=0 1"`
	Area string `json:"area" binding:"required,min=1"`
	ItemNum *int32 `json:"item_num" binding:"required,min=1"`
	PostStatus *int16 `json:"post_status" binding:"required,oneof=0 1"`
	Negotiable *int16 `json:"negotiable" binding:"required,oneof=0 1"`
	Images string `json:"images" binding:"required"`
}

type CreateNewPostResponse struct{
	PostId int64 `json:"post_id"`
}


func (server *Server) createNewPost(ctx *gin.Context){
	var req CreateNewPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 1. todo: below line will be used later for authentication
	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)


	// 2. todo: upload images
		// 1. get each one a unique uuid as name and save them locally
		// 2. save those uuids in a slice and store them in database


	arg := db.CreateNewPostParams{
		Title: req.Title,
		Content: req.Content,
		TotalPrice: req.TotalPrice,
		// todo: get post userId from the payload
		// PostUserID: todo
		DeliveryType: *req.DeliveryType,
		Area: pgtype.Text{
			String: req.Area,
			Valid: true,
		},
		ItemNum: *req.ItemNum,
		PostStatus: *req.PostStatus,
		Negotiable: *req.Negotiable,
		Images: req.Images,
	} 


	post, err := server.store.CreateNewPost(ctx, arg)
	if err != nil{
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := CreateNewPostResponse{
		PostId: post.ID,
	}
	ctx.JSON(http.StatusOK, rsp)
}



type GetPostListRequest struct{
	PageSize *int32 `form:"page_size" binding:"required"`
	PageNum *int32 `form:"page_num" binding:"required"`
}

type GetPostListResponse struct{
	ID int64 `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Images string `json:"images"`
	TotalPrice string `json:"total_price"`
	Area pgtype.Text `json:"area"`
}

func (server *Server) getPostList(ctx *gin.Context){
	var req GetPostListRequest;
	if err := ctx.ShouldBind(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetPostListParams{
		Limit: int32(*req.PageSize),
		Offset: (*req.PageNum - 1) * (*req.PageSize),
	}

	postList, err := server.store.GetPostList(ctx, arg)
	if err != nil{
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	var rsp []GetPostListResponse
	for i := range postList{
		listItem := GetPostListResponse{
			ID: postList[i].ID,
			Title: postList[i].Title,
			Content: postList[i].Content,
			Images: postList[i].Images,
			TotalPrice: postList[i].TotalPrice,
			Area: postList[i].Area,
		}
		rsp = append(rsp, listItem)
	}

	ctx.JSON(http.StatusOK, rsp)
}

type GetPostDetailWithOutAuthRequest struct{
	ID int64 `form:"id" binding:"required"`
}

type GetPostDetailWithOutAuthResponse struct{
	Postid       int64       `json:"postid"`
	Title        string      `json:"title"`
	Content      string      `json:"content"`
	TotalPrice   string      `json:"total_price"`
	PostUserID   int64       `json:"post_user_id"`
	DeliveryType int16       `json:"delivery_type"`
	Area         pgtype.Text `json:"area"`
	ItemNum      int32       `json:"item_num"`
	PostStatus   int16       `json:"post_status"`
	Negotiable   int16       `json:"negotiable"`
	Images       string      `json:"images"`
	CreatedAt    pgtype.Date `json:"created_at"`
	UpdatedAt    pgtype.Date `json:"updated_at"`
	UserID       pgtype.Int8 `json:"user_id"`
	FullName     pgtype.Text `json:"full_name"`
	Email        pgtype.Text `json:"email"`
	Phone        pgtype.Text `json:"phone"`
	Gender       pgtype.Int2 `json:"gender"`
	Avatar       pgtype.Text `json:"avatar"`
}

func (server *Server) getPostDetailInfoWithOutAuth(ctx *gin.Context){
	var req GetPostDetailWithOutAuthRequest
	if err := ctx.ShouldBind(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	postWithUserDetail, err := server.store.GetPostAndRelatedUser(ctx, req.ID)
	if err != nil{
		if err == db.ErrRecordNotFound{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// todo: get all images
	ctx.JSON(http.StatusOK, postWithUserDetail)
}


type GetPostDetailInfoWithRequest struct{
	ID int64 `form:"id" binding:"required"`
}

type GetPostDetailInfoWithAuthResponse struct{
	PostId       int64       `json:"postid"`
	Title        string      `json:"title"`
	Content      string      `json:"content"`
	TotalPrice   string      `json:"total_price"`
	PostUserID   int64       `json:"post_user_id"`
	DeliveryType int16       `json:"delivery_type"`
	Area         pgtype.Text `json:"area"`
	ItemNum      int32       `json:"item_num"`
	PostStatus   int16       `json:"post_status"`
	Negotiable   int16       `json:"negotiable"`
	Images       string      `json:"images"`
	CreatedAt    pgtype.Date `json:"created_at"`
	UpdatedAt    pgtype.Date `json:"updated_at"`
	UserID       pgtype.Int8 `json:"user_id"`
	FullName     pgtype.Text `json:"full_name"`
	Email        pgtype.Text `json:"email"`
	Phone        pgtype.Text `json:"phone"`
	Gender       pgtype.Int2 `json:"gender"`
	Avatar       pgtype.Text `json:"avatar"`
} 

func (server *Server) getPostDetailInfoWithAuth(ctx *gin.Context){
	var req GetPostDetailWithOutAuthRequest
	if err := ctx.ShouldBind(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	postWithUserDetail, err := server.store.GetPostAndRelatedUser(ctx, req.ID)
	if err != nil{
		if err == db.ErrRecordNotFound{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// todo: get all images
	ctx.JSON(http.StatusOK, postWithUserDetail)
}

type UpdatePostInfoRequest struct{
	PostId int64 `json:"post_id" binding:"required,min=1"`
	Title string `json:"title" binding:"required,min=6"`
	Content string `json:"content" binding:"required,min=10,max=2048"`
	TotalPrice string `json:"total_price" binding:"required,min=1"`
	DeliveryType *int16 `json:"delivery_type" binding:"required,oneof=0 1"`
	Area string `json:"area" binding:"required,min=1"`
	ItemNum *int32 `json:"item_num" binding:"required,min=1"`
	PostStatus *int16 `json:"post_status" binding:"required,oneof=0 1"`
	Negotiable *int16 `json:"negotiable" binding:"required,oneof=0 1"`
	Images string `json:"images" binding:"required"`
}
type UpdatePostInfoResponse struct{
	PostId int64 `json:"post_id"`
}
func (server *Server)updatePostInfo(ctx *gin.Context){
	var req UpdatePostInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 1. todo: authorization. get user id from payload, 
		// then find that post and check the corresponding userid

	
	// 2. check if post exist
	_, err := server.store.GetPost(ctx, req.PostId)
	if err != nil{
		if err == db.ErrRecordNotFound{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	// 3. update the post info

	arg := db.UpdatePostParams{
		ID: req.PostId,
		Title: req.Title,
		Content: req.Content,
		TotalPrice: req.TotalPrice,
		DeliveryType: *req.DeliveryType,
		Area: pgtype.Text{
			String: req.Area,
			Valid: true,
		},
		ItemNum: *req.ItemNum,
		PostStatus: *req.PostStatus,
		Negotiable: *req.Negotiable,
		Images: req.Images,
	}

	post, err := server.store.UpdatePost(ctx, arg)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, UpdatePostInfoResponse{PostId: post.ID})
}

type DeletePostRequest struct{
	PostId int64 `json:"post_id" binding:"required,min=1"`
}

type DeletePostResponse struct{
	Msg	string `json:"msg"`
}
func (server *Server) deletePostInfo(ctx * gin.Context){
	var req DeletePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 1. todo: get userId from payload
	
	// 2. get post and compare post.userid with userid in payload

	// 3. if ok delete
	err := server.store.DeletePost(ctx, req.PostId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return	
	}
	ctx.JSON(http.StatusOK, DeletePostResponse{Msg: "delete user success"})
}


type GetInterestListRequest struct{
	PostId int64 `form:"post_id" binding:"required,min=1"`
}
type GetInterestListResponse struct{
	RecordID int64       `json:"record_id"`
	UserID   pgtype.Int8 `json:"user_id"`
	Username pgtype.Text `json:"username"`
	Avatar   pgtype.Text `json:"avatar"`
	Gender   pgtype.Int2 `json:"gender"`
}
func (server *Server) GetInterestList(ctx *gin.Context){
	var req GetInterestListRequest
	if err := ctx.ShouldBind(&req); err !=nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	interestUsers, err := server.store.GetPostInterestList(ctx, req.PostId)
	if err!=nil{
		if err == db.ErrRecordNotFound{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, interestUsers)
}

type InterestPostRequest struct{
	PostId int64 `json:"post_id" binding:"required,min=1"`
}

type InterestPostResponse struct{
	Msg string `json:"msg"`
}

func (server *Server) InterestPost(ctx *gin.Context){
	var req InterestPostRequest
	if err:=ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// todo: authenticate user

	arg := db.CreateInterestRecordParams{
		PostID: req.PostId,
		InterestedUserID: 1, // todo,
	}
	err := server.store.CreateInterestRecord(ctx, arg)

	if err != nil {
		if err == db.ErrUniqueViolation{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp:=InterestPostResponse{
		Msg: "Interested",
	}
	ctx.JSON(http.StatusOK, rsp)
}


type UnInterestPostRequest struct{
	PostId int64 `json:"post_id" binding:"required,min=1"`
}

type UnInterestPostResponse struct{
	Msg string `json:"msg"`
}

func (server *Server) UnInterestPost(ctx *gin.Context){
	var req InterestPostRequest
	if err:=ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// todo: authenticate user

	arg := db.DeleteInterestRecordParams{
		PostID: req.PostId,
		InterestedUserID: 1, // todo,
	}
	err := server.store.DeleteInterestRecord(ctx, arg)

	if err != nil {
		if err == db.ErrUniqueViolation{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp:=InterestPostResponse{
		Msg: "UnInterested",
	}
	ctx.JSON(http.StatusOK, rsp)
}