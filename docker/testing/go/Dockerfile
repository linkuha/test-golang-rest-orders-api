FROM golang:1.19-alpine

MAINTAINER Alex Pudich <pudichas@gmail.com>

COPY ./docker/common/wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod 555 /usr/local/bin/wait-for-it

WORKDIR /app
COPY . .

RUN set -ex \
    && apk --no-cache add postgresql-client git \
    && go mod download && go mod tidy \
    && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
    && wget -q -O - https://git.io/JLbXn | sh -s \
    && wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.50.1 \
    && rm -rf /var/cache/apk/*


