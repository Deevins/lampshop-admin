
-- name: GetUserByUsername :one
SELECT username, hashed_password, created_at
FROM users
WHERE username = $1;

-- name: ListCategories :many
SELECT id, name
FROM categories
ORDER BY name;

-- name: CreateCategory :exec
INSERT INTO categories (name)
VALUES ($1)
RETURNING id, name;

-- name: DeleteCategory :exec
DELETE
FROM categories
WHERE id = $1;

-- name: ListAttributeOptionsByCategory :many
SELECT id, category_id, key, label, type
FROM attribute_options
WHERE category_id = $1
ORDER BY id;

-- name: CreateAttributeOption :exec
INSERT INTO attribute_options (category_id, key, label, type)
VALUES ($1, $2, $3, $4)
RETURNING id, category_id, key, label, type;

-- name: DeleteAttributeOption :exec
DELETE
FROM attribute_options
WHERE id = $1;
