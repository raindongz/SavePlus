
-- name: CreateInterestRecord :exec
INSERT INTO interest_info (
  post_id,
  interested_user_id
) VALUES (
  $1, $2
);

-- name: GetInterestRecordByUserIdAndPostId :one
SELECT id 
FROM interest_info 
WHERE post_id = $1 AND interested_user_id = $2 LIMIT 1;


-- name: DeleteInterestRecord :exec
DELETE FROM interest_info
WHERE id = $1;

-- name: GetInterestListByUserID :many
SELECT i.*,p.* FROM interest_info i
JOIN post_info p
ON i.interested_user_id = p.post_user_id
WHERE interested_user_id = $1 ;

-- name: GetMyPostList :many
SELECT * FROM post_info p where p.post_user_id = $1 AND p.deleted_flag=0;