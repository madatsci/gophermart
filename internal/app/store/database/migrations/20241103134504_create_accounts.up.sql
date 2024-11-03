SET statement_timeout = 0;

--bun:split

CREATE TABLE accounts (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    current_points_total int NOT NULL DEFAULT 0,
    withdrawn_total int NOT NULL DEFAULT 0,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

--bun:split

CREATE UNIQUE INDEX accounts_user_id_idx ON accounts(user_id);

--bun:split

ALTER TABLE accounts ADD CONSTRAINT user_id_constraint FOREIGN KEY (user_id) REFERENCES users(id);
