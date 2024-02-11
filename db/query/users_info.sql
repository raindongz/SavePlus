-- name: CreateNewUser :one
INSERT INTO users_info (
  username,
  hashed_password,
  full_name,
  email,
  phone,
  gender,
  avatar,
  deleted_flag
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, 0
) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users_info 
WHERE id = $1 and deleted_flag = 0 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users_info
WHERE email = $1 LIMIT 1;

-- name: UpdateUserInfo :one
UPDATE users_info
SET 
  username = $2,
  full_name = $3,
  email = $4,
  phone = $5,
  gender = $6,
  avatar = $7
WHERE id = $1
RETURNING *;


-- name: DeleteUser :exec
Update users_info
SET
    deleted_flag = 1
WHERE id = $1;
