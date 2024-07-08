-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  id, firstname,lastname,email,age,created
) VALUES (
  $1, $2, $3, $4, $5 ,$6
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
  set   firstname = COALESCE($2, firstname),
  lastname = COALESCE($3, lastname),
  email = COALESCE($4, email),
  age = COALESCE($5, age)
WHERE id = $1
RETURNING *;