-- name: CreateUser :one
INSERT INTO users     (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(),      NOW(),      NOW(),    $1,              $2)
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET
  email = COALESCE(sqlc.narg('email'), email),
  hashed_password = COALESCE(sqlc.narg('hashed_password'), hashed_password),
  updated_at = NOW()
WHERE id = $1
RETURNING
  id, created_at, updated_at, email, is_chirpy_red;

-- name: UpgradeUserToChirpyRed :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;
