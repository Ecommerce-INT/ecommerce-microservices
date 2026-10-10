FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY pkg/common ./pkg/common
COPY auth-service ./auth-service
COPY go.work ./go.work

WORKDIR /app/auth-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /auth-service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /auth-service /auth-service

EXPOSE 8088

ENTRYPOINT ["/auth-service"]
