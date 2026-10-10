FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY auth-service/go.mod auth-service/go.sum ./auth-service/
RUN cd auth-service && go mod download

COPY shared ./shared
COPY auth-service ./auth-service

WORKDIR /app/auth-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /auth-service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /auth-service /auth-service

EXPOSE 8088

ENTRYPOINT ["/auth-service"]
