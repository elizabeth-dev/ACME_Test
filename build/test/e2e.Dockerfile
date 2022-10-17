FROM golang:1.19-alpine as builder

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download -x

COPY . .

ENV CGO_ENABLED=0
CMD ["go", "test", "./test/e2e/..."]
