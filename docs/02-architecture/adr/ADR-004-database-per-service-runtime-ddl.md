# ADR-004 — Database-per-service, schema tạo lúc chạy bằng `CREATE TABLE IF NOT EXISTS`

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

| Thuộc tính | Giá trị |
| :--- | :--- |
| Trạng thái | **Accepted (đang hiện diện trong code)** |
| Phạm vi | 10 service có Postgres (trừ search-service) |
| Liên quan | [data-dictionary](../../06-data/data-dictionary.md), [go-service-architecture §4.3](../software/go-service-architecture.md#43-khởi-tạo-schema), [bun-service-architecture §5.3](../software/bun-service-architecture.md#53-initdb) |

## Bối cảnh

Microservice cần độc lập dữ liệu. Dự án muốn `docker compose up` chạy được ngay mà không cần bước migrate riêng.

## Quyết định

- **Một Postgres 16 instance**, mỗi service một **database logic** riêng (`authservice`, `productservice`, `orderservice`, `paymentservice`, `inventoryservice`, `shippingservice`, `ratingservice`, `mediaservice`, `promotionservice`), tạo bằng script init của container Postgres. Ngoại lệ: notification-service dùng DB chung `ecommerce`.
- Không service nào truy cập DB của service khác; không có FK xuyên service (ví dụ `orders.product_id` chỉ là số).
- Schema được tạo **lúc khởi động service** bằng DDL viết tay `CREATE TABLE IF NOT EXISTS` (`InitSchema` ở Go, `initDb()` ở Bun); lỗi chỉ log/bỏ qua. Không có công cụ migration, không bảng version. Drizzle `schema.ts` chỉ dùng cho truy vấn, không sinh migration.

## Hệ quả tích cực

- Khởi động "zero-step": stack chạy được ngay sau `docker compose up`.
- Ranh giới dữ liệu rõ ràng giữa service; thay đổi schema một service không ảnh hưởng service khác.

## Hệ quả tiêu cực

- **Không thể tiến hóa schema an toàn**: `IF NOT EXISTS` không thêm cột/đổi kiểu cho bảng đã có ([ISSUE-28](../../07-plan/known-issues-tech-debt.md)).
- DDL lặp hai nơi ở Bun (`db.ts` và `schema.ts`) dễ lệch.
- Lỗi DDL bị che (Go `_ = repo.InitSchema(ctx)`); DDL chỉ chạy một lần lúc start.
- DB thừa/không dùng: `favouriteservice`, `taxservice` ([ISSUE-31](../../07-plan/known-issues-tech-debt.md)); notification dùng DB chung ([ISSUE-33](../../07-plan/known-issues-tech-debt.md)).
- Không có nhất quán xuyên service (không saga/outbox): tồn kho hai nguồn ([ISSUE-26](../../07-plan/known-issues-tech-debt.md)), order/payment rời rạc ([ISSUE-23](../../07-plan/known-issues-tech-debt.md)).
- Chung một instance Postgres → chung điểm lỗi và tài nguyên.

## Bằng chứng trong code

- `deploy/postgres/init/create-all-databases.sql` — `CREATE DATABASE` cho từng service (+ `keycloak`).
- `docker-compose.yml:186,207,226,248,270,289,308,327,352` — `SPRING_DATASOURCE_URL` trỏ DB riêng; `:393-395` — notification dùng `ecommerce`.
- `services/product-service/internal/repository/product_repo.go:30-68` — `InitSchema` (log notice).
- `services/order-service/main.go:72` — `_ = repo.InitSchema(ctx)`.
- `services/rating-service/src/db.ts:46-67` — `initDb()`; `services/rating-service/src/schema.ts:4-16` — định nghĩa Drizzle trùng lặp.

## Issue liên quan

[ISSUE-26](../../07-plan/known-issues-tech-debt.md), ISSUE-28, ISSUE-31, ISSUE-33, ISSUE-36.
