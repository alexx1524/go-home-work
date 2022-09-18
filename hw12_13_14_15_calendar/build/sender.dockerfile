# Собираем в гошке
FROM golang:1.16.2 as build

ENV BIN_FILE /opt/calendar/sender
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
-ldflags "$LDFLAGS" \
-o ${BIN_FILE} cmd/sender/main.go

# На выходе тонкий образ
FROM alpine:3.9 as production

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="Calendar sender"
LABEL MAINTAINERS="alexx1524@gmail.com"

ENV BIN_FILE "/opt/calendar/sender"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/sender_config.toml
COPY ./configs/sender_config.toml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
