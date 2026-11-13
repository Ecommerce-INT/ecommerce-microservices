# ADR-001 — Backend đa ngôn ngữ: Go + Bun/TypeScript

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

| Thuộc tính | Giá trị |
| :--- | :--- |
| Trạng thái | **Accepted (đang hiện diện trong code)** |
| Phạm vi | 11 backend service |
| Liên quan | [go-service-architecture](../software/go-service-architecture.md), [bun-service-architecture](../software/bun-service-architecture.md) |

## Bối cảnh

Hệ thống có 11 service với mức độ phức tạp khác nhau. Các service "lõi giao dịch" (auth, product, order, payment, inventory, shipping, search) cần kết nối Postgres/Kafka/OpenSearch ổn định; các service "phụ trợ" (promotion/tax, rating, media, notification) chủ yếu là CRUD, upload file, gửi email. Dấu vết trong code (biến `SPRING_DATASOURCE_URL`, path `/actuator/health`, `ErrorVm`, route `/v3/api-docs`) cho thấy hệ thống được viết lại từ một bản Java/Spring trước đó.

## Quyết định

- 7 service viết bằng **Go 1.27** với `chi/v5`, `bob` + `pgx/v5`, quản lý chung bằng `go.work` và module `com.ecommerce/shared`.
- 4 service viết bằng **TypeScript chạy trên Bun 1.4** với Hono, Drizzle ORM (`bun-sql`), Zod, dùng API tích hợp sẵn của Bun (`Bun.SQL`, `Bun.S3Client`).
- Hai nhóm giữ **cùng hợp đồng bề ngoài**: route kép prefix + root, `/actuator/health`, `X-Correlation-Id`, logic giải quyết DB URL giống nhau, CORS `*`.

## Hệ quả tích cực

- Go cho binary tĩnh nhỏ (`CGO_ENABLED=0`, alpine) và mô hình đồng thời tốt cho consumer Kafka/OpenSearch.
- Bun giảm boilerplate cho CRUD và cung cấp sẵn S3/SQL client, không cần thêm dependency.
- Frontend (SvelteKit trên Bun) và service Bun dùng chung ngôn ngữ/công cụ (Zod, Bun test).

## Hệ quả tiêu cực

- **Không chia sẻ code xuyên ngôn ngữ**: middleware correlation, giải quyết DSN, định dạng health được viết hai lần; định dạng lỗi khác nhau giữa hai nhóm ([ISSUE-29](../../07-plan/known-issues-tech-debt.md)).
- Service Bun không dùng module chung nào (mỗi service một `package.json`/lockfile, lặp `db.ts`), service Go lặp `buildPostgresDSN` và `PageResponse`.
- Hai toolchain build/test/CI riêng (`go-check` và `bun-check` trong CI; `make go-*` và `make bun-*`).
- Kỹ năng đội cần trải trên 2 hệ sinh thái; quy ước đặt tên/phân trang lệch nhau (`page/size` vs `pageNo/pageSize`).

## Bằng chứng trong code

- `go.work:1-12` — 8 module Go.
- `Makefile:15-19` — `GO_MODULES`, `BUN_SERVICES`.
- `services/product-service/go.mod` — chi, bob, pgx, `replace com.ecommerce/shared => ../shared`.
- `services/rating-service/package.json:11-19` — hono, drizzle-orm, zod.
- `services/media-service/src/storage.ts:1,9-16` — `Bun.S3Client`; `services/rating-service/src/db.ts:1,36-44` — `Bun.SQL` + drizzle.
- `services/product-service/main.go:26-54` vs `services/rating-service/src/db.ts:4-32` — cùng logic DSN viết hai lần.
- `services/shared/middleware/middleware.go:25-36` vs `services/rating-service/src/index.ts:21-25` — correlation id hai bản.

## Issue liên quan

[ISSUE-29](../../07-plan/known-issues-tech-debt.md), [ISSUE-37](../../07-plan/known-issues-tech-debt.md), [ISSUE-46](../../07-plan/known-issues-tech-debt.md).
