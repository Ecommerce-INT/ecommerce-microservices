FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY search-service/go.mod search-service/go.sum ./search-service/
RUN cd search-service && go mod download

COPY shared ./shared
COPY search-service ./search-service

WORKDIR /app/search-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /search-service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /search-service /search-service

EXPOSE 8094

ENTRYPOINT ["/search-service"]
