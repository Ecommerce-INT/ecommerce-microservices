FROM golang:alpine AS builder

WORKDIR /app

COPY pkg/common ./pkg/common
COPY order-service ./order-service
COPY go.work ./go.work

WORKDIR /app/order-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /order-service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /order-service /order-service

EXPOSE 8084

ENTRYPOINT ["/order-service"]
