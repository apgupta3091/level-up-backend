-- name: ListModules :many
SELECT * FROM modules
ORDER BY order_index ASC;

-- name: GetModuleBySlug :one
SELECT * FROM modules
WHERE slug = $1
LIMIT 1;

-- name: GetModuleByID :one
SELECT * FROM modules
WHERE id = $1
LIMIT 1;
