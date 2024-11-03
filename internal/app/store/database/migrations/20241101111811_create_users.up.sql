SET statement_timeout = 0;

--bun:split

CREATE TABLE users (
    id uuid PRIMARY KEY,
    login character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE UNIQUE INDEX users_login_idx ON users(login);
