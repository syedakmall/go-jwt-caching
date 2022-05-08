-- name: GetUser :one

SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: ListUsers :many

SELECT *
FROM users
ORDER BY name;

-- name: CreateUser :one

INSERT INTO users (name, password)
VALUES ($1,
        $2) RETURNING *;

-- name: DeleteUser :exec

DELETE
FROM users
WHERE id = $1;

-- name: GetUserByName :one

SELECT *
FROM users
WHERE name = $1
LIMIT 1;