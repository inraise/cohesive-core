CREATE TABLE IF NOT EXISTS households (
    id         UUID         PRIMARY KEY     DEFAULT uuid_generate_v4(),
    version    INT                 NOT NULL DEFAULT 1,
    name       VARCHAR(100)        NOT NULL CHECK(char_length(name) BETWEEN 1 AND 100),
    created_at TIMESTAMPTZ         NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ         NOT NULL DEFAULT now(),

    CHECK(
        created_at <= updated_at
    )
);

CREATE TABLE IF NOT EXISTS household_members (
    id           UUID        PRIMARY KEY     DEFAULT uuid_generate_v4(),
    household_id UUID        NOT NULL        REFERENCES households(id) ON DELETE CASCADE,
    user_id      UUID        NOT NULL        REFERENCES users(id) ON DELETE CASCADE,
    role         VARCHAR(20) NOT NULL DEFAULT 'member' CHECK(role IN ('owner', 'admin', 'member')),
    joined_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(household_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_household_members_user_id ON household_members(user_id);
CREATE INDEX IF NOT EXISTS idx_household_members_household_id ON household_members(household_id);

CREATE TABLE IF NOT EXISTS household_invites (
    id           UUID        PRIMARY KEY     DEFAULT uuid_generate_v4(),
    household_id UUID        NOT NULL        REFERENCES households(id) ON DELETE CASCADE,
    code         VARCHAR(32) NOT NULL UNIQUE,
    created_by   UUID        NOT NULL        REFERENCES users(id) ON DELETE CASCADE,
    expires_at   TIMESTAMPTZ NOT NULL,
    max_uses     INT,
    use_count    INT         NOT NULL DEFAULT 0,
    revoked_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

    CHECK(
        max_uses IS NULL OR use_count <= max_uses
    )
);

CREATE INDEX IF NOT EXISTS idx_household_invites_household_id ON household_invites(household_id);