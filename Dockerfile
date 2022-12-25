FROM golang:1.19 as build

WORKDIR /app
COPY . .

RUN set -x \
  && apk --no-cache add git \
  && go mod tidy && go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -o apisrv cmd/api/main.go

FROM alpine

MAINTAINER Alex Pudich <pudichas@gmail.com>

RUN apk --no-cache add postgresql-client \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY /docker/common/wait-for-it.sh ./
RUN chmod +x wait-for-it.sh

COPY --from=build /project/apisrv ./
COPY --from=build /.env ./
COPY --from=build /config/config.yaml ./config/config.yaml

CMD ["/app/apisrv"]

