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
WHERE id = $1 and deleted_flag = 0 LIMIT 1;

-- name: GetPostAndRelatedUser :one
SELECT
pi.id as postId,
pi.title ,
pi.content,
pi.total_price,
pi.post_user_id,
pi.delivery_type,
pi.area,
pi.item_num,
pi.post_status,
pi.negotiable,
pi.images,
pi.created_at,
pi.updated_at,

ui.id as user_id,
ui.full_name,
ui.email,
ui.phone,
ui.gender,
ui.avatar
from post_info pi left join users_info ui on pi.post_user_id = ui.id
WHERE pi.id = $1 and pi.deleted_flag = 0 and ui.deleted_flag = 0;

-- name: GetPostListNoAuth :many
SELECT * FROM post_info
WHERE 
    deleted_flag = 0
ORDER BY updated_at desc
LIMIT $1
OFFSET $2;

-- name: GetPostListAuth :many
SELECT p.id as post_id,
       p.title as post_title,
       p.content as post_content,
       p.images as post_images,
       p.total_price as price,
       p.area,
       i.id as liked

FROM post_info p
    left join interest_info i on i.interested_user_id = $1 AND p.id = i.post_id
WHERE
    p.deleted_flag = 0
ORDER BY p.updated_at desc
LIMIT $2
OFFSET $3;

-- name: GetPostInterestList :many
select
ii.id as record_id,
ui.id as user_id,
ui.username,
ui.avatar,
ui.gender
from interest_info ii left join users_info ui on ii.interested_user_id = ui.id
where ii.post_id = $1 and ui.deleted_flag = 0;

-- name: UpdatePost :one
UPDATE post_info
SET 
title = $2,
content = $3,
total_price = $4,
delivery_type = $5,
area = $6,
item_num = $7,
post_status = $8,
negotiable = $9,
images = $10
WHERE id = $1 and deleted_flag = 0
RETURNING *;

-- name: DeletePost :exec
UPDATE post_info
SET
  deleted_flag = 1
WHERE id = $1;

-- name: SoldPost :exec
UPDATE post_info
SET
  post_status = 1
WHERE id = $1;