# Kiến trúc — Mục lục theo arc42 & C4

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> Điểm vào của thư mục `02-architecture/`. Ánh xạ toàn bộ tài liệu kiến trúc vào **12 mục của arc42** và **4 mức của mô hình C4**. Mọi tài liệu mô tả **code hiện tại** (as-is); kiến trúc đích tương lai được lưu trữ riêng tại [`target-revamp/`](target-revamp/README.md). Tài liệu hạ tầng (k8s, CI/CD, observability stack) sẽ được bổ sung ở bộ riêng.

## 1. Cấu trúc thư mục

```
02-architecture/
├── README.md                     ← tài liệu này (arc42 index + §4 Solution strategy)
├── quality-scenarios.md          arc42 §10 — cây chất lượng & kịch bản
├── system/                       Kiến trúc hệ thống
│   ├── system-context.md         C4 L1
│   ├── container-view.md         C4 L2
│   ├── component-view.md         C4 L3 (11 service + 2 app SvelteKit)
│   ├── runtime-flows.md          Sequence diagram RF-01 … RF-09
│   ├── runtime-view-local.md     docker compose (góc nhìn ứng dụng)
│   ├── integration-messaging.md  HTTP & Kafka
│   ├── data-architecture.md      Database-per-service, sở hữu dữ liệu
│   ├── security-architecture.md  Định danh, xác thực, phân quyền, phiên
│   └── threat-model.md           STRIDE, DFD, ranh giới tin cậy
├── software/                     Kiến trúc phần mềm (C4 L3–L4)
│   ├── go-service-architecture.md
│   ├── bun-service-architecture.md
│   ├── frontend-architecture.md
│   ├── api-conventions.md
│   └── cross-cutting-concerns.md
├── adr/                          ADR-001 … ADR-008
└── target-revamp/                Đề xuất kiến trúc đích dài hạn (Better-Auth, NATS JetStream, PG18)
```

## 2. Ánh xạ arc42

| § | Mục arc42 | Tài liệu | Ghi chú |
| :---: | :--- | :--- | :--- |
| 1 | Introduction & goals | [01-brd/brd.md](../01-brd/brd.md), [01-brd/non-functional-requirements.md](../01-brd/non-functional-requirements.md), [00-overview/project-overview.md](../00-overview/project-overview.md), [01-brd/actors-and-permissions.md](../01-brd/actors-and-permissions.md) | Yêu cầu nghiệp vụ BR-01…BR-28, stakeholder, mục tiêu chất lượng (tóm tắt ở [quality-scenarios.md §1](quality-scenarios.md#1-mục-tiêu-chất-lượng-ưu-tiên)) |
| 2 | Architecture constraints | [§3 dưới đây](#3-ràng-buộc-kiến-trúc-arc42-2), [ADR-001](adr/ADR-001-polyglot-go-bun.md) | Ràng buộc quan sát từ code & cấu hình |
| 3 | Context & scope | [system/system-context.md](system/system-context.md) | C4 L1: actor, Keycloak, SMTP, RustFS |
| 4 | Solution strategy | [§4 dưới đây](#4-solution-strategy-arc42-4) | Viết tại README này |
| 5 | Building block view | L1: [system/container-view.md](system/container-view.md) · L2: [system/component-view.md](system/component-view.md) · L3: [software/go-service-architecture.md](software/go-service-architecture.md), [software/bun-service-architecture.md](software/bun-service-architecture.md), [software/frontend-architecture.md](software/frontend-architecture.md) · chi tiết từng khối: [03-services/](../03-services/) | Whitebox container → component → package/file |
| 6 | Runtime view | [system/runtime-flows.md](system/runtime-flows.md), [system/integration-messaging.md](system/integration-messaging.md) | SSO, refresh phiên, catalog, checkout, CRUD admin, sự kiện, upload, search |
| 7 | Deployment view | [system/runtime-view-local.md](system/runtime-view-local.md) | Chỉ docker compose; k8s/k3d/GitOps thuộc bộ tài liệu hạ tầng (sẽ bổ sung) |
| 8 | Cross-cutting concepts | [software/cross-cutting-concerns.md](software/cross-cutting-concerns.md), [system/security-architecture.md](system/security-architecture.md), [system/data-architecture.md](system/data-architecture.md), [software/api-conventions.md](software/api-conventions.md), [06-data/data-privacy-pii.md](../06-data/data-privacy-pii.md) | Cấu hình, logging, lỗi, CORS, validation, i18n, bảo mật, dữ liệu cá nhân |
| 9 | Architecture decisions | [adr/README.md](adr/README.md) | ADR-001 … ADR-008 (suy ngược từ code) |
| 10 | Quality requirements | [quality-scenarios.md](quality-scenarios.md), [01-brd/non-functional-requirements.md](../01-brd/non-functional-requirements.md) | 44 kịch bản, 7 đạt |
| 11 | Risks & technical debt | [07-plan/known-issues-tech-debt.md](../07-plan/known-issues-tech-debt.md), [system/threat-model.md](system/threat-model.md) | ISSUE-xx; 38 threat STRIDE |
| 12 | Glossary | [00-overview/glossary.md](../00-overview/glossary.md) | |

## 3. Ràng buộc kiến trúc (arc42 §2)

Các ràng buộc dưới đây **quan sát được** trong repo (không có tài liệu ràng buộc gốc).

| Loại | Ràng buộc | Bằng chứng |
| :--- | :--- | :--- |
| Kỹ thuật | Backend Go (chi, bob/pgx) cho service lõi; Bun + Hono + Drizzle cho service phụ trợ | `services/*/go.mod`, `services/*/package.json`; [ADR-001](adr/ADR-001-polyglot-go-bun.md) |
| Kỹ thuật | Frontend SvelteKit (adapter-bun) dùng tính năng experimental `remoteFunctions` | `frontend/apps/web/vite.config.ts:12-18` (ISSUE-47) |
| Kỹ thuật | Một PostgreSQL 16 instance, một database mỗi service, một user DB chung | `deploy/postgres/init/create-all-databases.sql`; [ADR-004](adr/ADR-004-database-per-service-runtime-ddl.md) |
| Kỹ thuật | IdP là Keycloak 26 (OIDC); gateway là Apache APISIX 3.15 ở chế độ standalone YAML | `docker-compose.yml:122-170`; `deploy/apisix/config.yaml:5-8` |
| Kỹ thuật | Messaging Kafka 3.9 KRaft 1 node, PLAINTEXT | `docker-compose.yml:41-60` |
| Kỹ thuật | Tương thích di sản Spring: path `/actuator/health`, `SPRING_DATASOURCE_URL`, route `/v3/api-docs`, `PageResponse` kiểu Spring | `services/*/main.go` `buildPostgresDSN`; `deploy/apisix/apisix.yaml:103-141` |
| Tổ chức | Môi trường chạy tham chiếu cho ứng dụng là `docker compose`; production chạy k8s qua APISIX Ingress Controller (ngoài phạm vi) | `deploy/apisix/config.yaml:1-4` |
| Tổ chức | Monorepo: `services/` (go.work) + `frontend/` (Bun workspaces) | `go.work`, `frontend/package.json` |
| Quy ước | Ngôn ngữ UI vi (mặc định) / en; tài liệu tiếng Việt | Paraglide; [conventions-and-ids.md](../00-overview/conventions-and-ids.md) |

## 4. Solution strategy (arc42 §4)

Chiến lược giải pháp **như đang hiện thực** — tóm tắt các quyết định nền tảng và cách chúng hiện thực các mục tiêu chất lượng.

### 4.1 Quyết định công nghệ

| Quyết định | Lý do (suy từ code) | Hệ quả hiện tại |
| :--- | :--- | :--- |
| Microservice theo domain (11 service), mỗi service 1 DB | Tách quyền sở hữu dữ liệu, deploy độc lập | Không join/transaction xuyên domain; tồn kho và dòng hàng bị phân mảnh (ISSUE-26) |
| Polyglot Go + Bun | Go cho service lõi hiệu năng; Bun/TS cho service CRUD nhỏ | Hai bộ quy ước lỗi/health/validation (ISSUE-29) — [ADR-001](adr/ADR-001-polyglot-go-bun.md) |
| SvelteKit SSR làm **BFF** | Token không ra trình duyệt; SEO cho storefront | BFF là nơi duy nhất điều phối nghiệp vụ nhiều service — [ADR-006](adr/ADR-006-sveltekit-bff-httponly-session.md) |
| Keycloak + auth-service trung gian SSO bằng ticket | Không lộ code/token qua URL frontend | Thêm một bước trung gian có state in-memory (ISSUE-06) — [ADR-002](adr/ADR-002-keycloak-sso-ticket-exchange.md) |
| Xác thực tập trung tại APISIX | Service đơn giản, không cần thư viện OIDC | Gateway là lớp bảo vệ duy nhất; không có phân quyền (ISSUE-01, ISSUE-02) — [ADR-003](adr/ADR-003-gateway-authentication.md) |
| Schema tạo lúc khởi động | Không cần công cụ migration | Không đổi schema an toàn (ISSUE-28) — [ADR-004](adr/ADR-004-database-per-service-runtime-ddl.md) |
| Kafka cho sự kiện; OpenSearch cho tìm kiếm | Tách đồng bộ; full-text | Chưa có luồng nào chạy đầu-cuối (ISSUE-17, ISSUE-19) — [ADR-005](adr/ADR-005-kafka-async-events.md), [ADR-007](adr/ADR-007-opensearch-catalog-search.md) |
| Giỏ hàng client-side (`localStorage`) | Không cần backend cho khách vãng lai | Giỏ không đồng bộ giữa thiết bị — [ADR-008](adr/ADR-008-client-side-cart-localstorage.md) |

### 4.2 Phân rã & tổ chức code

- **Backend Go**: phân lớp `handler → service → repository` với DI thủ công trong `main.go`, module `services/shared` cho config/DB/middleware/response ([go-service-architecture.md](software/go-service-architecture.md)).
- **Backend Bun**: `index.ts` (Hono routes) → `service.ts` → `db.ts`/`schema.ts` (Drizzle), validate bằng Zod ([bun-service-architecture.md](software/bun-service-architecture.md)).
- **Frontend**: monorepo `apps/web`, `apps/admin`, `packages/lib` (session, API client, schema, cart), `packages/ui` (shadcn-svelte); dữ liệu nạp qua remote function chạy trên server ([frontend-architecture.md](software/frontend-architecture.md)).
- **Tích hợp**: không có HTTP service-to-service; BFF điều phối đồng bộ, Kafka cho bất đồng bộ ([integration-messaging.md](system/integration-messaging.md)).

### 4.3 Cách đạt mục tiêu chất lượng (và khoảng trống)

| Mục tiêu | Cách tiếp cận trong code | Khoảng trống chính |
| :--- | :--- | :--- |
| Security | Cookie `httpOnly`; xác thực JWT tại gateway; guard admin ở BFF; Zod ở BFF/Bun | Không RBAC backend, cookie `session_user` không ký, catalog ghi public, service không verify JWT — [threat-model.md](system/threat-model.md) |
| Reliability | Graceful shutdown, `Recoverer`, logout best-effort | Không transaction/outbox/idempotency; health luôn UP |
| Maintainability | Cấu trúc đồng nhất, shared module, lint/format | Không migration, nhiều định dạng response, không contract test |
| Performance | SSR, pool DB, service stateless | Endpoint list không phân trang; payment publish đồng bộ |
| Operability | Correlation id từng service, metrics gateway | Không truyền correlation id xuyên tầng, không metrics service |
| Usability | i18n vi/en, responsive, trang lỗi | Admin không đăng nhập được (ISSUE-100), không đăng ký (ISSUE-71) |

Chi tiết đo lường: [quality-scenarios.md](quality-scenarios.md).

## 5. Ánh xạ mô hình C4

| Mức C4 | Câu hỏi trả lời | Tài liệu | Phần tử chính |
| :--- | :--- | :--- | :--- |
| **L1 — System Context** | Hệ thống phục vụ ai, phụ thuộc hệ thống ngoài nào? | [system/system-context.md](system/system-context.md) | Khách hàng, Quản trị viên, Dev/API client; Keycloak, SMTP, RustFS |
| **L2 — Container** | Gồm những tiến trình/kho dữ liệu nào, giao tiếp ra sao? | [system/container-view.md](system/container-view.md), [system/runtime-view-local.md](system/runtime-view-local.md) | web, admin, APISIX, 11 service, PostgreSQL, Kafka, OpenSearch, RustFS, Redis, Keycloak |
| **L3 — Component** | Bên trong mỗi container có những component nào? | [system/component-view.md](system/component-view.md); sơ đồ component frontend tại [software/frontend-architecture.md §11](software/frontend-architecture.md#11-component-diagram) | Handler, service, repository, consumer/producer, client, middleware, store, remote function, hooks |
| **L4 — Code** | Package/file/class cụ thể | [software/go-service-architecture.md](software/go-service-architecture.md) (package diagram), [software/bun-service-architecture.md](software/bun-service-architecture.md), [03-services/*.md](../03-services/) §5 "Cấu trúc code", [06-data/data-dictionary.md](../06-data/data-dictionary.md) | File `.go`/`.ts`, bảng DB |
| Bổ trợ — Dynamic | Luồng chạy qua nhiều container | [system/runtime-flows.md](system/runtime-flows.md) | RF-01 … RF-09 |
| Bổ trợ — Deployment | Container chạy ở đâu | [system/runtime-view-local.md](system/runtime-view-local.md) (+ tài liệu hạ tầng sau) | docker compose |

## 6. Đọc theo vai trò

| Vai trò | Lộ trình đọc |
| :--- | :--- |
| Kiến trúc sư mới | System context → Container → Component → §4 Solution strategy → ADR → Quality scenarios → Threat model |
| Security reviewer | Security architecture → Threat model → Data privacy PII → Known issues §1 |
| Developer backend | Component view (service của mình) → Go/Bun architecture → API conventions → Cross-cutting |
| Developer frontend | Component view §12–13 → Frontend architecture → Runtime flows RF-01…RF-05 |
