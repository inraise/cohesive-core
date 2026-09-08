CREATE TABLE IF NOT EXISTS shopping_lists (
    id           UUID         PRIMARY KEY     DEFAULT uuid_generate_v4(),
    household_id UUID         NOT NULL        REFERENCES households(id) ON DELETE CASCADE,
    version      INT                  NOT NULL DEFAULT 1,
    name         VARCHAR(100)         NOT NULL CHECK(char_length(name) BETWEEN 1 AND 100),
    created_by   UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
 
    CHECK(
        created_at <= updated_at
    )
);
 
CREATE INDEX IF NOT EXISTS idx_shopping_lists_household_id ON shopping_lists(household_id);
 
CREATE TABLE IF NOT EXISTS shopping_list_items (
    id               UUID         PRIMARY KEY     DEFAULT uuid_generate_v4(),
    shopping_list_id UUID         NOT NULL        REFERENCES shopping_lists(id) ON DELETE CASCADE,
    version          INT                  NOT NULL DEFAULT 1,
    name             VARCHAR(200)         NOT NULL CHECK(char_length(name) BETWEEN 1 AND 200),
    quantity         VARCHAR(50),
    is_purchased     BOOLEAN              NOT NULL DEFAULT false,
    created_by       UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
 
    CHECK(
        created_at <= updated_at
    )
);
 
CREATE INDEX IF NOT EXISTS idx_shopping_list_items_list_id ON shopping_list_items(shopping_list_id);