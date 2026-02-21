-- name: GetLessonsByModule :many
SELECT * FROM lessons
WHERE module_id = $1
ORDER BY order_index ASC;

-- name: GetLessonBySlug :one
SELECT * FROM lessons
WHERE module_id = $1 AND slug = $2
LIMIT 1;

-- name: GetLessonByID :one
SELECT * FROM lessons
WHERE id = $1
LIMIT 1;
