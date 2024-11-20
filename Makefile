.PHONY: build
build:
	cd cmd/gophermart && go build -o gophermart *.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: run
run:
	./cmd/gophermart/gophermart -d 'postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable'

.PHONY: test_unit
test_unit:
	go test -v -cover ./internal/app/...
	go test -v -cover ./pkg/...

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