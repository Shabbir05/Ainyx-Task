-- name: CreateUser :execresult
INSERT INTO users (name, dob)
VALUES (?, ?);

-- name: GetUserByID :one
SELECT id, name, dob
FROM users
WHERE id = ?
LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET name = ?,
    dob  = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: ListUsers :many
SELECT id, name, dob
FROM users
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;
