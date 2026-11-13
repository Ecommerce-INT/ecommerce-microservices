# Kiến trúc đích & Đề xuất tái cấu trúc (Target Revamp Architecture)

> **Trạng thái:** `TARGET / PROPOSAL` (Kiến trúc đích dài hạn)  
> **Lưu ý quan trọng:** Các tài liệu trong thư mục này mô tả **định hướng tái cấu trúc trong tương lai** (chuyển dịch sang Better-Auth, NATS JetStream, PostgreSQL 18, OpenSearch 3.9 nâng cao, CNPG...).  
> Nếu bạn cần tra cứu **kiến trúc mã nguồn thực tế đang chạy hiện tại (As-Is)**, vui lòng xem [02-architecture/README.md](../README.md).

---

## 1. Danh mục tài liệu trong thư mục này

1. [**ARCHITECTURE.md**](./ARCHITECTURE.md):
   - Đặc tả kiến trúc đích tổng thể: 11 microservices (7 Go + 4 Bun), SvelteKit SSR BFF, Apache APISIX Edge, NATS JetStream, OpenSearch 3.9, PostgreSQL 18.
   - Sơ đồ C4 Level 0 đến Level 2, các ràng buộc chất lượng (NFRs, không dùng JVM runtime, cold-start < 100ms, idle RAM < 1.5GB).
2. [**architecture-revamp-plan.md**](./architecture-revamp-plan.md):
   - Kế hoạch chi tiết tái cấu trúc từ nền tảng Spring Boot / Next.js cũ sang Go / Bun / SvelteKit / Better-Auth / NATS.
   - Đánh giá các điểm lỗi lịch sử (Frankenstein split, $N+1$ waterfalls, nano-services) và ma trận di chuyển chi tiết.

---

## 2. Liên kết liên quan

- 🧭 [Kiến trúc hiện tại (As-Is arc42 & C4)](../README.md)
- ⚠️ [Danh sách lỗi & nợ kỹ thuật hiện tại (Known Issues)](../../07-plan/known-issues-tech-debt.md)
- 🗺️ [Lộ trình hoàn thiện hệ thống hiện tại P0 → P5](../../07-plan/project-plan.md)

