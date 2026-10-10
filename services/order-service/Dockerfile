FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY order-service/go.mod order-service/go.sum ./order-service/
RUN cd order-service && go mod download

COPY shared ./shared
COPY order-service ./order-service

WORKDIR /app/order-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /order-service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /order-service /order-service

EXPOSE 8084

ENTRYPOINT ["/order-service"]
