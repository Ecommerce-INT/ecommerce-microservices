# ADR-008 — Giỏ hàng phía client lưu localStorage

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

| Thuộc tính | Giá trị |
| :--- | :--- |
| Trạng thái | **Accepted (đang hiện diện trong code)** |
| Phạm vi | `frontend/apps/web` (cart store, checkout), `frontend/packages/lib/src/cart.ts` |
| Liên quan | [FEAT-08](../../04-features/FEAT-08-shopping-cart.md), [FEAT-09](../../04-features/FEAT-09-checkout-place-order.md), [frontend-architecture §8](../software/frontend-architecture.md#8-state-phía-client) |

## Bối cảnh

order-service có API `carts` nhưng mô hình dữ liệu (cart chứa các order một sản phẩm, `user_id` integer) không khớp với danh tính Keycloak và nhu cầu giỏ hàng nhiều dòng. Khách vãng lai cũng cần giỏ hàng trước khi đăng nhập.

## Quyết định

- Giỏ hàng là **state thuần client**: class `CartStore` (Svelte 5 runes) singleton, lưu `localStorage['ecommerce-cart']` dạng `{ state: { items }, version: 0 }` (tương thích định dạng zustand-persist của app cũ), nạp trong `onMount` của root layout.
- Mỗi dòng giữ **snapshot `Product`** (`{product, quantity}`); tổng tiền tính bằng `computeTotals` (phí ship 30 000, miễn phí từ 500 000) — hàm dùng chung với server.
- Không đồng bộ với backend: API `/api/carts` của order-service **không** được frontend dùng.
- Khi checkout, BFF nhận `[{productId, quantity}]`, **đọc lại giá** từ product-service và tính lại tổng bằng cùng `computeTotals` trước khi tạo order/payment.

## Hệ quả tích cực

- Giỏ hàng hoạt động cho khách chưa đăng nhập, không tốn request backend.
- Giá thanh toán không tin dữ liệu client (tính lại ở BFF).
- Một định nghĩa tiền tệ cho UI và server.

## Hệ quả tiêu cực

- Giỏ không theo người dùng giữa thiết bị/trình duyệt; mất khi xóa dữ liệu trình duyệt.
- Snapshot sản phẩm có thể cũ (giá/ tồn kho hiển thị khác giá khi thanh toán); không kiểm/giữ tồn kho ([ISSUE-26](../../07-plan/known-issues-tech-debt.md)).
- Sản phẩm mock ở trang chủ có `productId` không tồn tại → checkout lỗi ([ISSUE-38](../../07-plan/known-issues-tech-debt.md)); quick-add ghi `quantity: stock` vào snapshot ([ISSUE-42](../../07-plan/known-issues-tech-debt.md)).
- Order backend chỉ lưu `productId` của dòng đầu tiên và tổng tiền ([ISSUE-20](../../07-plan/known-issues-tech-debt.md)); API carts của order-service thành mã chết.
- Phí ship và ngưỡng hard-code ở frontend, không có ở backend ([ISSUE-21](../../07-plan/known-issues-tech-debt.md)).

## Bằng chứng trong code

- `frontend/apps/web/src/lib/stores/cart.svelte.ts:4,25-91` — `STORAGE_KEY`, `CartStore`, `persist()`.
- `frontend/apps/web/src/routes/+layout.svelte:15` — `onMount(() => cart.hydrate())`.
- `frontend/packages/lib/src/cart.ts:1-26` — `computeTotals`, hằng phí ship.
- `frontend/packages/lib/src/types.ts:137-141` — `CartItem` "never synced to the backend".
- `frontend/apps/web/src/routes/checkout/checkout.remote.ts:18-35` — đọc lại giá, tạo order với sản phẩm đầu tiên.
- `frontend/packages/lib/src/server/api.ts:78-110` — interface `Api` không có `carts`.

## Issue liên quan

[ISSUE-20](../../07-plan/known-issues-tech-debt.md), ISSUE-21, ISSUE-22, ISSUE-23, ISSUE-26, ISSUE-38, ISSUE-42.
