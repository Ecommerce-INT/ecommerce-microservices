# 🛒 E-Commerce Cloud-Native Microservices Platform

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: light)" srcset="https://socialify.git.ci/Ecommerce-INT/ecommerce-microservices/image?description=1&descriptionEditable=%E2%9A%A1%EF%B8%8F%2011%20Microservices%20%E2%80%A2%20Golang%201.27%20%E2%80%A2%20Bun%20%E2%80%A2%20Kubernetes&font=Inter&forks=1&language=1&owner=1&pattern=Floating%20Cogs&pulls=1&stargazers=1&theme=Light"/>
    <source media="(prefers-color-scheme: dark)" srcset="https://socialify.git.ci/Ecommerce-INT/ecommerce-microservices/image?description=1&descriptionEditable=%E2%9A%A1%EF%B8%8F%2011%20Microservices%20%E2%80%A2%20Golang%201.27%20%E2%80%A2%20Bun%20%E2%80%A2%20Kubernetes&font=Inter&forks=1&language=1&owner=1&pattern=Floating%20Cogs&pulls=1&stargazers=1&theme=Dark"/>
    <img alt="ecommerce-microservices" src="https://socialify.git.ci/Ecommerce-INT/ecommerce-microservices/image?description=1&descriptionEditable=%E2%9A%A1%EF%B8%8F%2011%20Microservices%20%E2%80%A2%20Golang%201.27%20%E2%80%A2%20Bun%20%E2%80%A2%20Kubernetes&font=Inter&forks=1&language=1&owner=1&pattern=Floating%20Cogs&pulls=1&stargazers=1&theme=Auto"/>
  </picture>
</p>

<p align="center">
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License">
  </a>
  <img src="https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white" alt="Go 1.27">
  <img src="https://img.shields.io/badge/Bun-1.x-f472b6?logo=bun&logoColor=white" alt="Bun 1.x">
  <img src="https://img.shields.io/badge/Architecture-11_Microservices-brightgreen" alt="11 Microservices">
</p>

---

## 📋 Overview

**ecommerce-microservices** is a production-grade, ultra-lightweight e-commerce ecosystem architected with **11 cloud-native microservices**:
- **7 Core Services** implemented in **Golang 1.27** (Fiber v2, pgx, Kafka, Elasticsearch).
- **4 Edge & Engagement Services** implemented in **Bun / TypeScript** (Hono, postgres.js, AWS S3 SDK, KafkaJS).
- Cắt giảm **>90% tài nguyên RAM** (tổng chỉ ~250 MB idle so với 6–8 GB của dàn Spring Boot JVM cũ).
- Cold start sub-millisecond (<50ms).

---

## 🛠️ Architecture & 11 Target Microservices

| # | Microservice | Runtime / Stack | Port | Responsibilities & Domain |
| :- | :--- | :--- | :--- | :--- |
| **1** | **`auth-service`** | Go 1.27 · Fiber v2 · pgx | `8088` | Keycloak OIDC callback, SSO ticket exchange, User & Role CRUD. |
| **2** | **`product-service`** | Go 1.27 · Fiber v2 · pgx | `8086` | Catalog, Categories, Brands, **User Favourites / Wishlists** (gộp từ `favourite-service`). |
| **3** | **`order-service`** | Go 1.27 · Fiber v2 · pgx | `8084` | Shopping Cart, Order creation, Checkout lifecycle, **OrderItem tracking**. |
| **4** | **`payment-service`** | Go 1.27 · Fiber v2 · pgx · Kafka | `8085` | Payment transactions, Kafka event publisher, payment status callbacks. |
| **5** | **`inventory-service`** | Go 1.27 · Fiber v2 · pgx | `8082` | Batch stock queries, stock reservations, warehouse ledger. |
| **6** | **`shipping-service`** | Go 1.27 · Fiber v2 · pgx | `8087` | Shipping fee quotes, carrier management, shipment lifecycle. |
| **7** | **`search-service`** | Go 1.27 · Fiber v2 · Elasticsearch 8 | `8094` | Full-text catalog search, auto-complete suggestions, Kafka consumer sync. |
| **8** | **`promotion-service`** | Bun 1.x · Hono · postgres.js | `8093` | Coupons, discount rules, **Tax Classes & Tax Rates calculation** (gộp từ `tax-service`). |
| **9** | **`rating-service`** | Bun 1.x · Hono · postgres.js | `8089` | Product reviews, verified purchases, star rating aggregates. |
| **10** | **`media-service`** | Bun 1.x · Hono · AWS S3 SDK | `8083` | Multipart file upload, RustFS (S3) stream, presigned upload URLs. |
| **11** | **`notification-service`** | Bun 1.x · Hono · KafkaJS | `8090` | Transactional email dispatches, Kafka notification consumer, in-app inbox. |

---

## 📁 Project Structure

```
ecommerce-microservices/
├── auth-service/                   # Port 8088 (Go 1.27)
├── product-service/                # Port 8086 (Go 1.27 - Includes Favourites)
├── order-service/                  # Port 8084 (Go 1.27)
├── payment-service/                # Port 8085 (Go 1.27)
├── inventory-service/              # Port 8082 (Go 1.27)
├── shipping-service/               # Port 8087 (Go 1.27)
├── search-service/                 # Port 8094 (Go 1.27)
├── promotion-service/              # Port 8093 (Bun / TS - Includes Taxes)
├── rating-service/                 # Port 8089 (Bun / TS)
├── media-service/                  # Port 8083 (Bun / TS)
├── notification-service/           # Port 8090 (Bun / TS)
├── pkg/common/                     # Shared Go module (database, config, middleware, response)
├── deploy/apisix/                  # Apache APISIX gateway configuration
├── docker/postgres/                # Database initialization scripts
├── k8s/                            # Kubernetes manifests & Gateway routes
├── go.work                         # Go workspace for 7 Go services
├── docker-compose.yml              # 11 Microservices + Infra orchestration
└── README.md
```

---

## 🚀 Running Locally

### 1. Run all with Docker Compose:
```bash
docker compose up -d --build
```

### 2. Run Go Services Locally:
```bash
# Verify all Go services build
go build ./auth-service/... ./product-service/... ./order-service/... ./payment-service/... ./inventory-service/... ./shipping-service/... ./search-service/...

# Run all Go tests
go test ./...
```

### 3. Run Bun Services Locally:
```bash
cd promotion-service && bun test && cd ..
cd rating-service && bun test && cd ..
cd media-service && bun test && cd ..
cd notification-service && bun test && cd ..
```

---

## 📊 Health Checks & Observability

Mỗi microservice đều cung cấp endpoint kiểm tra trạng thái tương thích chuẩn Spring Boot Actuator:
- Root health check: `GET /actuator/health` $\to$ `{"status": "UP"}`
- Service context health check: `GET /{service-name}/actuator/health` $\to$ `{"status": "UP"}`
- Correlation ID: Tự động gắn và chuyển tiếp qua header `X-Correlation-Id`.

---

## 📄 License
MIT License.
