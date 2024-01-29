
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
