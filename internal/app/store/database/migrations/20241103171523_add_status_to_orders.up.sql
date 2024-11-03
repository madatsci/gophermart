SET statement_timeout = 0;

--bun:split

ALTER TABLE orders ADD COLUMN status character varying(255) NOT NULL DEFAULT 'NEW';
