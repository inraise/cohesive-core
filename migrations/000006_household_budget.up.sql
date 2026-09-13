CREATE TABLE IF NOT EXISTS household_transactions (
    id           UUID         PRIMARY KEY     DEFAULT uuid_generate_v4(),
    household_id UUID         NOT NULL        REFERENCES households(id) ON DELETE CASCADE,
    version      INT                  NOT NULL DEFAULT 1,
    type         VARCHAR(20)          NOT NULL CHECK(type IN ('deposit', 'expense')),
    amount       BIGINT               NOT NULL CHECK(amount > 0),
    description  VARCHAR(300),
    created_by   UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CHECK(
        created_at <= updated_at
    )
);

CREATE INDEX IF NOT EXISTS idx_household_transactions_household_id ON household_transactions(household_id);
CREATE INDEX IF NOT EXISTS idx_household_transactions_created_by ON household_transactions(created_by);