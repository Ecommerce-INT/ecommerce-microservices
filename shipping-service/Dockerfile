FROM golang:alpine AS builder

WORKDIR /app

COPY pkg/common ./pkg/common
COPY shipping-service ./shipping-service
COPY go.work ./go.work

WORKDIR /app/shipping-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /shipping-service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /shipping-service /shipping-service

EXPOSE 8087

ENTRYPOINT ["/shipping-service"]
