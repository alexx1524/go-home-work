FROM golang:1.16.2 as build

WORKDIR "/dbmigrations"

COPY ./migrations .

RUN go get -u github.com/pressly/goose/cmd/goose

CMD ["/go/bin/goose", "postgres", "postgres://pguser:pgpwd@calendar_postgres:5432/calendar?sslmode=disable", "up"]