# Thuật ngữ (Glossary)

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

## 1. Nghiệp vụ

| Thuật ngữ | Định nghĩa trong hệ thống này |
| :--- | :--- |
| **Storefront** | Website dành cho khách hàng (`frontend/apps/web`). |
| **Backoffice** | Website quản trị nội bộ (`frontend/apps/admin`). |
| **Guest** | Người truy cập chưa đăng nhập. Được xem catalog, chi tiết sản phẩm, dùng giỏ hàng (lưu ở trình duyệt). |
| **Customer** | Người dùng đã đăng nhập với role `USER`. Được đặt hàng và xem lịch sử đơn. |
| **Admin** | Người dùng có role `ADMIN` (hoặc `ROLE_ADMIN`) trong token. Được vào backoffice. |
| **Category (Danh mục)** | Nhóm sản phẩm; hỗ trợ danh mục cha (`parent_category_id`). |
| **Product (Sản phẩm)** | Mặt hàng: tiêu đề, ảnh, SKU (unique), đơn giá, số lượng, danh mục. Không có biến thể, thương hiệu, mô tả. |
| **SKU** | Mã định danh hàng hóa, cột `products.sku`, unique. |
| **Cart (Giỏ hàng)** | Hai khái niệm khác nhau: (1) **giỏ phía client** — lưu trong `localStorage` của storefront; (2) **bảng `carts`** trong order-service — bản ghi gắn `user_id`, đóng vai trò nhóm các order. Storefront **không** đồng bộ (1) lên (2). |
| **Order (Đơn hàng)** | Bản ghi trong `orders`: ngày, mô tả, `order_fee` (tổng tiền), **một** `product_id`, `cart_id`. Không có trạng thái. |
| **Order item** | Bảng `order_items` trong **shipping-service** (order_id, product_id, ordered_quantity). Không được storefront sử dụng. |
| **Payment (Thanh toán)** | Bản ghi trong `payments`: `order_id`, `user_id`, `is_payed`, `payment_status`. Hệ thống chỉ **ghi nhận**, không xử lý tiền với cổng thanh toán. |
| **Payment status** | Chuỗi tự do; frontend dùng `NOT_STARTED`, `IN_PROGRESS`, `COMPLETED`, `FAILED`, `CANCELLED`. |
| **COD / ONLINE** | Phương thức thanh toán chọn ở checkout. COD → `isPayed=false, NOT_STARTED`; ONLINE → `isPayed=true, IN_PROGRESS` (không có cổng thanh toán thật). |
| **Free shipping threshold** | Ngưỡng miễn phí vận chuyển: tạm tính ≥ 500.000 ₫ thì phí ship = 0, ngược lại 30.000 ₫ (`frontend/packages/lib/src/cart.ts`). |
| **Favourite (Yêu thích)** | Bộ ba (`user_id`, `product_id`, `like_date`) trong product-service. |
| **Inventory (Tồn kho)** | Bảng `inventory(product_name, quantity)` trong inventory-service; tách biệt với `products.quantity`. |
| **Promotion (Khuyến mãi)** | Bản ghi chương trình khuyến mãi / coupon trong promotion-service; chỉ được quản lý, không được áp dụng vào giá. |
| **Tax class / Tax rate** | Nhóm thuế và mức thuế (%) theo quốc gia trong promotion-service. |
| **Rating (Đánh giá)** | Số sao 1–5 kèm nội dung cho một sản phẩm. |
| **Media** | File (ảnh, tài liệu) lưu trên S3 kèm metadata trong `media`. |
| **Notification** | Thông báo in-app (`notifications`) và thông báo thanh toán (`payment_notifications`) trong notification-service. |

## 2. Kỹ thuật

| Thuật ngữ | Định nghĩa |
| :--- | :--- |
| **BFF (Backend-for-Frontend)** | Server SvelteKit đứng giữa trình duyệt và API Gateway; giữ token, gọi API, render SSR. |
| **SSR** | Server-Side Rendering — HTML được render ở server SvelteKit. |
| **Remote function** | Tính năng (experimental) của SvelteKit: `query` (đọc), `form` (submit form), `command` (hành động) — hàm chạy trên server, gọi từ component. File `*.remote.ts`. |
| **Runes** | Cú pháp reactivity của Svelte 5 (`$state`, `$derived`, `$effect`, `$props`). |
| **Paraglide** | Thư viện i18n biên dịch message thành hàm `m.key()`; locale `vi` (mặc định, không prefix URL) và `en` (`/en/...`). |
| **APISIX** | API Gateway; điểm vào duy nhất của API backend. |
| **Upstream** | Định nghĩa backend đích trong APISIX (vd upstream 2 = `product-service:8086`). |
| **openid-connect (plugin)** | Plugin APISIX xác thực Bearer JWT với Keycloak (`bearer_only`). |
| **Keycloak** | Identity Provider (OIDC); realm `ecommerce`, client `ecommerce-client`. |
| **OIDC / OAuth2 Authorization Code** | Luồng đăng nhập chuyển hướng tới Keycloak, nhận `code`, đổi lấy token. |
| **SSO ticket** | Mã ngẫu nhiên dùng một lần, sống 2 phút, do auth-service cấp sau callback để frontend đổi lấy token (`GET /api/v1/auth/session?ticket=`). |
| **State (OIDC)** | Mã chống CSRF trong luồng OIDC, lưu trong bộ nhớ auth-service 5 phút. |
| **Access token / Refresh token** | JWT do Keycloak cấp; BFF lưu trong cookie `access_token` / `refresh_token` (httpOnly). |
| **Correlation ID** | Header `X-Correlation-Id` sinh/truyền qua mọi service để truy vết request. |
| **Database-per-service** | Mỗi service sở hữu một database PostgreSQL riêng; không truy vấn chéo. |
| **InitSchema / initDb** | Hàm tạo bảng `CREATE TABLE IF NOT EXISTS` khi service khởi động (thay cho migration). |
| **ErrorVm** | Định dạng lỗi `{statusCode, title, detail}` dùng ở đa số service Go. |
| **ApiResponse** | Envelope `{success, code, message, data, errors, path, traceId, timestamp}` dùng ở auth-service. |
| **Page (Spring-style)** | Object phân trang `{content, totalElements, totalPages, size, number, first, last, empty…}`. |
| **CDC / Debezium** | Change Data Capture — search-service được viết để đọc message Debezium nhưng repo không có connector. |
| **Drizzle ORM** | ORM TypeScript dùng ở service Bun. |
| **bob / pgx** | Query builder và driver PostgreSQL dùng ở service Go. |
| **Hono** | Web framework nhẹ cho Bun. |
| **chi** | HTTP router cho Go. |
| **RustFS** | Object storage tương thích S3. |
| **Presigned URL** | URL tạm thời có chữ ký để truy cập object S3 (README có nhắc, code hiện tại không dùng — media-service stream file qua chính nó). |
