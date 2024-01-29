// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type InterestInfo struct {
	ID int64 `json:"id"`
	// related post id
	PostID int64 `json:"post_id"`
	// user interested to this post
	InterestedUserID int64       `json:"interested_user_id"`
	CreatedAt        pgtype.Date `json:"created_at"`
	UpdatedAt        pgtype.Date `json:"updated_at"`
}

type PostInfo struct {
	ID int64 `json:"id"`
	// post title
	Title string `json:"title"`
	// post content
	Content string `json:"content"`
	// post price, accurate to cent
	TotalPrice string `json:"total_price"`
	// user who posted this post
	PostUserID int64 `json:"post_user_id"`
	// 0: pick up. 1: mail
	DeliveryType int16 `json:"delivery_type"`
	// the area that the seller wants to trade
	Area pgtype.Text `json:"area"`
	// total items in this post
	ItemNum int32 `json:"item_num"`
	// 0: active, 1: sold, 2: inactive
	PostStatus int16 `json:"post_status"`
	// 0: not negotiable, 1: negotiable
	Negotiable int16 `json:"negotiable"`
	// post images, separated by comma
	Images string `json:"images"`
	// 0: active, 1: deleted
	DeletedFlag int16       `json:"deleted_flag"`
	CreatedAt   pgtype.Date `json:"created_at"`
	UpdatedAt   pgtype.Date `json:"updated_at"`
}

type TradingHistory struct {
	ID int64 `json:"id"`
	// related post id
	PostID int64 `json:"post_id"`
	// user who bought at least one item in the post
	SoldToUserID int64 `json:"sold_to_user_id"`
	// seller id save for later usage
	SellerID int64 `json:"seller_id"`
	// true transaction price
	Price       string      `json:"price"`
	DeletedFlag int16       `json:"deleted_flag"`
	CreatedAt   pgtype.Date `json:"created_at"`
	UpdatedAt   pgtype.Date `json:"updated_at"`
}

type UsersInfo struct {
	ID int64 `json:"id"`
	// unique username
	Username string `json:"username"`
	// encrypted password
	HashedPassword string `json:"hashed_password"`
	// lastname firstname
	FullName string `json:"full_name"`
	// unique email address
	Email string `json:"email"`
	// unique, including country code
	Phone pgtype.Text `json:"phone"`
	// 0: femail, 1: male, 2: other
	Gender int16 `json:"gender"`
	// user icon address
	Avatar pgtype.Text `json:"avatar"`
	// 0: active, 1: deleted
	DeletedFlag       int16       `json:"deleted_flag"`
	PasswordChangedAt pgtype.Date `json:"password_changed_at"`
	CreatedAt         pgtype.Date `json:"created_at"`
	UpdatedAt         pgtype.Date `json:"updated_at"`
}
