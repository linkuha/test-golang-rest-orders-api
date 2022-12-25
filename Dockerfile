FROM golang:1.19-alpine as build

WORKDIR /app
COPY . .

RUN set -x \
  && apk --no-cache add git \
  && go mod tidy && go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -tags automigrate -a -o apisrv cmd/api/main.go

FROM alpine

MAINTAINER Alex Pudich <pudichas@gmail.com>

RUN apk --no-cache add postgresql-client bash \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY /docker/common/wait-for-it.sh ./
RUN chmod +x wait-for-it.sh

COPY --from=build /apisrv ./
COPY --from=build /.env ./
COPY --from=build /config/config.yml ./config/config.yml
COPY --from=build /database ./

CMD ["/app/apisrv"]

