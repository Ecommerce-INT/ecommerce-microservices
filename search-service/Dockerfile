FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY pkg/common ./pkg/common
COPY search-service ./search-service
COPY go.work ./go.work

WORKDIR /app/search-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /search-service main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /search-service /search-service

EXPOSE 8094

ENTRYPOINT ["/search-service"]
