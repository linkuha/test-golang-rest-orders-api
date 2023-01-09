FROM golang:1.19-alpine as build

WORKDIR /app
COPY . .

RUN set -x \
  && apk --no-cache add git \
  && go mod tidy && go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -tags automigrate -a -o apisrv cmd/apisrv/main.go

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

FROM scratch

MAINTAINER Alex Pudich <pudichas@gmail.com>

WORKDIR /app

COPY --from=build /app/apisrv ./
COPY --from=build /wait /wait
COPY /.env ./
COPY /config/config.yml ./config/config.yml
COPY /database ./database/

CMD ["/app/apisrv"]

