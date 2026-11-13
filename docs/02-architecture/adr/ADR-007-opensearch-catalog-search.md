# ADR-007 — OpenSearch cho tìm kiếm catalog

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

| Thuộc tính | Giá trị |
| :--- | :--- |
| Trạng thái | **Accepted (đang hiện diện trong code)** — chỉ mục chưa được nạp dữ liệu |
| Phạm vi | search-service, OpenSearch 3.9 |
| Liên quan | [search-service](../../03-services/search-service.md), [FEAT-14](../../04-features/FEAT-14-catalog-search-suggest.md), [ADR-005](ADR-005-kafka-async-events.md) |

## Bối cảnh

Tìm kiếm full-text có fuzzy, gợi ý theo tiền tố và lọc theo giá/brand/category không phù hợp với truy vấn `LIKE` trên Postgres của product-service.

## Quyết định

- search-service (Go) dùng `opensearch-go/v5` truy vấn index **`product`**:
  - `GET /storefront/catalog-search`: `bool` query với `multi_match` (`name^3`, `brand^2`, `categories`, `fuzziness: AUTO`), `term` filter brand/category/attribute, `range` giá, phân trang `page/size`.
  - `GET /storefront/search_suggest`: `match_phrase_prefix` trên `name`.
- Index được nạp **bất đồng bộ** từ Kafka (sự kiện CDC sản phẩm), không gọi product-service.
- search-service không có DB riêng; OpenSearch chạy single-node, tắt security plugin trong compose.

## Hệ quả tích cực

- Tìm kiếm có xếp hạng, fuzzy, gợi ý; tách tải đọc khỏi DB giao dịch.
- Route search công khai qua gateway, không cần token.

## Hệ quả tiêu cực

- Chỉ mục **luôn rỗng** vì không có producer/connector CDC và schema field lệch (`id/name/slug/price` vs `product_id/product_title/price_unit`) ([ISSUE-19](../../07-plan/known-issues-tech-debt.md)).
- `sortType` bị bỏ qua, aggregations rỗng, lỗi OpenSearch trả trang rỗng thay vì lỗi ([ISSUE-34](../../07-plan/known-issues-tech-debt.md)).
- Frontend storefront hiện **không gọi** search-service (catalog dùng `/api/products?search=`), nên quyết định chưa có giá trị người dùng.
- Không có mapping/template index trong repo (dựa vào dynamic mapping).

## Bằng chứng trong code

- `services/search-service/main.go:27-49` — cấu hình client (`OPENSEARCH_URL`, `DiscoverNodesOnStart=false`).
- `services/search-service/internal/service/search_service.go:44-60,112` — query và index `product`; `:194-215` — autocomplete; `:241-268` — index/delete document.
- `services/search-service/internal/handler/search_handler.go:30-33,50-84` — endpoint & tham số.
- `services/search-service/internal/model/search.go:5-42` — `ProductDoc`, `ProductListGetVm`.
- `docker-compose.yml:69-97` — OpenSearch single-node, `DISABLE_SECURITY_PLUGIN=true`.
- `frontend/packages/lib/src/server/api.ts:129` — catalog BFF gọi `/api/products`, không phải search.

## Issue liên quan

[ISSUE-19](../../07-plan/known-issues-tech-debt.md), ISSUE-34, ISSUE-48.
