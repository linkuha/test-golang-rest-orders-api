REST API of orders management ([pre-employment test](./ISSUE_RU.md))

![GitHub Actions tests status](https://github.com/linkuha/test-golang-rest-orders-api/actions/workflows/go.yml/badge.svg)
![Dockerfile build status](https://github.com/linkuha/test-golang-rest-orders-api/actions/workflows/docker-image.yml/badge.svg)
[![TravisCI Build Status](https://app.travis-ci.com/linkuha/test-golang-rest-orders-api.svg?branch=main)](https://app.travis-ci.com/linkuha/test-golang-rest-orders-api)
[![License](https://badgen.net/badge/license/MIT/blue)](https://github.com/linkuha/test-golang-rest-orders-api/blob/main/LICENSE)
[![Telegram](https://badgen.net/badge/icon/telegram?icon=telegram&label=@linkuha)](https://t.me/linkuha)


DB: PostgreSQL

This is first REST api service in Go implemented by me.
Before that I was implemented from scratch in Go only backend daemon (concurrent parser, tuned by config and environment, 
scaled and pulling via docker-compose from registry).

## Functionality

You can use OpenAPI configuration for research from testing section below.

## Comments

* **Router**. Intentionally selected for routing framework GIN `gin-gonic/gin`. There is an alternative to `labstack/echo`. 
      A good `gorilla/mux` multiplexer has fewer features, and unfortunately was archived before the service was written (December 22).
      All of them have support from swagger (`swaggo/swag`).
* **Middleware**.
  * Implemented middleware for authorizing user requests to API via JWT token. There is another approach - storing the authorization session ID / token in cookies.
  * Added middleware for generating and saving response in header - request ID (`gin-contrib/requestid`)
  * Implemented middleware for logging request execution time and its status. Useful, since client and detailed internal errors are logged, with reference to request_id.
* **DB**.
  * Intentionally selected standard lib for communicating with the database (`database/sql`, `lib/pq` driver).
  * You can use the more powerful and convenient `jackc/pgx` and the convenient fluent query builder `masterminds/squirrel`.
  * Migrations are raised automatically when the application is launched via `golang-migrate/migrate/v4`.
  * There are a lot of ORM's (`go-gorm/gorm` or others), I don't use it because this is not a GO-friendly approach, we lose speed due to a lot of reflection.
* **Logic**. 
  * Logic place - in the user cases, but somewhere the repository is used directly in handlers, left it for todo.
  * The User entity intentionally does not correspond to the test issue. I decided to complicate, to separate the essence of the profile from the credits.
    There is no need to carry all the data with you every time, the profile can expand. In this case, the profile can be created in another step after registration.
  * In an amicable way, DTO should be passed to the handlers inputs, and then their data being carried should be mapped into entities.
      This will make it easier to learn OpenAPI, as now not all fields are required / ignored. But it also breeds a lot of extra code, left it for todo.
  * In the product repository, the input argument is not a model, but a DTO ProductUpdateInput, was made to be able to send only those fields, which we want to change.
      Alternatively, for other small entities simply overwrites all fields.
* **Validation**. 
  * Validation in models, `go-ozzo/ozzo-validation/v4` package is used in order not to duplicate boilerplate code.
  * For validation at the structure binding level from the request, standard tags are used - the GIN validator - `go-playground/validator`
* **Tests**.
  * Unit tests for the Entities layer (validation of models is checked).
  * Unit tests for the Repositories layer to work with the User entity (the database is mocked via `DATA-DOG/go-sqlmock`).
  * Unit tests for the Use Cases layer to work with the User entity (mocking repositories via `golang/mock`).
  * Integration tests for the Controllers layer to work with the Products entity.
    Testing of handlers will generate a lot of code, even taking into account the development of tabular tests.
  * I don't like that the real test database is not used for tests. Thus, the execution of SQL queries / transactions is not checked, left it for todo.
* **Errors**. Custom types of errors worked out, for the client - more common errors, correct response statuses. For debugging - internal, with the ability to get a stack.
* **Context**. 
  * Added context forwarding into the repositories layer to cancel database queries on failures.
  * Graceful shutdown.

## Development environment

Install GO 1.19 on local machine https://go.dev/doc/install

Install tools on local machine with go:
* swagger https://github.com/swaggo/swag for generating API documentation
* migrate-tool  https://github.com/golang-migrate/migrate for make migrations
```
go install github.com/swaggo/swag/cmd/swag@latest 
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Install tools manually into ./bin directory: mockgen, ngrok (optionally).

PS: see [Makefile](./Makefile) for handy useful short commands (e.g. for build, run, test), get help with `make help` command.

For example, generate dummy file for migrations, located in `database/migrations` directory: `make migrate-create`.

### Build for development

1. You can build and run app with help of Makefile: `make api-build-run`. 

It will generate swag documentation, build binary and run with local installed compiler.

You must provide own env params like database connection in [.env](./.env) file (based on .env.sample).

The app will parse this params via wide-known package Viper.

PS. Фиксацию зависимостей в ./vendor для git можно производить через `go mod vendor`. Тогда билдить командой `make api-build-run-vendor`.

2. You can build & run app and ready-to-use database (Postgres with enabled extension for UUID generating)
for tests as production-like version of app with Docker Compose: `make up`. 

In that case you must not provide env params for database, it will be overwritten in [docker-compose.yml](./docker-compose.yml).

Supporting of `DATABASE_URL` was used for ability of easy deploy on Heroku.

### Build for production and Deploy

1. Prepare configs

Provide own env params in [.env.production](./.env.production) file (based on .env.sample) with production params (e.g., disable debug modes (GIN, logger), logging in file).

For setup clean server, install Docker and pass SSH key to server will be using Ansible. 
* Modify next files for your specific variables (e.g. domain): `./provisioning/ansible/vars/*`. 
* And modify params of servers for deploy: `./provisioning/ansible/hosts.yml` (group of servers - app) based on hosts.yml.dist sample in same directory.

For deploy to server will be used Docker Registry and Makefile.

* Modify domain in gateway (nginx) image: `./docker/production/nginx-gateway/sites-enabled/api.conf`
* Add OpenSSH keypair to `./provisioning/ssh_keys/` and name it as: `test_server_openssh_private_key` and `test_server_openssh_public_key`.

PS. This key names using in Ansible configs for connection to server and setup soft on behalf of root. 
This also using for upload private key to deployment server for user "deploy", to make deploy not from root (but for simplify - here used the same keys).

* Export common vars of your deployment server and docker repo (replace for your `REGISTRY=DOMAIN[:PORT]`):
```
export REGISTRY=docker.pudich.ru
export IMAGE_TAG=1
export DEPLOY_HOST=37.228.116.222
export DEPLOY_PORT=22 
export DEPLOY_KEY=$(pwd)/provisioning/ssh_keys/test_server_openssh_private_key
chmod 600 ${DEPLOY_KEY}
```
It keeps in terminal session and usable for Makefile commands.
You can use constant IMAGE_TAG=1 for rebuild images and rewrite it, for no tracking version and disk space economy.

* Connect to your docker registry (own service / GitLab registry / DockerHub) on local:
```
docker login ${REGISTRY} -u <PASTE USERNAME>
password: <PASTE PASSWORD>
```

2. Link in your DNS zone's A record for used domain with server's IP for Certbot correct resolving.

3. Prepare server via Ansible

- `make ansible-build` - Build Ansible with configs.
- `make ansible-playbook-test` - Run test task.
- `make ansible-playbook-setup` - Setup new server. Run it once, or if soft installation tasks changed.
- `make ansible-add-deploy-ssh` - Setup SSH key to user "deploy".
- `make ansible-registry-login` - Login to Docker Registry at server on behalf of user "deploy". Interactive input...

4. Build images for Delivery into Docker Registry

Build - `make cd-build`, push to repo - `make cd-push`.

5. Deploy via Makefile command

For create new release version: `BUILD_NUMBER=1 make deploy`, for rollback: `BUILD_NUMBER=0 make rollback`.
It will pull images from Registry to new directory, switch symlink on it and start services via Docker Compose. 
Modify to your version number. Not to be confused `BUILD_NUMBER` with `IMAGE_TAG`.


### Deploy to Heroku

Install Heroku CLI

`heroku login` - grant permissions to CLI to account via browser.

Create app in https://dashboard.heroku.com/ and check is ID (e.g. my - test-golang-rest)

I connected my GitHub repo to Heroku app, you can fork and do it.

`heroku git:remote -a <ID>` - link directory with local git to heroku's git repo.

If you connect own git repo, use section Deploy -> Manual deploy / Automatics, for pull changes from original repo.

If not - commit changes locally and use `git push heroku master`.

Add environment variables into section Settings -> Config Vars (like in your .env.production). 
But `DATABASE_URL` will be set automatically, if you install PostgreSQL addon on Heroku (section Resources) and connect it to this app.
Copy `DATABASE_URL` as `HEROKU_DSN` to .env.heroku based on [.env.heroku.sample](./.env.heroku.sample) for migration commands.
Don't set up `PORT` env, this is Heroku's responsibility.

Add env `BUILD_DIR=bin` to maintain the correct folder structure as on the local machine and Procfile declarations.

Don't forget make migrations on Heroku PostgreSQL.
Before that create extension on Heroku DB for that open section More -> Run console. Run psql terminal and do query:
```
~$ psql $DATABASE_URL
=> CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION
=> \dx
(check)
(press q)
=> exit
```

And now in local terminal run: `make heroku-migrate-up`. If conflict dirty - repeat after query in psql: 
```SQL
delete from schema_migrations where dirty != 'false';
```

Because this is not Dockerfile based app on heroku and app builds without tags by default (e.g. -tags automigrate = for include go:build commented file to compilation).

After deploy and app was built, check logs - button More -> View logs.

## Testing

You can test requests when service is running in swagger web from route `localhost:3000/swagger/index.html`.

You can share opened on local machine port to world via https://ngrok.com/ `ngrok http 3000`

Or you can import to Postman generated [OpenAPI 2.0 config](./docs/swagger.json) and check & execute allowed requests. 

Register with /auth/sign-up, login with /auth/sign-in and use received token for other requests: set up Authorize value: `Bearer <token>`

You can run **unit tests** with helping Makefile command: `make api-test`

Coverage: TODO fix

## Help

If error on Windows: listen tcp 0.0.0.0:<PORT>: bind: An attempt was made to access a socket in a way forbidden by its access permissions.
```
# PowerShell / cmd administrator:
net stop winnat
net start winnat
```

## TODO:
- add pagination to GetAll handlers - (headers Pagination-, pages fields in json, response 206 Partial Content)
or another variant - return collections of identifiers, than by get(ids) - return requested
- add requests throttling (aviddiviner/gin-limit, axiaoxin-com/ratelimiter middleware?) (response 429 Too Many Requests)
- add linters
  * dotenv-linter https://github.com/dotenv-linter/dotenv-linter
  * golangci-lint https://golangci-lint.run/usage/install/
  * sql linters
  * dockerfile linters (hadolint)
- test commands and deploy from OS linux


## Demo

Now my demos are running (working while I pay for it):

* https://go-rest-test.pudich.ru/swagger/index.html - deployed on clean Selectel's VPS with help of Ansible and my own docker registry.
* https://test-golang-rest.herokuapp.com/swagger/index.html#/ - deployed to Heroku with PostgreSQL paid addon.

It works :)
