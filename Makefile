.DEFAULT_GOAL := help

# env for migrate commands. NOTE! this affected to os.Getenv
include .env
# TODO for testing DB
include .env.testing
include .env.heroku
export

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

up: docker-up
down: docker-down
restart: down up
rebuild: swag-api-v1 docker-build restart

docker-up: ## docker up
	docker-compose up -d
docker-down: ## ..down and clean containers
	docker-compose down --remove-orphans
docker-down-clear: ## ..down and clean with volumes
	docker-compose down -v --remove-orphans
docker-pull: ## ..pull
	docker-compose pull
docker-build: ## ..build
	docker-compose build

swag-api-v1: ## generate OpenAPI for api v1
	swag init -g ./internal/delivery/httpserver/v1/router.go -o ./docs

api-run: ## run app
	go mod tidy && go mod download
	CGO_ENABLED=0 go run -tags automigrate cmd/apisrv/main.go -logdir ./log --

api-build-run: ## build OpenAPI, build and run app
	make swag-api-v1
	CGO_ENABLED=0 go build -o bin/apisrv cmd/apisrv/main.go
	./bin/apisrv -logdir ./log --

api-build-run-vendor: ## build OpenAPI, build with modules from vendor and run app
	make swag-api-v1
	go env -w CGO_ENABLED=0
	go build -mod vendor -o bin/apisrv cmd/apisrv/main.go
	./bin/apisrv -logdir ./log --

migrate-create: ## create dummy migrations file
	migrate create -ext sql -dir database/migrations 'migrate_name'

migrate-up: ## migrations apply
	migrate -path database/migrations -database '$(DATABASE_URL)' up

migrate-down: ## migrations rollback
	migrate -path database/migrations -database '$(DATABASE_URL)' down

heroku-migrate-up: ## migrations apply on heroku
	migrate -path database/migrations -database '$(HEROKU_DB_DSN)' up

heroku-migrate-down: ## migrations apply on heroku
	migrate -path database/migrations -database '$(HEROKU_DB_DSN)' down

api-test: ## run tests for api
	go test -v ./...

mock-testify:
	docker run --rm -v ${PWD}:/src -w /src/ vektra/mockery --keeptree --all

mock-testify-wincmd:
	docker run --rm -v %cd%:/src -w /src/ vektra/mockery --keeptree --all

mock-gomock:
	./bin/mockgen -source=internal/domain/service/passwordEncryptor.go \
	      -destination=internal/domain/service/mocks/mock_passwordEncryptor.go

MOCKS_DESTINATION=internal/mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
mocks: internal/domain/repository/*/* ## generate mocks to ./mocks (for repositories)
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do ./bin/mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done
.PHONY: mocks

cover: ## make coverage report
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

linter-hadolint: ## run linter for dockerfiles
	git ls-files --exclude='Dockerfile*' --ignored | xargs docker run --rm -i hadolint/hadolint

cd-build: api-build-image gateway-build-image db-build-image

db-build-image:
	docker --log-level=debug build --pull --file=docker/production/postgresql/Dockerfile --tag=${REGISTRY}/go-rest-api-postgres:${IMAGE_TAG} ./
gateway-build-image:
	docker --log-level=debug build --pull --file=docker/production/nginx-gateway/Dockerfile --tag=${REGISTRY}/go-rest-api-nginx:${IMAGE_TAG} ./
api-build-image:
	docker --log-level=debug build --pull --file=Dockerfile --tag=${REGISTRY}/go-rest-api-golang:${IMAGE_TAG} ./

cd-push: api-push-image gateway-push-image db-push-image

db-push-image:
	docker push ${REGISTRY}/go-rest-api-postgres:${IMAGE_TAG}
gateway-push-image:
	docker push ${REGISTRY}/go-rest-api-nginx:${IMAGE_TAG}
api-push-image:
	docker push ${REGISTRY}/go-rest-api-golang:${IMAGE_TAG}

ansible-build:
	docker-compose build ansible
ansible-console:
	docker-compose run --rm ansible sh
ansible-test-hosts:
	docker-compose run --rm ansible ansible all -a "/bin/echo Hello, world!"
ansible-playbook-test:
	docker-compose run --rm ansible ansible-playbook run_test.yml
ansible-add-deploy-ssh:
	docker-compose run --rm ansible ansible-playbook run_add_deploy_ssh.yml
ansible-registry-login:
	docker-compose run --rm ansible ansible-playbook run_docker_login.yml
ansible-playbook-setup:
	docker-compose run --rm ansible ansible-playbook run_setup.yml
ansible-site-facts:
	docker-compose run --rm ansible ansible app -m setup

deploy:
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'rm -rf api_${BUILD_NUMBER}'
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'mkdir api_${BUILD_NUMBER}'
	scp -o StrictHostKeyChecking=no -P ${DEPLOY_PORT} -i ${DEPLOY_KEY} docker-compose-production.yml deploy@${DEPLOY_HOST}:api_${BUILD_NUMBER}/docker-compose.yml
	scp -o StrictHostKeyChecking=no -P ${DEPLOY_PORT} -i ${DEPLOY_KEY} .env.production deploy@${DEPLOY_HOST}:api_${BUILD_NUMBER}/.env
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'cd api_${BUILD_NUMBER} && echo "REGISTRY=${REGISTRY}" >> .env'
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'cd api_${BUILD_NUMBER} && echo "IMAGE_TAG=${IMAGE_TAG}" >> .env'
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'cd api_${BUILD_NUMBER} && docker-compose pull'
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'docker network inspect api-orders >/dev/null 2>&1 || docker network create --driver=bridge --subnet=172.30.25.0/27 --gateway=172.30.25.1 api-orders'
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'cd api_${BUILD_NUMBER} && docker-compose up --build -d postgres'
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'cd api_${BUILD_NUMBER} && docker-compose run --rm api /wait'
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'cd api_${BUILD_NUMBER} && docker-compose up --build --remove-orphans -d'
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'rm -rf go-rest-api'
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'ln -sr api_${BUILD_NUMBER} go-rest-api'

rollback:
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'cd api_${BUILD_NUMBER} && docker-compose pull'
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'cd api_${BUILD_NUMBER} && docker-compose up --build --remove-orphans -d'
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'rm -rf go-rest-api'
	ssh -o StrictHostKeyChecking=no -p ${DEPLOY_PORT} -i ${DEPLOY_KEY} deploy@${DEPLOY_HOST} 'ln -sr api_${BUILD_NUMBER} go-rest-api'
