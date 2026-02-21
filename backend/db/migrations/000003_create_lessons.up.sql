CREATE TABLE lessons (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id         UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    title             TEXT NOT NULL,
    slug              TEXT NOT NULL,
    content           TEXT NOT NULL DEFAULT '',
    order_index       INTEGER NOT NULL,
    estimated_minutes INTEGER NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (module_id, slug)
);

CREATE INDEX idx_lessons_module_id ON lessons (module_id);

CREATE TRIGGER set_lessons_updated_at
    BEFORE UPDATE ON lessons
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();
