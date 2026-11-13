# Lịch sử thay đổi tài liệu (Documentation Changelog)

Tài liệu này ghi nhận toàn bộ các mốc phát hành và cập nhật của bộ tài liệu **E-Commerce Microservices Platform**.

---

## [1.0.0] - 2026-11-12 (Bản Final)

### Bổ sung (Added)
- **Kiểm soát & Quy chuẩn:** Hoàn thiện `document-control.md`, `glossary.md`, `conventions-and-ids.md`.
- **Yêu cầu & Phân tích (BRD):** Bổ sung use cases, actors & permissions matrix, non-functional requirements (NFRs), Requirements Traceability Matrix (RTM).
- **Kiến trúc (Arc42 & C4):** Bổ sung C4 Level 3 (`component-view.md`), quality scenarios (`quality-scenarios.md`), threat model STRIDE (`threat-model.md`), chính sách dữ liệu cá nhân Nghị định 13/2023 (`data-privacy-pii.md`).
- **Đặc tả API & Sự kiện:** Khởi tạo bộ OpenAPI 3.1 YAML cho 11 microservices và AsyncAPI 3.0 YAML cho Kafka events; bổ sung error catalog và event catalog.
- **Kế hoạch & Đánh giá hiện trạng:** Hoàn thiện audit source code với 107 open issues (`known-issues-tech-debt.md`), lộ trình 6 giai đoạn P0 → P5 (`project-plan.md`), sprint backlog và test plan chi tiết.
- **Tổ chức thư mục:** Hợp nhất bộ tài liệu vào thư mục chính thức `docs/` của dự án; lưu trữ các đề xuất kiến trúc đích dài hạn vào `02-architecture/target-revamp/`.

---

## [0.5.0] - 2026-11-12

### Bổ sung (Added)
- Khởi tạo tài liệu chi tiết cho 11 Backend microservices và 2 Frontend SvelteKit BFF (`03-services/`).
- Khởi tạo danh mục 28 tính năng `FEAT-01` đến `FEAT-28` (`04-features/`).
- Tổng hợp từ điển dữ liệu toàn hệ thống (`06-data/data-dictionary.md`) và bảng tra cứu API (`05-api/api-catalog.md`).
- Bổ sung 8 ADRs (`02-architecture/adr/`).
- Nhận diện các vấn đề từ `ISSUE-52` đến `ISSUE-102`.

---

## [0.1.0] - 2026-11-12

### Khởi tạo (Initial Draft)
- Khảo sát mã nguồn hiện tại của toàn bộ 11 service và frontend monorepo.
- Xây dựng khung tài liệu, quy ước định danh ID.
- Lập danh sách 51 vấn đề kỹ thuật và bảo mật ban đầu (`ISSUE-01` đến `ISSUE-51`).

