.DEFAULT_GOAL := hello

include .env
export

up: docker-up
down: docker-down
restart: down up

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down --remove-orphans

docker-down-clear:
	docker-compose down -v --remove-orphans

docker-pull:
	docker-compose pull

docker-build:
	docker-compose build

swag-api-v1:
	swag init -g ./internal/delivery/httpserver/v1/router.go -o ./docs

api-run:
	go mod tidy && go mod download
	CGO_ENABLED=0 go run -tags automigrate cmd/api/main.go -logdir ./log --

api-build-run:
	CGO_ENABLED=0 go build -o apisrv cmd/api/main.go
	./apisrv -logdir ./log --

api-test:
	go test -v ./...

mock:
	docker run -v "$PWD"/project:/src -w /src vektra/mockery --all

migrate-create:
	migrate create -ext sql -dir database/migrations 'migrate_name'

migrate-up:
	migrate -path database/migrations -database '$(DATABASE_URL)' up

migrate-down:
	migrate -path database/migrations -database '$(DATABASE_URL)' down

linter-hadolint:
	git ls-files --exclude='Dockerfile*' --ignored | xargs docker run --rm -i hadolint/hadolint

hello:
	@echo "Welcome to commands shortcuts"

