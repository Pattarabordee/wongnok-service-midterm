# wongnok-service

## 1. Install PostgreSQL

### Write compose.yml

```yaml
services:
  postgres:
    image: postgres:17
    ports:
      - 5432:5432
    volumes:
      - ./data:/var/lib/postgresql/data
    restart: always
    environment:
      - POSTGRES_PASSWORD=212224
```

### Run compose up

```sh
docker compose up -d
```

## 2. Setup environment

### Write .env file

```sh
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://postgres:pass2word@localhost:5432/wongnok
GOOSE_MIGRATION_DIR=./migrations
```

## 3. Migration

```sh
goose up
```

## 4. Testing

```sh
go test ./internal/...
```

## 5. Start server

```sh
go run cmd/server/main.go
```

## 6. Swagger

[Swagger](http://localhost:8000/swagger/index.html)

## 7. Authorization

[Login / Register](http://localhost:8000/api/v1/login)

[Logout](http://localhost:8000/api/v1/logout)
