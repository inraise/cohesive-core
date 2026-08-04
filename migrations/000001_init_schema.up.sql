CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id            UUID         PRIMARY KEY     DEFAULT uuid_generate_v4(),
    version       INT                 NOT NULL DEFAULT 1,
    email         VARCHAR(255) UNIQUE NOT NULL CHECK(char_length(email) BETWEEN 5 AND 100),
    password_hash VARCHAR(255)        NOT NULL,
    first_name    VARCHAR(100)        NOT NULL CHECK(char_length(first_name) BETWEEN 1 AND 100),
    last_name     VARCHAR(100)                 CHECK(char_length(last_name) BETWEEN 1 AND 100),
    age           INT                          CHECK(age >= 0 AND age <= 130),
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ,

    CHECK(
        created_at <= updated_at
    )
);
