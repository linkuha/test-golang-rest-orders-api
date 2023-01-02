.DEFAULT_GOAL := hello

include .env
export

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

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
	swag-api-v1
	CGO_ENABLED=0 go build -o apisrv cmd/api/main.go
	./apisrv -logdir ./log --

migrate-create:
	migrate create -ext sql -dir database/migrations 'migrate_name'

migrate-up:
	migrate -path database/migrations -database '$(DATABASE_URL)' up

migrate-down:
	migrate -path database/migrations -database '$(DATABASE_URL)' down

api-test:
	go test -v ./...

mock-testify:
	docker run --rm -v ${PWD}:/src -w /src/ vektra/mockery --keeptree --all

mock-testify-wincmd:
	docker run --rm -v %cd%:/src -w /src/ vektra/mockery --keeptree --all

mock-gomock:
	mockgen.exe -source=internal/domain/service/passwordEncryptor.go \
	      -destination=internal/domain/service/mocks/mock_passwordEncryptor.go

MOCKS_DESTINATION=internal/mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
mocks: internal/domain/repository/*/*
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do ./mockgen.exe -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done
.PHONY: mocks

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

linter-hadolint:
	git ls-files --exclude='Dockerfile*' --ignored | xargs docker run --rm -i hadolint/hadolint

hello:
	@echo "Welcome to commands shortcuts"

