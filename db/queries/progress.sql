-- name: MarkLessonComplete :exec
INSERT INTO user_lesson_progress (user_id, lesson_id)
VALUES ($1, $2)
ON CONFLICT (user_id, lesson_id) DO NOTHING;

-- name: GetCompletedLessonIDs :many
SELECT lesson_id FROM user_lesson_progress
WHERE user_id = $1;

-- name: GetCompletedLessonCountByModule :one
SELECT COUNT(ulp.lesson_id)
FROM user_lesson_progress ulp
JOIN lessons l ON l.id = ulp.lesson_id
WHERE ulp.user_id = $1 AND l.module_id = $2;
