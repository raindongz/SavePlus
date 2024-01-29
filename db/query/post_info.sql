-- name: CreateNewPost :one
INSERT INTO post_info (
title,
content,
total_price,
post_user_id,
delivery_type,
area,
item_num,
post_status,
negotiable,
images,
deleted_flag
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 0
) RETURNING *;

-- name: GetPost :one
SELECT * FROM post_info 
WHERE id = $1 LIMIT 1;

-- name: GetPostList :many
SELECT * FROM post_info
WHERE 
    deleted_flag = 0
ORDER BY updated_at desc
LIMIT $1
OFFSET $2;

-- name: UpdatePost :one
UPDATE post_info
SET 
title = $2,
content = $3,
total_price = $4,
post_user_id = $5,
delivery_type = $6,
area = $7,
item_num = $8,
post_status = $9,
negotiable = $10,
images = $11
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
UPDATE post_info
SET
  deleted_flag = 1
WHERE id = $1;