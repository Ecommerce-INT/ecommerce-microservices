# Yêu cầu phi chức năng (NFR)

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> Liệt kê các thuộc tính chất lượng **đang được hiện thực trong code ứng dụng** và mức độ đáp ứng. Không có SLO/SLA được đo lường trong repo; các con số dưới đây là **cấu hình thực tế**, không phải kết quả benchmark.

## 1. Bảo mật

| ID | Yêu cầu | Hiện thực | Trạng thái |
| :--- | :--- | :--- | :--- |
| NFR-SEC-01 | Token không bị JavaScript phía trình duyệt truy cập | BFF lưu `access_token`, `refresh_token`, `session_user`, `session_expires` trong cookie `httpOnly`, `SameSite=Lax`, `path=/` (`frontend/packages/lib/src/server/session.ts`) | 🟡 Không lộ token cho JS; nhưng thiếu `secure` (ISSUE-11) và cookie `session_user` không ký, bị dùng để quyết định quyền (ISSUE-66) |
| NFR-SEC-02 | Mọi API nghiệp vụ cần xác thực | APISIX `openid-connect` `bearer_only` cho đa số route; catalog ghi và POST rating để public | 🟡 ISSUE-102, ISSUE-03 |
| NFR-SEC-03 | Phân quyền theo vai trò | Chỉ ở admin UI | ⛔ ISSUE-01, ISSUE-100 |
| NFR-SEC-04 | Chống open redirect sau đăng nhập | `safeRedirect()` ở BFF chỉ cho path nội bộ; auth-service **không** kiểm `redirect_uri` | 🟡 ISSUE-05 |
| NFR-SEC-05 | Chống CSRF trong OIDC | `state` ngẫu nhiên 32 byte, dùng một lần, TTL 5 phút | 🟡 Callback vẫn chấp nhận state thiếu/hết hạn (ISSUE-52); in-memory (ISSUE-06); không PKCE (ISSUE-133) |
| NFR-SEC-06 | Ticket SSO dùng một lần, thời hạn ngắn | TTL 2 phút, xóa sau khi đổi | ✅ |
| NFR-SEC-07 | Giới hạn tần suất request | APISIX `limit-count`: 1200 request / 60 giây / IP, trả 429 (đếm in-memory mỗi node gateway) | 🟡 Key là IP của BFF → thực chất là trần chung cho mọi người dùng (ISSUE-61) |
| NFR-SEC-08 | Giá không tin client | `placeOrder` lấy lại giá từ product-service | ✅ (nhưng order API tin `orderFee` — ISSUE-21) |
| NFR-SEC-09 | Không lưu dữ liệu thẻ | Không có xử lý thẻ trong code | ✅ |
| NFR-SEC-10 | Crawler không index trang riêng tư | `robots.txt` chặn `/cart`, `/checkout`, `/orders`, `/auth/` (web); `Disallow: /` (admin) | ✅ |

## 2. Hiệu năng & khả năng mở rộng

| ID | Yêu cầu | Hiện thực | Trạng thái |
| :--- | :--- | :--- | :--- |
| NFR-PERF-01 | Render trang nhanh, SEO | SSR SvelteKit; dữ liệu nạp qua remote `query` phía server | ✅ |
| NFR-PERF-02 | Kết nối DB được pool | Go: pgx pool max 25 / min 2, lifetime 1h, idle 30m (`services/shared/database/postgres.go`). Bun: max 20, idle 30s, connect timeout 5s, tắt prepared statement (`src/db.ts`) | ✅ |
| NFR-PERF-03 | Danh sách lớn được phân trang | Có ở `/users/all`, `/carts/all`, `/orders/all`, `/payments/all`, promotion, tax, rating; **không** có ở `/api/products`, `/api/categories`, `/api/orders` mà frontend dùng | 🟡 ISSUE-48 |
| NFR-PERF-04 | Scale ngang service | Service stateless trừ auth-service (SSO store in-memory) | 🟡 ISSUE-06 |
| NFR-PERF-05 | Cache | Redis có trong stack nhưng không được dùng | ⛔ |

## 3. Độ tin cậy

| ID | Yêu cầu | Hiện thực | Trạng thái |
| :--- | :--- | :--- | :--- |
| NFR-REL-01 | Tắt service an toàn | Go: graceful shutdown 10 giây | ✅ |
| NFR-REL-02 | Không crash khi panic | chi `middleware.Recoverer` | ✅ |
| NFR-REL-03 | Health check phản ánh trạng thái thật | Health luôn UP, không kiểm phụ thuộc | ⛔ ISSUE-36 |
| NFR-REL-04 | Không mất dữ liệu khi lỗi giữa chừng | Không transaction xuyên service, không outbox, không idempotency key | ⛔ ISSUE-23, ISSUE-24 |
| NFR-REL-05 | Fail-fast khi thiếu DB | Service vẫn khởi động và trả dữ liệu rỗng/stub | ⛔ ISSUE-28 |
| NFR-REL-06 | Đăng xuất best-effort | Lỗi gọi logout vẫn xóa cookie | ✅ |
| NFR-REL-07 | Hết phiên xử lý mềm | Refresh lỗi → xóa cookie, coi là guest, không redirect vòng lặp | ✅ |

## 4. Khả năng quan sát

| ID | Yêu cầu | Hiện thực | Trạng thái |
| :--- | :--- | :--- | :--- |
| NFR-OBS-01 | Truy vết request | Header `X-Correlation-Id` sinh/echo ở mọi service (Go + Bun) | 🟡 Không xuyên suốt: gateway không sinh, BFF không gửi, service không log, không vào Kafka (ISSUE-98) |
| NFR-OBS-02 | Log request | chi `middleware.Logger`; Hono `logger()` | ✅ |
| NFR-OBS-03 | Metrics | APISIX plugin `prometheus` (port 9091); service không tự xuất metrics | 🟡 |

## 5. Khả dụng & trải nghiệm

| ID | Yêu cầu | Hiện thực | Trạng thái |
| :--- | :--- | :--- | :--- |
| NFR-UX-01 | Đa ngôn ngữ | Paraglide vi (mặc định) / en; kiểm tra parity key | ✅ |
| NFR-UX-02 | Responsive | Tailwind 4; admin có menu slide-out mobile | ✅ |
| NFR-UX-03 | Thông báo lỗi thân thiện | `+error.svelte` (404 vs lỗi chung); toast `svelte-sonner`; không dùng `alert()` | ✅ |
| NFR-UX-04 | Accessibility | Dùng bits-ui (có ARIA); chưa audit Lighthouse | 🟡 |

## 6. Khả năng bảo trì

| ID | Yêu cầu | Hiện thực | Trạng thái |
| :--- | :--- | :--- | :--- |
| NFR-MNT-01 | Cấu trúc code nhất quán | Go: handler/service/repository; Bun: index/service/db/schema | ✅ |
| NFR-MNT-02 | Dùng chung code | `services/shared` (Go), `frontend/packages/lib` & `ui` | ✅ |
| NFR-MNT-03 | Quản lý schema DB | Không migration | ⛔ ISSUE-28 |
| NFR-MNT-04 | Định dạng API thống nhất | 3+ định dạng | ⛔ ISSUE-29 |
| NFR-MNT-05 | Kiểm thử tự động | Unit test router & lib; không integration/E2E | 🟡 ISSUE-46 |
| NFR-MNT-06 | Kiểm tra định dạng/lint | `gofmt`, `go vet`, Prettier, ESLint, `svelte-check`, locale parity — chạy qua `make check`, `make web-check` | ✅ |

## 7. Tài nguyên runtime (cấu hình compose)

| Thành phần | Image | Ghi chú |
| :--- | :--- | :--- |
| Service Go | build từ `services/<svc>/Dockerfile` (context `./services`) | Static binary |
| Service Bun | `oven/bun:1.4-alpine` | `bun run src/index.ts` |
| Frontend | `ghcr.io/ecommerce-int/frontend-{web,admin}:latest` | adapter-bun |

Ngân sách tài nguyên, autoscaling, HA sẽ được mô tả trong tài liệu hạ tầng.
