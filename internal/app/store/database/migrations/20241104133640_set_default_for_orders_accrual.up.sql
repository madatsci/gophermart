SET statement_timeout = 0;

--bun:split

ALTER TABLE orders ALTER COLUMN accrual SET DEFAULT 0;
