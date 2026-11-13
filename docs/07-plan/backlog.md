# Product Backlog

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> Backlog ưu tiên theo epic → story, suy ra từ [project-plan.md](project-plan.md). Ưu tiên MoSCoW: **M**ust · **S**hould · **C**ould. Mỗi story liên kết FEAT/ISSUE để truy vết. Acceptance criteria chi tiết của hành vi hiện tại nằm trong tài liệu FEAT; tiêu chí dưới đây là **cho thay đổi mới**.

## EPIC-01 — Ổn định hợp đồng Frontend ↔ Backend (P0)

| ID | Story | Ưu tiên | Liên kết | Acceptance criteria |
| :--- | :--- | :---: | :--- | :--- |
| US-0101 | Là khách, tôi muốn danh sách sản phẩm phân trang và lọc đúng theo danh mục/từ khóa/sắp xếp | M | FEAT-06, ISSUE-48 | `GET /api/products?page=1&size=12&categoryId=2&sort=priceUnit,asc&search=abc` trả `Page` đúng `totalElements`; UI hiển thị đúng số trang |
| US-0102 | Là admin, tôi muốn danh sách danh mục/sản phẩm/đơn trong backoffice có dữ liệu thật | M | FEAT-22..25, ISSUE-48 | Dashboard hiển thị tổng số khớp DB |
| US-0103 | Là admin, tôi đăng nhập được vào backoffice | M | FEAT-01, FEAT-21, ISSUE-100, ISSUE-101 | User có role ADMIN vào `/dashboard`; user USER vào `/access-denied` |
| US-0104 | Là người dùng, đăng nhập SSO thành công với Keycloak thật | M | FEAT-01, ISSUE-04 | Không còn mock token; lỗi Keycloak hiển thị trang lỗi callback |
| US-0105 | Là khách, sau đăng nhập tôi không bị coi là guest khi profile lỗi | M | FEAT-03, ISSUE-67, ISSUE-69 | Phiên giữ được qua nhiều chu kỳ refresh |
| US-0106 | Là admin, danh mục của sản phẩm được lưu đúng | M | FEAT-23, ISSUE-60 | Tạo SP với danh mục X → GET trả `category.categoryId = X` |
| US-0107 | Là dev, mọi route gateway đều tới đúng service | M | ISSUE-13, ISSUE-15 | Smoke test tất cả route trong gateway-routing.md trả khác 404/502 |
| US-0108 | Là khách, sản phẩm trên trang chủ đặt hàng được | S | FEAT-05, ISSUE-38 | Thêm SP từ trang chủ → checkout thành công |
| US-0109 | Là khách mới, tôi tự tạo được tài khoản | M | FEAT-02, ISSUE-71 | Nút "Tạo tài khoản" dẫn tới form/trang đăng ký hoạt động; đăng nhập được ngay sau đó |

## EPIC-02 — Bảo mật & phân quyền (P1)

| ID | Story | Ưu tiên | Liên kết | Acceptance criteria |
| :--- | :--- | :---: | :--- | :--- |
| US-0201 | API quản trị chỉ chấp nhận ADMIN | M | ISSUE-01, ISSUE-102 | Token USER gọi `POST /api/products` → 403; không token → 401 |
| US-0202 | Khách chỉ thấy đơn/giỏ/thanh toán của mình | M | FEAT-10, ISSUE-49, ISSUE-76 | User A không đọc được order của user B (403/404) |
| US-0203 | Không thể giả mạo cookie phiên | M | ISSUE-66, ISSUE-68 | Sửa `session_user.roles` không vào được backoffice; gọi `/_app/remote/...` không có role → 403 |
| US-0204 | Redirect sau đăng nhập chỉ tới domain cho phép | M | ISSUE-05 | `redirect_uri` ngoài allowlist → 400 |
| US-0205 | Đăng nhập ổn định khi chạy nhiều replica auth-service | S | ISSUE-06 | Kiểm thử 2 replica, 100 lần login không lỗi |
| US-0206 | Chỉ khách đã đăng nhập mới được đánh giá | S | FEAT-15, ISSUE-03 | POST rating không token → 401 |
| US-0207 | Email chỉ gửi tới chủ tài khoản liên quan | M | FEAT-19, ISSUE-09, ISSUE-64 | Email thanh toán đi tới email của `userId` |
| US-0208 | Người dùng đổi được mật khẩu | S | FEAT-04, ISSUE-07 | Đổi mật khẩu → đăng nhập bằng mật khẩu mới |
| US-0209 | Cookie an toàn ở production | S | ISSUE-11 | `Secure` bật khi HTTPS; tên cookie khác nhau giữa web/admin |
| US-0210 | Service từ chối request không đi qua gateway hoặc token không hợp lệ | M | ISSUE-02, ISSUE-134 | Gọi trực tiếp `:8086` không token → 401; data store không publish ra host |
| US-0211 | Callback SSO từ chối `state` không hợp lệ; dùng PKCE | M | FEAT-01, ISSUE-52, ISSUE-133 | Callback với state lạ → lỗi `invalid_state`, không cấp ticket |
| US-0212 | Đăng ký chỉ tạo role USER; user tạo được thật trên Keycloak | M | FEAT-02, ISSUE-08, ISSUE-53 | Gửi `roles:["ADMIN"]` khi signup → vẫn chỉ USER; Keycloak down → lỗi rõ ràng, không id giả |
| US-0213 | Keycloak được hardening | S | ISSUE-133, ISSUE-136 | Brute-force bị khóa; password policy; auth-service dùng service account |
| US-0214 | Không lộ PII/secret qua log, lỗi, API public | S | ISSUE-93, ISSUE-132, ISSUE-137, ISSUE-138 | Response 500 không chứa SQL; log không chứa `ticket`/`code`; `size` bị giới hạn |

## EPIC-03 — Đơn hàng, thanh toán, tồn kho (P2)

| ID | Story | Ưu tiên | Liên kết | Acceptance criteria |
| :--- | :--- | :---: | :--- | :--- |
| US-0301 | Đơn hàng lưu đủ mọi sản phẩm và số lượng | M | FEAT-09, ISSUE-20 | Giỏ 3 dòng → đơn có 3 `order_items` |
| US-0302 | Đơn hàng lưu địa chỉ giao và phương thức thanh toán | M | FEAT-09, ISSUE-22 | Admin xem chi tiết đơn thấy địa chỉ |
| US-0303 | Đơn có trạng thái và admin cập nhật được | M | FEAT-25, ISSUE-41 | Đổi trạng thái theo state machine; khách thấy trạng thái mới |
| US-0304 | Không bán quá tồn kho | M | FEAT-12, ISSUE-26 | Đặt SP hết hàng → lỗi rõ ràng, đơn không tạo |
| US-0305 | Bấm đặt hàng 2 lần không tạo 2 đơn | M | ISSUE-23 | Cùng idempotency key → cùng orderId |
| US-0306 | Tổng tiền đơn do server tính | M | ISSUE-21 | `orderFee` client gửi bị bỏ qua |
| US-0307 | Lỗi thanh toán trùng trả 409 | S | FEAT-11, ISSUE-25 | |
| US-0308 | Lọc giá và lọc mã đơn trên toàn bộ dữ liệu | S | ISSUE-40 | |
| US-0309 | Thanh toán ONLINE không bị ghi "đã trả" khi chưa nhận tiền | M | FEAT-09, FEAT-11, ISSUE-75 | Đơn ONLINE tạo payment `isPayed=false`; vẫn ghi nhận được payment thật sau đó |
| US-0310 | Đơn hàng gắn với người đặt | M | FEAT-09, FEAT-10, ISSUE-76 | `orders.user_id` = `sub` của token |

## EPIC-04 — Sự kiện & tìm kiếm (P3)

| ID | Story | Ưu tiên | Liên kết | Acceptance criteria |
| :--- | :--- | :---: | :--- | :--- |
| US-0401 | Khách nhận email khi thanh toán thành công | S | FEAT-11, FEAT-19, ISSUE-17, ISSUE-18 | Thanh toán COMPLETED → 1 email + 1 `payment_notifications`; COD không phát |
| US-0402 | Người dùng mới nhận email chào mừng | C | FEAT-02, FEAT-19 | Đăng ký → email |
| US-0403 | Sản phẩm mới tìm thấy được trên search | S | FEAT-14, ISSUE-19 | Tạo SP → xuất hiện trong `catalog-search` < 5 giây |
| US-0404 | Kết quả tìm kiếm sắp xếp và có facet | C | ISSUE-34 | |

## EPIC-05 — Nền tảng kỹ thuật (P4)

| ID | Story | Ưu tiên | Liên kết |
| :--- | :--- | :---: | :--- |
| US-0501 | Migration có version cho mọi service | M | ISSUE-28 |
| US-0502 | Một định dạng response/lỗi chuẩn | S | ISSUE-29, ISSUE-65 |
| US-0503 | Readiness kiểm DB/Kafka/S3 | S | ISSUE-36 |
| US-0504 | Validate input mọi endpoint | S | ISSUE-35, ISSUE-62 |
| US-0505 | Bộ test integration + contract + E2E | M | ISSUE-46 |
| US-0506 | Dọn env thừa, dead code, cập nhật README | C | ISSUE-37, ISSUE-43, ISSUE-45 |
| US-0507 | Media an toàn và dùng được làm ảnh public | S | ISSUE-32, ISSUE-63 |

## EPIC-06 — Mở rộng tính năng (P5)

| ID | Story | Ưu tiên | Liên kết |
| :--- | :--- | :---: | :--- |
| US-0601 | Khách xem và viết đánh giá trên trang chi tiết | S | FEAT-15 |
| US-0602 | Khách lưu sản phẩm yêu thích | C | FEAT-13 |
| US-0603 | Ô tìm kiếm có gợi ý | S | FEAT-14 |
| US-0604 | Áp mã giảm giá khi checkout | S | FEAT-16 |
| US-0605 | Thuế được tính trong tổng đơn | C | FEAT-17 |
| US-0606 | Admin quản lý người dùng & role | S | FEAT-04, FEAT-26 |
| US-0607 | Admin quản lý tồn kho | S | FEAT-12, FEAT-26 |
| US-0608 | Admin quản lý khuyến mãi, đánh giá | C | FEAT-16, FEAT-15 |
| US-0609 | Admin upload ảnh sản phẩm | S | FEAT-18, FEAT-23 |
| US-0610 | Menu danh mục động | S | FEAT-05, ISSUE-44 |
| US-0611 | Giữ locale trên mọi link; định dạng theo locale | C | FEAT-27, ISSUE-39, ISSUE-83 |
| US-0613 | Trang chi tiết sản phẩm chuẩn SEO, sitemap có sản phẩm | C | FEAT-28, ISSUE-82 |
| US-0614 | Header storefront dùng được trên mobile; đăng nhập quay về trang đang xem | S | FEAT-01, FEAT-06, ISSUE-120, ISSUE-121 |
| US-0615 | Thông báo lỗi checkout rõ ràng theo nguyên nhân | S | FEAT-09, ISSUE-122 |
| US-0612 | Giỏ hàng đồng bộ server (đa thiết bị) | C | FEAT-08 |
