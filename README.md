# go-musthave-diploma-tpl

Шаблон репозитория для индивидуального дипломного проекта курса «Go-разработчик»

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без
   префикса `https://`) для создания модуля

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.

# Development

## Run Database

The database can be started via docker container using this command:

```bash
docker run --name gophermart-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=gophermart -p 5432:5432 -d postgres
```

## Build App

```bash
cd cmd/gophermart && go build -o gophermart *.go
```

or

```bash
make build
```

## Run App

Database connection is required to run the app. Some examples of how you can run the app (see Configuration below):

### Run app with database URI

```bash
./cmd/gophermart/gophermart -d 'postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable'
```

or

```bash
make run
```

## Configuration

App can be configured via flags and/or environment variables. If both flag and environment variable are set for the same parameter, environment variable prevails.

### `-a`, `RUN_ADDRESS`
Address and port to run server in the form of host:port.

### `-d`, `DATABASE_URI`
Database URI.

### `-r`, `ACCRUAL_SYSTEM_ADDRESS`
Accrual system address.

## Migrations

Migrations are implemented with [bun](https://bun.uptrace.dev/guide/migrations.html). You can run migrations using CLI app.

### CLI command

Create migration:

```bash
go run ./cmd/gophermart-cli/main.go -d 'postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable' db create_sql <migration_name>
```

Run migrations:

```bash
go run ./cmd/gophermart-cli/main.go -d 'postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable' db migrate
```

Rollback migrations:

```bash
go run ./cmd/gophermart-cli/main.go -d 'postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable' db rollback
```

### Makefile

Create migration:

```bash
make name="create_users" db_create_sql
```

Run migrations:

```bash
make db_migrate
```

Rollback migrations:

```bash
make db_rollback
```

# API Examples

## Register A New User

```bash
curl -i -X POST http://localhost:8080/api/user/register \
   -H "Content-Type: application/json" \
   -d '{
      "login":"john_doe",
      "password":"my_secret_password"
   }'

# Response:
HTTP/1.1 201 Created
Date: Fri, 01 Nov 2024 14:07:47 GMT
Content-Length: 0
```
