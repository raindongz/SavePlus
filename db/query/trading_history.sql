-- name: CreateTradingRecord :one
INSERT INTO trading_history (
  post_id, 
  sold_to_user_id,
  seller_id, 
  price, 
  deleted_flag
) VALUES (
  $1, $2, $3, $4, 0
) RETURNING *;


-- name: GetRecord :one
SELECT * FROM trading_history 
WHERE id = $1 LIMIT 1;

-- name: GetPurchaseByUserId :many
SELECT t.*,p.* FROM trading_history as t
LEFT JOIN post_info as p
ON t.post_id==p.post_id
WHERE sold_to_user_id = $1;

-- name: GetRecordList :many
SELECT * FROM trading_history
WHERE 
    deleted_flag = 0
ORDER BY updated_at desc
LIMIT $1
OFFSET $2;

-- name: DeleteTradingRecord :exec
Update trading_history
SET
    deleted_flag = $2
WHERE id = $1;