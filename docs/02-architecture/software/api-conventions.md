# Quy ước API (as-is)

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> Tài liệu này mô tả **các quy ước đang thực sự tồn tại** trong code — kể cả khi chúng không nhất quán. Danh mục endpoint đầy đủ: [`api-catalog.md`](../../05-api/api-catalog.md); quy tắc rewrite tại gateway: [`gateway-routing.md`](../../05-api/gateway-routing.md).

## 1. Tóm tắt nhanh

| Khía cạnh | Thực tế | Nhất quán? |
| :--- | :--- | :---: |
| Nhóm URL | `/api/...` (Go, notification), `/api/v1/...` (auth), `/backoffice/...` & `/storefront/...` (Bun, search), `/medias` (media) | ❌ |
| Prefix context | Mỗi service phục vụ cả `/<svc>/<path>` và `/<path>` | ✅ (mẫu chung) |
| Response thành công | `ApiResponse` envelope (auth) · JSON thô (còn lại) · text/plain (notification email) | ❌ ([ISSUE-29](../../07-plan/known-issues-tech-debt.md)) |
| Response lỗi | `ApiResponse{success:false}` (auth) · `ErrorVm` (Go) · `{statusCode,message}` / `{message}` / `{error}` (Bun) | ❌ |
| Phân trang | `page`/`size` (0-based) với `PageResponse` kiểu Spring · `pageNo`/`pageSize` với VM riêng | ❌ |
| JSON field | camelCase (trừ token Keycloak snake_case, vài key lệch như `fullname`) | Gần ✅ |
| ID | số nguyên (SERIAL/BIGSERIAL) trong path; user id Keycloak là UUID string | ❌ |
| Correlation | Header `X-Correlation-Id` (echo/sinh ở mọi service) | ✅ |
| Health | `GET /actuator/health` + `/<svc>/actuator/health` | ✅ |
| Versioning | Chỉ auth có `/v1` | ❌ |

## 2. URL

### 2.1 Các họ path

| Họ | Service | Ví dụ (path trong service) | Ghi chú |
| :--- | :--- | :--- | :--- |
| `/api/v1/{auth,users,roles}` | auth | `GET /api/v1/users/me` | `services/auth-service/internal/handler/auth_handler.go:45-56` |
| `/api/<resource-số-nhiều>` | product, order, payment, shipping, favourite | `GET /api/products/{productId}` | Tài nguyên REST cơ bản |
| `/api/inventory` | inventory | `GET /api/inventory?productName=a,b` | số ít |
| `/backoffice/<resource>` | promotion, tax, rating | `GET /backoffice/promotions` | quản trị |
| `/storefront/<resource>` | rating, search | `GET /storefront/catalog-search` | cửa hàng, có `search_suggest` (snake_case) |
| `/medias` | media | `GET /medias/{id}/file/{fileName}` | không có `/api` |
| `/api/{email,notifications,payment-notifications}` | notification | `POST /api/email/sendSimpleMail` | động từ camelCase trong path |

Phong cách động từ trong path (di sản Spring): `/api/orders/existOrderId`, `/api/payments/getOrder/{orderId}`, `/api/favourites/find`, `/api/favourites/delete`, `/api/shippings/find`, `/api/email/sendMail` — xem [`api-catalog.md`](../../05-api/api-catalog.md).

### 2.2 Prefix kép & gateway

Mỗi service đăng ký route cả **có** prefix context (`/product`, `/order`, `/payment`, `/inventory`, `/shipping`, `/favourite`, `/search`, `/auth`, `/rating`, `/promotion`, `/tax`, `/media`, `/notification`) và **không** prefix ([`go-service-architecture.md` §6.1](go-service-architecture.md#61-đăng-ký-route-kép-prefix--root), [`bun-service-architecture.md` §4.2](bun-service-architecture.md#42-mount-route-kép-prefix--root)). Gateway dùng hai kiểu:

| Kiểu route APISIX | Ví dụ | Path tới service |
| :--- | :--- | :--- |
| Path công khai không prefix + `proxy-rewrite` thêm prefix | `/api/products/*` → `/product/api/products/*` | `deploy/apisix/apisix.yaml:71-82` |
| Path có prefix, không rewrite | `/promotion/backoffice/promotions/*`, `/media/medias/*`, `/notification/api/...` | `apisix.yaml:280-351` |

Hệ quả: client (BFF) phải biết họ path nào cần prefix. Các path gốc không khớp `/*` bị 404 ([ISSUE-15](../../07-plan/known-issues-tech-debt.md)); route tới upstream không tồn tại ([ISSUE-13](../../07-plan/known-issues-tech-debt.md)).

### 2.3 Method

- `PUT /<resource>` (id trong body) **và** `PUT /<resource>/{id}` cùng tồn tại ở product/order/payment (`services/product-service/internal/handler/product_handler.go:145-146`).
- Không dùng `PATCH`.
- `GET`/`DELETE` có body ở favourites/shipping (`/find`, `/delete`) ([ISSUE-35](../../07-plan/known-issues-tech-debt.md)).
- `POST` tạo mới trả **201** ở product/category, promotion, tax; **200** ở order, payment, shipping, rating, media, payment-notifications.

## 3. Định dạng response thành công

### 3.1 `ApiResponse<T>` envelope — chỉ auth-service

Định nghĩa: `services/shared/response/api_response.go:12-21`; helper `OK`, `Message`, `Error` (`api_response.go:30-65`).

```json
{
  "success": true,
  "code": "OK",
  "message": "User updated successfully",
  "data": { "id": 7, "fullname": "Nguyen Van A", "username": "a", "email": "a@x.vn", "gender": "MALE", "phone": "", "avatar": "" },
  "traceId": "9f0c2a7e5b1d4c3e8a6f0b2d4c6e8a0f",
  "timestamp": "2026-11-12T03:04:05.123456789Z"
}
```

- `traceId` = correlation id của request; `timestamp` UTC RFC 3339 nano.
- `GET /api/v1/auth/session` và `POST /api/v1/auth/refresh` bọc **token Keycloak snake_case** (`access_token`, `refresh_token`, `expires_in`, `refresh_expires_in`, `token_type`, `scope`) trong `data` (`services/auth-service/internal/model/user.go:54-61`).
- BFF tự gỡ `data` (`frontend/packages/lib/src/server/api.ts:71-75`).

### 3.2 JSON thô — các service Go còn lại & Bun

```json
// GET /api/products/42 (product-service)
{
  "productId": 42,
  "productTitle": "Áo thun",
  "imageUrl": "",
  "sku": "TS-01",
  "priceUnit": 199000,
  "quantity": 10,
  "category": { "categoryId": 3, "categoryTitle": "Thời trang", "imageUrl": "" }
}
```

Biến thể:

| Biến thể | Ví dụ | Nguồn |
| :--- | :--- | :--- |
| Mảng thô | `GET /api/products`, `/api/categories`, `/api/orders`, `/api/payments`, `/medias?ids=` | `product_handler.go:150-157` |
| Bọc `{ "collection": [...] }` | favourites, shippings list | `services/product-service/internal/model/favourite.go:17-19`, `services/shipping-service/internal/model/shipping.go:18-20` |
| Giá trị nguyên thủy | `DELETE` → `true`; `average-star` → `4.5`; `existOrderId` → `true/false` | `product_handler.go:235`, `services/rating-service/src/index.ts:136-147` |
| Body rỗng | `DELETE` promotion (200), tax/media (204) | `services/promotion-service/src/index.ts:121,169` |
| `ResponeStatusVm` | `{ "title", "message", "statusCode": "200" }` khi xóa rating | `services/rating-service/src/service.ts:118-125` |
| text/plain | `"Mail Sent Successfully..."` | `services/notification-service/src/index.ts:42-72` |

Các DTO Go dùng `omitempty` không đồng đều: order/payment dùng con trỏ + `omitempty` (trường có thể biến mất), product/shipping luôn trả đủ trường.

## 4. Định dạng lỗi

### 4.1 `ErrorVm` — product, order, payment, inventory, shipping, search

`services/shared/response/response.go:6-34`:

```json
HTTP/1.1 404 Not Found
{ "statusCode": "404", "title": "Not found", "detail": "not found: product with id 42 not found" }
```

Chỉ có 3 helper: `BadRequest` (400), `NotFound` (404), `InternalError` (500). Không có 401/403/409/422. `statusCode` là **chuỗi**.

### 4.2 `ApiResponse` lỗi — auth-service

```json
HTTP/1.1 409 Conflict
{
  "success": false,
  "code": "CONFLICT",
  "message": "username already exists",
  "path": "/auth/api/v1/auth/signup",
  "traceId": "…",
  "timestamp": "2026-11-12T03:04:05Z"
}
```

Mã `code` đang dùng (`services/auth-service/internal/handler/auth_handler.go`):

| `code` | HTTP | Khi nào |
| :--- | :--- | :--- |
| `BAD_REQUEST` | 400 | body sai, thiếu ticket/refreshToken/username/role |
| `VALIDATION_FAILED` | 400 | signup thiếu username/password/email |
| `REGISTRATION_FAILED` | 400 | lỗi đăng ký khác |
| `INVALID_ID` | 400 | id không phải số |
| `UNAUTHORIZED` | 401 | ticket hết hạn, refresh lỗi |
| `AUTH_TOKEN_INVALID` | 401 | `/users/me` không có claims |
| `NOT_FOUND` | 404 | user không tồn tại |
| `CONFLICT` | 409 | trùng username/email/phone; gán/thu hồi role xung đột |
| `INTERNAL_ERROR` | 500 | lỗi khác |

### 4.3 Bun services

```json
HTTP/1.1 400 Bad Request
{ "statusCode": "400", "message": "Validation error", "errors": [ { "code": "too_small", "path": ["star"], "message": "…" } ] }
```

| Service / nhóm route | Format lỗi |
| :--- | :--- |
| rating, promotion (promotions), media | `{statusCode: "<code>", message[, errors]}` |
| tax routes | `{message}`; exception không bắt → `Internal Server Error` text/plain (ISSUE-94) |
| notification | `{error}` (500), body rỗng 404/500, `"Error: …"` text/plain, hoặc nuốt lỗi trả `[]`/`false` |

### 4.4 Lỗi tại gateway

Lỗi do APISIX sinh ra có định dạng riêng của APISIX/OpenResty: 401 từ `openid-connect` khi thiếu/sai token, 429 từ `limit-count`, 404/502/503 khi không khớp route hoặc upstream lỗi (`deploy/apisix/apisix.yaml:46-61`). Không có 403 theo role ([ISSUE-01](../../07-plan/known-issues-tech-debt.md)).

### 4.5 Bảng tổng hợp mã lỗi theo service

| Service | 400 | 401 | 404 | 409 | 500 | Ghi chú |
| :--- | :---: | :---: | :---: | :---: | :---: | :--- |
| auth | ✅ | ✅ | ✅ | ✅ | ✅ | `ApiResponse` |
| product | ✅ | — | ✅ (GET by id) | — | ✅ | DELETE id không tồn tại vẫn 200 |
| order | ✅ | — | ✅ | — | ✅ | |
| payment | ✅ | — | ✅ (+ "đã thanh toán" trả **404**) | — | ✅ | [ISSUE-25](../../07-plan/known-issues-tech-debt.md) |
| inventory | — | — | — | — | ✅ | |
| shipping | ✅ | — | ✅ | — | ✅ | |
| search | — | — | — | — | ✅ | lỗi OpenSearch thường trả trang rỗng ([ISSUE-34](../../07-plan/known-issues-tech-debt.md)) |
| promotion/tax | ✅ | — | ✅ | — | ✅ | |
| rating | ✅ | — | — | — | ✅ | average-star lỗi → `0.0` 200 |
| media | ✅ | — | ✅ | — | ✅ | `GET /medias?ids=` rỗng → 404 |
| notification | — | — | ✅ | — | ✅ | |

Mọi 500 đều đưa nguyên `err.Error()` / `err.message` ra client (ISSUE-93).

## 5. Phân trang & sắp xếp

### 5.1 Kiểu Spring `PageResponse` (`page`/`size`)

Endpoint: `GET /api/v1/users/all` (bọc trong `ApiResponse.data`), `GET /api/orders/all`, `GET /api/carts/all`, `GET /api/payments/all`.

| Query param | Default | Ghi chú |
| :--- | :--- | :--- |
| `page` | `0` (0-based; < 0 → 0) | `services/order-service/internal/handler/order_handler.go:207-208` |
| `size` | `10` (≤ 0 → 10, **không có max**) | `services/order-service/internal/service/order_service.go:103-108` |
| `sortBy` | `orderId` / `cartId` / `paymentId` / `id` | Whitelist qua `switch` → tên cột (an toàn SQL injection) — `services/order-service/internal/repository/order_repository.go:339-348` |
| `sortOrder` | `asc` (`DESC` không phân biệt hoa thường) | `order_repository.go:349` |

```json
{
  "content": [ { "orderId": 5, "orderDate": "2026-11-12T10:00:00", "orderFee": 530000, "productId": 42 } ],
  "totalElements": 31, "totalPages": 4, "size": 10, "number": 0,
  "numberOfElements": 10, "first": true, "last": false, "empty": false
}
```

Định nghĩa: `services/order-service/internal/model/order.go:36-46` (lặp ở payment, auth).

### 5.2 Kiểu `pageNo`/`pageSize` (Bun, search)

| Endpoint | Param (default) | Response |
| :--- | :--- | :--- |
| `GET /backoffice/promotions` | `pageNo=0`, `pageSize=5`, `promotionName`, `couponCode`, `startDate`, `endDate` (bị bỏ qua) | `{ promotionDetailVmList, totalElements, totalPages }` (`services/promotion-service/src/types.ts:48-52`) |
| `GET /backoffice/tax-classes/paging`, `/tax-rates/paging` | `pageNo=0`, `pageSize=10` | `{ taxClasses\|taxRates, pageNo, pageSize, totalElements, totalPages, isLast }` |
| `GET /backoffice/ratings`, `/storefront/ratings/products/{id}` | `pageNo=0`, `pageSize=5` | `{ ratingList, totalElements, totalPages }` |
| `GET /storefront/catalog-search` | `page=0`, `size=12`, `keyword`, `brand`, `category`, `attribute`, `minPrice`, `maxPrice`, `sortType` (bị bỏ qua) | `{ products, pageNo, pageSize, totalElements, totalPages, isLast, aggregations }` (`services/search-service/internal/model/search.go:34-42`) |

### 5.3 Kỳ vọng của BFF

BFF gửi `page`, `size`, `sort=productId,desc` (định dạng Spring `field,dir`), `categoryId`, `search` tới `GET /api/products`, `/api/categories`, `/api/orders` và đọc `PaginatedResponse` (`frontend/packages/lib/src/types.ts:12-20`, `frontend/packages/lib/src/schemas.ts:32-41`) — nhưng các endpoint này trả mảng thô và bỏ qua mọi query param ([ISSUE-48](../../07-plan/known-issues-tech-debt.md)). Không service nào hiểu tham số `sort=field,dir`.

## 6. Kiểu dữ liệu

### 6.1 ID

| Thực thể | Kiểu | Nguồn |
| :--- | :--- | :--- |
| product, category, order, cart, payment, inventory | `int` / `SERIAL` (32-bit) | `services/product-service/internal/repository/product_repo.go:36,45` |
| user (auth DB) | `int64` / `BIGINT IDENTITY` | `services/auth-service/internal/repository/user_repository.go:42` |
| rating, promotion, tax, media, notification | `bigserial` mode `number` | `services/rating-service/src/schema.ts:5` |
| User từ token (`sub`) | UUID string | `services/shared/middleware/middleware.go:55-56` |
| `userId` ở order/cart/payment/favourite | số nguyên — **không khớp** UUID Keycloak | [ISSUE-49](../../07-plan/known-issues-tech-debt.md) |
| `createdBy` (rating), `userId` (notification) | string | `services/notification-service/src/schema.ts:6` |
| Khóa phức hợp trong path | favourite `/{userId}/{productId}/{likeDate}`, shipping `/{orderId}/{productId}` | `favourite_handler.go:35`, `shipping_handler.go:37` |

Path param luôn parse bằng `strconv.Atoi`/`parseInt`; sai → 400.

### 6.2 Tiền

| Nơi | Kiểu wire | Kiểu lưu |
| :--- | :--- | :--- |
| `priceUnit` (product) | JSON number (`float64`) | `DECIMAL(12,2)` |
| `orderFee` (order) | JSON number (`float64`) | `DOUBLE PRECISION` |
| `amount` (payment-notifications) | chuỗi (Drizzle `numeric`) | `NUMERIC` |
| `discountAmount` (promotion) | số nguyên | `integer` |
| `rate` (tax) | number (`Number(row.rate)`) | `NUMERIC(10,4)` |

Không có trường currency; frontend mặc định định dạng **VND** (`frontend/packages/lib/src/format.ts:1-6`), ngưỡng free-ship 500 000 và phí 30 000 hard-code (`frontend/packages/lib/src/cart.ts:1-5`). Xem [ISSUE-27](../../07-plan/known-issues-tech-debt.md).

### 6.3 Ngày giờ

| Trường | Định dạng | Nguồn |
| :--- | :--- | :--- |
| `orderDate` | `YYYY-MM-DDTHH:MM:SS` **không timezone** (`to_char`) | `services/order-service/internal/repository/order_repository.go:248` |
| `likeDate` (favourite) | `2006-01-02T15:04:05` không timezone | `services/product-service/internal/repository/favourite_repo.go:60` |
| `createdAt/updatedAt` (shipping), `createdOn` (search) | RFC 3339 (`time.Time`) | `services/shipping-service/internal/model/shipping.go:14-15` |
| `timestamp` (ApiResponse) | RFC 3339 nano UTC | `api_response.go:41` |
| Bun (`createdOn`, `startDate`…) | ISO 8601 `toISOString()` / Date JSON | `services/rating-service/src/service.ts:16` |
| Query `createdFrom/createdTo` (rating) | ISO 8601 | `services/rating-service/src/index.ts:67-68` |

## 7. Header

| Header | Hướng | Ý nghĩa |
| :--- | :--- | :--- |
| `Authorization: Bearer <JWT>` | BFF → gateway → service | Gateway xác minh (route bảo vệ); service đọc claims không verify |
| `X-Correlation-Id` | vào/ra | Echo nếu có, sinh mới nếu không (Go: 32 hex; Bun: UUID). BFF hiện **không gửi** header này (ISSUE-98) |
| `Content-Type: application/json` | | Mọi body JSON; media dùng `multipart/form-data` (field `multipartFile` hoặc `file`) |
| `Location` | tax POST | `/tax-classes/{id}` — thiếu prefix `/backoffice` |
| `Content-Disposition: attachment` | media download | |

## 8. Xác thực theo route

Không có cơ chế khai báo ở service; mức bảo vệ do APISIX quyết định (`openid-connect bearer_only`). Public: `/api/v1/auth/*`, `/api/products*`, `/api/categories*`, `/storefront/*`, `/<svc>/v3/api-docs/*`. Còn lại yêu cầu token hợp lệ, **không** kiểm role. Lưu ý route `product-catalog` không có plugin xác thực nên **cả `POST/PUT/DELETE`** `/api/products*`, `/api/categories*` cũng công khai (ISSUE-102). Chi tiết: [`security-architecture.md`](../system/security-architecture.md), [`gateway-routing.md`](../../05-api/gateway-routing.md), [ADR-003](../adr/ADR-003-gateway-authentication.md).

## 9. Health endpoint

| Path | Service |
| :--- | :--- |
| `GET /actuator/health` | tất cả 11 service |
| `GET /<svc>/actuator/health` | tất cả (`/favourite/...` ở product, `/tax/...` ở promotion) |
| `GET /health` | auth, order, payment |

Response luôn `200 {"status":"UP","components":{"db":{"status":"UP"}}}` (search: `opensearch`; media thêm `storage`) — `services/shared/response/api_response.go:68-75`. Không kiểm phụ thuộc thật ([ISSUE-36](../../07-plan/known-issues-tech-debt.md)). Tên `actuator` là di sản Spring Boot. Gateway không route các path health ra ngoài.

## 10. OpenAPI

Không service nào phục vụ OpenAPI/Swagger; các route `/<svc>/v3/api-docs/*` ở gateway là route chết ([ISSUE-16](../../07-plan/known-issues-tech-debt.md)).

## 11. Danh sách bất nhất & issue

| Bất nhất | Issue |
| :--- | :--- |
| 3+ định dạng response/lỗi | [ISSUE-29](../../07-plan/known-issues-tech-debt.md), ISSUE-94 |
| List product/category/order không phân trang, bỏ qua query | [ISSUE-48](../../07-plan/known-issues-tech-debt.md) |
| `page/size` vs `pageNo/pageSize`; `sort=field,dir` vs `sortBy/sortOrder` vs `sortType` | [ISSUE-48](../../07-plan/known-issues-tech-debt.md), [ISSUE-34](../../07-plan/known-issues-tech-debt.md) |
| GET/DELETE có body | [ISSUE-35](../../07-plan/known-issues-tech-debt.md) |
| 404 cho "đã thanh toán" | [ISSUE-25](../../07-plan/known-issues-tech-debt.md) |
| Tiền float / decimal / string | [ISSUE-27](../../07-plan/known-issues-tech-debt.md) |
| user id integer vs UUID | [ISSUE-49](../../07-plan/known-issues-tech-debt.md) |
| `skuCode` vs `productName` | [ISSUE-51](../../07-plan/known-issues-tech-debt.md) |
| Profile `fullname` (BE) vs `fullName` (FE), thiếu `roles` | ISSUE-100 |
| Product payload `categoryId` phẳng (FE) vs `category.categoryId` (BE) | ISSUE-60 |
| Lộ thông điệp lỗi nội bộ trong 500 | ISSUE-93 |
| Endpoint ghi product/category công khai tại gateway | ISSUE-102 |
