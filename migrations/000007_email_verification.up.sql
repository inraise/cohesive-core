ALTER TABLE users ADD COLUMN IF NOT EXISTS is_verified BOOLEAN NOT NULL DEFAULT false;
 
CREATE TABLE IF NOT EXISTS email_verifications (
    id         UUID         PRIMARY KEY     DEFAULT uuid_generate_v4(),
    user_id    UUID         NOT NULL        REFERENCES users(id) ON DELETE CASCADE,
    email      VARCHAR(100) NOT NULL,
    token_hash VARCHAR(64)  NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ  NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);
 
CREATE INDEX IF NOT EXISTS idx_email_verifications_user_id ON email_verifications(user_id);