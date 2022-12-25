Install tools on local:
* swagger https://github.com/swaggo/swag
* migrate-tool  https://github.com/golang-migrate/migrate
```
go install github.com/swaggo/swag/cmd/swag@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

This will be installed in docker image:
* migrate-tool  https://github.com/golang-migrate/migrate
* dotenv-linter https://github.com/dotenv-linter/dotenv-linter
* golangci-lint https://golangci-lint.run/usage/install/

