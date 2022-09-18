FROM golang:1.16.2

WORKDIR /go/src

COPY . .

RUN go mod download
RUN go install github.com/onsi/ginkgo/v2/ginkgo

CMD go test ./testing/...