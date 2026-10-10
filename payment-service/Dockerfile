FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY pkg/common ./pkg/common
COPY payment-service ./payment-service
COPY go.work ./go.work

WORKDIR /app/payment-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /payment-service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /payment-service /payment-service

EXPOSE 8085

ENTRYPOINT ["/payment-service"]
