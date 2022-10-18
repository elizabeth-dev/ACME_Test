FROM golang:1.19-alpine as builder

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app ./cmd/user

RUN wget -qO/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.13/grpc_health_probe-linux-amd64 && \
    chmod +x /grpc_health_probe

# Prod container
FROM gcr.io/distroless/static as prod

COPY --from=builder /app /bin/app
COPY --from=builder /grpc_health_probe /bin/grpc_health_probe

USER 10001:10001

CMD ["/bin/app"]
