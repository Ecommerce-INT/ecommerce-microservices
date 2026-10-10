FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY pkg/common ./pkg/common
COPY inventory-service ./inventory-service
COPY go.work ./go.work

WORKDIR /app/inventory-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /inventory-service main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /inventory-service /inventory-service

EXPOSE 8082

ENTRYPOINT ["/inventory-service"]
