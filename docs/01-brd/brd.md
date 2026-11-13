# Business Requirements Document (BRD)

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

| Thuộc tính | Giá trị |
| :--- | :--- |
| Sản phẩm | E-Commerce Microservices Platform (Storefront + Backoffice) |
| Phiên bản tài liệu | 1.0 |
| Ngày | 2026-11-12 |
| Cơ sở | Đọc ngược từ source code hiện tại (reverse-engineered BRD) |
| Tài liệu liên quan | [Actors & quyền](actors-and-permissions.md) · [Quy trình nghiệp vụ](business-processes.md) · [NFR](non-functional-requirements.md) · [Tính năng](../04-features/README.md) |

> **Lưu ý cách đọc:** BRD này được xây dựng từ hệ thống đang có. Mỗi yêu cầu nghiệp vụ (`BR-xx`) mô tả **nhu cầu nghiệp vụ mà code đang phục vụ**, kèm **trạng thái đáp ứng thực tế** và các khoảng trống (`ISSUE-xx`). Yêu cầu không có dấu vết trong code thì **không** được đưa vào (ngoại trừ phần "Ngoài phạm vi hiện tại").

---

## 1. Bối cảnh kinh doanh

Doanh nghiệp bán lẻ trực tuyến (thị trường chính Việt Nam — giá VND, giao diện mặc định tiếng Việt, có tiếng Anh) cần:

1. Một **website bán hàng** cho khách: xem sản phẩm theo danh mục, thêm vào giỏ, đặt hàng với thanh toán khi nhận hàng (COD) hoặc trực tuyến, theo dõi đơn đã đặt.
2. Một **website quản trị** cho nhân viên: theo dõi số liệu tổng quan, quản lý sản phẩm và danh mục, xem đơn hàng.
3. Một **nền tảng backend mở rộng được** cho các nghiệp vụ hỗ trợ: tài khoản & phân quyền, thanh toán, tồn kho, vận chuyển, tìm kiếm, khuyến mãi & thuế, đánh giá, media, thông báo.

Hệ thống được tổ chức dưới dạng microservice, mỗi nghiệp vụ một service, để các nhóm phát triển độc lập.

## 2. Mục tiêu kinh doanh

| ID | Mục tiêu | Chỉ số đo (gợi ý) | Mức đáp ứng hiện tại |
| :--- | :--- | :--- | :--- |
| BO-1 | Cho phép khách hàng tự mua hàng trực tuyến 24/7 | Tỷ lệ checkout thành công | 🟡 Luồng đặt hàng có UI hoàn chỉnh; backend chỉ lưu đơn tối giản (1 sản phẩm/đơn, không trạng thái). |
| BO-2 | Trải nghiệm song ngữ, thân thiện SEO | Trang được index, Lighthouse SEO | ✅ i18n vi/en, sitemap, hreflang, robots. |
| BO-3 | Nhân viên tự quản lý catalog không cần dev | Thời gian cập nhật sản phẩm | 🟡 CRUD sản phẩm/danh mục có UI; phụ thuộc hợp đồng API ([ISSUE-48](../07-plan/known-issues-tech-debt.md#6-hợp-đồng-api-frontend--backend)). |
| BO-4 | Đăng nhập tập trung, an toàn | 0 token trong JS trình duyệt | 🟡 Token chỉ nằm trong cookie httpOnly; nhưng thiếu phân quyền backend ([ISSUE-01](../07-plan/known-issues-tech-debt.md#1-bảo-mật--xác-thực)). |
| BO-5 | Nền tảng cho khuyến mãi, đánh giá, thông báo | Số chương trình KM, số đánh giá | 🟡 Có API quản lý dữ liệu, chưa tích hợp vào luồng mua hàng/UI. |

## 3. Stakeholder

| Stakeholder | Quan tâm chính | Tương tác với hệ thống |
| :--- | :--- | :--- |
| Khách hàng (Guest/Customer) | Tìm, mua sản phẩm dễ dàng, đúng giá, theo dõi đơn | Storefront |
| Nhân viên vận hành / Quản trị | Quản lý catalog, xem đơn, số liệu | Backoffice |
| Chủ doanh nghiệp / PO | Doanh thu, mức độ hoàn thiện tính năng | Báo cáo, dashboard |
| Đội phát triển | Kiến trúc rõ ràng, dễ mở rộng | Repo, tài liệu này |
| Nhà cung cấp định danh (Keycloak) | Cấu hình realm, client, role | Tích hợp OIDC |
| Nhà cung cấp email (SMTP) | Gửi email giao dịch | notification-service |

## 4. Phạm vi

### 4.1 Trong phạm vi (có trong code)

- Định danh: đăng nhập SSO qua Keycloak, đăng ký qua API, làm mới phiên, đăng xuất, quản lý user & role (API).
- Storefront: trang chủ, danh mục sản phẩm, chi tiết sản phẩm, giỏ hàng, checkout 3 bước, lịch sử đơn.
- Backoffice: kiểm soát truy cập admin, dashboard, quản lý sản phẩm, danh mục, xem đơn hàng.
- Backend hỗ trợ (API, chưa có UI): thanh toán, tồn kho, yêu thích, tìm kiếm, đánh giá, khuyến mãi, thuế, media, thông báo/email, dòng hàng vận chuyển.
- Nền tảng: đa ngôn ngữ vi/en, SEO, trang lỗi.

### 4.2 Ngoài phạm vi hiện tại (không có trong code)

| Hạng mục | Ghi chú |
| :--- | :--- |
| Tích hợp cổng thanh toán thật (Stripe/VNPay/MoMo), webhook | Thanh toán chỉ được ghi nhận nội bộ. |
| Tích hợp đối tác vận chuyển, tính phí ship theo địa chỉ, tracking | Phí ship là hằng số phía frontend. |
| Áp dụng mã giảm giá / thuế vào giá đơn hàng | Dữ liệu KM/thuế có nhưng không được dùng khi checkout. |
| Giữ hàng (reservation), trừ kho khi đặt | Inventory chỉ đọc. |
| Vòng đời đơn hàng (trạng thái, hủy, hoàn tiền) | Bảng `orders` không có trạng thái. |
| Biến thể sản phẩm, thương hiệu, mô tả sản phẩm | Model sản phẩm tối giản. |
| Giỏ hàng lưu phía server, đồng bộ đa thiết bị | Giỏ chỉ ở `localStorage`. |
| UI quản lý user, tồn kho, vận chuyển, cấu hình | Placeholder ở backoffice. |
| Hạ tầng (k8s, CI/CD, monitoring) | Thuộc bộ tài liệu hạ tầng. |

## 5. Yêu cầu nghiệp vụ (Business Requirements)

Mức ưu tiên (MoSCoW) phản ánh tầm quan trọng nghiệp vụ; trạng thái phản ánh code hiện tại. Chi tiết hành vi, luồng, acceptance criteria nằm trong tài liệu FEAT tương ứng.

### 5.1 Định danh & truy cập

| ID | Yêu cầu | Ưu tiên | Trạng thái | FEAT | Khoảng trống chính |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **BR-01** | Khách hàng và quản trị viên đăng nhập bằng tài khoản tập trung (SSO) và đăng xuất an toàn; sau đăng nhập quay về đúng trang đang xem. | Must | 🟡 | [FEAT-01](../04-features/FEAT-01-sso-login-logout.md) | Lỗi `redirect_uri` + mock token ([ISSUE-04](../07-plan/known-issues-tech-debt.md#1-bảo-mật--xác-thực)); open redirect ([ISSUE-05](../07-plan/known-issues-tech-debt.md#1-bảo-mật--xác-thực)); store in-memory ([ISSUE-06](../07-plan/known-issues-tech-debt.md#1-bảo-mật--xác-thực)). |
| **BR-02** | Người dùng mới có thể đăng ký tài khoản với username, email, mật khẩu; không trùng username/email/số điện thoại. | Must | 🟡 | [FEAT-02](../04-features/FEAT-02-user-registration.md) | Chỉ có API, không có form; role không đồng bộ sang Keycloak ([ISSUE-08](../07-plan/known-issues-tech-debt.md#1-bảo-mật--xác-thực)). |
| **BR-03** | Phiên đăng nhập được duy trì liên tục, tự làm mới khi sắp hết hạn; hết hạn thì coi như đăng xuất. | Must | 🟡 | [FEAT-03](../04-features/FEAT-03-session-refresh.md) | `session_user` không được gia hạn khi refresh ([ISSUE-69](../07-plan/known-issues-tech-debt.md)); phụ thuộc FEAT-01 tạo được phiên (ISSUE-67). |
| **BR-04** | Quản trị viên quản lý người dùng (xem, sửa, xóa) và gán/thu hồi vai trò. | Should | 🟡 | [FEAT-04](../04-features/FEAT-04-user-role-management.md) | Chỉ API; đổi mật khẩu là stub ([ISSUE-07](../07-plan/known-issues-tech-debt.md#1-bảo-mật--xác-thực)); không kiểm role ([ISSUE-01](../07-plan/known-issues-tech-debt.md#1-bảo-mật--xác-thực)). |

### 5.2 Mua sắm (Storefront)

| ID | Yêu cầu | Ưu tiên | Trạng thái | FEAT | Khoảng trống chính |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **BR-05** | Trang chủ giới thiệu banner, deal, flash sale, sản phẩm nổi bật, danh mục, thương hiệu để kích thích mua hàng. | Should | 🟡 | [FEAT-05](../04-features/FEAT-05-home-page.md) | Toàn bộ nội dung mock ([ISSUE-38](../07-plan/known-issues-tech-debt.md#5-frontend)). |
| **BR-06** | Khách duyệt sản phẩm theo danh mục, tìm theo từ khóa, sắp xếp theo mới nhất/giá, lọc theo khoảng giá, phân trang. | Must | 🟡 | [FEAT-06](../04-features/FEAT-06-product-catalog.md) | Backend không hỗ trợ phân trang/lọc ([ISSUE-48](../07-plan/known-issues-tech-debt.md#6-hợp-đồng-api-frontend--backend)); lọc giá chỉ trên trang hiện tại ([ISSUE-40](../07-plan/known-issues-tech-debt.md#5-frontend)). |
| **BR-07** | Khách xem chi tiết sản phẩm: ảnh, giá, tình trạng còn hàng, chọn số lượng, thêm giỏ / mua ngay. | Must | ✅ | [FEAT-07](../04-features/FEAT-07-product-detail.md) | Không có mô tả sản phẩm ([ISSUE-50](../07-plan/known-issues-tech-debt.md#6-hợp-đồng-api-frontend--backend)). |
| **BR-08** | Khách (kể cả chưa đăng nhập) quản lý giỏ hàng: thêm, đổi số lượng, xóa, xem tạm tính, phí ship, tổng; giỏ được giữ khi tải lại trang. | Must | ✅ | [FEAT-08](../04-features/FEAT-08-shopping-cart.md) | Giỏ chỉ ở trình duyệt. |
| **BR-09** | Khách đã đăng nhập đặt hàng: nhập thông tin giao hàng, chọn COD/Online, xem lại, xác nhận; giá được tính lại phía server. | Must | 🟡 | [FEAT-09](../04-features/FEAT-09-checkout-place-order.md) | Chỉ lưu sản phẩm đầu tiên ([ISSUE-20](../07-plan/known-issues-tech-debt.md#4-nghiệp-vụ--dữ-liệu-backend)); không lưu địa chỉ ([ISSUE-22](../07-plan/known-issues-tech-debt.md#4-nghiệp-vụ--dữ-liệu-backend)); không idempotency ([ISSUE-23](../07-plan/known-issues-tech-debt.md#4-nghiệp-vụ--dữ-liệu-backend)). |
| **BR-10** | Khách xem danh sách đơn hàng của mình. | Must | 🟡 | [FEAT-10](../04-features/FEAT-10-order-history.md) | API trả mọi đơn của mọi người ([ISSUE-49](../07-plan/known-issues-tech-debt.md#6-hợp-đồng-api-frontend--backend)); không có trạng thái ([ISSUE-41](../07-plan/known-issues-tech-debt.md#5-frontend)). |

### 5.3 Nghiệp vụ thương mại (Backend)

| ID | Yêu cầu | Ưu tiên | Trạng thái | FEAT | Khoảng trống chính |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **BR-11** | Hệ thống ghi nhận thanh toán cho đơn hàng, không cho thanh toán trùng một đơn đã trả, và thông báo khi thanh toán thành công. | Must | 🟡 | [FEAT-11](../04-features/FEAT-11-payment.md) | Sự kiện thành công phát sai ([ISSUE-18](../07-plan/known-issues-tech-debt.md#3-messaging-kafka)); topic lệch ([ISSUE-17](../07-plan/known-issues-tech-debt.md#3-messaging-kafka)). |
| **BR-12** | Hệ thống cho biết sản phẩm còn hàng hay không. | Must | 🟡 | [FEAT-12](../04-features/FEAT-12-inventory-stock-check.md) | Khóa theo tên sản phẩm, hai nguồn tồn kho ([ISSUE-26](../07-plan/known-issues-tech-debt.md#4-nghiệp-vụ--dữ-liệu-backend)). |
| **BR-13** | Khách lưu sản phẩm yêu thích. | Could | 🟡 | [FEAT-13](../04-features/FEAT-13-favourites.md) | Không có UI; route gateway hỏng ([ISSUE-13](../07-plan/known-issues-tech-debt.md#2-api-gateway-apisix)). |
| **BR-14** | Khách tìm kiếm full-text có sai chính tả, lọc theo thương hiệu/danh mục/giá, có gợi ý khi gõ. | Should | 🟡 | [FEAT-14](../04-features/FEAT-14-catalog-search-suggest.md) | Index không được nạp ([ISSUE-19](../07-plan/known-issues-tech-debt.md#3-messaging-kafka)); không có UI. |
| **BR-15** | Khách đánh giá sản phẩm (1–5 sao + nội dung); hiển thị điểm trung bình; quản trị xóa đánh giá xấu. | Should | 🟡 | [FEAT-15](../04-features/FEAT-15-ratings-reviews.md) | Cho phép đánh giá ẩn danh ([ISSUE-03](../07-plan/known-issues-tech-debt.md#1-bảo-mật--xác-thực)); không UI. |
| **BR-16** | Quản trị tạo/sửa/xóa chương trình khuyến mãi, mã coupon, giới hạn sử dụng, thời gian hiệu lực. | Should | 🟡 | [FEAT-16](../04-features/FEAT-16-promotion-management.md) | Không áp dụng vào giá ([ISSUE-30](../07-plan/known-issues-tech-debt.md#4-nghiệp-vụ--dữ-liệu-backend)). |
| **BR-17** | Quản trị cấu hình nhóm thuế và mức thuế theo quốc gia/vùng; hệ thống tra cứu % thuế. | Could | 🟡 | [FEAT-17](../04-features/FEAT-17-tax-management.md) | Bỏ qua vùng/zip ([ISSUE-31](../07-plan/known-issues-tech-debt.md#4-nghiệp-vụ--dữ-liệu-backend)); route gateway hỏng local ([ISSUE-13](../07-plan/known-issues-tech-debt.md#2-api-gateway-apisix)). |
| **BR-18** | Hệ thống lưu trữ và phục vụ file media (ảnh sản phẩm…). | Should | 🟡 | [FEAT-18](../04-features/FEAT-18-media-management.md) | Không giới hạn kích thước/loại ([ISSUE-32](../07-plan/known-issues-tech-debt.md#4-nghiệp-vụ--dữ-liệu-backend)); chưa gắn với sản phẩm. |
| **BR-19** | Hệ thống gửi email giao dịch, thông báo in-app, email chào mừng. | Should | 🟡 | [FEAT-19](../04-features/FEAT-19-notifications-email.md) | Người nhận hard-code ([ISSUE-09](../07-plan/known-issues-tech-debt.md#1-bảo-mật--xác-thực)); không lọc user ([ISSUE-33](../07-plan/known-issues-tech-debt.md#4-nghiệp-vụ--dữ-liệu-backend)). |
| **BR-20** | Hệ thống ghi nhận các dòng hàng (sản phẩm, số lượng) của mỗi đơn để phục vụ giao hàng. | Should | 🟡 | [FEAT-20](../04-features/FEAT-20-shipping-order-items.md) | Không được checkout sử dụng. |

### 5.4 Quản trị (Backoffice)

| ID | Yêu cầu | Ưu tiên | Trạng thái | FEAT | Khoảng trống chính |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **BR-21** | Chỉ người có vai trò quản trị mới truy cập được backoffice; người khác thấy trang từ chối truy cập. | Must | 🟡 | [FEAT-21](../04-features/FEAT-21-admin-access-control.md) | Chỉ chặn ở UI; API không chặn ([ISSUE-01](../07-plan/known-issues-tech-debt.md#1-bảo-mật--xác-thực)). |
| **BR-22** | Dashboard hiển thị tổng số sản phẩm, đơn, danh mục và các đơn/sản phẩm mới nhất. | Should | 🟡 | [FEAT-22](../04-features/FEAT-22-admin-dashboard.md) | Đọc `.content` trên mảng thô ([ISSUE-48](../07-plan/known-issues-tech-debt.md#6-hợp-đồng-api-frontend--backend)); không vào được do ISSUE-100; không có doanh thu. |
| **BR-23** | Quản trị tìm, xem, thêm, sửa, xóa sản phẩm. | Must | 🟡 | [FEAT-23](../04-features/FEAT-23-admin-product-mgmt.md) | Phụ thuộc [ISSUE-48](../07-plan/known-issues-tech-debt.md#6-hợp-đồng-api-frontend--backend). |
| **BR-24** | Quản trị thêm, sửa, xóa danh mục. | Must | 🟡 | [FEAT-24](../04-features/FEAT-24-admin-category-mgmt.md) | Phụ thuộc [ISSUE-48](../07-plan/known-issues-tech-debt.md#6-hợp-đồng-api-frontend--backend). |
| **BR-25** | Quản trị xem danh sách và chi tiết đơn hàng. | Must | 🟡 | [FEAT-25](../04-features/FEAT-25-admin-order-mgmt.md) | Không có trạng thái/dòng hàng; lọc phía client. |
| **BR-26** | Quản trị quản lý người dùng, tồn kho, vận chuyển, cấu hình hệ thống qua giao diện. | Could | ⛔ | [FEAT-26](../04-features/FEAT-26-admin-placeholder-modules.md) | Chỉ là placeholder. |

### 5.5 Nền tảng

| ID | Yêu cầu | Ưu tiên | Trạng thái | FEAT | Khoảng trống chính |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **BR-27** | Giao diện hỗ trợ tiếng Việt (mặc định) và tiếng Anh, chuyển đổi được, URL thân thiện. | Must | ✅ | [FEAT-27](../04-features/FEAT-27-i18n.md) | Một số link không giữ locale ([ISSUE-39](../07-plan/known-issues-tech-debt.md#5-frontend)). |
| **BR-28** | Website thân thiện SEO (sitemap, robots, canonical, hreflang) và có trang lỗi thân thiện. | Should | ✅ | [FEAT-28](../04-features/FEAT-28-seo-error-pages.md) | — |

## 6. Quy tắc nghiệp vụ toàn cục (đang được code thực thi)

| ID | Quy tắc | Nơi thực thi |
| :--- | :--- | :--- |
| RULE-01 | Đơn vị tiền tệ hiển thị là VND (định dạng `vi-VN`) ở mọi locale — `formatPrice` mặc định `currency='VND'`; nhánh USD chỉ dùng khi truyền tham số khác, hiện không nơi nào truyền. Ngày giờ định dạng `vi-VN`. | `frontend/packages/lib/src/format.ts:1-16` |
| RULE-02 | Phí vận chuyển = 0 nếu tạm tính ≥ 500.000; ngược lại 30.000. Giỏ rỗng → phí 0. | `frontend/packages/lib/src/cart.ts` |
| RULE-03 | Số lượng mỗi dòng khi checkout: 1–99. | `frontend/packages/lib/src/schemas.ts` (`checkoutLineSchema`) |
| RULE-04 | Thông tin giao hàng: họ tên ≥ 2 ký tự; SĐT 8–15 ký tự gồm số, `+`, `-`, khoảng trắng; địa chỉ ≥ 5 ký tự; thành phố, quận bắt buộc; ghi chú ≤ 500 ký tự. | `shippingSchema` |
| RULE-05 | Giá khi đặt hàng được lấy lại từ product-service phía server, không tin giá trong giỏ. | `frontend/apps/web/src/routes/checkout/checkout.remote.ts` |
| RULE-06 | COD → thanh toán `NOT_STARTED`, `isPayed=false`; ONLINE → `IN_PROGRESS`, `isPayed=true`. | `checkout.remote.ts` |
| RULE-07 | Không tạo thanh toán mới cho đơn đã có bản ghi `is_payed=true`. | payment-service |
| RULE-08 | Username, email, số điện thoại người dùng là duy nhất. | auth-service (`users` UNIQUE) |
| RULE-09 | SKU sản phẩm là duy nhất. | product-service (`products.sku` UNIQUE) |
| RULE-10 | Số sao đánh giá là số nguyên 1–5. | rating-service (Zod) |
| RULE-11 | Sản phẩm "còn hàng" khi tồn tại bản ghi inventory với `quantity > 0`. | inventory-service |
| RULE-12 | Chỉ người có role `ADMIN` hoặc `ROLE_ADMIN` vào được backoffice (UI). | `frontend/packages/lib/src/server/session.ts` (`isAdmin`), `frontend/apps/admin/src/hooks.server.ts` |
| RULE-13 | Sau đăng nhập chỉ chuyển hướng tới đường dẫn nội bộ (bắt đầu `/`, không `//`). | `frontend/packages/lib/src/server/redirect.ts` |

## 7. Giả định & ràng buộc

| ID | Nội dung |
| :--- | :--- |
| ASM-01 | Một deployment phục vụ một doanh nghiệp (không multi-tenant). |
| ASM-02 | Keycloak là nguồn sự thật cho thông tin đăng nhập và role trong token. |
| ASM-03 | Mọi request API từ bên ngoài đều đi qua APISIX; service tin tưởng gateway đã xác thực token. |
| CON-01 | Backend chỉ dùng Go và Bun/TypeScript; frontend SvelteKit chạy trên Bun. |
| CON-02 | Mỗi service sở hữu database riêng; không truy vấn chéo database. |
| CON-03 | Trình duyệt không giữ token trong JavaScript (chỉ cookie httpOnly do BFF quản lý). |

## 8. Rủi ro nghiệp vụ

| ID | Rủi ro | Hậu quả | Liên quan |
| :--- | :--- | :--- | :--- |
| RSK-01 | Khách xem được đơn hàng của người khác | Vi phạm quyền riêng tư | ISSUE-49 |
| RSK-02 | User thường gọi được API quản trị | Sửa/xóa catalog, user trái phép | ISSUE-01 |
| RSK-03 | Đơn nhiều sản phẩm chỉ lưu 1 sản phẩm, không lưu địa chỉ | Không thể giao hàng đúng | ISSUE-20, ISSUE-22 |
| RSK-04 | Bán quá số lượng tồn kho | Hủy đơn, mất uy tín | ISSUE-26 |
| RSK-05 | Email thanh toán gửi sai người | Lộ dữ liệu | ISSUE-09 |
| RSK-06 | Đăng nhập thất bại hoặc phiên giả | Không mua được hàng | ISSUE-04 |

Kế hoạch xử lý: [project-plan.md](../07-plan/project-plan.md).

## 9. Truy vết

Ma trận BR ↔ FEAT ↔ Service ↔ Route ↔ Issue: [04-features/README.md](../04-features/README.md) và [conventions-and-ids.md §3](../00-overview/conventions-and-ids.md#3-sổ-đăng-ký-business-requirement--feature).
