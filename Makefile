# make name="create_users" db_create_sql
.PHONY: db_create_sql
db_create_sql:
	go run ./cmd/gophermart-cli/main.go -d 'postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable' db create_sql $(name)

.PHONY: db_migrate
db_migrate:
	go run ./cmd/gophermart-cli/main.go -d 'postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable' db migrate

.PHONY: db_rollback
db_rollback:
	go run ./cmd/gophermart-cli/main.go -d 'postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable' db rollback