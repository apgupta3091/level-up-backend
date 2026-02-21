CREATE TYPE submission_status AS ENUM (
    'pending',
    'reviewed',
    'approved',
    'needs_revision'
);

CREATE TABLE submissions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id   UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    github_url      TEXT NOT NULL,
    written_answers TEXT NOT NULL DEFAULT '',
    status          submission_status NOT NULL DEFAULT 'pending',
    feedback        TEXT,
    submitted_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_at     TIMESTAMPTZ,
    UNIQUE (assignment_id, user_id)
);

CREATE INDEX idx_submissions_user_id ON submissions (user_id);
CREATE INDEX idx_submissions_assignment_id ON submissions (assignment_id);
CREATE INDEX idx_submissions_status ON submissions (status);
