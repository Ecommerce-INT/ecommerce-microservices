FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy shared library, service and workspace
COPY pkg/common ./pkg/common
COPY tax-service ./tax-service
COPY go.work ./go.work

WORKDIR /app/tax-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /tax-service main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /tax-service /tax-service

EXPOSE 8091

ENTRYPOINT ["/tax-service"]
