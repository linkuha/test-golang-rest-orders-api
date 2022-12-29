REST API of orders management (testing task)

Works with database: PostgreSQL

This is first REST api service in Go implemented by me.
Before that I was implemented from scratch in Go only backend daemon (concurrent parser, tuned by config and environment, 
scaled and pulling via docker-compose from registry).

## Functionality

You can use OpenAPI configuration for research from testing section below.

## Comments

* Намерено выбран для маршрутизации фреймворк GIN, т.к. обладает хорошоей производительностью и прост.
Другой хороший gorilla/mux к сожалению переведён в архив до начала написания сервиса (декабрь 22), поэтому не был выбран.
* Намерено выбрана стандартная либа для общения с БД (database/sql, драйвер lib/pq). 
Можно использовать более мощный и удобный jackc/pgx и удобный fluent билдер запросов (masterminds/squirrel).
Но это тестовый пример, этого достаточно.
* Не использую ORM (gorm или другие), т.к. это не очень GO-friendly подход, тяреем по скорости (много рефлексии)
* Сущность user намерено не соответствует ТЗ. Решил усложнить, отделить сущность профиля от кредов. 
Нет нужды таскать с собой все данные, профиль может расширяться. В этом кейсе профиль можно создать другим этапом после регистрации.
Апи закрыто. Есть другой подход - хранение ID авторизационной сессии в куках, но мне он меньше нравится
* Реализован middleware для авторизации запросов пользователя к API через JWT токен
* Валидация в моделях, используется пакет go-ozzo/ozzo-validation/v4, чтобы не дублировать похожий код.
* Для валидации на уровне биндинга структуры из запроса используется стандартные теги валидации - валидатора GIN - go-playground/validator
* В репозитории продукта входной аргумент - не модель, а DTO ProductUpdateInput, был сделан для возможности присылать только те поля,
которые мы хотим изменить. Новый для меня подход, раньше просто перезаписывал все поля.
* Будет написано по 3 теста на слой. Сейчас написаны только для сущностей (проверяется валидация модельки), остальное в TODO

## Development environment

Install GO 1.19 on local machine https://go.dev/doc/install

Install tools on local machine:
* swagger https://github.com/swaggo/swag for generating API documentation
* migrate-tool  https://github.com/golang-migrate/migrate for make migrations
```
go install github.com/swaggo/swag/cmd/swag@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

PS: see [Makefile](./Makefile) for handy useful short commands (e.g. for build/run/test)

### Build

1. You can build and run app with help of Makefile: `make api-build-run`. 
It will generate swag documentation, build binary and run with local installed compiler. 
You must provide own env params like database connection in [.env](./.env) file. 
App parse this params via Viper wide-known packed.

2. You can build & run app and ready-to-use database (Postgres with enabled extension for UUID generating)
for tests as production-similar version of app with Docker Compose: `make up`. 
In that case you must not provide env params for database, it will be overwritten in [docker-compose.yml](./docker-compose.yml).
Supporting of DATABASE_URL was used for ability of easy deploy on Heroku. 

## Testing

You can test requests when service is running in swagger web from route `localhost:3000/swagger/index.html`.

You can share opened on local machine port to world via https://ngrok.com/ `ngrok http 3000`

Or you can import to Postman generated [OpenAPI 2.0 config](./docs/swagger.json) and check & execute allowed requests. 


You can run **unit tests** with helping Makefile command: `make api-test`

Coverage: `todo`

## Deploy

todo

## TODO:
- добавить теги binding в структуры моделей для GIN парсинга
- перепроверить проброс ошибок, правильность их обертывания, создать новые типы доменных ошибок если нужно
- добавить проброс контекста в слой репозиториев, чтобы отменять запросы к БД при сбоях
- внедрить ответ со статусом http.StatusUnprocessableEntity
- написать тесты 3-х репозиториев через моки/тестовое хранилище
- написать тесты 3-х хэндлеров апи 
- перепроверить всё апи через postman, фиксить баги если будут
- добавить middleware для CORS (опционально)
- добавить middleware для генерации и сохранении в хэдэр ответа ID запроса (опционально)
- добавить middleware для записи ID запроса в лог, для связки с ошибками и отладки времени генерации ответа (опционально)
- задеплоить сервис в heroku для демо
- сделать provisioning конфигурации ansible чтоб на любую чистую виртуалку задеплоить по команде можно было для демо (опционально)
- поправить мультистадийный докер образ на scratch
- протестить линтеры

### Future 
add and setup linters:
* dotenv-linter https://github.com/dotenv-linter/dotenv-linter
* golangci-lint https://golangci-lint.run/usage/install/
* линтеры sql ?
* hadolint для dockerfiles
