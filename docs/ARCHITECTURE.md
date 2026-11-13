# Kiến trúc hệ thống — Architecture Guide

Tài liệu kiến trúc của dự án được tổ chức rõ ràng thành hai phần theo hiện trạng và định hướng:

---

### 1. 🏗️ Kiến trúc hiện tại (As-Is Architecture — Phản ánh 100% Codebase)

> **Điểm vào chính:** 👉 **[`02-architecture/README.md`](02-architecture/README.md)**

Tài liệu mô tả chi tiết, chính xác mã nguồn hiện tại đang chạy (7 Go services + 4 Bun services + 2 SvelteKit apps + APISIX + Keycloak + Kafka):
- **Kiến trúc hệ thống (C4 & arc42):**
  - [System Context (C4 Level 1)](02-architecture/system/system-context.md)
  - [Container View (C4 Level 2)](02-architecture/system/container-view.md)
  - [Component View (C4 Level 3)](02-architecture/system/component-view.md)
  - [Runtime Flows (Sequence diagrams)](02-architecture/system/runtime-flows.md)
  - [Security Architecture](02-architecture/system/security-architecture.md) & [Threat Model](02-architecture/system/threat-model.md)
  - [Data Architecture](02-architecture/system/data-architecture.md)
- **Kiến trúc phần mềm:**
  - [Go Service Architecture](02-architecture/software/go-service-architecture.md)
  - [Bun/Hono Service Architecture](02-architecture/software/bun-service-architecture.md)
  - [Frontend SvelteKit BFF Architecture](02-architecture/software/frontend-architecture.md)
  - [API Conventions](02-architecture/software/api-conventions.md)
- **Quyết định kiến trúc:** [ADR-001 đến ADR-008](02-architecture/adr/README.md)

---

### 2. 🎯 Kiến trúc đích & Đề xuất tái cấu trúc (Target Revamp Architecture)

> **Điểm vào:** 👉 **[`02-architecture/target-revamp/README.md`](02-architecture/target-revamp/README.md)**

Các tài liệu nghiên cứu và đề xuất định hướng tương lai (chuyển đổi sang NATS JetStream, Better-Auth, PostgreSQL 18, Cilium...):
- [**ARCHITECTURE.md (Bản đặc tả kiến trúc đích)**](02-architecture/target-revamp/ARCHITECTURE.md)
- [**architecture-revamp-plan.md (Kế hoạch tái cấu trúc)**](02-architecture/target-revamp/architecture-revamp-plan.md)
