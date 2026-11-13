# Error Catalog — mọi phản hồi lỗi của backend, gateway và BFF

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> Liệt kê **mọi phản hồi lỗi** mà 11 service (7 Go + 4 Bun), gateway APISIX và BFF SvelteKit (web, admin) có thể trả, với thông điệp **nguyên văn trong code**, điều kiện kích hoạt, vị trí `file:line`, cách frontend xử lý và issue liên quan. Quy ước định dạng tổng quát: [`api-conventions.md` §4](../02-architecture/software/api-conventions.md#4-định-dạng-lỗi). Đặc tả từng service: [`openapi/`](openapi/).
>
> Đường dẫn trong bảng tính từ thư mục service (vd `internal/handler/auth_handler.go` = `services/auth-service/internal/handler/auth_handler.go`) trừ khi ghi đầy đủ.

## 1. Định dạng lỗi đang tồn tại

| # | Kiểu | Ai dùng | Hình dạng | Content-Type | Nguồn |
| :--- | :--- | :--- | :--- | :--- | :--- |
| F1 | `ApiResponse` lỗi | auth-service | `{success:false, code, message, errors?, path, traceId, timestamp}` | `application/json` | `services/shared/response/api_response.go:12-21,55-65` |
| F2 | `ErrorVm` | product, order, payment, inventory, shipping, search | `{statusCode:"400"\|"404"\|"500", title, detail}` — `title` cố định `Bad request` / `Not found` / `Internal server error` | `application/json` | `services/shared/response/response.go:6-34` |
| F3 | Bun `{statusCode, message[, errors]}` | rating, media, promotion (nhóm promotion) | `statusCode` là **chuỗi**; `errors` = mảng issue Zod v4 | `application/json` | vd `services/rating-service/src/index.ts:124` |
| F4 | Bun `{message}` | promotion (nhóm tax) | chỉ `message` | `application/json` | `services/promotion-service/src/index.ts:141-251` |
| F5 | Bun hỗn hợp | notification | text `Error: …`; `{error}`; body rỗng; `false`; hoặc nuốt lỗi trả `[]` | nhiều loại | `services/notification-service/src/index.ts:42-185` |
| F6 | Mặc định framework | Go (chi), Bun (Hono) | chi: 404 `404 page not found` (text/plain), 405 body rỗng, panic → `Recoverer` 500 body rỗng; Hono: 404 `404 Not Found`, exception không bắt → 500 `Internal Server Error` (text/plain) | text/plain | chi mặc định + `chimiddleware.Recoverer` (vd `services/product-service/main.go:82`); Hono mặc định (không có `app.onError` / `app.notFound` ở service nào) |
| F7 | Gateway APISIX | mọi request qua `:9080` | 404 `{"error_msg":"404 Route Not Found"}`; 401 (`openid-connect`) + `WWW-Authenticate`; 429 (`limit-count`); 5xx upstream — body do APISIX/OpenResty sinh | JSON/HTML | `deploy/apisix/apisix.yaml:46-61` |
| F8 | BFF | web, admin | `ApiError(status, message, body)` phía server; ra trình duyệt: trang `+error.svelte`, `{:catch}` trong component, toast, redirect 303 | HTML / JSON remote | `frontend/packages/lib/src/server/api.ts:18-28,62-69` |

Ví dụ:

```json
// F1 — auth-service
{ "success": false, "code": "NOT_FOUND", "message": "User not found", "path": "/auth/api/v1/users/me",
  "traceId": "9f1c…", "timestamp": "2026-11-12T03:04:05Z" }
// F2 — product-service
{ "statusCode": "404", "title": "Not found", "detail": "not found: product with id 42 not found" }
// F3 — rating-service (Zod)
{ "statusCode": "400", "message": "Validation error",
  "errors": [ { "origin": "number", "code": "too_big", "maximum": 5, "inclusive": true, "path": ["star"], "message": "Too big: expected number to be <=5" } ] }
// F4 — promotion-service (tax)
{ "message": "Tax class not found" }
```

Đặc điểm chung: không có mã lỗi máy đọc ngoài auth-service (`code`); mọi 500 trả nguyên `err.Error()` / `err.message` (ISSUE-93); không có 403 theo role ở bất kỳ tầng backend nào (ISSUE-01); 400 không trả lỗi theo field (trừ Zod ở rating/promotion).

---

## 2. Gateway APISIX (F7)

| HTTP | Thông điệp / body | Điều kiện | Nguồn | Frontend xử lý | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 404 | `{"error_msg":"404 Route Not Found"}` | Không route nào khớp — gồm path gốc collection của route `/<svc>/.../*` (vd `POST /media/medias`, `GET /promotion/backoffice/promotions`) | `deploy/apisix/apisix.yaml:62-375` | Không gọi các path này | ISSUE-15 |
| 401 | Body do plugin sinh; header `WWW-Authenticate: Bearer realm="ecommerce"` (thêm `error="invalid_token"` khi token sai) | Route có `openid-connect` (`bearer_only: true`) mà thiếu/sai/hết hạn token | `deploy/apisix/apisix.yaml:142-375` | `ApiError(401)`; hook chủ động refresh trước khi hết hạn (`frontend/apps/web/src/hooks.server.ts:35-46`); nếu vẫn 401 → lỗi chung | ISSUE-01 (không có 403) |
| 429 | Body mặc định APISIX (không cấu hình `rejected_msg`) | > 1200 request / 60 s cho cùng `remote_addr` (thực tế là IP container BFF) | `deploy/apisix/apisix.yaml:55-61` | Như mọi lỗi khác (không retry/backoff) | ISSUE-61 |
| 502/503 | Trang lỗi mặc định APISIX/OpenResty | Upstream không tồn tại (`upstream_id` 7, 13) hoặc container service không chạy | `deploy/apisix/apisix.yaml:250-267,352-375` | Lỗi chung | ISSUE-13 |
| 504 | Trang lỗi mặc định | Upstream không phản hồi trong timeout mặc định (không cấu hình timeout) | `deploy/apisix/apisix.yaml:1-45` | Lỗi chung | — |
| (CORS) | Thiếu header `Access-Control-Allow-Origin` | Origin không thuộc allowlist (vd `http://localhost:3001`) — chỉ ảnh hưởng nếu trình duyệt gọi trực tiếp | `deploy/apisix/apisix.yaml:50-54` | Không áp dụng (BFF gọi server-side) | ISSUE-14 |

> Định dạng chính xác của body 401/429/5xx phụ thuộc phiên bản APISIX (`apache/apisix:3.15.0-debian`) và chưa được kiểm chứng runtime; client không nên parse các body này.

---

## 3. auth-service (Go, F1)

Handler: `internal/handler/auth_handler.go`. Sentinel: `internal/service/user_service.go:13-20`.

| HTTP | `code` | `message` (nguyên văn) | Endpoint / điều kiện | Nguồn | Frontend xử lý | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 400 | `BAD_REQUEST` | `Invalid request body` | `POST /auth/signup` body không phải JSON | `:102-104` | Không dùng (không có UI signup) | ISSUE-71 |
| 400 | `VALIDATION_FAILED` | `Username, password and email are required` | signup thiếu `username`/`password`/`email` | `:107-109` | — | |
| 409 | `CONFLICT` | `username already exists: <u>` · `email already exists: <e>` · `phone already exists: <p>` | Trùng trong DB `authservice` | `:114-116`; `internal/service/user_service.go:40,48,57` | — | ISSUE-55 |
| 400 | `REGISTRATION_FAILED` | `err.Error()` — vd `failed to create keycloak user: <body> (status <n>)`, `db not connected` | Lỗi khác khi đăng ký (Keycloak ≥ 400, DB…) — kể cả lỗi DB cũng là **400** | `:118`; `internal/service/keycloak_client.go:230`; `internal/repository/user_repository.go:123` | — | ISSUE-08, ISSUE-93 |
| 302 | — | Redirect `SSO_DEFAULT_FRONTEND_REDIRECT?error=<error>` | `GET /auth/callback` có `?error=` từ Keycloak | `:142-145` | Callback page hiển thị `auth_callbackErrorTitle` "Đăng nhập không thành công" + `auth_callbackErrorText` (`frontend/apps/web/src/routes/auth/callback/+page.server.ts:15-18`) | |
| (không lỗi) | — | Mock token `mock-access-token` | Đổi `code` thất bại → vẫn cấp ticket | `:154-162` | Ticket đổi được nhưng token giả → 401 ở gateway về sau | ISSUE-04, ISSUE-52 |
| 400 | `BAD_REQUEST` | `Ticket is required` | `GET /auth/session` thiếu `ticket` | `:171-173` | `exchangeTicket` ném `ApiError` → callback `failed: true` (`+page.server.ts:34-36`) | |
| 401 | `UNAUTHORIZED` | `Invalid or expired ticket` | Ticket không tồn tại/hết hạn (2 phút)/đã dùng | `:176-179` | Như trên | ISSUE-06 |
| 400 | `BAD_REQUEST` | `Refresh token is required` | `POST /auth/refresh` body sai/thiếu `refreshToken` | `:187-189` | `hooks.server.ts` bắt lỗi → `clearSession` (người dùng thành khách) | |
| 401 | `UNAUTHORIZED` | `err.Error()` — vd `keycloak token error: <body Keycloak> (status 400)` | Keycloak từ chối refresh | `:193-195`; `internal/service/keycloak_client.go:281` | Như trên (`frontend/apps/web/src/hooks.server.ts:42-45`) | ISSUE-93 |
| 200 | — | `Logout successful` (luôn) | `POST /auth/logout` — lỗi Keycloak bị nuốt | `:201-208` | `auth.remote.ts` bỏ qua lỗi, luôn xóa cookie | |
| 400 | `INVALID_ID` | `User ID must be numeric` | `{id}`/`{userId}` không phải số (PUT/DELETE/GET user, assign/revoke/get roles) | `:217,247,286,348,375,402` | Không dùng | |
| 400 | `BAD_REQUEST` | `Invalid request body` | `PUT /users/{id}` body sai | `:222-224` | — | |
| 400 | `BAD_REQUEST` | `username query parameter is required` | `GET /users` thiếu `username` | `:265-267` | — | |
| 400 | `BAD_REQUEST` | `Role name cannot be empty` | assign/revoke body text rỗng | `:354-356,381-383` | — | ISSUE-08 |
| 401 | `AUTH_TOKEN_INVALID` | `Invalid or missing token` | `GET /users/me` không có `preferred_username`/`sub` (token parse không verify) | `:323-325` | `profile()` lỗi → `user = null` (`+page.server.ts:27-31`) | ISSUE-02, ISSUE-67 |
| 404 | `NOT_FOUND` | `User not found` | user không có trong DB `authservice` (`/users/me`, get/put/delete) | `:230,253,273,293,331` | `/users/me` 404 → phiên không có user → bị coi là khách | ISSUE-101, ISSUE-67 |
| 409 | `CONFLICT` | `Role already assigned or user/role not found` | assign không chèn được | `:364-366` | — | |
| 409 | `CONFLICT` | `Role not present on user` | revoke không xóa được | `:391-393` | — | |
| 500 | `INTERNAL_ERROR` | `err.Error()` — vd `db not connected`, lỗi pgx | Lỗi DB ở các handler user/role | `:233,256,276,296,311,334,361,388,408` | Lỗi chung | ISSUE-28, ISSUE-93 |
| 200 | `OK` | `Password change request processed` (stub) | `PUT /users/me/password` luôn thành công, không đổi gì | `:240-242` | — | ISSUE-07 |

---

## 4. product-service (Go, F2) — gồm favourites

Handler: `internal/handler/product_handler.go` (P), `internal/handler/favourite_handler.go` (F). Sentinel `ErrNotFound = "not found"` (`internal/service/product_service.go:12`).

| HTTP | `detail` (nguyên văn) | Endpoint / điều kiện | Nguồn | Frontend xử lý | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 400 | `Invalid categoryId` | `{categoryId}` không phải số (GET/PUT/DELETE) | P `:63,109,128` | Admin: toast `common_noData` "Chưa có dữ liệu" khi xóa/lưu lỗi (`frontend/apps/admin/src/routes/categories/+page.svelte:43-44,160-161`) | ISSUE-83, ISSUE-115 |
| 400 | `Invalid productId` | `{productId}` không phải số | P `:162,208,227` | Web: chi tiết SP `{:catch}` → "Không tìm thấy sản phẩm" (`frontend/apps/web/src/routes/products/[id]/+page.svelte:159-167`) | ISSUE-115 |
| 400 | `Invalid request body` | Body category/product không phải JSON / sai kiểu | P `:81,95,114,180,194,213` | Admin toast `common_noData` (`frontend/apps/admin/src/routes/products/+page.svelte:203-204`) | ISSUE-35 |
| 404 | `not found: category with id <n> not found` | GET category không tồn tại | P `:68-70`; `internal/service/product_service.go:34` | (không dùng `categories.get`) | |
| 404 | `not found: product with id <n> not found` | GET product không tồn tại | P `:167-169`; `product_service.go:73` | Web: "Không tìm thấy sản phẩm" (status thật bị mất — trang trả 200) ; checkout: `placeOrder` lỗi → toast `checkout_orderFailed` | ISSUE-82, ISSUE-38 |
| 500 | `err.Error()` (pgx) | Lỗi DB ở list/get/create/update/delete | P `:54,72,86,100,119,133,153,171,185,199,218,232` | Web catalog `{:catch}` → `productgrid_loadError` "Không thể tải sản phẩm. Vui lòng thử lại."; admin `{:catch}` → `common_noData` | ISSUE-93, ISSUE-55 (UNIQUE `sku`) |
| — | (không có 404) | PUT/DELETE id không tồn tại vẫn 200 | P `:92-136,191-235` | — | ISSUE-65 |
| 500 | `err.Error()` | `GET /favourites` lỗi DB | F `:48-50` | Không dùng | ISSUE-13 |
| 400 | `Invalid request payload` | Body favourite sai (find/save/update/delete) | F `:58,106,125` | — | ISSUE-35 |
| 400 | `Invalid userId` · `Invalid productId` | Path favourite không phải số | F `:77,82,140,145` | — | |
| 404 | `Favourite not found` | Không có bản ghi | F `:67-69,96-98` | — | |
| 500 | `invalid likeDate: <parse error>` | `likeDate` sai định dạng khi tìm/xóa | F `:63-65,92-94,155-157`; `internal/repository/favourite_repo.go:93,152` | — | ISSUE-74 |

---

## 5. order-service (Go, F2)

Handler `internal/handler/order_handler.go`. Sentinel `cart not found`, `order not found` (`internal/service/order_service.go:13-14`).

| HTTP | `detail` | Endpoint / điều kiện | Nguồn | Frontend xử lý | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 400 | `Invalid cartId` | `{cartId}` không phải số | `:90,147,172` | Không dùng carts | |
| 400 | `Invalid orderId` | `{orderId}` không phải số; `existOrderId` thiếu/sai `orderId` | `:224,243,296,321` | — | |
| 400 | `Invalid request body` | Body cart/order sai | `:109,124,153,258,273,302` | Checkout: toast `checkout_orderFailed` "Đặt hàng thất bại. Vui lòng thử lại." (`frontend/apps/web/src/routes/checkout/+page.svelte:61-62`) | |
| 400 | `cartId is required for update` | `PUT /carts` không có `cartId` | `:127-129` | — | |
| 400 | `orderId is required for update` | `PUT /orders` không có `orderId` | `:276-278` | — | |
| 404 | `cart not found` | Cart không tồn tại (GET/PUT) | `:96-98,134-136,159-161` | — | |
| 404 | `order not found` | Order không tồn tại (GET/PUT) | `:230-232,283-285,308-310` | — | |
| 500 | `err.Error()` | Lỗi DB (list, paged, save, delete, exists…) | `:67,81,100,115,138,163,177,201,215,234,249,264,287,312,326` | Web `/orders` `{:catch}` → `common_error` "Đã xảy ra lỗi. Vui lòng thử lại." (`frontend/apps/web/src/routes/orders/+page.svelte:80-88`); admin `{:catch}` → `common_noData`; checkout → toast | ISSUE-93, ISSUE-24 |

---

## 6. payment-service (Go, F2)

Handler `internal/handler/payment_handler.go`. Sentinel `payment not found`, `Order has already been paid.` (`internal/service/payment_service.go:14-15`).

| HTTP | `detail` | Endpoint / điều kiện | Nguồn | Frontend xử lý | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 400 | `Invalid paymentId` | `{paymentId}` không phải số | `:81,149,174` | — | |
| 400 | `Invalid orderId` | `GET /getOrder/{orderId}` không phải số | `:99-101` | — | ISSUE-25 |
| 400 | `Invalid request body` | Body sai (POST/PUT) | `:111,130,155` | Checkout toast `checkout_orderFailed` | |
| **404** | `Order has already been paid.` | `POST /api/payments` cho order đã có payment `is_payed = true` | `:117-119`; `internal/service/payment_service.go:88-96` | Checkout toast chung — order đã tạo nhưng payment thất bại (không bù trừ) | ISSUE-25, ISSUE-75, ISSUE-23 |
| 404 | `payment not found` | PUT/GET payment không tồn tại | `:87-89,136-138,161-163` | — | |
| 500 | `paymentId is required for update` | `PUT /api/payments` thiếu `paymentId` (lẽ ra 400) | `:139-141`; `payment_service.go:118-121` | — | ISSUE-93 |
| 500 | `err.Error()` | Lỗi DB | `:58,72,91,121,140,165,179` | Checkout toast | ISSUE-93 |
| (không lỗi) | — | Gửi Kafka thất bại bị nuốt, vẫn 200 | `payment_service.go:110-113` | — | ISSUE-18, ISSUE-57 |

---

## 7. inventory-service (Go, F2)

| HTTP | `detail` | Endpoint / điều kiện | Nguồn | Frontend xử lý | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 500 | `err.Error()` | Lỗi DB ở `GET /api/inventory`, `GET /api/inventory/{products}` | `internal/handler/inventory_handler.go:58-61,69-72` | `checkStock` không được gọi | ISSUE-51, ISSUE-93 |
| — | (không có 400) | Thiếu `productName` → mảng rỗng 200 | `inventory_handler.go:53-63` | — | ISSUE-26 |

## 8. shipping-service (Go, F2)

Handler `internal/handler/shipping_handler.go`. Sentinel `not found` (`internal/service/shipping_service.go:12`).

| HTTP | `detail` | Endpoint / điều kiện | Nguồn | Frontend xử lý | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 400 | `Invalid orderId` · `Invalid productId` | Path không phải số | `:62,67,135,140` | Không dùng | |
| 400 | `Invalid request body` | Body sai (find, save, update, delete) | `:86,105,120,155` | — | ISSUE-81 |
| 404 | `not found: order item not found for orderId: <o>, productId: <p>` | Không có dòng | `:73-75,92-94`; `internal/service/shipping_service.go:32` | — | |
| 500 | `err.Error()` | Lỗi DB | `:51,77,96,111,126,146,161` | — | ISSUE-93 |

## 9. search-service (Go, F2)

| HTTP | `detail` | Endpoint / điều kiện | Nguồn | Frontend xử lý | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 500 | `err.Error()` | Chỉ khi encode query JSON lỗi (thực tế không xảy ra) | `internal/handler/search_handler.go:78-81,89-92`; `internal/service/search_service.go:107-109,207-209` | Không được web gọi | ISSUE-78 |
| 200 (che lỗi) | Trang rỗng | OpenSearch lỗi/không kết nối → `products: []` | `internal/service/search_service.go:119-125,219-222` | — | ISSUE-34 |
| — | (không có 400) | `page`/`size`/`minPrice`/`maxPrice` sai → dùng default/bỏ qua | `search_handler.go:39-48,62-76` | — | |

---

## 10. rating-service (Bun, F3)

Nguồn `services/rating-service/src/index.ts`.

| HTTP | Body (nguyên văn) | Endpoint / điều kiện | Nguồn | Frontend xử lý | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 400 | `{statusCode:"400", message:"Invalid productId"}` | `GET /storefront/ratings/products/{productId}` không phải số | `:105-107` | Không được frontend gọi | |
| 400 | `{statusCode:"400", message:"Validation error", errors:[…Zod]}` | `POST /storefront/ratings` sai schema (`star` ngoài 1–5, thiếu `productId`…) | `:122-125`; `src/types.ts:3-8` | — | ISSUE-62 |
| 400 | `{statusCode:"400", message:"Invalid ID"}` | `DELETE /backoffice/ratings/{id}` không phải số | `:91-93` | — | |
| 500 | `{statusCode:"500", message:<err.message>}` | Body không phải JSON, lỗi DB, `createdFrom/To` sai, `pageNo` âm | `:83,97,114,131` | — | ISSUE-93, ISSUE-111 |
| 200 (che lỗi) | `0` | `average-star` với id sai hoặc lỗi DB | `:139-146` | — | ISSUE-65 |
| 200 (sai nghĩa) | `{title:"Delete Rating", message:"The request has been processed successfully", statusCode:"200"}` | Xóa id không tồn tại | `src/service.ts:118-125` | — | ISSUE-65 |
| 404 | `404 Not Found` (text) | Path không khớp (Hono mặc định) | — | — | |

## 11. media-service (Bun, F3)

Nguồn `services/media-service/src/index.ts`.

| HTTP | Body | Endpoint / điều kiện | Nguồn | Frontend xử lý | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 400 | `{statusCode:"400", message:"File is required"}` | `POST /medias` không có `multipartFile`/`file` (hoặc body không phải multipart) | `:46-48` | Không được frontend gọi | ISSUE-32 |
| 400 | `{statusCode:"400", message:"Invalid ID"}` | `{id}` không phải số (GET/DELETE/file) | `:63-65,77-79,117-119` | — | |
| 400 | `{statusCode:"400", message:"ids query param is required"}` | `GET /medias` thiếu `ids` | `:94-96` | — | |
| 404 | `{statusCode:"404", message:"Media not found"}` | `GET /medias/{id}` không tồn tại | `:81-83` | — | |
| 404 | `{statusCode:"404", message:"No medias found"}` | `GET /medias?ids=` không có bản ghi nào | `:103-105` | — | |
| 404 | `{statusCode:"404", message:"File not found"}` | Sai `fileName`, không có object S3 | `:121-123`; `src/service.ts:101-118` | — | ISSUE-63 |
| 500 | `{statusCode:"500", message:<err.message>}` | Lỗi S3/DB; có thể cả khi tên file chứa ký tự ngoài Latin-1 lúc set `Content-Disposition` | `:55,69,86,108,128` | — | ISSUE-93, ISSUE-113, ISSUE-65f |
| 204 (sai nghĩa) | body rỗng | Xóa id không tồn tại | `:66-67` | — | ISSUE-65 |

## 12. promotion-service (Bun, F3 + F4)

Nguồn `services/promotion-service/src/index.ts`.

| HTTP | Body | Endpoint / điều kiện | Nguồn | Frontend xử lý | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 400 | `{statusCode:"400", message:"Invalid ID"}` | `{promotionId}` không phải số | `:74,118` | Không được frontend gọi | |
| 400 | `{statusCode:"400", message:"Validation error", errors:[…Zod]}` | POST/PUT promotion sai schema | `:89-91,104-106`; `src/types.ts:3-23` | — | ISSUE-62 |
| 404 | `{statusCode:"404", message:"Promotion not found"}` | GET promotion không tồn tại | `:77` | — | |
| 500 | `{statusCode:"500", message:<err.message>}` | Body không phải JSON, lỗi DB, ngày sai, **PUT id không tồn tại** (TypeError) | `:66,80,95,110,123`; `src/service.ts:112-136` | — | ISSUE-65c, ISSUE-93 |
| 200 (che lỗi) | trang rỗng | Lỗi DB khi list promotion / phân trang tax class | `src/service.ts:73-79`; `src/tax-service.ts:96-105` | — | ISSUE-28 |
| 200 (sai nghĩa) | body rỗng | DELETE promotion id không tồn tại | `:120-121` | — | ISSUE-65d |
| 400 | `{message:"Invalid tax class ID"}` | `{id}` tax class không phải số | `:141,163,174` | — | |
| 400 | `{message:"Invalid tax rate ID"}` | `{id}` tax rate không phải số | `:215,237,248` | — | |
| 400 | `{message:"Invalid payload"}` | Body tax class / tax rate sai schema (không trả chi tiết) | `:155,166,229,240` | — | ISSUE-62 |
| 400 | `{message:"taxClassId and countryId are required"}` | `tax-percent` thiếu/sai tham số | `:191-193` | — | |
| 400 | `{message:"taxClassIds and countryId are required"}` | `location-based-batch` thiếu/sai tham số | `:203-205` | — | |
| 404 | `{message:"Tax class not found"}` | Tax class không tồn tại (GET/PUT/DELETE) | `:143,168,176` | — | |
| 404 | `{message:"Tax rate not found"}` | Tax rate không tồn tại (GET/PUT/DELETE) | `:217,242,250` | — | |
| 500 | `Internal Server Error` (text/plain) | Exception không bắt ở route tax: body không phải JSON, trùng tên tax class (UNIQUE), FK `tax_class_id` sai/đang được dùng, lỗi DB | `:130-253` (không try/catch); `src/db.ts:73-86` | — | ISSUE-94, ISSUE-110 |

## 13. notification-service (Bun, F5)

Nguồn `services/notification-service/src/index.ts`.

| HTTP | Body | Endpoint / điều kiện | Nguồn | Frontend xử lý | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 200 (che lỗi) | `Error while Sending Mail: <err.message>` (text) | SMTP lỗi khi gửi | `src/mailer.ts:43-46` | Không được frontend gọi | ISSUE-65a, ISSUE-117 |
| 200 | `Mail Sent Successfully (Mock)` | SMTP chưa cấu hình | `src/mailer.ts:30-33` | — | |
| 500 | `Error: <err.message>` (text/plain) | Body email không phải JSON | `:48,58,70` | — | ISSUE-94 |
| 200 (che lỗi) | `[]` | Lỗi DB khi list notifications / payment-notifications | `:79-81,132-134` | — | ISSUE-65 |
| 404 | body rỗng | `GET /notifications/{id}`, `GET /payment-notifications/{id}` không tồn tại | `:88,146` | — | ISSUE-94 |
| 500 | body rỗng | Lỗi DB (vd id không phải số → `NaN`) khi GET theo id / DELETE payment-notification | `:91,149,183` | — | ISSUE-94 |
| 500 | `{error:<err.message>}` | Lỗi insert notification / payment-notification (body không phải JSON…) | `:110,169` | — | ISSUE-93 |
| 500 | `false` | Lỗi DB khi xóa notification | `:120` | — | ISSUE-94 |
| 200/204 (sai nghĩa) | `true` / rỗng | Xóa id không tồn tại | `:118,181` | — | ISSUE-65d |
| (Kafka) | chỉ log | Lỗi xử lý message — không phản hồi HTTP | `src/kafka.ts:57-59,62-64` | — | ISSUE-80 |

---

## 14. BFF SvelteKit (web, admin — F8)

### 14.1 Cách `ApiError` được tạo

`request()` (`frontend/packages/lib/src/server/api.ts:39-76`): mọi response không `ok` → `ApiError(status, message, body)` với `message` = `body.message` nếu là chuỗi không rỗng, ngược lại `"<status> <statusText>"`. Hệ quả theo định dạng:

| Định dạng backend | `ApiError.message` | Ghi chú |
| :--- | :--- | :--- |
| F1 auth | `message` (vd `User not found`) | |
| F2 ErrorVm | `"404 Not Found"`, `"500 Internal Server Error"` | `detail` bị bỏ qua (ISSUE-115) |
| F3 Bun | `message` | |
| F4/F5/F6/F7 | `"<status> <statusText>"` (hoặc `message` nếu có) | body text được giữ trong `ApiError.body` |

Không có `handleError` hook ở web/admin → `ApiError` ném ra từ remote function/`load` bị SvelteKit coi là lỗi bất ngờ: trình duyệt nhận **500** `{ message: "Internal Error" }`, status gốc bị mất (ISSUE-115).

### 14.2 Lỗi do BFF tự sinh

| Loại | Status / hành vi | Thông điệp | Điều kiện | Nguồn | Issue |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `error()` | 401 | `Your session has expired. Please sign in again.` | `placeOrder` không có access token | `frontend/apps/web/src/routes/checkout/checkout.remote.ts:16` | ISSUE-83 (tiếng Anh cứng) |
| Redirect | 303 → `/login?redirect=<path>` | — | Web `/checkout`, `/orders` khi chưa đăng nhập | `frontend/apps/web/src/routes/checkout/+page.server.ts:8`, `frontend/apps/web/src/routes/orders/+page.server.ts:8` | ISSUE-66 |
| Redirect | 303 → `/login?redirect=<path>` | — | Admin: mọi path không public khi chưa đăng nhập | `frontend/apps/admin/src/hooks.server.ts:67-70` | ISSUE-68 |
| Redirect | 303 → `/access-denied` | Trang `access_denied_title` "Không có quyền truy cập" | Admin: user không có role `ADMIN` trong cookie | `frontend/apps/admin/src/hooks.server.ts:72-74` | ISSUE-100 (luôn xảy ra) |
| Callback thất bại | 200 trang lỗi | `auth_callbackErrorTitle` "Đăng nhập không thành công" · `auth_callbackErrorText` "Phiên đăng nhập đã hết hạn hoặc không hợp lệ. Vui lòng thử lại." | Thiếu `ticket`, có `error`, hoặc `exchangeTicket` lỗi | `frontend/apps/web/src/routes/auth/callback/+page.server.ts:14-37`; `+page.svelte:16-22` | |
| Refresh thất bại | Xóa cookie, tiếp tục như khách | — | `auth.refresh` ném lỗi | `frontend/apps/web/src/hooks.server.ts:42-45` (admin tương tự) | ISSUE-69 |
| Validate Zod (remote) | 400 từ SvelteKit (`invalid`) | Thông điệp Zod tiếng Anh, vd `Please enter your full name`, `Your cart is empty` | Input remote `query/form/command` không qua schema | `frontend/packages/lib/src/schemas.ts:32-82` | ISSUE-83 |

### 14.3 Trang lỗi & hiển thị

| Nơi | Hiển thị (vi) | Khi nào | Nguồn |
| :--- | :--- | :--- | :--- |
| `+error.svelte` (web) | 404: "Không tìm thấy trang" / "Trang có thể đã được di chuyển hoặc bị xoá."; khác: "Đã xảy ra lỗi" / "Dịch vụ có thể tạm thời không kết nối được. Vui lòng thử lại."; tiêu đề `{status} \| EzBuy`; nút "Về trang chủ", "Thử lại" | Route không tồn tại; lỗi trong `load` | `frontend/apps/web/src/routes/+error.svelte:9-28`; `frontend/apps/web/messages/vi.json:139-145` |
| `+error.svelte` (admin) | Như web, tiêu đề `{status} \| EzBuy Admin` | Như trên | `frontend/apps/admin/src/routes/+error.svelte:9-28`; `frontend/apps/admin/messages/vi.json:50-56` |
| Web catalog `{:catch}` | "Không thể tải sản phẩm. Vui lòng thử lại." | `getProducts` lỗi | `frontend/apps/web/src/routes/products/+page.svelte:241-249` |
| Web sidebar danh mục `{:catch}` | `—` | `getCategories` lỗi | `frontend/apps/web/src/routes/products/+page.svelte:112-113` |
| Web chi tiết SP `{:catch}` | "Không tìm thấy sản phẩm" — cho **mọi** lỗi (kể cả 500, 429) | `getProduct` lỗi | `frontend/apps/web/src/routes/products/[id]/+page.svelte:159-167` |
| Web đơn hàng `{:catch}` | "Đã xảy ra lỗi. Vui lòng thử lại." | `getOrders` lỗi | `frontend/apps/web/src/routes/orders/+page.svelte:80-88` |
| Web checkout | Toast "Đặt hàng thất bại. Vui lòng thử lại." — cho mọi lỗi (401, 404 "already paid", 500…) | `placeOrder` ném lỗi | `frontend/apps/web/src/routes/checkout/+page.svelte:61-62` |
| Admin products/categories/orders/dashboard `{:catch}` | "Chưa có dữ liệu" (`common_noData`) — hiển thị như trạng thái rỗng | Query lỗi | `frontend/apps/admin/src/routes/products/+page.svelte:182-184`, `categories/+page.svelte:139`, `orders/+page.svelte:125-127`, `dashboard/+page.svelte:135-137` |
| Admin toast lưu/xóa | "Chưa có dữ liệu" | `saveProduct`/`deleteProduct`/`saveCategory`/`deleteCategory` lỗi | `frontend/apps/admin/src/routes/products/+page.svelte:48-49,203-204`; `categories/+page.svelte:43-44,160-161` |

---

## 15. Ma trận tổng hợp status theo service

| Service | 400 | 401 | 404 | 409 | 429 (gateway) | 500 | 5xx gateway | Định dạng |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :--- |
| auth | ✅ | ✅ (service + gateway) | ✅ | ✅ | ✅ | ✅ | ✅ | F1 |
| product | ✅ | — (route public) / ✅ favourites | ✅ | — | ✅ | ✅ | ✅ (favourite upstream 7) | F2 |
| order | ✅ | ✅ gateway | ✅ | — | ✅ | ✅ | ✅ | F2 |
| payment | ✅ | ✅ gateway | ✅ (gồm "already paid") | — | ✅ | ✅ | ✅ | F2 |
| inventory | — | ✅ gateway | — | — | ✅ | ✅ | ✅ | F2 |
| shipping | ✅ | ✅ gateway | ✅ | — | ✅ | ✅ | ✅ | F2 |
| search | — | — | — | — | ✅ | (hiếm) | ✅ | F2 |
| rating | ✅ | ✅ gateway (backoffice) | — | — | ✅ | ✅ | ✅ | F3 |
| media | ✅ | ✅ gateway | ✅ | — | ✅ | ✅ | ✅ | F3 |
| promotion/tax | ✅ | ✅ gateway | ✅ | — | ✅ | ✅ (JSON / text) | ✅ (tax upstream 13) | F3 + F4 |
| notification | — | ✅ gateway | ✅ (rỗng) | — | ✅ | ✅ (text / JSON / rỗng / `false`) | ✅ | F5 |

---

## 16. Đề xuất mô hình lỗi thống nhất (**PROPOSAL** — chưa hiện thực)

> Phần này là **đề xuất**, không mô tả code hiện tại. Liên quan ISSUE-29, ISSUE-93, ISSUE-94, ISSUE-115.

### 16.1 Định dạng: RFC 9457 Problem Details

```http
HTTP/1.1 409 Conflict
Content-Type: application/problem+json
X-Correlation-Id: 6f1c2a…

{
  "type": "https://errors.ecommerce.local/payment/order-already-paid",
  "title": "Order has already been paid",
  "status": 409,
  "code": "PAYMENT_ORDER_ALREADY_PAID",
  "detail": "Order 34 already has a completed payment.",
  "instance": "/api/payments",
  "traceId": "6f1c2a…",
  "errors": [ { "field": "orderId", "code": "ALREADY_PAID", "message": "…" } ]
}
```

| Field | Bắt buộc | Quy tắc |
| :--- | :---: | :--- |
| `status` | ✅ | Số nguyên, bằng HTTP status |
| `code` | ✅ | `UPPER_SNAKE`, ổn định, có prefix miền (`AUTH_`, `CATALOG_`, `ORDER_`, `PAYMENT_`, `TAX_`, `MEDIA_`, `RATING_`, `NOTIFICATION_`, `GATEWAY_`, `COMMON_`) — client và i18n dựa vào field này |
| `title` | ✅ | Mô tả ngắn cố định theo `code` (tiếng Anh) |
| `detail` | — | Thông điệp an toàn cho người dùng; **không** chứa lỗi SQL/stack (thay ISSUE-93) |
| `traceId` | ✅ | = `X-Correlation-Id` (gateway sinh bằng plugin `request-id`, BFF truyền tiếp — ISSUE-98) |
| `errors[]` | — | Lỗi theo field cho 400/422 (`field`, `code`, `message`) — map từ Zod issue / validate Go |

### 16.2 Bảng ánh xạ đề xuất

| Tình huống hiện tại | Hiện trả | Đề xuất |
| :--- | :--- | :--- |
| Body không phải JSON (mọi service) | 400 / 500 / text | 400 `COMMON_MALFORMED_BODY` |
| Validate input (Zod, thiếu field) | 400 nhiều dạng | 400 `COMMON_VALIDATION_FAILED` + `errors[]` |
| id path không phải số | 400 `Invalid ID`… | 400 `COMMON_INVALID_ID` |
| Không tìm thấy | 404 / 500 (PUT promotion) / 200 (DELETE) | 404 `<DOMAIN>_NOT_FOUND`; DELETE id không tồn tại → 404 hoặc 204 nhất quán |
| Trùng khóa (username, sku, tax class name) | 409 / 500 / 400 | 409 `<DOMAIN>_DUPLICATE` (map lỗi Postgres `23505`) |
| FK sai/đang được dùng | 500 | 409 `<DOMAIN>_IN_USE` / 422 `<DOMAIN>_REFERENCE_NOT_FOUND` (map `23503`) |
| Order đã thanh toán | **404** | 409 `PAYMENT_ORDER_ALREADY_PAID` |
| Token thiếu/sai (gateway) | 401 body APISIX | 401 `GATEWAY_UNAUTHENTICATED` (plugin `response-rewrite`/custom error page) |
| Thiếu role (chưa có) | — | 403 `GATEWAY_FORBIDDEN` / `<DOMAIN>_FORBIDDEN` (cần ISSUE-01) |
| Rate limit | 429 body APISIX | 429 `GATEWAY_RATE_LIMITED` + `Retry-After` (cấu hình `rejected_msg`) |
| Phụ thuộc lỗi (DB, OpenSearch, SMTP, S3) | 500 lộ lỗi / 200 che lỗi | 503 `<DOMAIN>_DEPENDENCY_UNAVAILABLE`, log chi tiết server-side; **không** trả 200 rỗng |
| Lỗi không lường trước | 500 lộ `err.message` / text | 500 `COMMON_INTERNAL_ERROR`, `detail` chung |

### 16.3 Hiện thực đề xuất (theo tầng)

| Tầng | Việc cần làm |
| :--- | :--- |
| `services/shared/response` (Go) | Thay `ErrorVm` + `ApiResponse` lỗi bằng `Problem`; helper `response.Problem(w, r, status, code, detail, errs...)`; `Recoverer` trả Problem |
| Bun | `app.onError` + `app.notFound` dùng chung (package nội bộ), `safeParse` trả `errors[]`; bọc `c.req.json()` |
| APISIX | `request-id` plugin; `rejected_msg` cho `limit-count`; trang lỗi JSON cho 401/404/5xx |
| BFF | `ApiError` đọc `code`/`detail`/`traceId`; `handleError` hook giữ status + `code`; UI map `code` → key i18n thay vì thông điệp chung |

---

## 17. Issue liên quan

| ID | Mức độ | Tóm tắt |
| :--- | :--- | :--- |
| ISSUE-29 | Low | Nhiều định dạng response/lỗi |
| ISSUE-93 | Medium | 500 lộ thông điệp lỗi nội bộ |
| ISSUE-94 | Low | Route tax/notification trả text/plain hoặc body rỗng |
| ISSUE-65 | Low | Mã lỗi sai / che lỗi ở service Bun |
| ISSUE-25 | Medium | "Order has already been paid." trả 404 |
| ISSUE-110 | Medium | (mới) Vi phạm UNIQUE/FK ở tax → 500; `Location` sai |
| ISSUE-111 | Medium | (mới) Phân trang Bun không giới hạn / không validate |
| ISSUE-113 | Medium | (mới) `Content-Disposition` không escape tên file |
| ISSUE-115 | Low | (mới) BFF làm mất status/thông điệp lỗi backend |
