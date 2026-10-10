FROM golang:alpine AS builder

WORKDIR /app

COPY pkg/common ./pkg/common
COPY product-service ./product-service
COPY go.work ./go.work

WORKDIR /app/product-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /product-service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /product-service /product-service

EXPOSE 8086

ENTRYPOINT ["/product-service"]
