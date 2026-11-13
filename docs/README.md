# Tài liệu ứng dụng — E-Commerce Microservices Platform

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> **Phạm vi:** toàn bộ phần **ứng dụng** (11 backend service + 2 frontend SvelteKit + package dùng chung). Tài liệu hạ tầng (Kubernetes, CI/CD, GitOps, observability) sẽ được bổ sung ở bộ riêng.
>
> **Nguyên tắc cốt lõi:** mọi nội dung mô tả **đúng code hiện tại (as-is)**. Chức năng còn thiếu hoặc có lỗi được ghi trung thực kèm mã `ISSUE-xx`. Các tài liệu cũ trong `docs/` (`ARCHITECTURE.md`, `architecture-revamp-plan.md`) mô tả **kiến trúc đích khác** và **không** phải nguồn của bộ tài liệu này.
>
> **Cập nhật:** 2026-11-12 · **Ngôn ngữ:** Tiếng Việt · **Sơ đồ:** Mermaid

---

## 1. Bắt đầu từ đâu?

| Bạn là… | Đọc theo thứ tự |
| :--- | :--- |
| **Người mới vào dự án** | [Tổng quan dự án](00-overview/project-overview.md) → [Container view](02-architecture/system/container-view.md) → [Developer guide](08-guides/developer-guide.md) |
| **BA / PO** | [BRD](01-brd/brd.md) → [Quy trình nghiệp vụ](01-brd/business-processes.md) → [Danh mục tính năng](04-features/README.md) → [Plan](07-plan/project-plan.md) |
| **Kiến trúc sư / Tech lead** | [Kiến trúc hệ thống](#4-kiến-trúc) → [Kiến trúc phần mềm](#4-kiến-trúc) → [ADR](02-architecture/adr/README.md) → [Known issues](07-plan/known-issues-tech-debt.md) |
| **Backend dev** | [Service docs](#5-service) → [API catalog](05-api/api-catalog.md) → [Data dictionary](06-data/data-dictionary.md) → [Coding conventions](08-guides/coding-conventions.md) |
| **Frontend dev** | [Frontend architecture](02-architecture/software/frontend-architecture.md) → [frontend-web](03-services/frontend-web.md) / [frontend-admin](03-services/frontend-admin.md) / [packages](03-services/frontend-packages.md) |
| **QA** | [Danh mục tính năng](04-features/README.md) (Acceptance criteria trong từng FEAT) → [Test plan](07-plan/test-plan.md) |

## 2. Cấu trúc thư mục

```
docs/
├── README.md                 ← Cổng tra cứu tài liệu chính (tài liệu này)
├── ARCHITECTURE.md           ← Điều hướng kiến trúc: Hiện tại (As-Is) vs Đích (Target)
├── 00-overview/              Tổng quan, thuật ngữ, quy ước & sổ đăng ký ID
├── 01-brd/                   Business Requirements Document & Quy trình nghiệp vụ
├── 02-architecture/          Kiến trúc hệ thống & phần mềm (arc42, C4, ADRs)
│   ├── system/               System Context, Container, Component, Runtime, Security, Data
│   ├── software/             Kiến trúc Go, Bun/Hono, Frontend BFF, API conventions
│   ├── adr/                  Architecture Decision Records (ADR-001 … ADR-008)
│   └── target-revamp/        Kiến trúc đích tương lai (ARCHITECTURE.md, revamp plan)
├── 03-services/              Mỗi service / app một tài liệu chi tiết (11 backend + 2 frontend)
├── 04-features/              Mỗi tính năng một tài liệu (FEAT-01 … FEAT-28 + screen catalog)
├── 05-api/                   API catalog, Gateway routing, OpenAPI specs, AsyncAPI spec
├── 06-data/                  Data dictionary toàn hệ thống & Data privacy PII
├── 07-plan/                  Known issues (107 issues), Project plan (P0-P5), Backlog, Test plan
├── 08-guides/                Developer guide, Coding conventions, Rollout guides
└── web/                      Tài liệu & tasklist quá trình rewrite SvelteKit frontend
```

## 3. Tổng quan & BRD

| Tài liệu | Nội dung |
| :--- | :--- |
| [00-overview/project-overview.md](00-overview/project-overview.md) | Hệ thống làm gì, tech stack, danh sách service/app, port, trạng thái tổng thể |
| [00-overview/glossary.md](00-overview/glossary.md) | Thuật ngữ nghiệp vụ & kỹ thuật |
| [00-overview/conventions-and-ids.md](00-overview/conventions-and-ids.md) | Quy ước viết, trạng thái, sổ đăng ký BR/FEAT/ISSUE/ADR, template |
| [01-brd/brd.md](01-brd/brd.md) | BRD: bối cảnh, mục tiêu, stakeholder, phạm vi, yêu cầu nghiệp vụ BR-01…BR-28 |
| [01-brd/actors-and-permissions.md](01-brd/actors-and-permissions.md) | Actor & ma trận quyền (thiết kế vs thực thi thực tế) |
| [01-brd/business-processes.md](01-brd/business-processes.md) | Quy trình nghiệp vụ end-to-end |
| [01-brd/non-functional-requirements.md](01-brd/non-functional-requirements.md) | Yêu cầu phi chức năng đang được hiện thực |

## 4. Kiến trúc

**Kiến trúc hệ thống** (`02-architecture/system/`)

| Tài liệu | Nội dung |
| :--- | :--- |
| [system-context.md](02-architecture/system/system-context.md) | C4 Level 1 — hệ thống và các actor / hệ thống ngoài |
| [container-view.md](02-architecture/system/container-view.md) | C4 Level 2 — các container runtime và giao tiếp |
| [runtime-flows.md](02-architecture/system/runtime-flows.md) | Sequence diagram các luồng chính |
| [integration-messaging.md](02-architecture/system/integration-messaging.md) | Giao tiếp đồng bộ / bất đồng bộ (Kafka) |
| [security-architecture.md](02-architecture/system/security-architecture.md) | Xác thực, phân quyền, phiên, bề mặt tấn công |
| [data-architecture.md](02-architecture/system/data-architecture.md) | Database-per-service, quyền sở hữu dữ liệu |
| [runtime-view-local.md](02-architecture/system/runtime-view-local.md) | Môi trường chạy docker compose (góc nhìn ứng dụng) |

**Kiến trúc phần mềm** (`02-architecture/software/`)

| Tài liệu | Nội dung |
| :--- | :--- |
| [go-service-architecture.md](02-architecture/software/go-service-architecture.md) | Kiến trúc phân lớp các service Go |
| [bun-service-architecture.md](02-architecture/software/bun-service-architecture.md) | Kiến trúc các service Bun / Hono |
| [frontend-architecture.md](02-architecture/software/frontend-architecture.md) | SvelteKit BFF, monorepo frontend |
| [api-conventions.md](02-architecture/software/api-conventions.md) | Quy ước URL, response, phân trang, lỗi |
| [cross-cutting-concerns.md](02-architecture/software/cross-cutting-concerns.md) | Cấu hình, logging, lỗi, CORS, validation, test |
| [adr/README.md](02-architecture/adr/README.md) | Các quyết định kiến trúc (ADR-001 … ADR-008) |
| [target-revamp/](02-architecture/target-revamp/README.md) | 🎯 **Kiến trúc đích & Đề xuất tái cấu trúc** (NATS JetStream, Better-Auth, PG 18) |

## 5. Service

| Service / App | Stack | Port | Tài liệu |
| :--- | :--- | :---: | :--- |
| auth-service | Go | 8088 | [auth-service.md](03-services/auth-service.md) |
| product-service | Go | 8086 | [product-service.md](03-services/product-service.md) |
| order-service | Go | 8084 | [order-service.md](03-services/order-service.md) |
| payment-service | Go | 8085 | [payment-service.md](03-services/payment-service.md) |
| inventory-service | Go | 8082 | [inventory-service.md](03-services/inventory-service.md) |
| shipping-service | Go | 8087 | [shipping-service.md](03-services/shipping-service.md) |
| search-service | Go | 8094 | [search-service.md](03-services/search-service.md) |
| promotion-service | Bun / TS | 8093 | [promotion-service.md](03-services/promotion-service.md) |
| rating-service | Bun / TS | 8089 | [rating-service.md](03-services/rating-service.md) |
| media-service | Bun / TS | 8083 | [media-service.md](03-services/media-service.md) |
| notification-service | Bun / TS | 8090 | [notification-service.md](03-services/notification-service.md) |
| services/shared | Go module | — | [shared-go-module.md](03-services/shared-go-module.md) |
| Storefront (web) | SvelteKit | 3000 | [frontend-web.md](03-services/frontend-web.md) |
| Backoffice (admin) | SvelteKit | 3001 | [frontend-admin.md](03-services/frontend-admin.md) |
| frontend/packages | TS / Svelte | — | [frontend-packages.md](03-services/frontend-packages.md) |

## 6. Tính năng

Danh mục đầy đủ, trạng thái và ma trận truy vết: **[04-features/README.md](04-features/README.md)**.

| Nhóm | Tính năng |
| :--- | :--- |
| Identity | [FEAT-01 SSO login/logout](04-features/FEAT-01-sso-login-logout.md) · [FEAT-02 Đăng ký](04-features/FEAT-02-user-registration.md) · [FEAT-03 Làm mới phiên](04-features/FEAT-03-session-refresh.md) · [FEAT-04 User & role](04-features/FEAT-04-user-role-management.md) |
| Storefront | [FEAT-05 Trang chủ](04-features/FEAT-05-home-page.md) · [FEAT-06 Catalog](04-features/FEAT-06-product-catalog.md) · [FEAT-07 Chi tiết SP](04-features/FEAT-07-product-detail.md) · [FEAT-08 Giỏ hàng](04-features/FEAT-08-shopping-cart.md) · [FEAT-09 Checkout](04-features/FEAT-09-checkout-place-order.md) · [FEAT-10 Lịch sử đơn](04-features/FEAT-10-order-history.md) |
| Commerce backend | [FEAT-11 Thanh toán](04-features/FEAT-11-payment.md) · [FEAT-12 Tồn kho](04-features/FEAT-12-inventory-stock-check.md) · [FEAT-13 Yêu thích](04-features/FEAT-13-favourites.md) · [FEAT-14 Tìm kiếm](04-features/FEAT-14-catalog-search-suggest.md) · [FEAT-15 Đánh giá](04-features/FEAT-15-ratings-reviews.md) · [FEAT-16 Khuyến mãi](04-features/FEAT-16-promotion-management.md) · [FEAT-17 Thuế](04-features/FEAT-17-tax-management.md) · [FEAT-18 Media](04-features/FEAT-18-media-management.md) · [FEAT-19 Thông báo](04-features/FEAT-19-notifications-email.md) · [FEAT-20 Vận chuyển](04-features/FEAT-20-shipping-order-items.md) |
| Backoffice | [FEAT-21 Kiểm soát truy cập](04-features/FEAT-21-admin-access-control.md) · [FEAT-22 Dashboard](04-features/FEAT-22-admin-dashboard.md) · [FEAT-23 Sản phẩm](04-features/FEAT-23-admin-product-mgmt.md) · [FEAT-24 Danh mục](04-features/FEAT-24-admin-category-mgmt.md) · [FEAT-25 Đơn hàng](04-features/FEAT-25-admin-order-mgmt.md) · [FEAT-26 Module chờ phát triển](04-features/FEAT-26-admin-placeholder-modules.md) |
| Platform | [FEAT-27 i18n](04-features/FEAT-27-i18n.md) · [FEAT-28 SEO & trang lỗi](04-features/FEAT-28-seo-error-pages.md) |

## 7. API, dữ liệu, kế hoạch, hướng dẫn

| Tài liệu | Nội dung |
| :--- | :--- |
| [05-api/api-catalog.md](05-api/api-catalog.md) | Toàn bộ endpoint của 11 service, endpoint frontend thực sự gọi |
| [05-api/gateway-routing.md](05-api/gateway-routing.md) | Cấu hình APISIX: upstream, route, plugin |
| [06-data/data-dictionary.md](06-data/data-dictionary.md) | Mọi database, bảng, cột |
| [07-plan/known-issues-tech-debt.md](07-plan/known-issues-tech-debt.md) | Danh sách lỗi & nợ kỹ thuật (ISSUE-xx) |
| [07-plan/project-plan.md](07-plan/project-plan.md) | Lộ trình theo giai đoạn |
| [07-plan/backlog.md](07-plan/backlog.md) | Backlog ưu tiên (epic → story) |
| [07-plan/test-plan.md](07-plan/test-plan.md) | Hiện trạng & kế hoạch kiểm thử |
| [08-guides/developer-guide.md](08-guides/developer-guide.md) | Chạy local, thêm endpoint/route, troubleshooting |
| [08-guides/coding-conventions.md](08-guides/coding-conventions.md) | Quy ước code Go / Bun / SvelteKit |
| [08-guides/keycloak-postgres-rollout.md](08-guides/keycloak-postgres-rollout.md) | Hướng dẫn cấu hình môi trường Keycloak & PostgreSQL |
| [web/README.md](web/README.md) | Ghi chú & backlog quá trình chuyển đổi SvelteKit 3 |

## 8. Bảo trì tài liệu

- Khi sửa code làm thay đổi hành vi, cập nhật **FEAT** và **service doc** tương ứng trong cùng PR.
- Khi sửa một `ISSUE-xx`, đánh dấu "Đã xử lý" trong [known-issues-tech-debt.md](07-plan/known-issues-tech-debt.md) và cập nhật trạng thái FEAT liên quan.
- ID mới (BR / FEAT / ISSUE / ADR) phải được đăng ký trong [conventions-and-ids.md](00-overview/conventions-and-ids.md).
