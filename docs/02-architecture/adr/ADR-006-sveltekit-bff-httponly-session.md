# ADR-006 — SvelteKit BFF, phiên bằng cookie httpOnly, không token trong JS trình duyệt

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

| Thuộc tính | Giá trị |
| :--- | :--- |
| Trạng thái | **Accepted (đang hiện diện trong code)** |
| Phạm vi | `frontend/apps/web`, `frontend/apps/admin`, `frontend/packages/lib/src/server` |
| Liên quan | [frontend-architecture](../software/frontend-architecture.md), [FEAT-01](../../04-features/FEAT-01-sso-login-logout.md), [FEAT-03](../../04-features/FEAT-03-session-refresh.md), [ADR-002](ADR-002-keycloak-sso-ticket-exchange.md) |

## Bối cảnh

Lưu access/refresh token trong `localStorage` hoặc biến JS khiến token bị lộ khi có XSS. Gateway nằm ở origin khác frontend, gọi trực tiếp từ trình duyệt cần CORS và quản lý token ở client.

## Quyết định

- Mỗi app SvelteKit là một **Backend-for-Frontend**: mọi dữ liệu lấy qua SSR `load` hoặc **remote functions** (`query`/`form`/`command`) chạy trên server; server gọi APISIX bằng `API_INTERNAL_URL`.
- Token Keycloak lưu trong cookie `access_token`, `refresh_token` (`httpOnly`, `sameSite=lax`, `path=/`), kèm `session_expires` và `session_user` (profile JSON).
- `hooks.server.ts` đọc cookie, **tự refresh** khi access token còn < 30 s hoặc đã mất, gán `event.locals.user`.
- Trình duyệt chỉ chạm gateway một lần: điều hướng tới `API_PUBLIC_URL/api/v1/auth/login` để bắt đầu SSO.

## Hệ quả tích cực

- JavaScript trình duyệt không đọc được token; giảm tác động XSS.
- Không cần CORS cho luồng dữ liệu; gateway có thể chỉ nghe nội bộ.
- Validation Zod dùng chung chạy trên server trước khi gọi backend; BFF tính lại tổng tiền khi checkout.

## Hệ quả tiêu cực

- Mọi request đi qua server SvelteKit (thêm một hop, server cần scale theo lưu lượng); rate limit theo IP ở gateway thấy một IP duy nhất (ISSUE-61).
- `secure` không đặt tường minh; hai app cùng tên cookie → chia sẻ phiên trên cùng host ([ISSUE-11](../../07-plan/known-issues-tech-debt.md)).
- `session_user` (chứa roles) **không ký** → người dùng có thể sửa cookie để qua guard admin (ISSUE-66).
- Phụ thuộc tính năng experimental `remoteFunctions` ([ISSUE-47](../../07-plan/known-issues-tech-debt.md)).
- Khi `/users/me` lỗi, cookie `session_user` bị xóa → đăng nhập "thành công" nhưng không có phiên (ISSUE-100).

## Bằng chứng trong code

- `frontend/packages/lib/src/server/session.ts:4-11,32-49` — tên cookie, thuộc tính, `setSession`.
- `frontend/packages/lib/src/server/env.ts:11-27` — `API_INTERNAL_URL`, `API_PUBLIC_URL`.
- `frontend/apps/web/src/hooks.server.ts:32-53` — refresh & `locals.user`.
- `frontend/apps/web/src/lib/server/api.ts:10-13` — client request-scoped lấy token từ cookie.
- `frontend/apps/web/vite.config.ts:13-15` — `experimental.remoteFunctions`.
- `frontend/apps/web/.env.example` — "The browser never calls it".

## Issue liên quan

[ISSUE-11](../../07-plan/known-issues-tech-debt.md), ISSUE-14, ISSUE-47, ISSUE-100, ISSUE-66, ISSUE-61.
