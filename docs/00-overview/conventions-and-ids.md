# Quy ước tài liệu & Sổ đăng ký ID

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> Tài liệu này là **nguồn chuẩn cho mọi ID** (BR, FEAT, ISSUE, ADR) dùng trong `docs-all/`. Khi thêm ID mới, cập nhật file này trước.

## 1. Nguyên tắc viết

| # | Nguyên tắc |
| :--- | :--- |
| P-1 | **Mô tả đúng code hiện tại (as-is).** Source code là nguồn sự thật. Không mô tả kiến trúc đích trong `docs/ARCHITECTURE.md` / `docs/architecture-revamp-plan.md`. |
| P-2 | Viết **tiếng Việt**, giữ nguyên thuật ngữ kỹ thuật tiếng Anh (endpoint, handler, repository, cookie, token…). |
| P-3 | Sơ đồ dùng **Mermaid** trong block ` ```mermaid `. |
| P-4 | Mỗi khẳng định quan trọng có tham chiếu source dạng `` `path/to/file.go:123` `` (đường dẫn tính từ gốc repo). |
| P-5 | **Không** chép giá trị secret (mật khẩu, client secret, access key) từ `.env`, `ecommerce-realm.json`, manifest k8s. Chỉ ghi **tên biến** và trỏ tới file. |
| P-6 | Chức năng chưa có hoặc có lỗi được ghi trung thực với trạng thái và liên kết `ISSUE-xx`. |
| P-7 | Phạm vi hiện tại: **ứng dụng** (backend service + frontend). Hạ tầng (k8s, CI/CD, GitOps, observability stack) sẽ có bộ tài liệu riêng. |

## 2. Trạng thái triển khai

| Ký hiệu | Ý nghĩa |
| :--- | :--- |
| ✅ **Đã triển khai** | Có code chạy được đầu-cuối (UI → API → DB) hoặc API đầy đủ theo mô tả. |
| 🟡 **Một phần** | Có code nhưng thiếu luồng, có lỗi chặn, hoặc chỉ có backend/chỉ có frontend. |
| ⛔ **Chưa có** | Chỉ có placeholder / chưa có code. |

Mức độ nghiêm trọng của issue: **Critical** (lỗ hổng bảo mật / mất dữ liệu / chặn luồng chính) · **High** (sai nghiệp vụ, tính năng không dùng được) · **Medium** (sai lệch, thiếu ràng buộc) · **Low** (code smell, tài liệu, cấu hình thừa).

## 3. Sổ đăng ký Business Requirement ↔ Feature

Mỗi yêu cầu nghiệp vụ `BR-xx` ánh xạ 1-1 tới một tài liệu tính năng `FEAT-xx` trong [`04-features/`](../04-features/README.md).

| BR | FEAT | Tên tính năng | App / Service chính | Trạng thái |
| :--- | :--- | :--- | :--- | :--- |
| BR-01 | [FEAT-01](../04-features/FEAT-01-sso-login-logout.md) | Đăng nhập SSO / Đăng xuất | web, admin, auth-service, Keycloak | 🟡 |
| BR-02 | [FEAT-02](../04-features/FEAT-02-user-registration.md) | Đăng ký tài khoản | auth-service | 🟡 |
| BR-03 | [FEAT-03](../04-features/FEAT-03-session-refresh.md) | Duy trì & làm mới phiên | web, admin, auth-service | 🟡 |
| BR-04 | [FEAT-04](../04-features/FEAT-04-user-role-management.md) | Quản lý người dùng & vai trò (API) | auth-service | 🟡 |
| BR-05 | [FEAT-05](../04-features/FEAT-05-home-page.md) | Trang chủ storefront | web | 🟡 |
| BR-06 | [FEAT-06](../04-features/FEAT-06-product-catalog.md) | Danh mục sản phẩm (lọc/sắp xếp/phân trang) | web, product-service | 🟡 |
| BR-07 | [FEAT-07](../04-features/FEAT-07-product-detail.md) | Chi tiết sản phẩm | web, product-service | ✅ |
| BR-08 | [FEAT-08](../04-features/FEAT-08-shopping-cart.md) | Giỏ hàng | web (client-side) | ✅ |
| BR-09 | [FEAT-09](../04-features/FEAT-09-checkout-place-order.md) | Checkout & đặt hàng | web, order-service, payment-service | 🟡 |
| BR-10 | [FEAT-10](../04-features/FEAT-10-order-history.md) | Lịch sử đơn hàng của khách | web, order-service | 🟡 |
| BR-11 | [FEAT-11](../04-features/FEAT-11-payment.md) | Ghi nhận thanh toán | payment-service | 🟡 |
| BR-12 | [FEAT-12](../04-features/FEAT-12-inventory-stock-check.md) | Kiểm tra tồn kho | inventory-service | 🟡 |
| BR-13 | [FEAT-13](../04-features/FEAT-13-favourites.md) | Sản phẩm yêu thích | product-service | 🟡 |
| BR-14 | [FEAT-14](../04-features/FEAT-14-catalog-search-suggest.md) | Tìm kiếm full-text & gợi ý | search-service | 🟡 |
| BR-15 | [FEAT-15](../04-features/FEAT-15-ratings-reviews.md) | Đánh giá sản phẩm | rating-service | 🟡 |
| BR-16 | [FEAT-16](../04-features/FEAT-16-promotion-management.md) | Quản lý khuyến mãi | promotion-service | 🟡 |
| BR-17 | [FEAT-17](../04-features/FEAT-17-tax-management.md) | Quản lý thuế | promotion-service | 🟡 |
| BR-18 | [FEAT-18](../04-features/FEAT-18-media-management.md) | Quản lý media (upload/stream) | media-service | 🟡 |
| BR-19 | [FEAT-19](../04-features/FEAT-19-notifications-email.md) | Thông báo & email | notification-service | 🟡 |
| BR-20 | [FEAT-20](../04-features/FEAT-20-shipping-order-items.md) | Dòng hàng vận chuyển | shipping-service | 🟡 |
| BR-21 | [FEAT-21](../04-features/FEAT-21-admin-access-control.md) | Kiểm soát truy cập backoffice | admin | 🟡 |
| BR-22 | [FEAT-22](../04-features/FEAT-22-admin-dashboard.md) | Dashboard quản trị | admin | 🟡 |
| BR-23 | [FEAT-23](../04-features/FEAT-23-admin-product-mgmt.md) | Quản lý sản phẩm (admin) | admin, product-service | 🟡 |
| BR-24 | [FEAT-24](../04-features/FEAT-24-admin-category-mgmt.md) | Quản lý danh mục (admin) | admin, product-service | 🟡 |
| BR-25 | [FEAT-25](../04-features/FEAT-25-admin-order-mgmt.md) | Quản lý đơn hàng (admin) | admin, order-service | 🟡 |
| BR-26 | [FEAT-26](../04-features/FEAT-26-admin-placeholder-modules.md) | Module Users / Inventory / Shipping / Settings (admin) | admin | ⛔ |
| BR-27 | [FEAT-27](../04-features/FEAT-27-i18n.md) | Đa ngôn ngữ vi / en | web, admin | ✅ |
| BR-28 | [FEAT-28](../04-features/FEAT-28-seo-error-pages.md) | SEO & trang lỗi | web, admin | ✅ |

> Trạng thái ở bảng trên là đánh giá ban đầu; tài liệu FEAT tương ứng là nơi chốt cuối cùng. Nếu có thay đổi, cập nhật cả hai.

## 4. Sổ đăng ký Issue

Danh sách đầy đủ, mô tả chi tiết và đề xuất xử lý: [`07-plan/known-issues-tech-debt.md`](../07-plan/known-issues-tech-debt.md). Khoảng ID theo nhóm:

| Khoảng | Nhóm |
| :--- | :--- |
| ISSUE-01 … ISSUE-12 | Bảo mật & xác thực |
| ISSUE-13 … ISSUE-16 | API Gateway (APISIX) |
| ISSUE-17 … ISSUE-19 | Messaging (Kafka) |
| ISSUE-20 … ISSUE-37 | Nghiệp vụ & dữ liệu backend |
| ISSUE-38 … ISSUE-47 | Frontend |
| ISSUE-48 … ISSUE-51 | Hợp đồng API Frontend ↔ Backend |
| ISSUE-52 … ISSUE-59 | Service Go & data |
| ISSUE-60 … ISSUE-69 | Service Bun, frontend, API |
| ISSUE-70 … ISSUE-89 | Tài liệu tính năng |
| ISSUE-90 … ISSUE-99 | Kiến trúc phần mềm & hướng dẫn |
| ISSUE-100 … ISSUE-109 | Nhóm biên soạn chính (hợp đồng profile/role, catalog không xác thực) |
| ISSUE-110 … ISSUE-119 | OpenAPI/AsyncAPI, error catalog, cấu hình |
| ISSUE-120 … ISSUE-129 | Use case, RTM, screen catalog |
| ISSUE-130 … ISSUE-139 | C4 L3, quality scenarios, threat model, PII |
| ISSUE-140 trở đi | Dành cho cập nhật sau phiên bản 1.0 |

ID bị đánh dấu **trùng** (giữ lại để không tái sử dụng): ISSUE-54 → 102, ISSUE-90 → 100/101/67, ISSUE-91 → 66, ISSUE-92 → 60, ISSUE-95 → 102, ISSUE-99 → 61, ISSUE-116 → 133, ISSUE-118 → 135. ID không sử dụng: 70, 72, 73, 77, 84–89, 103–109, 119, 124–129.

## 5. Sổ đăng ký ADR

Xem [`02-architecture/adr/README.md`](../02-architecture/adr/README.md). ADR ghi lại các **quyết định đang hiện diện trong code**, không phải đề xuất tương lai.

## 6. Template tài liệu tính năng (`04-features/FEAT-xx-*.md`)

```markdown
# FEAT-xx — <Tên tính năng>

| Thuộc tính | Giá trị |
| :--- | :--- |
| ID | FEAT-xx |
| Business Requirement | BR-xx |
| Trạng thái | ✅ / 🟡 / ⛔ + 1 câu giải thích |
| Actor | Guest / Customer / Admin / Hệ thống |
| Ứng dụng / Service | ... |
| Issue liên quan | ISSUE-xx, ... |

## 1. Mô tả & mục tiêu nghiệp vụ
## 2. User stories
## 3. Tiền điều kiện / Hậu điều kiện
## 4. Luồng nghiệp vụ
### 4.1 Luồng chính
### 4.2 Luồng thay thế
### 4.3 Luồng ngoại lệ
## 5. Sequence diagram
## 6. Màn hình & route (UI)
## 7. API sử dụng
## 8. Dữ liệu
## 9. Business rules & validation
## 10. Xử lý lỗi
## 11. i18n (key chính)
## 12. Acceptance criteria (Given / When / Then — theo hành vi hiện tại)
## 13. Giới hạn & vấn đề đã biết
## 14. Test hiện có
## 15. Tham chiếu source
```

Mục không áp dụng thì ghi "Không áp dụng" kèm lý do (ví dụ: tính năng chỉ có backend → mục 6, 11 ghi "Chưa có UI").

## 7. Template tài liệu service (`03-services/*.md`)

```markdown
# <service-name>

| Thuộc tính | Giá trị |  (Ngôn ngữ/stack, port, database, đường dẫn code, prefix gateway, tính năng phục vụ FEAT-xx)

## 1. Tổng quan & trách nhiệm
## 2. Tech stack & thư viện
## 3. Cấu hình (biến môi trường: tên · default · ý nghĩa — KHÔNG ghi giá trị secret)
## 4. Phụ thuộc (DB, Kafka topic, OpenSearch, S3, Keycloak, SMTP, service khác)
## 5. Cấu trúc code (cây thư mục + vai trò từng file)
## 6. API endpoints (bảng đầy đủ: method, path trong service, path qua gateway, auth, mô tả, request, response, mã lỗi)
## 7. Data model (ERD Mermaid + bảng cột)
## 8. Business logic chi tiết
## 9. Sự kiện (produce / consume)
## 10. Health check & observability
## 11. Test
## 12. Vấn đề đã biết (ISSUE-xx)
## 13. Tham chiếu source
```
