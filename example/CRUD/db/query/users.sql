-- name: CreateUser :one
INSERT INTO users (
  username, full_name, email
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;