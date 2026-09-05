CREATE TABLE IF NOT EXISTS tasks (
    id           UUID         PRIMARY KEY     DEFAULT uuid_generate_v4(),
    household_id UUID         NOT NULL        REFERENCES households(id) ON DELETE CASCADE,
    version      INT                  NOT NULL DEFAULT 1,
    title        VARCHAR(200)         NOT NULL CHECK(char_length(title) BETWEEN 1 AND 200),
    description  TEXT,
    status       VARCHAR(20)          NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'in_progress', 'done')),
    assigned_to  UUID                 REFERENCES users(id) ON DELETE SET NULL,
    created_by   UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CHECK(
        created_at <= updated_at
    )
);

CREATE INDEX IF NOT EXISTS idx_tasks_household_id ON tasks(household_id);
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to ON tasks(assigned_to);