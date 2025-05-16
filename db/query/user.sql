-- name: CreateUser :one
INSERT INTO users(
    username, name, email, phone, password
) VALUES(
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: LoginUser :one
SELECT username, password
FROM users where username = $1;

-- name: GetUser :one
SELECT *
FROM users where username = $1;

-- name: UpdatePassword :exec
UPDATE users
SET password = $1
WHERE username = $2;