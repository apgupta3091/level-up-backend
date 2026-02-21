CREATE TABLE skills (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id   UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    skill_name  TEXT NOT NULL,
    order_index INTEGER NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_skills_module_id ON skills (module_id);

CREATE TABLE user_skill_progress (
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    skill_id     UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, skill_id)
);

CREATE INDEX idx_user_skill_progress_user_id ON user_skill_progress (user_id);
