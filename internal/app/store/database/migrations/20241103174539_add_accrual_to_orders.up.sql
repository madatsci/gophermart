SET statement_timeout = 0;

--bun:split

ALTER TABLE orders ADD COLUMN accrual decimal;
