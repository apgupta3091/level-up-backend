-- name: GetSkillsByModule :many
SELECT * FROM skills
WHERE module_id = $1
ORDER BY order_index ASC;

-- name: GetSkillByID :one
SELECT * FROM skills
WHERE id = $1
LIMIT 1;

-- name: MarkSkillComplete :exec
INSERT INTO user_skill_progress (user_id, skill_id)
VALUES ($1, $2)
ON CONFLICT (user_id, skill_id) DO NOTHING;

-- name: GetCompletedSkillIDs :many
SELECT skill_id FROM user_skill_progress
WHERE user_id = $1;
