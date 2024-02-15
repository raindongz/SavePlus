package apis

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/nioliu/commons/log"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/randongz/save_plus/db/sqlc"
	"github.com/randongz/save_plus/token"
)

type CreateNewPostRequest struct {
	Title        string `json:"title" binding:"required,min=6"`
	UserId       int64  `json:"user_id" binding:"required,min=1"`
	Content      string `json:"content" binding:"required,min=10,max=2048"`
	TotalPrice   string `json:"total_price" binding:"required,min=1"`
	DeliveryType *int16 `json:"delivery_type" binding:"required,oneof=0 1"`
	Area         string `json:"area" binding:"required,min=1"`
	ItemNum      *int32 `json:"item_num" binding:"required,min=1"`
	PostStatus   *int16 `json:"post_status" binding:"required,oneof=0 1"`
	Negotiable   *int16 `json:"negotiable" binding:"required,oneof=0 1"`
	Images       string `json:"images" binding:"required"`
}

type CreateNewPostResponse struct {
	PostId       int64       `json:"post_id"`
	Title        string      `json:"title"`
	Content      string      `json:"content"`
	TotalPrice   string      `json:"total_price"`
	PostUserID   int64       `json:"post_user_id"`
	DeliveryType int16       `json:"delivery_type"`
	Area         string      `json:"area"`
	ItemNum      int32       `json:"item_num"`
	PostStatus   int16       `json:"post_status"`
	Negotiable   int16       `json:"negotiable"`
	Images       string      `json:"images"`
	CreatedAt    pgtype.Date `json:"created_at"`
	UpdatedAt    pgtype.Date `json:"updated_at"`
	FullName     string      `json:"full_name"`
	Email        string      `json:"email"`
	Phone        string      `json:"phone"`
	Gender       int16       `json:"gender"`
	Avatar       pgtype.Text `json:"avatar"`
}

func (server *Server) createNewPost(ctx *gin.Context) {
	var req CreateNewPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.ErrorWithCtxFields(ctx, "bind json failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 1. authentication
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userId, err := strconv.Atoi(authPayload.Uid)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "user id convertion failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if userId != int(req.UserId) {
		log.ErrorWithCtxFields(ctx, "unauthorized request:")
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized request")))
		return
	}

	// 2. create new post in databse
	arg := db.CreateNewPostParams{
		Title:        req.Title,
		Content:      req.Content,
		TotalPrice:   req.TotalPrice,
		PostUserID:   req.UserId,
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
			log.ErrorWithCtxFields(ctx, "duplicated record", zap.Error(err))
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		log.ErrorWithCtxFields(ctx, "create new post internal error: ", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 3. get related user info by userid in post
	user, err := server.store.GetUserById(ctx, int64(userId))
	if err != nil {
		if err == db.ErrRecordNotFound {
			log.ErrorWithCtxFields(ctx, "record not found", zap.Error(err))
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.ErrorWithCtxFields(ctx, "internal server error", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 4. put user info in response alone with post info just created
	rsp := CreateNewPostResponse{
		PostId:       post.ID,
		PostUserID:   user.ID,
		Title:        post.Title,
		Content:      post.Content,
		TotalPrice:   post.TotalPrice,
		DeliveryType: post.DeletedFlag,
		Area:         post.Area.String,
		ItemNum:      post.ItemNum,
		PostStatus:   post.PostStatus,
		Negotiable:   post.Negotiable,
		Images:       post.Images,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
		Gender:       user.Gender,
		Avatar:       user.Avatar,
		FullName:     user.FullName,
		Email:        user.Email,
		Phone:        user.Phone.String,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type GetPostListRequest struct {
	PageSize *int32 `form:"page_size" binding:"required,min=1"`
	PageNum  *int32 `form:"page_num" binding:"required,min=1"`
}

type GetPostListResponse struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Images     string `json:"images"`
	TotalPrice string `json:"total_price"`
	Area       string `json:"area"`
}

func (server *Server) getPostList(ctx *gin.Context) {
	// 1. check if params in request is valid
	var req GetPostListRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.ErrorWithCtxFields(ctx, "getPostList request params not valid", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 2. get post list from database
	arg := db.GetPostListParams{
		Limit:  int32(*req.PageSize),
		Offset: (*req.PageNum - 1) * (*req.PageSize),
	}
	postList, err := server.store.GetPostList(ctx, arg)
	if err != nil {
		if err == db.ErrRecordNotFound {
			log.ErrorWithCtxFields(ctx, "getPostList: no record found: ", zap.Error(err))
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.ErrorWithCtxFields(ctx, "getpostList internal server error: ", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 3. parse database model to response json object
	var rsp []GetPostListResponse
	for i := range postList {
		listItem := GetPostListResponse{
			ID:         postList[i].ID,
			Title:      postList[i].Title,
			Content:    postList[i].Content,
			Images:     postList[i].Images,
			TotalPrice: postList[i].TotalPrice,
			Area:       postList[i].Area.String,
		}
		rsp = append(rsp, listItem)
	}

	ctx.JSON(http.StatusOK, rsp)
}

type GetPostDetailWithOutAuthRequest struct {
	ID int64 `form:"id" binding:"required,min=1"`
}

type GetPostDetailWithOutAuthResponse struct {
	Postid       int64       `json:"postid"`
	Title        string      `json:"title"`
	Content      string      `json:"content"`
	TotalPrice   string      `json:"total_price"`
	DeliveryType int16       `json:"delivery_type"`
	Area         string      `json:"area"`
	ItemNum      int32       `json:"item_num"`
	PostStatus   int16       `json:"post_status"`
	Negotiable   int16       `json:"negotiable"`
	Images       string      `json:"images"`
	CreatedAt    pgtype.Date `json:"created_at"`
	UpdatedAt    pgtype.Date `json:"updated_at"`
	Gender       pgtype.Int2 `json:"gender"`
	Avatar       pgtype.Text `json:"avatar"`
}

func (server *Server) getPostDetailInfoWithOutAuth(ctx *gin.Context) {
	// 1. check if params in request body is valid
	var req GetPostDetailWithOutAuthRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.ErrorWithCtxFields(ctx, "request params not valid", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 2. get models from database
	postWithUserDetail, err := server.store.GetPostAndRelatedUser(ctx, req.ID)
	if err != nil {
		if err == db.ErrRecordNotFound {
			log.ErrorWithCtxFields(ctx, "db record not found:", zap.Error(err))
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.ErrorWithCtxFields(ctx, "internalserver error: ", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 3. only get the first img
	firstImg := strings.Split(postWithUserDetail.Images, ",")[0]

	// 4. construct response
	rsp := GetPostDetailWithOutAuthResponse{
		Postid:       postWithUserDetail.Postid,
		Title:        postWithUserDetail.Title,
		Content:      postWithUserDetail.Content,
		TotalPrice:   postWithUserDetail.TotalPrice,
		DeliveryType: postWithUserDetail.DeliveryType,
		Area:         postWithUserDetail.Area.String,
		ItemNum:      postWithUserDetail.ItemNum,
		PostStatus:   postWithUserDetail.PostStatus,
		Negotiable:   postWithUserDetail.Negotiable,
		Images:       firstImg,
		CreatedAt:    postWithUserDetail.CreatedAt,
		UpdatedAt:    postWithUserDetail.UpdatedAt,
		Gender:       postWithUserDetail.Gender,
		Avatar:       postWithUserDetail.Avatar,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type GetPostDetailInfoWithAuthRequest struct {
	PostId int64 `form:"post_id" binding:"required,min=1"`
	UserId int64 `form:"user_id" binding:"required,min=1"`
}

type GetPostDetailInfoWithAuthResponse struct {
	PostId       int64       `json:"postid"`
	Title        string      `json:"title"`
	Content      string      `json:"content"`
	TotalPrice   string      `json:"total_price"`
	PostUserID   int64       `json:"post_user_id"`
	DeliveryType int16       `json:"delivery_type"`
	Area         string      `json:"area"`
	ItemNum      int32       `json:"item_num"`
	PostStatus   int16       `json:"post_status"`
	Negotiable   int16       `json:"negotiable"`
	Images       string      `json:"images"`
	CreatedAt    pgtype.Date `json:"created_at"`
	UpdatedAt    pgtype.Date `json:"updated_at"`
	FullName     pgtype.Text `json:"full_name"`
	Email        pgtype.Text `json:"email"`
	Phone        pgtype.Text `json:"phone"`
	Gender       pgtype.Int2 `json:"gender"`
	Avatar       pgtype.Text `json:"avatar"`
}

func (server *Server) getPostDetailInfoWithAuth(ctx *gin.Context) {
	// 1. check request params
	var req GetPostDetailInfoWithAuthRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.ErrorWithCtxFields(ctx, "invalid request params", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 2.authorization
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userId, err := strconv.Atoi(authPayload.Uid)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "user id convertion failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if userId != int(req.UserId) {
		log.ErrorWithCtxFields(ctx, "unauthorized request:")
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized request")))
		return
	}

	// 2. get record from database
	postWithUserDetail, err := server.store.GetPostAndRelatedUser(ctx, req.PostId)
	if err != nil {
		if err == db.ErrRecordNotFound {
			log.ErrorWithCtxFields(ctx, "no record found in db", zap.Error(err))
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		log.ErrorWithCtxFields(ctx, "internal server error", zap.Error(err))
		return
	}

	// 4. construct response
	rsp := GetPostDetailInfoWithAuthResponse{
		PostId:       postWithUserDetail.Postid,
		Title:        postWithUserDetail.Title,
		Content:      postWithUserDetail.Content,
		TotalPrice:   postWithUserDetail.TotalPrice,
		DeliveryType: postWithUserDetail.DeliveryType,
		PostUserID:   postWithUserDetail.PostUserID,
		Area:         postWithUserDetail.Area.String,
		ItemNum:      postWithUserDetail.ItemNum,
		PostStatus:   postWithUserDetail.PostStatus,
		Negotiable:   postWithUserDetail.Negotiable,
		Images:       postWithUserDetail.Images,
		CreatedAt:    postWithUserDetail.CreatedAt,
		UpdatedAt:    postWithUserDetail.UpdatedAt,
		FullName:     postWithUserDetail.FullName,
		Email:        postWithUserDetail.Email,
		Phone:        postWithUserDetail.Phone,
		Gender:       postWithUserDetail.Gender,
		Avatar:       postWithUserDetail.Avatar,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type UpdatePostInfoRequest struct {
	PostId       int64  `json:"post_id" binding:"required,min=1"`
	PostUserID   int64  `json:"post_user_id" binding:"required,min=1"`
	Title        string `json:"title" binding:"required,min=6"`
	Content      string `json:"content" binding:"required,min=10,max=2048"`
	TotalPrice   string `json:"total_price" binding:"required,min=1"`
	DeliveryType *int16 `json:"delivery_type" binding:"required,oneof=0 1"`
	Area         string `json:"area" binding:"required,min=1"`
	ItemNum      *int32 `json:"item_num" binding:"required,min=1"`
	PostStatus   *int16 `json:"post_status" binding:"required,oneof=0 1"`
	Negotiable   *int16 `json:"negotiable" binding:"required,oneof=0 1"`
	Images       string `json:"images" binding:"required"`
}

type CreateOrUpdatePostResponse struct {
	PostId       int64       `json:"post_id"`
	Title        string      `json:"title"`
	Content      string      `json:"content"`
	TotalPrice   string      `json:"total_price"`
	PostUserID   int64       `json:"post_user_id"`
	DeliveryType int16       `json:"delivery_type"`
	Area         string      `json:"area"`
	ItemNum      int32       `json:"item_num"`
	PostStatus   int16       `json:"post_status"`
	Negotiable   int16       `json:"negotiable"`
	Images       string      `json:"images"`
	CreatedAt    pgtype.Date `json:"created_at"`
	UpdatedAt    pgtype.Date `json:"updated_at"`
	FullName     string      `json:"full_name"`
	Email        string      `json:"email"`
	Phone        string      `json:"phone"`
	Gender       int16       `json:"gender"`
	Avatar       pgtype.Text `json:"avatar"`
}

func (server *Server) updatePostInfo(ctx *gin.Context) {
	var req UpdatePostInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 1. authorization. get user id from payload,
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userId, err := strconv.Atoi(authPayload.Uid)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "user id convertion failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if userId != int(req.PostUserID) {
		log.ErrorWithCtxFields(ctx, "unauthorized request:")
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized request")))
		return
	}

	// 2. check if post exist
	postInfo, err := server.store.GetPost(ctx, req.PostId)
	if err != nil {
		if err == db.ErrRecordNotFound {
			log.ErrorWithCtxFields(ctx, "record not exist", zap.Error(err))
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.ErrorWithCtxFields(ctx, "internal server error", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 3. check if post belone to authenticated user
	if postInfo.PostUserID != int64(userId) {
		log.ErrorWithCtxFields(ctx, "unauthorized request:", zap.Error(errors.New("unauthorized request")))
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized request")))
		return
	}

	// 4. update the post info
	arg := db.UpdatePostParams{
		ID:           req.PostId,
		Title:        req.Title,
		Content:      req.Content,
		TotalPrice:   req.TotalPrice,
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

	post, err := server.store.UpdatePost(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 5. get related user info by userid in post
	user, err := server.store.GetUserById(ctx, int64(userId))
	if err != nil {
		if err == db.ErrRecordNotFound {
			log.ErrorWithCtxFields(ctx, "record not found", zap.Error(err))
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.ErrorWithCtxFields(ctx, "internal server error", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 6. put user info in response alone with the post info just created
	rsp := CreateNewPostResponse{
		PostId:       post.ID,
		PostUserID:   user.ID,
		Title:        post.Title,
		Content:      post.Content,
		TotalPrice:   post.TotalPrice,
		DeliveryType: post.DeletedFlag,
		Area:         post.Area.String,
		ItemNum:      post.ItemNum,
		PostStatus:   post.PostStatus,
		Negotiable:   post.Negotiable,
		Images:       post.Images,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
		Gender:       user.Gender,
		Avatar:       user.Avatar,
		FullName:     user.FullName,
		Email:        user.Email,
		Phone:        user.Phone.String,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type DeletePostRequest struct {
	PostId int64 `json:"post_id" binding:"required,min=1"`
	UserId int64 `json:"user_id" binding:"required,min=1"`
}

type DeletePostResponse struct {
	Msg string `json:"msg"`
}

func (server *Server) deletePostInfo(ctx *gin.Context) {
	var req DeletePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 1. authorization
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userId, err := strconv.Atoi(authPayload.Uid)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "user id convertion failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if userId != int(req.UserId) {
		log.ErrorWithCtxFields(ctx, "unauthorized request:")
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized request")))
		return
	}

	// 2. get post and compare post.userid with userid in payload
	postInfo, err := server.store.GetPost(ctx, req.PostId)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "record not exist", zap.Error(err))
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if postInfo.PostUserID != req.UserId {
		log.ErrorWithCtxFields(ctx, "unauthorized request", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// 3. if ok delete
	err = server.store.DeletePost(ctx, req.PostId)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "internal error", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, DeletePostResponse{Msg: "delete user success"})
}

type GetInterestListRequest struct {
	PostId int64 `form:"post_id" binding:"required,min=1"`
}
type GetInterestListResponse struct {
	RecordID int64       `json:"record_id"`
	UserID   pgtype.Int8 `json:"user_id"`
	Username pgtype.Text `json:"username"`
	Avatar   pgtype.Text `json:"avatar"`
	Gender   pgtype.Int2 `json:"gender"`
}

func (server *Server) GetInterestList(ctx *gin.Context) {
	var req GetInterestListRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.ErrorWithCtxFields(ctx, "invalid request params", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	interestUsers, err := server.store.GetPostInterestList(ctx, req.PostId)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "get post info error", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, interestUsers)
}

type InterestPostRequest struct {
	PostId int64 `json:"post_id" binding:"required,min=1"`
}

type InterestPostResponse struct {
	Msg string `json:"msg"`
}

func (server *Server) InterestPost(ctx *gin.Context) {
	var req InterestPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.ErrorWithCtxFields(ctx, "invalid request params", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 1. get user id from token
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userId, err := strconv.Atoi(authPayload.Uid)
	if err != nil {
		log.ErrorWithCtxFields(ctx, "user id convert failed", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// 2. if user already in interest then uninterest, if no record found then add interest record
	argForInterestInfo := db.GetInterestRecordByUserIdAndPostIdParams{
		PostID:           req.PostId,
		InterestedUserID: int64(userId),
	}
	interestRecordId, err := server.store.GetInterestRecordByUserIdAndPostId(ctx, argForInterestInfo)
	if err == nil {
		// already interested, remove record
		err := server.store.DeleteInterestRecord(ctx, interestRecordId)
		if err != nil {
			log.ErrorWithCtxFields(ctx, "delete interest record error", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		rsp := InterestPostResponse{
			Msg: "uninterest",
		}
		ctx.JSON(http.StatusOK, rsp)
	} else {
		if err == db.ErrRecordNotFound {
			// not interested yet, do interest
			arg := db.CreateInterestRecordParams{
				PostID:           req.PostId,
				InterestedUserID: int64(userId),
			}
			err = server.store.CreateInterestRecord(ctx, arg)
			if err != nil {
				log.ErrorWithCtxFields(ctx, "create interest record error", zap.Error(err))
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			rsp := InterestPostResponse{
				Msg: "Interested",
			}
			ctx.JSON(http.StatusOK, rsp)
		}

	}
}

// private method for create new post and update new post response
// func (server *Server) getUserInfoForCreateAndUpdatePostResponse(ctx *gin.Context, userId int, post db.PostInfo)(CreateOrUpdatePostResponse, error){
// 	var rsp CreateOrUpdatePostResponse
// 	user, err := server.store.GetUserById(ctx, int64(userId))
// 	if err != nil {
// 		if err == db.ErrRecordNotFound {
// 			log.ErrorWithCtxFields(ctx, "record not found", zap.Error(err))
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return rsp, err
// 		}
// 		log.ErrorWithCtxFields(ctx, "internal server error", zap.Error(err))
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return rsp, err
// 	}

// 	// 6. put user info in response alone with post info just created
// 	rsp = CreateOrUpdatePostResponse{
// 		PostId:       post.ID,
// 		PostUserID:   user.ID,
// 		Title:        post.Title,
// 		Content:      post.Content,
// 		TotalPrice:   post.TotalPrice,
// 		DeliveryType: post.DeletedFlag,
// 		Area:         post.Area.String,
// 		ItemNum:      post.ItemNum,
// 		PostStatus:   post.PostStatus,
// 		Negotiable:   post.Negotiable,
// 		Images:       post.Images,
// 		CreatedAt:    post.CreatedAt,
// 		UpdatedAt:    post.UpdatedAt,
// 		Gender:       user.Gender,
// 		Avatar:       user.Avatar,
// 		FullName:     user.FullName,
// 		Email:        user.Email,
// 		Phone:        user.Phone.String,
// 	}
// 	return rsp, nil
// }
