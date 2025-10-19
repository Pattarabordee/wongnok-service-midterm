create-postgres:
	docker run --name db-postgres -e POSTGRES_PASSWORD=pass2word -p 5432:5432 -d postgres:17-alpine

start-postgres:
	docker start db-postgres

stop-postgres:
	docker stop db-postgres

migrate-rating-table:
	goose -dir migrations create add_ratings_table sql

migrate-up:
	goose up

migrate-down:
	goose down

migrate-down-all:
    goose reset

migrate-status:
	goose status
	
mockery-install:
	go install github.com/vektra/mockery/v3@v3.5.0

mockery:
	mockery

run-test:
	go test ./internal/... 

// craete folder coverage if not exists
run-test-coverage:
	go test -short -v -coverprofile coverage/cover.out ./internal/...

coverage-html:
	go tool cover -html=coverage/cover.out -o coverage/cover.html

swagger:
	swag init \
	-g ./cmd/server/main.go \
	-o ./cmd/server/docs \
	--parseInternal
