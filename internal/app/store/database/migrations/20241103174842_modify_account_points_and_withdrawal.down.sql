SET statement_timeout = 0;

--bun:split

ALTER TABLE accounts
ALTER COLUMN current_points_total TYPE int,
ALTER COLUMN withdrawn_total TYPE int;
