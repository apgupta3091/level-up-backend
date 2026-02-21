-- name: CreateSubmission :one
INSERT INTO submissions (assignment_id, user_id, github_url, written_answers)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSubmissionsByUser :many
SELECT * FROM submissions
WHERE user_id = $1
ORDER BY submitted_at DESC;

-- name: GetSubmissionByID :one
SELECT * FROM submissions
WHERE id = $1
LIMIT 1;

-- name: ListPendingSubmissions :many
SELECT * FROM submissions
WHERE status = 'pending'
ORDER BY submitted_at ASC;

-- name: ReviewSubmission :one
UPDATE submissions
SET
    status      = $2,
    feedback    = $3,
    reviewed_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetSubmissionByAssignmentAndUser :one
SELECT * FROM submissions
WHERE assignment_id = $1 AND user_id = $2
LIMIT 1;
