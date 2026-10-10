FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY shipping-service/go.mod shipping-service/go.sum ./shipping-service/
RUN cd shipping-service && go mod download

COPY shared ./shared
COPY shipping-service ./shipping-service

WORKDIR /app/shipping-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /shipping-service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /shipping-service /shipping-service

EXPOSE 8087

ENTRYPOINT ["/shipping-service"]
