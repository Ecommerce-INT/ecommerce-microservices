# 05-api — Tài liệu API & sự kiện

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> Mục lục thư mục `05-api/`: hợp đồng HTTP của 11 service, định tuyến gateway, sự kiện Kafka và danh mục lỗi. Mọi tài liệu mô tả **code hiện tại** (P-1 trong [`conventions-and-ids.md`](../00-overview/conventions-and-ids.md)); điểm lệch được gắn `ISSUE-xx` ([`known-issues-tech-debt.md`](../07-plan/known-issues-tech-debt.md)).

## 1. Tài liệu dạng đọc

| Tài liệu | Nội dung |
| :--- | :--- |
| [`api-catalog.md`](api-catalog.md) | Bảng mọi endpoint của 11 service (path service ↔ path gateway, auth, request/response) + endpoint frontend thực sự gọi + sai lệch hợp đồng |
| [`gateway-routing.md`](gateway-routing.md) | Phân tích APISIX local: upstream, global rule, `openid-connect`, bảng 33 route, ma trận public/protected |
| [`event-catalog.md`](event-catalog.md) | Danh mục topic Kafka: producer, consumer, payload, ví dụ, ngữ nghĩa giao nhận, trạng thái |
| [`error-catalog.md`](error-catalog.md) | Mọi phản hồi lỗi của 11 service, gateway và BFF (thông điệp nguyên văn, `file:line`, xử lý ở frontend) + đề xuất mô hình lỗi thống nhất |

## 2. Đặc tả máy đọc

### 2.1 OpenAPI 3.1 — `openapi/`

| Service | Ngôn ngữ | Port | Spec | Route gateway chính |
| :--- | :--- | :--- | :--- | :--- |
| auth-service | Go | 8088 | [`openapi/auth-service.yaml`](openapi/auth-service.yaml) | `auth-public`, `auth-users`, `auth-roles` |
| product-service (+ favourites) | Go | 8086 | [`openapi/product-service.yaml`](openapi/product-service.yaml) | `product-catalog`, `favourite` (hỏng — ISSUE-13) |
| order-service | Go | 8084 | [`openapi/order-service.yaml`](openapi/order-service.yaml) | `order-carts`, `order-orders` |
| payment-service | Go | 8085 | [`openapi/payment-service.yaml`](openapi/payment-service.yaml) | `payment` |
| inventory-service | Go | 8082 | [`openapi/inventory-service.yaml`](openapi/inventory-service.yaml) | `inventory` |
| shipping-service | Go | 8087 | [`openapi/shipping-service.yaml`](openapi/shipping-service.yaml) | `shipping` |
| search-service | Go | 8094 | [`openapi/search-service.yaml`](openapi/search-service.yaml) | `search-storefront` |
| rating-service | Bun | 8089 | [`openapi/rating-service.yaml`](openapi/rating-service.yaml) | `rating-storefront`, `rating-backoffice` |
| media-service | Bun | 8083 | [`openapi/media-service.yaml`](openapi/media-service.yaml) | `media` |
| notification-service | Bun | 8090 | [`openapi/notification-service.yaml`](openapi/notification-service.yaml) | `notif-email`, `notif-notifications`, `notif-payment` |
| promotion-service (+ tax) | Bun | 8093 | [`openapi/promotion-service.yaml`](openapi/promotion-service.yaml) | `promotion-backoffice`, `tax-classes`/`tax-rates` (hỏng — ISSUE-13) |

**Quy ước chung của các spec**

| Quy ước | Ý nghĩa |
| :--- | :--- |
| `info.version: "1.0.0"`, mô tả tiếng Việt | Phiên bản tài liệu, không phải phiên bản API (service không version API) |
| `servers` | `http://localhost:9080` (gateway APISIX compose) + port trực tiếp của service (publish ra host — ISSUE-02). Path không có route gateway khai báo `servers` riêng ở cấp path |
| Path key | Path **tại gateway**; service luôn đăng ký thêm dạng có/không prefix context-root nên path vẫn hợp lệ khi gọi trực tiếp |
| `x-upstream-path` | Các path service thực nhận (sau `proxy-rewrite`, kèm alias) |
| `x-gateway-route` | `id` route trong `deploy/apisix/apisix.yaml`; `none — not exposed` khi không route nào khớp; ghi chú "hỏng" cho route trỏ upstream không tồn tại (ISSUE-13) hoặc không khớp path gốc (ISSUE-15) |
| `security` | `bearerAuth` khi route APISIX có `openid-connect`; `[]` khi public hoặc chỉ gọi trực tiếp |
| `x-required-role` | Role được **thực thi**; hiện luôn `none (ISSUE-01)` vì không tầng backend nào kiểm role |
| `x-known-issues` | Danh sách `ISSUE-xx` liên quan tới operation |
| Tag `health` | `/actuator/health` (+ alias) — tĩnh, không có route gateway (ISSUE-36) |
| Không có endpoint OpenAPI runtime | Route `/<svc>/v3/api-docs/*` trên gateway là route chết (ISSUE-16); spec ở đây là tài liệu tĩnh |

### 2.2 AsyncAPI 3.0 — `asyncapi/`

| Spec | Nội dung |
| :--- | :--- |
| [`asyncapi/kafka-events.yaml`](asyncapi/kafka-events.yaml) | Broker Kafka (compose + host), 5 channel (`SUCCESSFUL`, `paymentCompleted`, `status-payment-successful`, `profile-onboarding-topic`, `dbproduct.public.product`), operation send/receive theo service, consumer group, schema payload, `x-known-issues`, `x-status` (orphan-produce / orphan-consume) |

## 3. Cách xem / kiểm tra spec

Chạy từ thư mục gốc repo. Các lệnh cần Docker hoặc Node.js (`npx`); không có công cụ nào được cài sẵn trong repo.

### 3.1 Swagger UI (OpenAPI)

```bash
# Mở http://localhost:8081 — chọn spec trong ô "Explore"
docker run --rm -p 8081:8080 \
  -v "$PWD/docs-all/05-api/openapi:/usr/share/nginx/html/specs:ro" \
  -e URLS='[{"url":"./specs/auth-service.yaml","name":"auth"},{"url":"./specs/product-service.yaml","name":"product"},{"url":"./specs/order-service.yaml","name":"order"},{"url":"./specs/payment-service.yaml","name":"payment"},{"url":"./specs/inventory-service.yaml","name":"inventory"},{"url":"./specs/shipping-service.yaml","name":"shipping"},{"url":"./specs/search-service.yaml","name":"search"},{"url":"./specs/rating-service.yaml","name":"rating"},{"url":"./specs/media-service.yaml","name":"media"},{"url":"./specs/notification-service.yaml","name":"notification"},{"url":"./specs/promotion-service.yaml","name":"promotion"}]' \
  swaggerapi/swagger-ui
```

> "Try it out" với server `localhost:9080` cần access token Keycloak (nút **Authorize** → `bearerAuth`). Cách lấy token local: [`developer-guide.md`](../08-guides/developer-guide.md). Gọi từ trình duyệt tới gateway có thể bị CORS chặn (allowlist không có `localhost:8081` — ISSUE-14).

### 3.2 Redocly (OpenAPI — lint & tài liệu tĩnh)

```bash
# Lint một spec (rule mặc định "recommended"; cảnh báo về extension x-* là bình thường)
npx @redocly/cli lint docs-all/05-api/openapi/promotion-service.yaml

# Xem trước có hot-reload
npx @redocly/cli preview-docs docs-all/05-api/openapi/rating-service.yaml

# Xuất HTML tĩnh
npx @redocly/cli build-docs docs-all/05-api/openapi/media-service.yaml -o /tmp/media-service.html
```

### 3.3 AsyncAPI Studio / CLI

```bash
# Validate
npx @asyncapi/cli validate docs-all/05-api/asyncapi/kafka-events.yaml

# Xem trước trong trình duyệt (Studio local)
npx @asyncapi/cli start studio --file docs-all/05-api/asyncapi/kafka-events.yaml

# Hoặc dán nội dung file vào https://studio.asyncapi.com (không gửi dữ liệu nhạy cảm — file không chứa secret)
```

### 3.4 Kiểm tra cú pháp YAML nhanh (không cần cài thêm)

```bash
for f in docs-all/05-api/openapi/*.yaml docs-all/05-api/asyncapi/*.yaml; do
  ruby -ryaml -e 'YAML.load_file(ARGV[0]); puts "OK  #{ARGV[0]}"' "$f"
done
```

## 4. Bảo trì

| Khi | Cập nhật |
| :--- | :--- |
| Thêm/sửa handler | Spec của service + [`api-catalog.md`](api-catalog.md) + [`error-catalog.md`](error-catalog.md) nếu đổi mã lỗi |
| Sửa `deploy/apisix/apisix.yaml` | [`gateway-routing.md`](gateway-routing.md) + `x-gateway-route` / `security` trong spec |
| Thêm/sửa topic Kafka | [`asyncapi/kafka-events.yaml`](asyncapi/kafka-events.yaml) + [`event-catalog.md`](event-catalog.md) |
| Đổi biến môi trường | [`../08-guides/configuration-reference.md`](../08-guides/configuration-reference.md) |
