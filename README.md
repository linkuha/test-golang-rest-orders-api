REST API of orders management ([pre-employment test](./ISSUE.md))

[![License](https://badgen.net/badge/license/MIT/blue)](https://github.com/linkuha/test-golang-rest-orders-api/blob/main/LICENSE)
[![Telegram](https://badgen.net/badge/icon/telegram?icon=telegram&label=@linkuha)](https://t.me/linkuha)

DB: PostgreSQL

This is first REST api service in Go implemented by me.
Before that I was implemented from scratch in Go only backend daemon (concurrent parser, tuned by config and environment, 
scaled and pulling via docker-compose from registry).

## Functionality

You can use OpenAPI configuration for research from testing section below.

## Comments

* Намерено выбран для маршрутизации фреймворк Джин (GIN).
Хороший мультиплексор gorilla/mux имеет меньше возможностей, и к сожалению, переведён в архив до начала написания сервиса (декабрь 22), 
поэтому не был выбран. Существует альтернатива labstack/echo.
* Намерено выбрана стандартная либа для общения с БД (database/sql, драйвер lib/pq). 
Можно использовать более мощный и удобный jackc/pgx и удобный fluent билдер запросов (masterminds/squirrel).
Но это тестовый пример, этого достаточно.
* Миграции подымаются автоматически при запуске приложения через golang-migrate/migrate/v4. 
* Не использую ORM (gorm или другие), т.к. это не очень GO-friendly подход, тяреем по скорости (много рефлексии)
* Сущность user намерено не соответствует ТЗ. Решил усложнить, отделить сущность профиля от кредов. 
Нет нужды таскать с собой все данные, профиль может расширяться. В этом кейсе профиль можно создать другим этапом после регистрации.
* Реализован middleware для авторизации запросов пользователя к API через JWT токен. Есть другой подход - хранение ID авторизационной сессии в куках.
* Валидация в моделях, используется пакет go-ozzo/ozzo-validation/v4, чтобы не дублировать похожий код.
* Для валидации на уровне биндинга структуры из запроса используется стандартные теги - валидатора GIN - go-playground/validator
* В репозитории продукта входной аргумент - не модель, а DTO ProductUpdateInput, был сделан для возможности присылать только те поля,
которые мы хотим изменить. Как альтернатива, в других малых моделях просто перезаписываются все поля.
* По-хорошему, в хэндлеры должны передаваться на вход - DTO, а затем их несомые данные мапиться в сущности. 
Это упростит изучение OpenAPI, т.к. сейчас не все поля обязательны и игнорируются. Но также наплодит
много дополнительного кода. Поэтому здесь в основном просто передаются заведомо валидные сущности.
* Не во всех хэндлерах логика вынесена в юзкейс, а используется репозиторий напрямую - доделать.
* Написаны unit-тесты для слоя Entities (проверяется валидация моделей).
* Написаны unit-тесты для слоя Repositories для работы с сущностью User (база мокнута через DATA-DOG/go-sqlmock).
* Написаны unit-тесты для слоя Use Cases для работы с сущностью User (моки репозиториев через golang/mock).
* Написаны integration-тесты для слоя Controllers для работы с сущностью Products. 
Дальнейшее тестирование хендлеров - наплодит много кода, даже с учетом проработки табличных тестов.
* Не нравится, что для тестов не используется реальная тестовая БД. Таким образом не проверяется само выполнение запросов/транзакций.
* Проработаны кастомные типа ошибок, для клиента - более общие ошибки, корректные статусы ответов. 
Для дебага - внутренние, с возможностью получения стека. 
* Добавлен middleware для генерации и сохранении в хэдэр ответа - ID запроса (gin-contrib/requestid)
* Реализован middleware для логирования времени исполнения запроса, его статуса. 
Полезно, так как логируются клиентская и подробная внутренняя ошибки, с привязкой к request_id.
* Добавлен проброс контекста в слой репозиториев, чтобы отменять запросы к БД при сбоях. 
* Graceful shutdown.

## Development environment

Install GO 1.19 on local machine https://go.dev/doc/install

Install tools on local machine:
* swagger https://github.com/swaggo/swag for generating API documentation
* migrate-tool  https://github.com/golang-migrate/migrate for make migrations
```
go install github.com/swaggo/swag/cmd/swag@latest 
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

PS: see [Makefile](./Makefile) for handy useful short commands (e.g. for build/run/test), get help with `make help` command.

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
- добавить пагинацию для хэндлеров GetAll - возвращающих коллекции (заголовками Pagination-, полями json и ответ 206 Partial Content)
так же можно переделать на такой вариант - для getall вовзращать коллекцию идентификаторов. затем по get(ids) - возвращать нужные (опционально)
- добавить настройки троттлинга запросов к API (aviddiviner/gin-limit, axiaoxin-com/ratelimiter middleware?) (ответ 429 Too Many Requests)
- протестить линтеры 
- протестить команды и весь процесс деплоя с линукса
- запустить тесты в github actions

### Future 
add and setup linters:
* dotenv-linter https://github.com/dotenv-linter/dotenv-linter
* golangci-lint https://golangci-lint.run/usage/install/
* линтеры sql ?
* hadolint для dockerfiles


## Demo

Now my demos are running (working while I pay for it):

* https://go-rest-test.pudich.ru/swagger/index.html - deployed on clean Selectel's VPS with help of Ansible and my own docker registry.
* https://test-golang-rest.herokuapp.com/swagger/index.html#/ - deployed to Heroku with PostgreSQL paid addon.

It works :)
