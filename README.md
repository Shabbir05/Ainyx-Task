# Ainyx Task

A production-ready RESTful API built with **GoFiber** that manages users with `name` and `dob` (date of birth). Age is **never stored** — it is always calculated dynamically at runtime using Go's `time` package.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Tech Stack](#tech-stack)
3. [Project Structure](#project-structure)
4. [Prerequisites](#prerequisites)
5. [Local Setup (without Docker)](#local-setup-without-docker)
6. [Local Setup (with Docker)](#local-setup-with-docker)
7. [Running Migrations](#running-migrations)
8. [Regenerating SQLC Code](#regenerating-sqlc-code)
9. [API Endpoints](#api-endpoints)
10. [Running Unit Tests](#running-unit-tests)

---

## Project Overview

`user-api` exposes a CRUD REST API for user records. Key design decisions:

- **No stored age**: age is computed at request time by subtracting the user's DOB from today's date.
- **Layered architecture**: handler → service → repository → SQLC/MySQL.
- **Structured logging**: every request logged via [Uber Zap](https://github.com/uber-go/zap).
- **Input validation**: enforced via [go-playground/validator](https://github.com/go-playground/validator) with a custom `dob_date` rule.
- **Request tracing**: every response carries an `X-Request-ID` UUID header.

---

## Tech Stack

| Layer | Technology |
|---|---|
| HTTP framework | [GoFiber v2](https://gofiber.io/) |
| Database | MySQL 8.0 |
| Query layer | [SQLC](https://sqlc.dev/) |
| Logging | [Uber Zap](https://pkg.go.dev/go.uber.org/zap) |
| Validation | [go-playground/validator v10](https://github.com/go-playground/validator) |
| Config | [godotenv](https://github.com/joho/godotenv) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Containerisation | Docker + Docker Compose |

---

## Project Structure

```
.
├── cmd/server/main.go          # Application entry point
├── config/config.go            # Environment variable loading
├── db/
│   ├── migrations/             # SQL migration files (up + down)
│   └── sqlc/                   # SQLC config + generated Go code
├── internal/
│   ├── handler/                # HTTP handlers (parse → service → respond)
│   ├── logger/                 # Zap logger initialisation
│   ├── middleware/             # RequestID + RequestLogger middlewares
│   ├── models/                 # Request / response structs
│   ├── repository/             # DB calls via SQLC (interface + impl)
│   ├── routes/                 # Fiber route registration
│   ├── service/                # Business logic + age calculation
│   └── validator/              # Custom validator rules
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── sqlc.yaml
```

---

## Prerequisites

- **Go 1.22+** — https://go.dev/dl/
- **MySQL 8.0+** (for local, non-Docker setup)
- **Docker + Docker Compose** (for containerised setup)
- **golang-migrate CLI** — https://github.com/golang-migrate/migrate
- **SQLC CLI** — https://docs.sqlc.dev/en/latest/overview/install.html

---

## Local Setup (without Docker)

### 1. Clone and configure

```bash
git clone https://github.com/yourusername/user-api.git
cd user-api
cp .env.example .env
# Fill in DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, APP_PORT
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Run migrations

See [Running Migrations](#running-migrations) below.

### 4. Start the server

```bash
go run ./cmd/server
```

The API will be available at `http://localhost:3000`.

---

## Local Setup (with Docker)

### 1. Configure environment

```bash
cp .env.example .env
# Edit .env — DB_HOST will be overridden to "mysql" by docker-compose
```

### 2. Build and start all services

```bash
docker compose up --build
```

Docker Compose will:
1. Start MySQL with a healthcheck.
2. Wait until MySQL is healthy.
3. Build and start the Go application.

The API will be available at `http://localhost:3000`.

### 3. Run migrations inside Docker

```bash
# From the host, targeting the exposed MySQL port
migrate -path db/migrations \
        -database "mysql://root:secret@tcp(localhost:3306)/usersdb" up
```

### 4. Stop services

```bash
docker compose down
# To also remove the database volume:
docker compose down -v
```

---

## Running Migrations

Install the CLI first:

```bash
# macOS / Linux
brew install golang-migrate

# or via Go
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Apply all pending migrations (up)

```bash
migrate -path db/migrations \
        -database "mysql://root:secret@tcp(localhost:3306)/usersdb" up
```

### Roll back the last migration (down)

```bash
migrate -path db/migrations \
        -database "mysql://root:secret@tcp(localhost:3306)/usersdb" down 1
```

### Roll back all migrations

```bash
migrate -path db/migrations \
        -database "mysql://root:secret@tcp(localhost:3306)/usersdb" down
```

---

## Regenerating SQLC Code

Install the CLI:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Regenerate:

```bash
sqlc generate
```

This reads `sqlc.yaml`, parses `db/sqlc/query.sql` + the migration schema, and writes:

| File | Contents |
|---|---|
| `db/sqlc/db.go` | DBTX interface + Queries struct |
| `db/sqlc/models.go` | `User` model struct |
| `db/sqlc/query.sql.go` | All query implementations |

> **Note:** The generated files are committed to the repository for reproducibility. Regenerate only when you change `db/sqlc/query.sql` or the schema.

---

## API Endpoints

Base URL: `http://localhost:3000`

### Create a user

```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "dob": "1990-05-10"}'
```

**Response 201:**

```json
{ "id": 1, "name": "Alice", "dob": "1990-05-10" }
```

---

### Get a user by ID (includes age)

```bash
curl http://localhost:3000/users/1
```

**Response 200:**

```json
{ "id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35 }
```

---

### Update a user

```bash
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Updated", "dob": "1991-03-15"}'
```

**Response 200:**

```json
{ "id": 1, "name": "Alice Updated", "dob": "1991-03-15" }
```

---

### Delete a user

```bash
curl -X DELETE http://localhost:3000/users/1
```

**Response: 204 No Content**

---

### List users (paginated)

```bash
curl "http://localhost:3000/users?page=1&limit=10"
```

**Response 200:**

```json
{
  "data": [
    { "id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35 }
  ],
  "page": 1,
  "limit": 10,
  "total": 1
}
```

---

### Error Responses

| Scenario | Status | Body |
|---|---|---|
| Missing / malformed body | 400 | `{"error": "invalid request body"}` |
| Validation failure | 422 | `{"error": "validation failed", "details": {"name": "required"}}` |
| User not found | 404 | `{"error": "user not found"}` |
| Server error | 500 | `{"error": "internal server error"}` |

---

## Running Unit Tests

```bash
# Run all tests
go test ./...

# Run only the age-calculation unit tests with verbose output
go test -v ./internal/service/...

# Run with race detector
go test -race ./...
```

The `CalculateAge` unit tests live in [`internal/service/age_test.go`](internal/service/age_test.go) and cover:

- Birthday already passed this year
- Birthday has not yet occurred this year
- Birthday is today (age does not increment)
- Edge cases (leap years, newborns)

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `APP_PORT` | `3000` | Port the HTTP server listens on |
| `DB_HOST` | `localhost` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL user |
| `DB_PASSWORD` | — | MySQL password |
| `DB_NAME` | `usersdb` | MySQL database name |

Copy `.env.example` to `.env` and fill in values before running locally.
