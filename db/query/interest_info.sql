
-- name: CreateInterestRecord :one
INSERT INTO interest_info (
  post_id,
  interested_user_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetInterestList :many
SELECT * FROM interest_info 
WHERE post_id = $1;

-- name: DeleteInterestRecord :exec
DELETE FROM interest_info
WHERE id = $1;

-- name: GetInterestListByUserID :many
SELECT i.*,p.* FROM interest_info as i
JOIN post_info as p
ON i.interested_user_id == post_info.post_user_id
WHERE interested_user_id = $1 ;