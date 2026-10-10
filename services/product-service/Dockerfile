FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY product-service/go.mod product-service/go.sum ./product-service/
RUN cd product-service && go mod download

COPY shared ./shared
COPY product-service ./product-service

WORKDIR /app/product-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /product-service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /product-service /product-service

EXPOSE 8086

ENTRYPOINT ["/product-service"]
