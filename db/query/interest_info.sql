
-- name: CreateInterestRecord :exec
INSERT INTO interest_info (
  post_id,
  interested_user_id
) VALUES (
  $1, $2
);

-- name: DeleteInterestRecord :exec
DELETE FROM interest_info
WHERE post_id = $1 and interested_user_id = $2;

