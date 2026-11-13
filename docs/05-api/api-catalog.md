# API Catalog — toàn bộ endpoint của 11 service

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> Mô tả **as-is** từ handler trong `services/`. Cấu hình gateway chi tiết: [`gateway-routing.md`](gateway-routing.md). Tài liệu từng service: [`../03-services/`](../03-services/).

## 1. Quy ước chung

### 1.1 Cột "Gateway"

- **Path qua gateway** là path client gọi tới APISIX (`http://apisix:9080` nội bộ / `http://localhost:9080` từ host).
- **Auth**: `Public` = route không có `openid-connect`; `Bearer` = APISIX verify JWT (RS256, JWKS Keycloak) nhưng **không kiểm role** (ISSUE-01).
- ❌ = không có route APISIX local nào khớp (thiếu route, ISSUE-15 path gốc, hoặc upstream không tồn tại — ISSUE-13). Endpoint vẫn gọi được trực tiếp qua port service publish ra host (ISSUE-02).

### 1.2 Đăng ký route kép

Mọi service (Go và Bun) đăng ký route **hai lần**: dưới context-root của service (`/auth`, `/product`, `/favourite`, `/order`, `/payment`, `/inventory`, `/shipping`, `/search`, `/rating`, `/media`, `/notification`, `/promotion`, `/tax`) và ở gốc. Bảng dưới ghi path ở gốc; path có prefix tương đương. Gateway kiểu `/api/...` rewrite sang path có prefix.

### 1.3 Định dạng response

| Kiểu | Service | Cấu trúc |
| :--- | :--- | :--- |
| `ApiResponse<T>` envelope | auth-service | `{success, code, message?, data?, errors?[], path?, traceId?, timestamp}` — `services/shared/response/api_response.go:12-60` |
| JSON thô + `ErrorVm` | product, order, payment, inventory, shipping, search | Thành công: JSON thô; lỗi: `{statusCode:"400"\|"404"\|"500", title, detail}` — `services/shared/response/response.go:6-35` |
| `PageResponse<T>` (Go) | order `/all`, payment `/all`, auth `/users/all` | `{content[], totalElements, totalPages, size, number, numberOfElements, first, last, empty}` |
| Bun tùy service | rating, media, promotion, notification | Lỗi `{statusCode:"400", message}` hoặc `{message}` hoặc `{error}` hoặc text — xem từng mục |

Không đồng nhất (ISSUE-29). Health của mọi service là tĩnh `{"status":"UP",...}` (ISSUE-36), không có route gateway.

### 1.4 Middleware chung

- **Go** (chi): `Recoverer` → `Logger` → CORS `*` → `CorrelationID` → `UserClaims` (parse JWT **không verify**, đặt `sub`/`email`/`preferred_username` vào context, không bao giờ chặn — `services/shared/middleware/middleware.go:25-70`, ISSUE-02). Chỉ `GET /users/me` dùng danh tính.
- **Bun** (Hono): `logger` → CORS `*` → `X-Correlation-Id`. Chỉ rating-service đọc token (decode không verify, ISSUE-03).

---

## 2. auth-service (Go, `:8088`, upstream 1)

Handler `services/auth-service/internal/handler/auth_handler.go`; base `/api/v1/...` và `/auth/api/v1/...` (`:45-57`). Envelope `ApiResponse`. `UserResponse = {id, fullname, username, email, gender, phone, avatar}` — **không có `roles`**, key `fullname` chữ thường (`services/auth-service/internal/model/user.go:15-22`, ISSUE-100).

| Method | Path service | Path qua gateway | Auth | Mục đích | Request | Response / lỗi |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| POST | `/api/v1/auth/signup` | `/api/v1/auth/signup` | Public | Đăng ký user (Keycloak + DB) | `RegisterRequest {fullName, username*, password*, email*, gender, phone, avatar, roles[]}` | 200 message; 400 `BAD_REQUEST`/`VALIDATION_FAILED`/`REGISTRATION_FAILED`; 409 `CONFLICT` |
| GET | `/api/v1/auth/login` | `/api/v1/auth/login` | Public | Bắt đầu SSO | `?redirect_uri` (không kiểm allowlist — ISSUE-05) | 302 → Keycloak authorize |
| GET | `/api/v1/auth/callback` | `/api/v1/auth/callback` | Public | Callback OIDC | `?code&state&error` | 302 → frontend `?ticket=` / `?error=`; đổi code lỗi → lưu **mock token** (ISSUE-04) |
| GET | `/api/v1/auth/session` | `/api/v1/auth/session` | Public | Đổi ticket 1 lần lấy token | `?ticket` | 200 `data: {access_token, refresh_token, token_type, expires_in, refresh_expires_in, scope}`; 400; 401 |
| POST | `/api/v1/auth/refresh` | `/api/v1/auth/refresh` | Public | Làm mới token | `{refreshToken}` | 200 token; 400; 401 |
| POST | `/api/v1/auth/logout` | `/api/v1/auth/logout` | Public | Logout Keycloak | `{refreshToken}` | 200 luôn |
| GET | `/api/v1/users/me` | `/api/v1/users/me` | Bearer | Profile hiện tại (tra DB theo `preferred_username`, fallback `sub`) | — | 200 `UserResponse`; 401 `AUTH_TOKEN_INVALID`; 404 (user chỉ có ở Keycloak — ISSUE-101) |
| PUT | `/api/v1/users/me/password` | `/api/v1/users/me/password` | Bearer | **Stub** | — | 200 luôn (ISSUE-07) |
| GET | `/api/v1/users/all` | `/api/v1/users/all` | Bearer | Danh sách phân trang | `?page=0&size=10&sortBy=id&sortOrder=ASC` | 200 `PageResponse<UserResponse>` |
| GET | `/api/v1/users/{id}` | `/api/v1/users/{id}` | Bearer | Lấy user | — | 200; 400 `INVALID_ID`; 404 |
| PUT | `/api/v1/users/{id}` | `/api/v1/users/{id}` | Bearer | Cập nhật một phần | `{fullName?, email?, gender?, phone?, avatar?}` | 200; 400; 404 |
| DELETE | `/api/v1/users/{id}` | `/api/v1/users/{id}` | Bearer | Xóa user | — | 200; 400; 404 |
| GET | `/api/v1/users?username=` | `/api/v1/users?username=` | Bearer | Tìm theo username | `?username` | 200; 400; 404 |
| POST | `/api/v1/roles/users/{userId}/assign` | như service | Bearer | Gán role (chỉ DB, ISSUE-08) | Body **text thô** = tên role | 200; 400; 409 |
| POST | `/api/v1/roles/users/{userId}/revoke` | như service | Bearer | Thu hồi role | Body text | 200; 400; 409 |
| GET | `/api/v1/roles/users/{userId}` | như service | Bearer | Role của user | — | 200 `string[]` |
| GET | `/actuator/health`, `/auth/actuator/health`, `/health` | ❌ | — | Health | — | 200 tĩnh |

## 3. product-service (Go, `:8086`, upstream 2) — gồm favourites

Handler `services/product-service/internal/handler/product_handler.go` (base `/api/categories`, `/api/products` + prefix `/product`) và `favourite_handler.go` (base `/api/favourites` + prefix `/favourite`, `/product`). `ProductDto = {productId, productTitle, imageUrl, sku, priceUnit, quantity, category?}`; `CategoryDto = {categoryId, categoryTitle, imageUrl, parentCategory?}` (`services/product-service/internal/model/product.go:3-18`).

| Method | Path service | Path qua gateway | Auth | Mục đích | Request | Response / lỗi |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| GET | `/api/categories` | `/api/categories` | **Public** | Tất cả danh mục | query bị bỏ qua | 200 `CategoryDto[]` (mảng thô — ISSUE-48) |
| GET | `/api/categories/paging`, `/paging-and-sorting` | như service | Public | Như trên, **không phân trang** (ISSUE-35) | | 200 mảng |
| GET | `/api/categories/{categoryId}` | như service | Public | Chi tiết | | 200; 400; 404 |
| POST | `/api/categories` | như service | **Public** (ISSUE-102) | Tạo | `CategoryDto` | **201**; 400 |
| PUT | `/api/categories` | như service | **Public** | Cập nhật (id trong body) | `CategoryDto` | 200; 400 |
| PUT | `/api/categories/{categoryId}` | như service | **Public** | Cập nhật theo path | `CategoryDto` | 200; 400 |
| DELETE | `/api/categories/{categoryId}` | như service | **Public** | Xóa | | 200 `true` |
| GET | `/api/products` | `/api/products` | Public | Tất cả sản phẩm, **bỏ qua** `page/size/sort/categoryId/search` | | 200 `ProductDto[]` thô (ISSUE-48) |
| GET | `/api/products/{productId}` | như service | Public | Chi tiết | | 200; 400; 404 |
| POST | `/api/products` | như service | **Public** (ISSUE-102) | Tạo | `ProductDto` (category lồng nhau) | **201**; 400 |
| PUT | `/api/products` | như service | **Public** | Cập nhật (id trong body) | `ProductDto` | 200; 400 |
| PUT | `/api/products/{productId}` | như service | **Public** | Cập nhật theo path | `ProductDto` | 200; 400 |
| DELETE | `/api/products/{productId}` | như service | **Public** | Xóa | | 200 `true` |
| GET | `/api/favourites` | `/api/favourites` ❌ (ISSUE-13) | Bearer | Tất cả | | 200 `{collection: FavouriteDto[]}` |
| GET | `/api/favourites/find` | ❌ | Bearer | Tìm theo id ghép | **Body JSON trên GET** `{userId, productId, likeDate}` (ISSUE-35) | 200; 400; 404 |
| GET | `/api/favourites/{userId}/{productId}/{likeDate}` | ❌ | Bearer | Tìm theo path | | 200; 400; 404 |
| POST | `/api/favourites` | ❌ | Bearer | Tạo | `FavouriteDto` | 200 |
| PUT | `/api/favourites` | ❌ | Bearer | = POST | | 200 |
| DELETE | `/api/favourites/delete` | ❌ | Bearer | Xóa theo body | `{userId, productId, likeDate}` | 200 bool |
| DELETE | `/api/favourites/{userId}/{productId}/{likeDate}` | ❌ | Bearer | Xóa theo path | | 200 bool |
| GET | `/actuator/health`, `/product/actuator/health`, `/favourite/actuator/health` | ❌ | — | Health | | 200 tĩnh |

## 4. order-service (Go, `:8084`, upstream 3)

Handler `services/order-service/internal/handler/order_handler.go`; base `/api/carts`, `/api/orders` (+ prefix `/order`). `OrderDto = {orderId?, orderDate?, orderDesc?, orderFee?, productId?, product?, cart?}`; `CartDto = {cartId?, userId?, order?: OrderDto[], user?}` (`services/order-service/internal/model/order.go:19-35`). Không handler nào đọc danh tính (ISSUE-49).

| Method | Path service | Path qua gateway | Auth | Mục đích | Request | Response / lỗi |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| GET | `/api/carts/all` | như service | Bearer | Phân trang | `?page=0&size=10&sortBy=cartId&sortOrder=asc` | 200 `PageResponse<CartDto>` |
| GET | `/api/carts/{cartId}` | như service | Bearer | Chi tiết | | 200; 400; 404 |
| GET | `/api/carts` | như service | Bearer | Tất cả | | 200 mảng |
| POST | `/api/carts` | như service | Bearer | Tạo (insert order lồng — ISSUE-24) | `CartDto` | 200 |
| PUT | `/api/carts/{cartId}` | như service | Bearer | Cập nhật | `CartDto` | 200; 404 |
| PUT | `/api/carts` | như service | Bearer | Cập nhật (id trong body) | `CartDto` có `cartId` | 200; 400; 404 |
| DELETE | `/api/carts/{cartId}` | như service | Bearer | Xóa | | 200 `true` |
| GET | `/api/orders/existOrderId?orderId=` | như service | Bearer | Kiểm tồn tại | `?orderId` | 200 bool; 400 |
| GET | `/api/orders/all` | như service | Bearer | Phân trang (mọi user) | `?page&size&sortBy=orderId&sortOrder` | 200 `PageResponse<OrderDto>` |
| GET | `/api/orders/{orderId}` | như service | Bearer | Chi tiết | | 200; 400; 404 |
| GET | `/api/orders` | `/api/orders` | Bearer | **Tất cả đơn của mọi user**, bỏ qua query | | 200 `OrderDto[]` thô (ISSUE-48, ISSUE-49) |
| POST | `/api/orders` | `/api/orders` | Bearer | Tạo đơn (`orderFee` tin client — ISSUE-21) | `OrderDto` | 200 |
| PUT | `/api/orders/{orderId}`, `/api/orders` | như service | Bearer | Cập nhật | `OrderDto` | 200; 400; 404 |
| DELETE | `/api/orders/{orderId}` | như service | Bearer | Xóa | | 200 `true` |
| GET | `/actuator/health`, `/order/actuator/health`, `/health` | ❌ | — | Health | | 200 |

## 5. payment-service (Go, `:8085`, upstream 4)

Handler `services/payment-service/internal/handler/payment_handler.go`; base `/api/payments` (+ `/payment`). `PaymentDto = {paymentId?, isPayed?, paymentStatus?, orderId?, userId?, order?, user?}`.

| Method | Path service | Path qua gateway | Auth | Mục đích | Request | Response / lỗi |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| GET | `/api/payments/getOrder/{orderId}` | như service | Bearer | **Stub** (ISSUE-25) | | 200 (có thể `null`); 400 |
| GET | `/api/payments/all` | như service | Bearer | Phân trang | `?page&size&sortBy=paymentId&sortOrder` | 200 `PageResponse` |
| GET | `/api/payments/{paymentId}` | như service | Bearer | Chi tiết | | 200; 400; 404 |
| GET | `/api/payments` | như service | Bearer | Tất cả | | 200 mảng |
| POST | `/api/payments` | `/api/payments` | Bearer | Tạo; publish Kafka mỗi lần (ISSUE-18) | `PaymentDto` | 200; 400; **404** "Order has already been paid." (ISSUE-25) |
| PUT | `/api/payments/{paymentId}`, `/api/payments` | như service | Bearer | Cập nhật | `PaymentDto` | 200; 400; 404 |
| DELETE | `/api/payments/{paymentId}` | như service | Bearer | Xóa | | 200 `true` |
| GET | `/actuator/health`, `/payment/actuator/health`, `/health` | ❌ | — | Health | | 200 |

## 6. inventory-service (Go, `:8082`, upstream 5)

Handler `services/inventory-service/internal/handler/inventory_handler.go:25-75`; base `/api/inventory` (+ `/inventory`).

| Method | Path service | Path qua gateway | Auth | Mục đích | Request | Response / lỗi |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| GET | `/api/inventory` | `/api/inventory` | Bearer | Kiểm tồn kho theo tên | `?productName=` (lặp hoặc CSV) | 200 `[{productName, isInStock, quantity?}]`; 500 |
| GET | `/api/inventory/{products}` | như service | Bearer | Như trên, tên trong path (CSV) | | 200 mảng |
| GET | `/actuator/health`, `/inventory/actuator/health` | ❌ | — | Health | | 200 |

## 7. shipping-service (Go, `:8087`, upstream 6)

Handler `services/shipping-service/internal/handler/shipping_handler.go`; base `/api/shippings` (+ `/shipping`). Gateway chỉ có `/shipping/api/shippings/*` (không rewrite). `OrderItemDto = {orderId, productId, orderedQuantity, createdAt?, updatedAt?}`.

| Method | Path service | Path qua gateway | Auth | Mục đích | Request | Response / lỗi |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| GET | `/api/shippings` | `/shipping/api/shippings` ❌ (ISSUE-15) | Bearer | Tất cả | | 200 `{collection: OrderItemDto[]}` |
| GET | `/api/shippings/find` | `/shipping/api/shippings/find` | Bearer | Tìm theo **body JSON trên GET** | `{orderId, productId}` | 200; 400; 404 |
| GET | `/api/shippings/{orderId}/{productId}` | `/shipping/api/shippings/{o}/{p}` | Bearer | Tìm theo path | | 200; 400; 404 |
| POST | `/api/shippings` | ❌ (ISSUE-15) | Bearer | Tạo | `OrderItemDto` | 200 |
| PUT | `/api/shippings` | ❌ (ISSUE-15) | Bearer | Cập nhật | `OrderItemDto` | 200 |
| DELETE | `/api/shippings/delete` | `/shipping/api/shippings/delete` | Bearer | Xóa theo body | `{orderId, productId}` | 200 `true` |
| DELETE | `/api/shippings/{orderId}/{productId}` | `/shipping/api/shippings/{o}/{p}` | Bearer | Xóa theo path | | 200 `true` |
| GET | `/actuator/health`, `/shipping/actuator/health` | ❌ | — | Health | | 200 |

## 8. search-service (Go, `:8094`, upstream 11)

Handler `services/search-service/internal/handler/search_handler.go:25-90`; prefix `/search` và gốc.

| Method | Path service | Path qua gateway | Auth | Mục đích | Request | Response / lỗi |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| GET | `/storefront/catalog-search` | `/storefront/catalog-search` | Public | Tìm kiếm full-text | `?keyword&page=0&size=12&brand&category&attribute&minPrice&maxPrice&sortType=DEFAULT` (sortType bị bỏ qua — ISSUE-34) | 200 `ProductListGetVm {products[{id,name,slug,thumbnailId,price,isAllowedToOrder,isPublished,isFeatured,isVisibleIndividually,createdOn}], pageNo, pageSize, totalElements, totalPages, isLast, aggregations}`; lỗi OpenSearch → trang rỗng |
| GET | `/storefront/search_suggest` | `/storefront/search_suggest` | Public | Gợi ý | `?keyword` | 200 `{productNames:[{id,name}]}` |
| GET | `/actuator/health`, `/search/actuator/health` | ❌ | — | Health | | 200 |

Index luôn rỗng vì không có nguồn đồng bộ (ISSUE-19).

## 9. rating-service (Bun, `:8089`, upstream 10)

Chi tiết: [`rating-service.md`](../03-services/rating-service.md#6-api-endpoints).

| Method | Path service | Path qua gateway | Auth | Mục đích | Request | Response / lỗi |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| GET | `/storefront/ratings/products/{productId}` | như service | Public | Đánh giá theo SP | `?pageNo=0&pageSize=5` | 200 `{ratingList[], totalElements, totalPages}`; 400 |
| POST | `/storefront/ratings` | `/storefront/ratings` | **Public** (ISSUE-03) | Tạo | `{content, star 1–5, productId, productName}` | 200 `RatingVm`; 400 |
| GET | `/storefront/ratings/product/{productId}/average-star` | như service | Public | Điểm TB | | 200 number |
| GET | `/backoffice/ratings` | `/rating/backoffice/ratings` ❌ (ISSUE-15) | Bearer | Lọc | `?productName&customerName&message&createdFrom&createdTo&pageNo&pageSize` | 200 |
| DELETE | `/backoffice/ratings/{id}` | `/rating/backoffice/ratings/{id}` | Bearer | Xóa | | 200 `{title,message,statusCode}`; 400 |
| GET | `/actuator/health`, `/rating/actuator/health` | ❌ | — | Health | | 200 |

## 10. media-service (Bun, `:8083`, upstream 9)

Chi tiết: [`media-service.md`](../03-services/media-service.md#6-api-endpoints).

| Method | Path service | Path qua gateway | Auth | Mục đích | Request | Response / lỗi |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| POST | `/medias` | `/media/medias` ❌ (ISSUE-15) | Bearer | Upload | multipart `multipartFile`/`file`, `caption`, `fileNameOverride` | 200 `{id,caption,fileName,mediaType}`; 400 |
| GET | `/medias?ids=` | ❌ (ISSUE-15) | Bearer | Metadata nhiều id | `?ids=1,2` | 200 `MediaVm[]`; 400; 404 |
| GET | `/medias/{id}` | `/media/medias/{id}` | Bearer | Metadata | | 200 `{id,caption,fileName,mediaType,url}`; 400; 404 |
| GET | `/medias/{id}/file/{fileName}` | `/media/medias/{id}/file/{fileName}` | Bearer (ISSUE-63) | Stream file | | 200 stream; 404 |
| DELETE | `/medias/{id}` | `/media/medias/{id}` | Bearer | Xóa | | 204 |
| GET | `/actuator/health`, `/media/actuator/health` | ❌ | — | Health | | 200 |

## 11. notification-service (Bun, `:8090`, upstream 8)

Chi tiết: [`notification-service.md`](../03-services/notification-service.md#6-api-endpoints).

| Method | Path service | Path qua gateway | Auth | Mục đích | Request | Response / lỗi |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| POST | `/api/email/sendSimpleMail` | `/notification/api/email/sendSimpleMail` | Bearer | Gửi email (ISSUE-64) | `{recipient, subject, msgBody}` | 200 text (kể cả khi SMTP lỗi — ISSUE-65) |
| POST | `/api/email/sendMailWithAttachment` | `/notification/api/email/sendMailWithAttachment` | Bearer | Như trên, bỏ qua attachment | như trên | 200 text |
| POST | `/api/email/sendMail` | `/notification/api/email/sendMail` | Bearer | Gửi email từ query | `?to&subject&body` | 200 text |
| GET | `/api/notifications` | ❌ (ISSUE-15) | Bearer | Tất cả thông báo (ISSUE-33) | | 200 mảng |
| GET | `/api/notifications/{id}` | `/notification/api/notifications/{id}` | Bearer | Chi tiết | | 200; 404 |
| POST | `/api/notifications` | ❌ (ISSUE-15) | Bearer | Tạo | `{userId,title,message,type,status}` | 201 |
| DELETE | `/api/notifications/{id}` | `/notification/api/notifications/{id}` | Bearer | Xóa | | 200 `true` |
| GET | `/api/payment-notifications` | ❌ (ISSUE-15) | Bearer | Tất cả | | 200 mảng |
| GET | `/api/payment-notifications/{paymentId}` | `/notification/api/payment-notifications/{id}` | Bearer | Theo `payment_id` hoặc `id` | | 200; 404 |
| POST | `/api/payment-notifications` | ❌ (ISSUE-15) | Bearer | Tạo | `{paymentId,userId,orderId,amount,isPayed,paymentStatus}` | 200 |
| DELETE | `/api/payment-notifications/{paymentId}` | `/notification/api/payment-notifications/{id}` | Bearer | Xóa | | 204 |
| GET | `/actuator/health`, `/notification/actuator/health` | ❌ | — | Health | | 200 |

## 12. promotion-service (Bun, `:8093`, upstream 12) — promotion + tax

Chi tiết: [`promotion-service.md`](../03-services/promotion-service.md#6-api-endpoints). Route tax trên gateway local trỏ upstream 13 không tồn tại (ISSUE-13) → mọi path `/tax/...` ❌ trên compose.

| Method | Path service | Path qua gateway | Auth | Mục đích | Request | Response / lỗi |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| GET | `/backoffice/promotions` | `/promotion/backoffice/promotions` ❌ (ISSUE-15) | Bearer | List + lọc | `?pageNo=0&pageSize=5&promotionName&couponCode&startDate&endDate` | 200 `{promotionDetailVmList[], totalElements, totalPages}` |
| GET | `/backoffice/promotions/{id}` | `/promotion/backoffice/promotions/{id}` | Bearer | Chi tiết | | 200; 400; 404 |
| POST | `/backoffice/promotions` | ❌ (ISSUE-15) | Bearer | Tạo | `PromotionPostVm` | 201; 400 |
| PUT | `/backoffice/promotions` | ❌ (ISSUE-15) | Bearer | Cập nhật | `PromotionPutVm` | 200; 400; 500 nếu id không tồn tại |
| DELETE | `/backoffice/promotions/{id}` | `/promotion/backoffice/promotions/{id}` | Bearer | Xóa | | 200 |
| GET | `/backoffice/tax-classes/paging` | `/tax/backoffice/tax-classes/paging` ❌ (ISSUE-13) | Bearer | Phân trang | `?pageNo=0&pageSize=10` | 200 `TaxClassListVm` |
| GET | `/backoffice/tax-classes/{id}` | ❌ (ISSUE-13) | Bearer | Chi tiết | | 200; 400; 404 |
| GET | `/backoffice/tax-classes` | ❌ | Bearer | Tất cả | | 200 mảng |
| POST | `/backoffice/tax-classes` | ❌ | Bearer | Tạo | `{name}` | 201 + `Location` |
| PUT | `/backoffice/tax-classes/{id}` | ❌ | Bearer | Đổi tên | `{name}` | 204; 404 |
| DELETE | `/backoffice/tax-classes/{id}` | ❌ | Bearer | Xóa | | 204; 404 |
| GET | `/backoffice/tax-rates/paging` | ❌ | Bearer | Phân trang | | 200 `TaxRateListVm` |
| GET | `/backoffice/tax-rates/tax-percent` | ❌ | Bearer | Tra thuế suất | `?taxClassId*&countryId*&stateOrProvinceId&zipCode` | 200 number; 400 |
| GET | `/backoffice/tax-rates/location-based-batch` | ❌ | Bearer | Tra nhiều class | `?taxClassIds*&countryId*` | 200 mảng; 400 |
| GET | `/backoffice/tax-rates/{id}` | ❌ | Bearer | Chi tiết | | 200; 404 |
| GET | `/backoffice/tax-rates` | ❌ | Bearer | Tất cả | | 200 mảng |
| POST | `/backoffice/tax-rates` | ❌ | Bearer | Tạo | `{rate, countryId, stateOrProvinceId?, zipCode?, taxClassId}` | 201 |
| PUT | `/backoffice/tax-rates/{id}` | ❌ | Bearer | Cập nhật | như POST | 204; 404 |
| DELETE | `/backoffice/tax-rates/{id}` | ❌ | Bearer | Xóa | | 204; 404 |
| GET | `/actuator/health`, `/promotion/actuator/health`, `/tax/actuator/health` | ❌ | — | Health | | 200 |

---

## 13. Endpoint frontend thực sự gọi

Toàn bộ lời gọi đi qua client `createApi()` tại `frontend/packages/lib/src/server/api.ts:116-161` (server-side, base `API_INTERNAL_URL`), cộng một redirect trình duyệt tới `API_PUBLIC_URL`. `request()` tự **bóc envelope** nếu body có field `data` (`api.ts:71-75`).

| Method client | HTTP | Path gateway | Service | Được dùng ở | Ghi chú hợp đồng |
| :--- | :--- | :--- | :--- | :--- | :--- |
| (redirect trình duyệt) | GET | `/api/v1/auth/login?redirect_uri=` | auth | `frontend/apps/web/src/routes/login/+page.server.ts:24`, `frontend/apps/admin/src/routes/login/+page.server.ts:21` | ISSUE-04, ISSUE-05 |
| `auth.exchangeTicket` | GET | `/api/v1/auth/session?ticket=` | auth | callback web/admin | ✅ khớp (bóc `data`) |
| `auth.profile` | GET | `/api/v1/users/me` | auth | callback web/admin | ❌ `fullname`/thiếu `roles` (ISSUE-100); 404 khi user chỉ ở Keycloak (ISSUE-101) |
| `auth.refresh` | POST | `/api/v1/auth/refresh` | auth | `hooks.server.ts` web/admin | ✅ |
| `auth.logout` | POST | `/api/v1/auth/logout` | auth | `auth.remote.ts` web/admin | ✅ |
| `products.list` | GET | `/api/products?page&size&sort&categoryId&search` | product | web catalog, admin products, dashboard | ❌ trả mảng thô, bỏ qua query (ISSUE-48) |
| `products.get` | GET | `/api/products/{id}` | product | web product detail, checkout | ✅ (trừ `description` — ISSUE-50) |
| `products.create` | POST | `/api/products` | product | admin | ⚠️ gửi `categoryId` phẳng — service đọc `category.categoryId` (ISSUE-60); public (ISSUE-102) |
| `products.update` | PUT | `/api/products/{id}` | product | admin | ⚠️ như trên |
| `products.remove` | DELETE | `/api/products/{id}` | product | admin | ✅ |
| `categories.list` | GET | `/api/categories?page&size` | product | web catalog, admin categories/products, dashboard | ❌ mảng thô (ISSUE-48) |
| `categories.get` | GET | `/api/categories/{id}` | product | (không dùng) | |
| `categories.create/update/remove` | POST/PUT/DELETE | `/api/categories[/{id}]` | product | admin | ✅ hình dạng khớp; public (ISSUE-102) |
| `orders.list` | GET | `/api/orders?page&size` | order | web `/orders`, admin `/orders`, dashboard | ❌ mảng thô (ISSUE-48), mọi user (ISSUE-49) |
| `orders.get` | GET | `/api/orders/{id}` | order | (không dùng) | |
| `orders.create` | POST | `/api/orders` | order | web checkout (`checkout.remote.ts:30-35`) | ⚠️ chỉ 1 `productId`, không `userId`/cart (ISSUE-20, ISSUE-49) |
| `payments.create` | POST | `/api/payments` | payment | web checkout (`checkout.remote.ts:37-41`) | ⚠️ không `userId`; COD vẫn phát sự kiện "thành công" (ISSUE-18) |
| `inventory.checkStock` | GET | `/api/inventory?skuCode=` | inventory | (không dùng) | ❌ service đọc `productName` (ISSUE-51) |

**Service không được frontend gọi**: rating, media, notification, promotion/tax, search, shipping, favourites, carts, auth users/roles (ngoài `/me`).

## 14. Sai lệch hợp đồng (tổng hợp)

| ID | Mức độ | Điểm sai lệch |
| :--- | :--- | :--- |
| ISSUE-48 | Critical | `GET /api/products`, `/api/categories`, `/api/orders` trả mảng thô; BFF đọc `.content` (Page) |
| ISSUE-49 | High | `/api/orders` không lọc theo user; `user_id` integer vs `sub` UUID |
| ISSUE-50 | Medium | `Product.description` không tồn tại ở backend |
| ISSUE-51 | Low | `checkStock` gửi `skuCode`, service đọc `productName` |
| ISSUE-60 | High | `ProductPayload.categoryId` phẳng vs `ProductDto.category.categoryId` lồng nhau → danh mục bị bỏ khi tạo/sửa sản phẩm |
| ISSUE-100 | Critical | `/users/me` thiếu `roles`, key `fullname` |
| ISSUE-101 | High | User Keycloak không có trong DB authservice → `/users/me` 404 |
| ISSUE-102 | Critical | Ghi catalog qua gateway không cần xác thực |
| ISSUE-29 | Low | Nhiều định dạng response/lỗi |
| — | Info | Type frontend `Rating`, `Notification`, `InventoryItem`, `Favourite` (`frontend/packages/lib/src/types.ts:108-135`) **không khớp** VM backend (vd `Rating.ratingValue` vs `star`, `Notification.isRead` vs `status`, `InventoryItem.skuCode` vs `productName/isInStock`) — hiện không dùng (ISSUE-43) |

## 15. Tham chiếu source

- `services/*/main.go` (đăng ký handler, port), `services/*/internal/handler/*.go`
- `services/{rating,media,notification,promotion}-service/src/index.ts`
- `services/shared/middleware/middleware.go`, `services/shared/response/*.go`
- `frontend/packages/lib/src/server/api.ts`, `frontend/packages/lib/src/types.ts`
- `frontend/apps/web/src/routes/checkout/checkout.remote.ts:15-44`, `frontend/apps/admin/src/routes/products/products.remote.ts:16-31`
- `deploy/apisix/apisix.yaml`
