// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: trading_history.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createTradingRecord = `-- name: CreateTradingRecord :one
INSERT INTO trading_history (
  post_id, 
  sold_to_user_id,
  seller_id, 
  price, 
  deleted_flag
) VALUES (
  $1, $2, $3, $4, 0
) RETURNING id, post_id, sold_to_user_id, seller_id, price, deleted_flag, created_at, updated_at
`

type CreateTradingRecordParams struct {
	PostID       int64  `json:"post_id"`
	SoldToUserID int64  `json:"sold_to_user_id"`
	SellerID     int64  `json:"seller_id"`
	Price        string `json:"price"`
}

func (q *Queries) CreateTradingRecord(ctx context.Context, arg CreateTradingRecordParams) (TradingHistory, error) {
	row := q.db.QueryRow(ctx, createTradingRecord,
		arg.PostID,
		arg.SoldToUserID,
		arg.SellerID,
		arg.Price,
	)
	var i TradingHistory
	err := row.Scan(
		&i.ID,
		&i.PostID,
		&i.SoldToUserID,
		&i.SellerID,
		&i.Price,
		&i.DeletedFlag,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteTradingRecord = `-- name: DeleteTradingRecord :exec
Update trading_history
SET
    deleted_flag = $2
WHERE id = $1
`

type DeleteTradingRecordParams struct {
	ID          int64 `json:"id"`
	DeletedFlag int16 `json:"deleted_flag"`
}

func (q *Queries) DeleteTradingRecord(ctx context.Context, arg DeleteTradingRecordParams) error {
	_, err := q.db.Exec(ctx, deleteTradingRecord, arg.ID, arg.DeletedFlag)
	return err
}

const getPurchaseByUserId = `-- name: GetPurchaseByUserId :many
SELECT t.id, t.post_id, t.sold_to_user_id, t.seller_id, t.price, t.deleted_flag, t.created_at, t.updated_at,p.id, p.title, p.content, p.total_price, p.post_user_id, p.delivery_type, p.area, p.item_num, p.post_status, p.negotiable, p.images, p.deleted_flag, p.created_at, p.updated_at FROM trading_history as t
LEFT JOIN post_info as p
ON t.post_id==p.post_id
WHERE sold_to_user_id = $1
`

type GetPurchaseByUserIdRow struct {
	ID            int64       `json:"id"`
	PostID        int64       `json:"post_id"`
	SoldToUserID  int64       `json:"sold_to_user_id"`
	SellerID      int64       `json:"seller_id"`
	Price         string      `json:"price"`
	DeletedFlag   int16       `json:"deleted_flag"`
	CreatedAt     pgtype.Date `json:"created_at"`
	UpdatedAt     pgtype.Date `json:"updated_at"`
	ID_2          pgtype.Int8 `json:"id_2"`
	Title         pgtype.Text `json:"title"`
	Content       pgtype.Text `json:"content"`
	TotalPrice    pgtype.Text `json:"total_price"`
	PostUserID    pgtype.Int8 `json:"post_user_id"`
	DeliveryType  pgtype.Int2 `json:"delivery_type"`
	Area          pgtype.Text `json:"area"`
	ItemNum       pgtype.Int4 `json:"item_num"`
	PostStatus    pgtype.Int2 `json:"post_status"`
	Negotiable    pgtype.Int2 `json:"negotiable"`
	Images        pgtype.Text `json:"images"`
	DeletedFlag_2 pgtype.Int2 `json:"deleted_flag_2"`
	CreatedAt_2   pgtype.Date `json:"created_at_2"`
	UpdatedAt_2   pgtype.Date `json:"updated_at_2"`
}

func (q *Queries) GetPurchaseByUserId(ctx context.Context, soldToUserID int64) ([]GetPurchaseByUserIdRow, error) {
	rows, err := q.db.Query(ctx, getPurchaseByUserId, soldToUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPurchaseByUserIdRow{}
	for rows.Next() {
		var i GetPurchaseByUserIdRow
		if err := rows.Scan(
			&i.ID,
			&i.PostID,
			&i.SoldToUserID,
			&i.SellerID,
			&i.Price,
			&i.DeletedFlag,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.Title,
			&i.Content,
			&i.TotalPrice,
			&i.PostUserID,
			&i.DeliveryType,
			&i.Area,
			&i.ItemNum,
			&i.PostStatus,
			&i.Negotiable,
			&i.Images,
			&i.DeletedFlag_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRecord = `-- name: GetRecord :one
SELECT id, post_id, sold_to_user_id, seller_id, price, deleted_flag, created_at, updated_at FROM trading_history 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRecord(ctx context.Context, id int64) (TradingHistory, error) {
	row := q.db.QueryRow(ctx, getRecord, id)
	var i TradingHistory
	err := row.Scan(
		&i.ID,
		&i.PostID,
		&i.SoldToUserID,
		&i.SellerID,
		&i.Price,
		&i.DeletedFlag,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRecordList = `-- name: GetRecordList :many
SELECT id, post_id, sold_to_user_id, seller_id, price, deleted_flag, created_at, updated_at FROM trading_history
WHERE 
    deleted_flag = 0
ORDER BY updated_at desc
LIMIT $1
OFFSET $2
`

type GetRecordListParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetRecordList(ctx context.Context, arg GetRecordListParams) ([]TradingHistory, error) {
	rows, err := q.db.Query(ctx, getRecordList, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []TradingHistory{}
	for rows.Next() {
		var i TradingHistory
		if err := rows.Scan(
			&i.ID,
			&i.PostID,
			&i.SoldToUserID,
			&i.SellerID,
			&i.Price,
			&i.DeletedFlag,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
