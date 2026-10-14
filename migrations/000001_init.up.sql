CREATE TABLE users (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username      TEXT        NOT NULL,
    password_hash TEXT        NOT NULL,
    balance       INTEGER     NOT NULL DEFAULT 1000,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT users_username_key    UNIQUE (username),
    CONSTRAINT users_username_length CHECK (char_length(username) BETWEEN 1 AND 64),
    CONSTRAINT users_balance_nonneg  CHECK (balance >= 0)
);

CREATE TABLE merch (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name          TEXT        NOT NULL,
    price         INTEGER     NOT NULL,

    CONSTRAINT merch_name_length CHECK (char_length(name) BETWEEN 1 AND 128),
    CONSTRAINT merch_price_nonneg CHECK (price > 0),
    CONSTRAINT merch_name_key UNIQUE (name)
);

INSERT INTO merch (name, price) VALUES
('t-shirt', 80),
('cup', 20),
('book', 50),
('pen', 10),
('powerbank', 200),
('hoody', 300),
('umbrella', 200),
('socks', 10),
('wallet', 50),
('pink-hoody', 500)
;
CREATE TABLE inventory (
    user_id       BIGINT NOT NULL,
    merch_id      BIGINT NOT NULL,
    quantity      INTEGER NOT NULL, CHECK (quantity >= 0),

    CONSTRAINT inventory_pkey PRIMARY KEY (user_id, merch_id),
    CONSTRAINT inventory_user_fk  FOREIGN KEY (user_id)  REFERENCES users (id),
    CONSTRAINT inventory_merch_fk FOREIGN KEY (merch_id) REFERENCES merch (id)

);

CREATE TABLE transfers (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    from_user_id  BIGINT NOT NULL,
    to_user_id    BIGINT NOT NULL,
    amount        INTEGER NOT NULL CHECK (amount > 0),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT transfers_from_user_fk FOREIGN KEY (from_user_id) REFERENCES users (id),
    CONSTRAINT transfers_to_user_fk   FOREIGN KEY (to_user_id)   REFERENCES users (id),
    CONSTRAINT transfers_from_to_check CHECK (from_user_id != to_user_id)
);

CREATE INDEX transfers_from_user_idx ON transfers (from_user_id);
CREATE INDEX transfers_to_user_idx   ON transfers (to_user_id);