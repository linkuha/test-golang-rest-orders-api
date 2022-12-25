FROM golang:1.19-alpine

MAINTAINER Alex Pudich <pudichas@gmail.com>

COPY ./docker/common/wait-for-it.sh /usr/local/bin/wait-for-it

RUN set -ex \
    && chmod 555 /usr/local/bin/wait-for-it \
    && apk --no-cache add postgresql-client git bash coreutils

WORKDIR /app
COPY . .

RUN go mod download && go mod tidy \
    && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
    && rm -rf /var/cache/apk/*

RUN CGO_ENABLED=0 GOOS=linux go build -tags automigrate -a -o apisrv cmd/api/main.go
