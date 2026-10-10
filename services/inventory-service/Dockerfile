FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY inventory-service/go.mod inventory-service/go.sum ./inventory-service/
RUN cd inventory-service && go mod download

COPY shared ./shared
COPY inventory-service ./inventory-service

WORKDIR /app/inventory-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /inventory-service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /inventory-service /inventory-service

EXPOSE 8082

ENTRYPOINT ["/inventory-service"]
