-- name: GetAssignmentByModuleID :one
SELECT * FROM assignments
WHERE module_id = $1
LIMIT 1;

-- name: GetAssignmentByID :one
SELECT * FROM assignments
WHERE id = $1
LIMIT 1;
