SET statement_timeout = 0;

--bun:split

CREATE TABLE transactions (
    id uuid PRIMARY KEY,
    account_id uuid NOT NULL,
    amount decimal NOT NULL,
    order_number character varying(255) NOT NULL,
    direction character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

--bun:split

ALTER TABLE transactions ADD CONSTRAINT account_id_constraint FOREIGN KEY (account_id) REFERENCES accounts(id);

--bun:split

CREATE INDEX transactions_account_direction_idx ON transactions(account_id, direction);
