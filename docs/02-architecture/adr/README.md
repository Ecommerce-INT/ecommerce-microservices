# Architecture Decision Records (ADR)

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> ADR trong thư mục này ghi lại các **quyết định kiến trúc đang hiện diện trong code** (được suy ra ngược từ source — không tìm thấy ADR gốc trong repo). Chúng **không** phải đề xuất tương lai; kiến trúc đích hay các đề xuất chưa triển khai không được phản ánh ở đây.

## Định dạng

Mỗi ADR gồm: **Trạng thái** · **Bối cảnh** · **Quyết định** · **Hệ quả tích cực** · **Hệ quả tiêu cực** · **Bằng chứng trong code** (`path:line`) · **Issue liên quan**.

Trạng thái dùng trong bộ này:

| Trạng thái | Ý nghĩa |
| :--- | :--- |
| Accepted (đang hiện diện trong code) | Quyết định đang được code thực thi |
| Superseded | Bị thay bằng ADR khác (chưa có) |
| Deprecated | Còn trong code nhưng có kế hoạch gỡ (chưa có) |

## Danh sách

| ID | Tiêu đề | Trạng thái | Phạm vi |
| :--- | :--- | :--- | :--- |
| [ADR-001](ADR-001-polyglot-go-bun.md) | Backend đa ngôn ngữ: Go + Bun/TypeScript | Accepted | Toàn bộ backend |
| [ADR-002](ADR-002-keycloak-sso-ticket-exchange.md) | Keycloak làm IdP, auth-service trung gian SSO bằng ticket | Accepted | auth-service, web, admin |
| [ADR-003](ADR-003-gateway-authentication.md) | Xác thực tập trung tại APISIX (`openid-connect bearer_only`), service tin gateway | Accepted | APISIX, mọi service |
| [ADR-004](ADR-004-database-per-service-runtime-ddl.md) | Database-per-service, schema tạo lúc chạy bằng `CREATE TABLE IF NOT EXISTS` | Accepted | 10 service có DB |
| [ADR-005](ADR-005-kafka-async-events.md) | Kafka cho sự kiện bất đồng bộ (payment → notification, product CDC → search) | Accepted | payment, notification, search |
| [ADR-006](ADR-006-sveltekit-bff-httponly-session.md) | SvelteKit BFF, phiên bằng cookie httpOnly, không token trong JS trình duyệt | Accepted | web, admin |
| [ADR-007](ADR-007-opensearch-catalog-search.md) | OpenSearch cho tìm kiếm catalog | Accepted | search-service |
| [ADR-008](ADR-008-client-side-cart-localstorage.md) | Giỏ hàng phía client lưu localStorage | Accepted | web |

## Thêm ADR mới

1. Đánh số tiếp theo (`ADR-009`…), cập nhật bảng trên và sổ đăng ký trong [`conventions-and-ids.md`](../../00-overview/conventions-and-ids.md#5-sổ-đăng-ký-adr).
2. Chỉ ghi `Accepted (đang hiện diện trong code)` khi code đã thực thi quyết định; đề xuất chưa làm thuộc `07-plan/`.
3. Mọi khẳng định phải có `path:line`.
