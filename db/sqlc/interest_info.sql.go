// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: interest_info.sql

package db

import (
	"context"
)

const createInterestRecord = `-- name: CreateInterestRecord :exec
INSERT INTO interest_info (
  post_id,
  interested_user_id
) VALUES (
  $1, $2
)
`

type CreateInterestRecordParams struct {
	PostID           int64 `json:"post_id"`
	InterestedUserID int64 `json:"interested_user_id"`
}

func (q *Queries) CreateInterestRecord(ctx context.Context, arg CreateInterestRecordParams) error {
	_, err := q.db.Exec(ctx, createInterestRecord, arg.PostID, arg.InterestedUserID)
	return err
}

const deleteInterestRecord = `-- name: DeleteInterestRecord :exec
DELETE FROM interest_info
WHERE id = $1
`

func (q *Queries) DeleteInterestRecord(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteInterestRecord, id)
	return err
}

const getInterestListByUserID = `-- name: GetInterestListByUserID :many
SELECT p.id, p.title, p.content, p.total_price, p.post_user_id, p.delivery_type, p.area, p.item_num, p.post_status, p.negotiable, p.images, p.deleted_flag, p.created_at, p.updated_at FROM post_info p
LEFT JOIN  interest_info i
ON p.id = i.post_id
WHERE i.interested_user_id = $1
`

func (q *Queries) GetInterestListByUserID(ctx context.Context, interestedUserID int64) ([]PostInfo, error) {
	rows, err := q.db.Query(ctx, getInterestListByUserID, interestedUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PostInfo{}
	for rows.Next() {
		var i PostInfo
		if err := rows.Scan(
			&i.ID,
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

const getInterestRecordByUserIdAndPostId = `-- name: GetInterestRecordByUserIdAndPostId :one
SELECT id 
FROM interest_info 
WHERE post_id = $1 AND interested_user_id = $2 LIMIT 1
`

type GetInterestRecordByUserIdAndPostIdParams struct {
	PostID           int64 `json:"post_id"`
	InterestedUserID int64 `json:"interested_user_id"`
}

func (q *Queries) GetInterestRecordByUserIdAndPostId(ctx context.Context, arg GetInterestRecordByUserIdAndPostIdParams) (int64, error) {
	row := q.db.QueryRow(ctx, getInterestRecordByUserIdAndPostId, arg.PostID, arg.InterestedUserID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getMyPostList = `-- name: GetMyPostList :many
SELECT id, title, content, total_price, post_user_id, delivery_type, area, item_num, post_status, negotiable, images, deleted_flag, created_at, updated_at FROM post_info p where p.post_user_id = $1 AND p.deleted_flag=0
`

func (q *Queries) GetMyPostList(ctx context.Context, postUserID int64) ([]PostInfo, error) {
	rows, err := q.db.Query(ctx, getMyPostList, postUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PostInfo{}
	for rows.Next() {
		var i PostInfo
		if err := rows.Scan(
			&i.ID,
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
