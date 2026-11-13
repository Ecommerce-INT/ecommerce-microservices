# ADR-003 — Xác thực tập trung tại APISIX, service tin gateway

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

| Thuộc tính | Giá trị |
| :--- | :--- |
| Trạng thái | **Accepted (đang hiện diện trong code)** |
| Phạm vi | APISIX (`deploy/apisix/apisix.yaml`), mọi backend service |
| Liên quan | [security-architecture](../system/security-architecture.md), [gateway-routing](../../05-api/gateway-routing.md), [shared-go-module](../../03-services/shared-go-module.md) |

## Bối cảnh

11 service, 2 ngôn ngữ. Việc nhúng thư viện OIDC/JWKS vào từng service tốn công và dễ lệch cấu hình. Mọi lưu lượng từ frontend đi qua một API gateway duy nhất (Apache APISIX 3.15).

## Quyết định

- APISIX xác minh JWT bằng plugin `openid-connect` với `bearer_only: true`, `use_jwks: true`, `token_signing_alg_values_expected: RS256`, discovery tới Keycloak realm `ecommerce` — gắn **theo route**.
- Route công khai (không plugin): `/api/v1/auth/*`, `/api/products*`, `/api/categories*`, `/storefront/ratings*`, `/storefront/catalog-search`, `/storefront/search_suggest`, `/<svc>/v3/api-docs/*`.
- Service **không** xác minh chữ ký; chỉ đọc claims không verify:
  - Go: middleware `UserClaims` (`ParseUnverified`) đặt `sub/email/preferred_username` vào context.
  - rating-service: `extractUser` base64-decode payload.
- Không có kiểm tra role/scope ở cả gateway lẫn service; phân quyền admin chỉ có ở frontend admin.
- Plugin global: `prometheus`, `cors`, `limit-count` (1200 req/60 s theo IP).

## Hệ quả tích cực

- Một điểm cấu hình xác thực; service đơn giản, không phụ thuộc Keycloak lúc runtime (trừ auth-service).
- Đổi IdP/thuật toán chỉ cần sửa gateway.

## Hệ quả tiêu cực

- **Không authorization theo role** → user thường có token gọi được API backoffice ([ISSUE-01](../../07-plan/known-issues-tech-debt.md)).
- Route `product-catalog` công khai cho **mọi method** → tạo/sửa/xóa sản phẩm, danh mục không cần token (ISSUE-102).
- Service tin header `Authorization` mà không verify; port service được publish trong compose → gọi trực tiếp giả mạo danh tính ([ISSUE-02](../../07-plan/known-issues-tech-debt.md), [ISSUE-03](../../07-plan/known-issues-tech-debt.md)).
- Sai sót cấu hình route trở thành lỗ hổng/404 (upstream 7/13 không tồn tại — [ISSUE-13](../../07-plan/known-issues-tech-debt.md); path gốc không khớp — [ISSUE-15](../../07-plan/known-issues-tech-debt.md)).
- Rate limit theo `remote_addr` trong khi mọi request đến từ BFF (ISSUE-61).

## Bằng chứng trong code

- `deploy/apisix/apisix.yaml:46-61` — global rules (prometheus, cors, limit-count).
- `deploy/apisix/apisix.yaml:63-102` — route công khai (auth-public, product-catalog, rating-storefront, search-storefront).
- `deploy/apisix/apisix.yaml:142-159` — mẫu route bảo vệ `openid-connect bearer_only`.
- `services/shared/middleware/middleware.go:44-69` — `UserClaims` parse không verify.
- `services/rating-service/src/index.ts:39-58` — `extractUser`.
- `services/product-service/main.go:89-90` — `r.Use(commonmw.CorrelationID)`, `r.Use(commonmw.UserClaims)` (không có middleware chặn).
- `frontend/apps/admin/src/hooks.server.ts:58-77` — guard role duy nhất trong hệ thống.

## Issue liên quan

[ISSUE-01](../../07-plan/known-issues-tech-debt.md), ISSUE-02, ISSUE-03, ISSUE-12, ISSUE-13, ISSUE-14, ISSUE-15, ISSUE-16, ISSUE-102, ISSUE-61.
