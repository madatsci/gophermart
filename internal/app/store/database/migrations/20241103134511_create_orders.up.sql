SET statement_timeout = 0;

--bun:split

CREATE TABLE orders (
    id uuid PRIMARY KEY,
    account_id uuid NOT NULL,
    number character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

--bun:split

CREATE UNIQUE INDEX orders_number_idx ON orders(number);

--bun:split

ALTER TABLE orders ADD CONSTRAINT account_id_constraint FOREIGN KEY (account_id) REFERENCES accounts(id);
