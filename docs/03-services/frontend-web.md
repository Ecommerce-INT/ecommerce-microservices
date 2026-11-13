# frontend-web (`@ecommerce/web`) — Storefront SvelteKit

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

| Thuộc tính | Giá trị |
| :--- | :--- |
| Ngôn ngữ / stack | TypeScript · SvelteKit 3 (remote functions, experimental) · Svelte 5 runes · Tailwind CSS 4 · Paraglide JS 2 · Bun runtime (`@sveltejs/adapter-bun`) |
| Port | `3000` (dev/preview: `frontend/apps/web/vite.config.ts:8-9`; production: `frontend/apps/web/Dockerfile:24,27`; compose `3000:3000`: `docker-compose.yml:416-417`) |
| Database | Không có. Giỏ hàng lưu ở `localStorage` trình duyệt; phiên lưu ở cookie httpOnly |
| Đường dẫn code | `frontend/apps/web` (+ helper dùng chung `frontend/packages/lib`, UI kit `frontend/packages/ui`) |
| Kiểu kiến trúc | BFF (Backend-for-Frontend): trình duyệt **không** gọi gateway, trừ bước redirect SSO; mọi API call chạy server-side qua remote function / load function |
| Upstream | API Gateway APISIX (`API_INTERNAL_URL`) → auth-service, product-service, order-service, payment-service |
| Image | `ghcr.io/ecommerce-int/frontend-web:latest` (`docker-compose.yml:413`) |
| Tính năng phục vụ | [FEAT-01](../04-features/FEAT-01-sso-login-logout.md) · [FEAT-03](../04-features/FEAT-03-session-refresh.md) · [FEAT-05](../04-features/FEAT-05-home-page.md) · [FEAT-06](../04-features/FEAT-06-product-catalog.md) · [FEAT-07](../04-features/FEAT-07-product-detail.md) · [FEAT-08](../04-features/FEAT-08-shopping-cart.md) · [FEAT-09](../04-features/FEAT-09-checkout-place-order.md) · [FEAT-10](../04-features/FEAT-10-order-history.md) · [FEAT-27](../04-features/FEAT-27-i18n.md) · [FEAT-28](../04-features/FEAT-28-seo-error-pages.md) |
| Issue liên quan | ISSUE-04, 05, 11, 100, 101, 20, 21, 22, 23, 38–44, 46, 47, 48, 49, 50, **66**, **67** |

---

## 1. Tổng quan & trách nhiệm

`frontend/apps/web` là ứng dụng storefront dành cho khách hàng (thương hiệu hiển thị "EzBuy" — `frontend/apps/web/src/lib/components/header.svelte:85`). Trách nhiệm theo code hiện tại:

| Trách nhiệm | Thực hiện ở | Ghi chú |
| :--- | :--- | :--- |
| Trang chủ marketing | `src/routes/+page.svelte` | **Toàn bộ dữ liệu là mock** từ `src/lib/data/home-mock.ts` (ISSUE-38) |
| Danh mục sản phẩm (tìm kiếm, lọc danh mục, sắp xếp, phân trang, lọc giá client-side) | `src/routes/products/*` | Phụ thuộc hợp đồng Page của backend (ISSUE-48) |
| Chi tiết sản phẩm, thêm vào giỏ, "Mua ngay" | `src/routes/products/[id]/*` | |
| Giỏ hàng client-side | `src/lib/stores/cart.svelte.ts`, `src/routes/cart/+page.svelte` | Không đồng bộ backend |
| Checkout 3 bước, tạo order + payment | `src/routes/checkout/*` | Yêu cầu đăng nhập |
| Lịch sử đơn hàng | `src/routes/orders/*` | Yêu cầu đăng nhập; backend trả mọi đơn (ISSUE-49) |
| Đăng nhập SSO (ticket exchange), đăng xuất | `src/routes/login`, `src/routes/auth/callback`, `src/lib/remote/auth.remote.ts` | |
| Duy trì & làm mới phiên (refresh token) | `src/hooks.server.ts:32-51` | |
| Đa ngôn ngữ vi/en (URL `/` = vi, `/en/...` = en) | Paraglide: `vite.config.ts:26-32`, `src/hooks.ts`, `src/hooks.server.ts:14-25` | |
| SEO: canonical, hreflang, sitemap, robots, trang lỗi | `src/routes/+layout.svelte:26-35`, `src/routes/sitemap.xml/+server.ts`, `static/robots.txt`, `src/routes/+error.svelte` | |

Ứng dụng **không** dùng các tính năng favourites, rating, search-service, promotion, inventory dù type/API client có khai báo một phần (ISSUE-43).

---

## 2. Tech stack & thư viện

Tất cả khai báo trong `devDependencies` (`frontend/apps/web/package.json:16-35`):

| Thư viện | Version (range) | Vai trò |
| :--- | :--- | :--- |
| `@sveltejs/kit` | `^3.0.0` | Framework; dùng `experimental.remoteFunctions` (`vite.config.ts:13-15`) |
| `svelte` | `^5.57.1` | Runes mode bắt buộc cho code dự án (`vite.config.ts:21-22`), `compilerOptions.experimental.async: true` (`vite.config.ts:16-19`) cho phép `{#await}` trên remote query khi SSR |
| `@sveltejs/vite-plugin-svelte` | `^7.2.0` | Plugin Vite |
| `@sveltejs/adapter-bun` | `^1.0.0` | Build ra server Bun (`vite.config.ts:24`); cấu hình adapter đặt trong `vite.config.ts`, **không có** `svelte.config.js` |
| `vite` | `^8.3.0` | Bundler/dev server |
| `@inlang/paraglide-js` | `^2.18.2` | i18n compile-time |
| `tailwindcss`, `@tailwindcss/vite` | `^4.3.0` | Styling |
| `tw-animate-css` | `^1.4.0` | Animation utility |
| `@fontsource-variable/geist` | `^5.3.0` | Font |
| `@lucide/svelte` | `^1.55.0` | Icon |
| `svelte-sonner` | `^1.2.0` | Toast (`toast.success/error`) |
| `zod` | `^4.6.5` | Validation schema cho remote function (import dưới alias `v`) |
| `typescript` | `^6.0.3` | |
| `svelte-check` | `^4.6.0` | Typecheck |
| `@types/bun` | `latest` | Type Bun (không pin version) |
| `@ecommerce/lib` | `workspace:*` | API client, session, schema, cart math, types (`frontend/packages/lib`) |
| `@ecommerce/ui` | `workspace:*` | Component UI dùng chung (Button, Input, Badge, Skeleton, EmptyState, DropdownMenu, Separator, Label, Textarea, Toaster) |

Scripts (`package.json:6-15`): `dev`/`build`/`preview` chạy Vite **dưới Bun** (`bun run --bun vite ...`); `check` compile Paraglide (`--strategy url cookie baseLocale`) rồi `svelte-kit sync && svelte-check`; `lint` = prettier + eslint. Không có script `test` (xem mục 14).

Subpath imports: `#lib` → `./src/lib/index.js`, `#lib/*` → `./src/lib/*` (`package.json:36-39`, `tsconfig.json:6-8`).

Docker (`Dockerfile`): build context là `frontend/` (`Dockerfile:3-4`); stage build `oven/bun:1`, `bun install --frozen-lockfile`, `bun run --filter @ecommerce/web build` (`Dockerfile:6-19`); stage runtime chỉ copy `apps/web/build`, chạy user `bun`, `CMD ["bun", "./build"]` (`Dockerfile:21-28`).

---

## 3. Cấu hình (biến môi trường)

| Biến | Default | Ý nghĩa | Nguồn |
| :--- | :--- | :--- | :--- |
| `API_INTERNAL_URL` | `http://localhost:9080` | Origin gateway APISIX dùng cho **mọi** lời gọi server-side; dấu `/` cuối bị cắt | `frontend/packages/lib/src/server/env.ts:11-14` |
| `API_PUBLIC_URL` | = `API_INTERNAL_URL` | Origin gateway mà **trình duyệt** truy cập được; chỉ dùng để dựng link SSO `/api/v1/auth/login` | `frontend/packages/lib/src/server/env.ts:24-27`, `src/routes/login/+page.server.ts:24` |
| `PORT` | `3000` | Port server Bun build | `Dockerfile:24`, `.env.example:10-11` |
| `HOST` | (mặc định của adapter-bun) | Interface lắng nghe | `.env.example:12` |
| `ORIGIN` | — | Origin công khai; cần khi chạy HTTP thuần/proxy để SvelteKit kiểm tra CSRF và dựng `url.origin` đúng (ảnh hưởng `redirect_uri` SSO và sitemap) | `.env.example:14-15` |
| `PROTOCOL_HEADER` | — | Header proxy chứa scheme (vd `x-forwarded-proto`) | `.env.example:16` |
| `NODE_ENV` | `production` (Docker) | | `Dockerfile:23` |

Giá trị trong docker compose (`docker-compose.yml:418-426`): `NODE_ENV=production`, `PORT=3000`, `API_INTERNAL_URL=http://apisix:9080`, `API_PUBLIC_URL=http://localhost:9080`, `ORIGIN=http://localhost:3000`; `depends_on: apisix` (`docker-compose.yml:427-429`). Ứng dụng web không dùng secret nào.

Cấu hình i18n tĩnh: `project.inlang/settings.json` (`baseLocale: "vi"`, `locales: ["vi","en"]`, `pathPattern: ./messages/{locale}.json` — dòng 7-11). Thư mục sinh `src/lib/paraglide` bị gitignore (`.gitignore`).

---

## 4. Phụ thuộc

### 4.1 Gateway & service backend

Mọi lời gọi đi qua `createApi()` (`frontend/packages/lib/src/server/api.ts:116-161`), gắn `accept: application/json`, `content-type: application/json` khi có body, và `authorization: Bearer <access_token>` khi có token (`api.ts:45-48`). Response dạng envelope `{ data }` được unwrap, payload thô được trả nguyên (`api.ts:71-75`); lỗi HTTP ném `ApiError(status, message, body)` (`api.ts:62-69`).

| Method | Path gateway | Hàm client | Dùng ở web | Route APISIX → service | Auth tại gateway |
| :--- | :--- | :--- | :--- | :--- | :--- |
| GET (browser) | `/api/v1/auth/login?redirect_uri=` | (link) | `login/+page.server.ts:23-24` | `auth-public` → auth-service `/auth/api/v1/auth/login` (`deploy/apisix/apisix.yaml:63-70`, `services/auth-service/internal/handler/auth_handler.go:66`) | Public |
| GET | `/api/v1/auth/session?ticket=` | `auth.exchangeTicket` (không token) | `auth/callback/+page.server.ts:21` | auth-service `Session` (`auth_handler.go:68`) | Public |
| GET | `/api/v1/users/me` | `auth.profile` | `auth/callback/+page.server.ts:28` | auth-service `GetCurrentUser` (`auth_handler.go:74`) | openid-connect bearer_only (`apisix.yaml:161-177`) |
| POST | `/api/v1/auth/refresh` body `{refreshToken}` | `auth.refresh` | `hooks.server.ts:39` | auth-service `Refresh` (`auth_handler.go:69`) | Public |
| POST | `/api/v1/auth/logout` body `{refreshToken}` | `auth.logout` | `lib/remote/auth.remote.ts:12` | auth-service `Logout` (`auth_handler.go:70`) | Public |
| GET | `/api/products?page&size&sort&categoryId&search` | `products.list` | `products/catalog.remote.ts:13` | `product-catalog` → product-service (`apisix.yaml:71-82`) | Public |
| GET | `/api/products/{id}` | `products.get` | `products/[id]/product.remote.ts:5`, `checkout/checkout.remote.ts:19` | product-service | Public |
| GET | `/api/categories?page=0&size=50` | `categories.list` | `products/catalog.remote.ts:16-19` | product-service | Public |
| GET | `/api/orders?page&size` | `orders.list` | `orders/orders.remote.ts:6-8` | `order-orders` → order-service (`apisix.yaml:197-213`) | bearer_only |
| POST | `/api/orders` | `orders.create` | `checkout/checkout.remote.ts:30-35` | order-service | bearer_only |
| POST | `/api/payments` | `payments.create` | `checkout/checkout.remote.ts:37-41` | `payment` → payment-service (`apisix.yaml:215-231`) | bearer_only |

Các hàm client khác (`products.create/update/remove`, `categories.*` ghi, `orders.get`, `inventory.checkStock`) chỉ được admin dùng hoặc không dùng (ISSUE-43, ISSUE-51).

**Hợp đồng dữ liệu kỳ vọng** (`frontend/packages/lib/src/types.ts`): `PaginatedResponse<T>` (`content`, `totalElements`, `totalPages`… — `types.ts:12-20`), `TokenResponse` snake_case (`types.ts:22-29`), `UserResponse` (`types.ts:36-45`), `Product` (`types.ts:65-74`, có `description` — ISSUE-50), `Order` (`types.ts:90-97`), `CreateOrderPayload` (`types.ts:158-163`), `CreatePaymentPayload` (`types.ts:165-169`). Backend thực tế trả mảng thô cho list → ISSUE-48.

### 4.2 Package nội bộ

| Package / module | Dùng gì |
| :--- | :--- |
| `@ecommerce/lib/server` (`frontend/packages/lib/src/server/index.ts`) | `createApi`, `apiPublicUrl`, `safeRedirect`, `session`, `setSession`, `clearSession`, `accessToken`, `refreshToken`, `shouldRefresh`, `toSessionUser` |
| `@ecommerce/lib/schemas` | `pageQuerySchema`, `productQuerySchema`, `shippingSchema`, `checkoutSchema` |
| `@ecommerce/lib/cart` | `computeTotals`, `FREE_SHIPPING_THRESHOLD` |
| `@ecommerce/lib/format` | `formatPrice`, `formatDate` |
| `@ecommerce/lib/types` | `Product`, `CartItem`, `SessionUser` |
| `@ecommerce/ui/*` | Component UI |

---

## 5. Cấu trúc code

```text
frontend/apps/web/
├── Dockerfile                 # build 2 stage Bun, chạy `bun ./build` port 3000
├── .env.example               # API_INTERNAL_URL, API_PUBLIC_URL, PORT, HOST, ORIGIN, PROTOCOL_HEADER
├── package.json               # deps, scripts, subpath imports #lib
├── vite.config.ts             # port 3000, sveltekit(remoteFunctions, async, runes, adapter-bun), paraglide
├── tsconfig.json              # extends $app/tsconfig, strict
├── eslint.config.js / prettier.config.js
├── project.inlang/settings.json   # baseLocale vi, locales [vi, en]
├── messages/{vi,en}.json      # 306 key mỗi file
├── static/robots.txt
└── src/
    ├── app.html               # <html lang/dir> thay bằng placeholder Paraglide
    ├── app.d.ts               # App.Locals { user: SessionUser | null }
    ├── hooks.server.ts        # handle = sequence(paraglide, session refresh)
    ├── hooks.ts               # reroute: bỏ prefix locale khỏi URL
    ├── lib/
    │   ├── index.ts           # rỗng (placeholder alias #lib)
    │   ├── assets/favicon.svg
    │   ├── components/
    │   │   ├── header.svelte        # top bar, đổi ngôn ngữ, search, menu user/đăng xuất, badge giỏ, nav danh mục cứng
    │   │   ├── footer.svelte        # cột danh mục cứng, hỗ trợ, liên hệ
    │   │   └── product-card.svelte  # thẻ sản phẩm + nút thêm nhanh vào giỏ
    │   ├── data/home-mock.ts  # dữ liệu mock trang chủ (banner, deal, sản phẩm id 101–405, brand)
    │   ├── remote/auth.remote.ts    # command logout
    │   ├── server/api.ts      # api(): createApi gắn access token của request hiện tại
    │   └── stores/cart.svelte.ts    # CartStore (runes) + localStorage
    └── routes/
        ├── +layout.server.ts  # trả { user: locals.user }
        ├── +layout.svelte     # Header/Footer/Toaster, cart.hydrate(), canonical + hreflang
        ├── +page.svelte       # trang chủ (mock)
        ├── +error.svelte      # trang lỗi 404 / lỗi chung
        ├── layout.css         # Tailwind + theme
        ├── sitemap.xml/+server.ts
        ├── login/{+page.server.ts,+page.svelte}
        ├── auth/callback/{+page.server.ts,+page.svelte}
        ├── products/{+page.svelte,catalog.remote.ts}
        ├── products/[id]/{+page.svelte,product.remote.ts}
        ├── cart/+page.svelte
        ├── checkout/{+page.server.ts,+page.svelte,checkout.remote.ts}
        └── orders/{+page.server.ts,+page.svelte,orders.remote.ts}
```

`src/lib/server/api.ts:10-13`: helper `api()` lấy `cookies` từ `getRequestEvent()` và tạo client với `accessToken(cookies)` — chỉ dùng được trong request context (remote function, load, hook).

---

## 6. Routes

Bảng tóm tắt (path vi; mọi route có biến thể `/en/...` nhờ reroute — mục 9.2):

| Route | Server load / endpoint | Guard | Remote functions | Gateway API | FEAT |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `/` | layout load | — | — | — (mock) | FEAT-05 |
| `/products` | layout load | — | `getProducts` (query), `getCategories` (query) | `GET /api/products`, `GET /api/categories` | FEAT-06 |
| `/products/[id]` | layout load | — | `getProduct` (query) | `GET /api/products/{id}` | FEAT-07 |
| `/cart` | layout load | — | — | — | FEAT-08 |
| `/checkout` | `checkout/+page.server.ts` | Phải đăng nhập | `placeOrder` (command) | `GET /api/products/{id}`, `POST /api/orders`, `POST /api/payments` | FEAT-09 |
| `/orders` | `orders/+page.server.ts` | Phải đăng nhập | `getOrders` (query) | `GET /api/orders` | FEAT-10 |
| `/login` | `login/+page.server.ts` | Đã đăng nhập → redirect | — | (link SSO) | FEAT-01 |
| `/auth/callback` | `auth/callback/+page.server.ts` | — | — | `GET /api/v1/auth/session`, `GET /api/v1/users/me` | FEAT-01 |
| `/sitemap.xml` | `sitemap.xml/+server.ts` (GET) | — | — | — | FEAT-28 |
| (mọi route) | `hooks.server.ts` | — | `logout` (command, từ header) | `POST /api/v1/auth/refresh`, `POST /api/v1/auth/logout` | FEAT-01, FEAT-03 |

Không có remote `form` nào trong app web; chỉ có `query` và `command`.

### 6.1 Layout gốc

- `+layout.server.ts:3-5`: trả `{ user: locals.user }` cho mọi trang.
- `+layout.svelte`:
  - `onMount(() => cart.hydrate())` (dòng 15) — giỏ chỉ đọc `localStorage` sau khi mount.
  - `<head>`: favicon, `<link rel="canonical" href={page.url.href}>` (dòng 28, gồm cả query string), `hreflang` cho từng locale + `x-default` trỏ bản `vi` (dòng 17-23, 29-34).
  - Khung: `<Header user={data.user}>`, `<main>`, `<Footer>`, `<Toaster richColors position="top-center">` (dòng 37-45).

### 6.2 `/` — Trang chủ (FEAT-05)

- Không có load function; toàn bộ dữ liệu import từ `#lib/data/home-mock` (`+page.svelte:8-22`): `heroSlides`, `heroSideBanners`, `quickDeals`, `categoryIcons`, `brands`, `heroSidebar`, `flashSaleProducts` (id 101–106), `hotProducts` (201–210), `phoneProducts` (301–305), `laptopProducts` (401–405), `phoneBrands`, `laptopBrands`, `nextFlashSaleDeadline()` (`home-mock.ts:58-450`).
- Các khối UI: sidebar danh mục + hero slides + side banners (`+page.svelte:104-151`), quick deals (153-163), lưới 10 icon danh mục (165-177), flash sale có countdown `setInterval` 1s tới `nextFlashSaleDeadline()` và thanh % đã bán (84-97, 179-215), "Sản phẩm hot" (217-227), 2 section điện thoại / laptop link `/products?categoryId=1|2` (64-81, 229-261), thương hiệu nổi bật (263-275), USP (58-63, 277-288).
- Nhãn danh mục dùng `m.category_c1..c10` cố định (`+page.svelte:34-57`) → ISSUE-44.
- `ProductCard` với sản phẩm mock cho phép thêm vào giỏ `productId` không tồn tại → `placeOrder` lỗi khi fetch lại giá (ISSUE-38).
- SEO: `<title>` = `home_metaTitle`, `<meta name="description">` = `home_metaDesc` (dòng 99-102).

### 6.3 `/products` — Danh mục (FEAT-06)

**Remote functions** (`src/routes/products/catalog.remote.ts`):

| Tên | Loại | Schema | Gọi |
| :--- | :--- | :--- | :--- |
| `getProducts` | `query` | `productQuerySchema` (`frontend/packages/lib/src/schemas.ts:37-41`): `page` coerce int ≥0 default 0; `size` coerce int 1–60 default 12; `sort` ∈ `productId,asc` · `productId,desc` · `priceUnit,asc` · `priceUnit,desc` (optional); `categoryId` chuỗi số / int dương / `''` → number \| undefined (`schemas.ts:21-30`); `search` trim, min 1 (optional) | `api().products.list(input)` → `GET /api/products?…` (dòng 13) |
| `getCategories` | `query` (không tham số) | — | `api().categories.list({ page: 0, size: 50 })`, trả `result.content` (dòng 16-19) |

**UI** (`src/routes/products/+page.svelte`):
- State đọc từ URL: `page`, `search`, `categoryId`, `sort` (chỉ nhận 3 giá trị trong `SORTS`) → `args = { page, size: 12, search, categoryId, sort }` (dòng 15-41). Dropdown chỉ có 3 lựa chọn: mới nhất (`productId,desc`), giá tăng, giá giảm (15-19).
- Sidebar danh mục từ `getCategories()` với skeleton / `—` khi lỗi (96-114).
- Lọc khoảng giá 5 mức (0–100k, 100k–500k, 500k–1tr, 1tr–5tr, ≥5tr) là state client `priceRange`, chỉ lọc `result.content` của **trang hiện tại** (22-28, 57-63, 120-131) → ISSUE-40. Không ghi vào URL.
- Ô tìm kiếm là `<form method="get">` (147-157); đổi sort dùng `goto(urlFor(...))` (162-171). `urlFor` giữ query hiện tại, xoá key rỗng, bọc `localizeHref` (45-53).
- Hiển thị tổng `products_count` theo `totalElements` (187-189); phân trang cửa sổ 5 trang + Prev/Next khi `totalPages > 1` (65-69, 208-239). Lỗi → `productgrid_loadError` + nút thử lại (241-250).
- Vì backend trả mảng thô, `result.content` là `undefined` → luôn rơi vào "không tìm thấy" (ISSUE-48; code có `?? []` ở dòng 186).
- `<title>` = `products_metaTitle`.

### 6.4 `/products/[id]` — Chi tiết sản phẩm (FEAT-07)

- Remote `getProduct = query(v.coerce.number().int().positive(), id => api().products.get(id))` (`product.remote.ts:5`) → `GET /api/products/{id}`.
- UI (`products/[id]/+page.svelte`): breadcrumb (26-34); skeleton khi chờ; ảnh hoặc 📦; SKU nếu có; giá + "đã gồm VAT"; badge còn hàng (`quantity`) / hết hàng (65-86).
- Bộ chọn số lượng: tối thiểu 1, tăng tối đa `product.quantity || 99` (88-109).
- "Thêm vào giỏ": `cart.add(product, quantity)` + toast (112-122). "Mua ngay": thêm vào giỏ rồi điều hướng `/checkout` (123-133). Cả hai disable khi `quantity <= 0`.
- Mô tả chỉ hiện khi `product.description` có giá trị (151-156) — backend không có field này (ISSUE-50).
- Lỗi (gồm 404 từ gateway) → khối `productdetail_notFound` + nút về danh sách (159-168); trang vẫn trả HTTP 200.
- Trang **không** đặt `<title>`/meta riêng (key `productdetail_metaTitle` tồn tại nhưng không được dùng).

### 6.5 `/cart` — Giỏ hàng (FEAT-08)

`src/routes/cart/+page.svelte` — thuần client, không gọi API:
- Chưa hydrate → 3 skeleton (27-32); giỏ rỗng → empty state + link `/products` (33-43).
- Mỗi dòng: ảnh/tiêu đề link `localizeHref('/products/{id}')`, đơn giá, nút −/+ gọi `cart.updateQuantity` (về 0 thì xoá), nút xoá `cart.remove`, thành tiền (54-119). Nút "Xoá tất cả" `cart.clear()` (47-52). Nút + **không** có giới hạn trên.
- Tóm tắt: tạm tính, phí ship (miễn phí nếu `shippingFee === 0`), gợi ý "mua thêm X để miễn phí" = `FREE_SHIPPING_THRESHOLD − totalPrice` (12, 122-148), tổng `grandTotal`, nút tới `/checkout` (152-159).
- Quy tắc tiền (`frontend/packages/lib/src/cart.ts:1-26`): `FREE_SHIPPING_THRESHOLD = 500_000`, `SHIPPING_FEE = 30_000`; phí = 0 nếu giỏ rỗng hoặc subtotal ≥ ngưỡng.

### 6.6 `/checkout` — Checkout (FEAT-09)

**Guard** (`checkout/+page.server.ts:5-11`): nếu `!locals.user` → `redirect(303, localizeHref('/login') + '?redirect=' + encodeURIComponent(<path đã bỏ locale + query>))`. Trả `{ user }`.

**Remote** `placeOrder = command(checkoutSchema, …)` (`checkout/checkout.remote.ts:13-44`):

| Bước | Hành vi | Dòng |
| :--- | :--- | :--- |
| Validate | `checkoutSchema` (`schemas.ts:43-64`): `shipping` = `fullName` trim ≥2; `phone` regex `^[0-9+\-\s]{8,15}$`; `address` ≥5; `city` ≥1; `district` ≥1; `note` ≤500 optional. `paymentMethod` ∈ `COD`/`ONLINE`. `items` mảng ≥1 phần tử `{ productId: int>0, quantity: int 1–99 }` | — |
| Kiểm phiên | Không có `access_token` → `error(401, 'Your session has expired. Please sign in again.')` | 14-16 |
| Lấy giá thật | `Promise.all(items.map(api.products.get))` — giá không lấy từ client | 18-19 |
| Tính tổng | `computeTotals` với `priceUnit` từ catalog | 21-26 |
| Tạo order | `POST /api/orders` `{ orderDate: now ISO, orderDesc: shipping.note \|\| 'EzBuy order', orderFee: totals.total, productId: products[0].productId }` — chỉ dòng đầu tiên được ghi (ISSUE-20); địa chỉ giao hàng không được gửi đi đâu (ISSUE-22) | 28-35 |
| Tạo payment | `POST /api/payments` `{ isPayed: paymentMethod !== 'COD', paymentStatus: COD ? 'NOT_STARTED' : 'IN_PROGRESS', orderId }` — `ONLINE` được đánh dấu `isPayed: true` dù không có cổng thanh toán nào | 37-41 |
| Trả về | `{ orderId, total }` | 43 |

Không có idempotency key, không bù trừ nếu payment lỗi sau khi order đã tạo (ISSUE-23); không kiểm tồn kho.

**UI wizard** (`checkout/+page.svelte`):
- State: `step` 1..4, `paymentMethod` mặc định `COD`, `shipping` (prefill `fullName` từ `data.user`) (20-32).
- Bước 1 – Giao hàng: form 6 trường; validate client bằng `shippingSchema.safeParse` (34); "Tiếp tục" chỉ sang bước 2 khi hợp lệ, ngược lại hiện `checkout_validationErr` (41-44, 110-157). Lưu ý: "Thành phố" là input text tự do dù placeholder là "Chọn…".
- Bước 2 – Thanh toán: chọn COD hoặc "Chuyển khoản" (`ONLINE`) (158-196).
- Bước 3 – Xác nhận: hiển thị lại thông tin, nút "Đặt hàng" gọi `submit()` (197-230). `submit` chặn double-click bằng cờ `placing` phía client (46-66).
- Thành công: lưu `orderId`, `cart.clear()`, `step = 4` → màn hình thành công với mã đơn, link `/orders` và `/products` (58-60, 231-245). Lỗi bất kỳ (validation 400, 401, lỗi gateway) → chỉ một toast chung `checkout_orderFailed` (61-62).
- Giỏ rỗng (sau hydrate, trước bước 4) → empty state (75-85). Sidebar tóm tắt dùng giá đang lưu trong giỏ (`cart.totalPrice/shippingFee/grandTotal`) — có thể khác tổng server tính nếu giá đã đổi (249-288).

### 6.7 `/orders` — Lịch sử đơn hàng (FEAT-10)

- Guard giống checkout (`orders/+page.server.ts:5-11`).
- Remote `getOrders = query(pageQuerySchema, ({page,size}) => api().orders.list({page,size}))` (`orders/orders.remote.ts:6-8`) → `GET /api/orders?page&size`. Comment "the gateway scopes them by the JWT" (dòng 5) **không đúng** — ISSUE-49.
- UI (`orders/+page.svelte`): `page` từ URL, `size: 10` (13-14); skeleton; empty state; mỗi đơn hiển thị badge **luôn** là `orders_processing` (ISSUE-41), `#orderId`, `formatDate(orderDate)`, `formatPrice(orderFee)`, `orderDesc` (40-54); phân trang Prev/Next khi `totalPages > 1` (56-78); lỗi → `common_error` + thử lại (80-86). Vì `orders.list` nhận mảng thô từ order-service, `result.content` là `undefined` → lỗi render rơi vào nhánh `{:catch}` (ISSUE-48).
- Lưu ý: bản thân remote query không kiểm `locals.user`; khi gọi không có token, gateway trả 401.

### 6.8 `/login` — Đăng nhập SSO (FEAT-01)

`login/+page.server.ts`:
1. `target = safeRedirect(url.searchParams.get('redirect'))` — chỉ chấp nhận path bắt đầu bằng `/` và không bắt đầu bằng `//`, ngược lại `/` (`frontend/packages/lib/src/server/redirect.ts:2-5`).
2. Đã đăng nhập (`locals.user`) → `redirect(303, localizeHref(target))` (dòng 12).
3. Ghi cookie `auth_redirect = target` (httpOnly, sameSite lax, path `/`, maxAge 600s) (dòng 16-21) — vì gateway chỉ cho redirect cố định về `/auth/callback`.
4. Dựng `ssoUrl = ${API_PUBLIC_URL}/api/v1/auth/login?redirect_uri=${encodeURIComponent(origin + '/auth/callback')}` (dòng 23-26). `redirect_uri` luôn **không** có prefix locale.

`login/+page.svelte`: hai nút (đăng nhập và "đăng ký qua SSO") cùng trỏ `data.ssoUrl` (19-21, 38-40); không có form username/password. Đăng ký thực chất là luồng Keycloak.

### 6.9 `/auth/callback` — Đổi ticket lấy phiên (FEAT-01)

`auth/callback/+page.server.ts`:
1. Đọc và xoá cookie `auth_redirect`, chuẩn hoá bằng `safeRedirect` (11-12).
2. Thiếu `ticket` hoặc có `error` → `clearSession`, trả `{ failed: true }` (14-18).
3. `createApi().auth.exchangeTicket(ticket)` → `GET /api/v1/auth/session?ticket=` **không** kèm bearer (21; `api.ts:154-155`).
4. Lấy profile bằng access token mới: `GET /api/v1/users/me` → `toSessionUser` (`session.ts:13-21`: `id, username, fullName, email, roles[].name`). Lỗi profile → `user = null` (26-31).
5. `setSession(cookies, tokens, user)` (33). Mọi lỗi exchange → `clearSession` + `{ failed: true }` (34-37).
6. Thành công → `redirect(303, localizeHref(target))` (39).

Lưu ý: comment dòng 23-25 nói "giữ phiên không có role" khi profile lỗi, nhưng `setSession(..., null)` **xoá** cookie `session_user` (`session.ts:44-45`) và `session()` trả `null` khi thiếu cookie này (`session.ts:70-72`) → người dùng bị coi là chưa đăng nhập dù có token → **ISSUE-67**. Trường hợp này xảy ra **thường xuyên** chứ không chỉ khi lỗi tạm thời: user tạo trên Keycloak không có trong DB `authservice` nên `/users/me` trả 404 (**ISSUE-101**). Khi profile trả về được, response `model.UserResponse` dùng key `fullname` (chữ thường) và không có `roles` dạng `{id,name}` (`services/auth-service/internal/model/user.go:16-22`) → `toSessionUser` cho ra `fullName` undefined và `roles: []` (**ISSUE-100**); ở web, header rơi về hiển thị `username`.

`auth/callback/+page.svelte`: `failed` → icon lỗi, `auth_callbackErrorTitle/Text`, nút thử lại về `/login`; ngược lại spinner `auth_callbackLoading` (thực tế ít khi thấy vì load redirect trước khi render).

### 6.10 Đăng xuất (command `logout`)

`src/lib/remote/auth.remote.ts:6-19`: nếu có refresh token → `POST /api/v1/auth/logout {refreshToken}` (best-effort, nuốt lỗi), rồi luôn `clearSession`. Header gọi `await logout(); await invalidateAll(); await goto(localizeHref('/'))` với cờ `loggingOut` chống gọi lặp (`header.svelte:35-45`). Giỏ hàng **không** bị xoá khi đăng xuất.

### 6.11 `/sitemap.xml`, `+error.svelte`, `robots.txt`

Xem mục 12.

---

## 7. Components

| Component | File | Props | Hành vi chính |
| :--- | :--- | :--- | :--- |
| Header | `src/lib/components/header.svelte` | `user?: SessionUser \| null` | Top bar: thông điệp ship, link theo dõi đơn (`/orders`), hotline, **bộ chuyển ngôn ngữ** dùng `localizeHref(page.url.pathname, { locale })` + `data-sveltekit-reload` (64-76; query string bị bỏ khi đổi ngôn ngữ). Logo "EzBuy". Form search `GET` → `/products?search=` (88-98). Đã đăng nhập → DropdownMenu: tên (`fullName \|\| username`), email, "Đơn hàng của tôi", "Giỏ hàng", "Đăng xuất" (101-132); chưa → nút "Đăng nhập" (133-138). Nút giỏ với badge `cart.totalItems` (chỉ khi hydrated) (18, 140-148). Nav 10 danh mục cứng `categoryId=1..10` (22-33, 153-163) — ISSUE-44 |
| Footer | `src/lib/components/footer.svelte` | — | Cột danh mục cứng `categoryId=1..5` (8-16), cột hỗ trợ (`/products`, `/orders`, `/cart`, `/login`) (17-25), liên hệ (địa chỉ/điện thoại/email từ i18n), copyright |
| ProductCard | `src/lib/components/product-card.svelte` | `product: Product \| ShowcaseProduct`, `class?` | Ảnh/📦, badge giảm giá hoặc "hết hàng", badge mock, giá + giá cũ, sao/đã bán (chỉ mock), "còn N" khi tồn ≤10. Nút "Thêm vào giỏ" `preventDefault` rồi `cart.add({...product, quantity: stock ?? 1}, 1)` + toast (19-33) — ISSUE-42. Link `href="/products/{id}"` không qua `localizeHref` (37) — ISSUE-39. Nút nằm trong thẻ `<a>` (interactive lồng nhau) |

UI primitive dùng từ `@ecommerce/ui`: `button`, `input`, `label`, `textarea`, `separator`, `badge`, `skeleton`, `empty-state`, `dropdown-menu`, `sonner` (Toaster).

---

## 8. Stores — `cart` (FEAT-08)

`src/lib/stores/cart.svelte.ts` — singleton `export const cart = new CartStore()` (dòng 91), dùng Svelte 5 runes (module `.svelte.ts`). Vì là singleton cấp module nhưng chỉ đọc/ghi `localStorage` ở client, không rò dữ liệu giữa các request SSR (trên server `items` luôn rỗng).

| Thành phần | Chi tiết | Dòng |
| :--- | :--- | :--- |
| Persistence key | `localStorage['ecommerce-cart']` | 4 |
| Định dạng ghi | `{ "state": { "items": CartItem[] }, "version": 0 }` (tương thích zustand-persist của app cũ) | 81-84 |
| Định dạng đọc | Chấp nhận `{state:{items}}`, `{items}` hoặc mảng trần; JSON lỗi → `[]` | 6-23 |
| `items` | `$state<CartItem[]>` — `CartItem = { product: Product, quantity }` (`types.ts:138-141`); lưu **snapshot cả object Product** (giá, ảnh, tồn) tại thời điểm thêm | 26 |
| `hydrated` | `$state(false)`; UI dùng để tránh nháy "giỏ rỗng" khi SSR | 27 |
| `totalItems` | `$derived` tổng `quantity` | 29 |
| `totals`, `totalPrice`, `shippingFee`, `grandTotal` | `$derived` từ `computeTotals` (dùng chung với server) | 30-37 |
| `hydrate()` | Đọc storage một lần (gọi ở `+layout.svelte:15`) | 40-44 |
| `add(product, quantity = 1)` | Cộng dồn nếu trùng `productId`, ngược lại push; **không giới hạn** theo tồn kho hay mức 99 của `checkoutSchema` | 46-54 |
| `remove(productId)` | Lọc bỏ dòng | 56-59 |
| `updateQuantity(productId, quantity)` | `≤0` → `remove`; ngược lại gán | 61-71 |
| `clear()` | Xoá hết (gọi sau đặt hàng thành công và nút "Xoá tất cả") | 73-76 |
| `persist()` (private) | Ghi storage; lỗi (đầy/bị chặn) bị bỏ qua, giỏ chỉ nằm trong bộ nhớ | 78-88 |

Hệ quả hành vi: số lượng một dòng có thể vượt 99 (bấm "Thêm" nhiều lần / nút + ở `/cart`), khi đó `placeOrder` bị từ chối ở bước validate và người dùng chỉ thấy toast chung `checkout_orderFailed` (xem mục 15, ghi chú). Giỏ không đồng bộ giữa thiết bị, không gắn với user, không xoá khi logout; giá trong giỏ có thể cũ (server tính lại khi đặt hàng).

---

## 9. Hooks

### 9.1 `src/hooks.server.ts`

`export const handle = sequence(handleParaglide, handleSession)` (dòng 53) — thứ tự: Paraglide chạy **trước**, session **sau**.

1. **`handleParaglide`** (14-25): `paraglideMiddleware(event.request, ({ request, locale }) => resolve({ ...event, request }, { transformPageChunk }))` — xác định locale theo strategy `url → cookie → baseLocale`, rồi thay `%paraglide.lang%` và `%paraglide.dir%` trong `app.html` (`src/app.html:2`) bằng locale và `getTextDirection(locale)`.
2. **`handleSession`** (32-51):
   - `current = session(event.cookies)` — đọc và `JSON.parse` cookie `session_user`; hợp lệ khi có `username` kiểu string (`session.ts:70-80`). Cookie này **không được ký** (xem ISSUE-66).
   - Nếu có `current` **và** `shouldRefresh(cookies)` (`session.ts:87-91`: không có `access_token`, hoặc thiếu `session_expires`, hoặc còn < 30 giây) và có `refresh_token` → `POST /api/v1/auth/refresh` → `setSession(cookies, tokens)` (không truyền `user` → giữ nguyên cookie profile, `session.ts:29-31`); lỗi → `clearSession` và `current = null` (không quay lại SSO).
   - `event.locals.user = current?.user ?? null`. Giá trị `roles` trong cookie luôn rỗng do ISSUE-100.
   - Khi không có cookie `session_user` thì **không** refresh dù vẫn còn `refresh_token` (liên quan ISSUE-67).
   - Refresh diễn ra trên mọi request (cả request remote function), nên token được xoay vòng trong suốt phiên — FEAT-03.

### 9.2 `src/hooks.ts`

`reroute = (request) => deLocalizeUrl(request.url).pathname` (dòng 4): `/en/products` được route tới `src/routes/products`, nên không cần thư mục `[[lang]]`.

---

## 10. Session & cookies

Hằng số và option: `frontend/packages/lib/src/server/session.ts:4-11`. Option chung `base = { path: '/', httpOnly: true, sameSite: 'lax' }` — **không** đặt `secure` tường minh (SvelteKit tự đặt `secure: true` trừ khi host là `localhost`) và tên cookie trùng với app admin (ISSUE-11).

| Cookie | Nội dung | maxAge | httpOnly | sameSite | Ghi bởi |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `access_token` | JWT access token Keycloak | `expires_in` (fallback 300s) | ✓ | lax | `setSession` (`session.ts:37`) |
| `refresh_token` | Refresh token | `refresh_expires_in` (fallback 1800s) | ✓ | lax | `session.ts:38` |
| `session_expires` | epoch giây hết hạn access token | = refresh maxAge | ✓ | lax | `session.ts:39-42` |
| `session_user` | JSON `SessionUser {id, username, fullName, email, roles[]}` — plaintext, không ký; `fullName`/`roles` thực tế rỗng do ISSUE-100 | = refresh maxAge | ✓ | lax | `session.ts:44-48` |
| `auth_redirect` | Path đích sau đăng nhập | 600s | ✓ | lax | `login/+page.server.ts:16-21`; xoá ở `auth/callback/+page.server.ts:12` |
| Cookie locale Paraglide | Locale đã chọn | do runtime Paraglide sinh (không có trong source repo) | — | — | `paraglideMiddleware` / runtime |

`clearSession` xoá 4 cookie phiên với `path: '/'` (`session.ts:51-55`). Access token không bao giờ tới JavaScript phía trình duyệt; mọi API call do server SvelteKit thực hiện. CSRF cho remote command dựa vào kiểm tra origin mặc định của SvelteKit (cần `ORIGIN` đúng) và `sameSite=lax`.

---

## 11. i18n (FEAT-27)

| Mục | Giá trị |
| :--- | :--- |
| Thư viện | `@inlang/paraglide-js` 2 — plugin Vite (`vite.config.ts:26-32`), plugin inlang `message-format@4`, `m-function-matcher@2` (`project.inlang/settings.json:3-6`) |
| Locale | `vi` (base), `en` (`settings.json:10-11`) |
| Strategy | `['url', 'cookie', 'baseLocale']` (`vite.config.ts:30`): URL không prefix = `vi`, `/en/...` = `en`; sau đó tới cookie; cuối cùng `vi` |
| Output | `src/lib/paraglide` (sinh lúc build, gitignore), có `.d.ts` |
| Sử dụng | `m.<key>()` từ `#lib/paraglide/messages.js`; `localizeHref`, `deLocalizeUrl`, `getLocale`, `locales`, `baseLocale` từ `#lib/paraglide/runtime.js` |
| Đổi ngôn ngữ | Link `localizeHref(pathname, { locale })` + full reload (`header.svelte:64-76`) |
| HTML | `<html lang dir>` thay trong `transformPageChunk` (`hooks.server.ts:19-22`) |

**Key** (`messages/vi.json`, `messages/en.json`): mỗi file **306 key**, hai file khớp nhau hoàn toàn. Nhóm theo prefix:

| Prefix | Số key | Prefix | Số key |
| :--- | :---: | :--- | :---: |
| `auth_` | 59 | `home_` | 43 |
| `footer_` | 37 | `checkout_` | 36 |
| `header_` | 25 | `products_` | 22 |
| `cart_` | 17 | `productdetail_` | 17 |
| `common_` | 13 | `category_` | 12 |
| `error_` | 7 | `orders_` | 7 |
| `productcard_` | 5 | `lang_` | 3 |
| `productgrid_` | 3 | | |

Đếm theo tham chiếu `m.<key>` trong `src/`: khoảng **113 key không được dùng** (43 `auth_*`, 27 `footer_*`, 15 `header_*`, 6 `home_*`, 5 `common_*`, 3 `lang_*`, và một số key lẻ như `cart_applyCoupon`, `checkout_paypal`, `productdetail_metaTitle`, `productdetail_reviews`) — ISSUE-43. Tham số message truyền dạng string (vd `m.cart_title({ count: String(n) })`).

Thông điệp lỗi trong `checkoutSchema`/`shippingSchema` là tiếng Anh cố định (`schemas.ts:44-63`) nhưng UI không hiển thị chúng (chỉ dùng `checkout_validationErr`/`checkout_orderFailed`).

---

## 12. SEO & trang lỗi (FEAT-28)

| Thành phần | Hành vi | Nguồn |
| :--- | :--- | :--- |
| Canonical | `page.url.href` (bao gồm query string) | `+layout.svelte:28` |
| hreflang | Một `<link rel="alternate">` cho mỗi locale + `x-default` → bản `vi` | `+layout.svelte:17-34` |
| `<title>` / meta | Trang chủ có title + description; `/products`, `/cart`, `/checkout`, `/orders`, `/login`, callback, trang lỗi có title; **`/products/[id]` không có title** | các `+page.svelte` |
| `sitemap.xml` | `GET` trả XML tĩnh gồm `/` và `/products` × 2 locale (4 URL) với origin từ request; không liệt kê trang sản phẩm; `content-type: application/xml` | `sitemap.xml/+server.ts:5-21` |
| `robots.txt` | `Allow: /`; `Disallow` `/cart`, `/checkout`, `/orders`, `/auth/` và bản `/en/...`; không chặn `/login`; không khai báo `Sitemap:` | `static/robots.txt:1-11` |
| Trang lỗi | `+error.svelte`: 404 → icon `SearchX` + `error_notFoundTitle/Hint`; khác → `AlertTriangle` + `error_title/hint`; nút về trang chủ + `location.reload()`; title `error_metaTitle({status})` | `src/routes/+error.svelte:9-29` |
| Preload | `data-sveltekit-preload-data="hover"` | `src/app.html:9` |

Lỗi tải sản phẩm ở `/products/[id]` được xử lý trong `{:catch}` của trang, không ném `error(404)`, nên không đi qua `+error.svelte` và trả HTTP 200.

---

## 13. Sequence diagrams

### 13.1 Đăng nhập SSO

```mermaid
sequenceDiagram
    autonumber
    actor U as Browser
    participant W as web (SvelteKit server)
    participant G as APISIX gateway
    participant A as auth-service
    participant K as Keycloak

    U->>W: GET /checkout (chưa có session_user)
    W-->>U: 303 /login?redirect=%2Fcheckout
    U->>W: GET /login?redirect=/checkout
    W->>W: safeRedirect(); set cookie auth_redirect (600s)
    W-->>U: Trang login, ssoUrl = API_PUBLIC_URL/api/v1/auth/login?redirect_uri=<origin>/auth/callback
    U->>G: GET /api/v1/auth/login?redirect_uri=...
    G->>A: /auth/api/v1/auth/login
    A-->>U: 302 tới Keycloak authorize
    U->>K: Đăng nhập
    K-->>A: callback (code)
    Note over A: ISSUE-04 / ISSUE-05
    A-->>U: 302 <origin>/auth/callback?ticket=...
    U->>W: GET /auth/callback?ticket=...
    W->>W: đọc + xoá auth_redirect
    W->>G: GET /api/v1/auth/session?ticket=... (không bearer)
    G->>A: Session
    A-->>W: TokenResponse
    W->>G: GET /api/v1/users/me (Bearer access_token)
    G->>A: GetCurrentUser
    A-->>W: UserResponse (roles)
    W->>W: setSession: access_token, refresh_token, session_expires, session_user
    W-->>U: 303 localizeHref('/checkout') + Set-Cookie
    Note over W: Profile 404 (ISSUE-101) → session_user bị xoá → vẫn "chưa đăng nhập" (ISSUE-67); profile OK nhưng roles=[] (ISSUE-100)
```

### 13.2 Checkout / đặt hàng

```mermaid
sequenceDiagram
    autonumber
    actor U as Browser
    participant S as cart store (localStorage)
    participant W as web server (remote command)
    participant G as APISIX
    participant P as product-service
    participant O as order-service
    participant Pay as payment-service

    U->>W: GET /checkout
    W->>W: hooks: refresh token nếu cần; guard locals.user
    W-->>U: Wizard (bước 1)
    U->>U: Bước 1 shippingSchema.safeParse → Bước 2 chọn COD/ONLINE → Bước 3 xác nhận
    U->>S: đọc cart.items
    U->>W: placeOrder({shipping, paymentMethod, items[{productId, quantity}]})
    W->>W: validate checkoutSchema (quantity 1..99, ≥1 dòng)
    alt không có access_token
        W-->>U: 401 → toast checkout_orderFailed
    end
    loop mỗi dòng
        W->>G: GET /api/products/{id}
        G->>P: lấy sản phẩm
        P-->>W: Product (priceUnit)
    end
    W->>W: computeTotals (ship 30.000 nếu < 500.000)
    W->>G: POST /api/orders {orderDate, orderDesc, orderFee=total, productId=dòng đầu}
    G->>O: tạo order (ISSUE-20, ISSUE-21)
    O-->>W: Order {orderId}
    W->>G: POST /api/payments {isPayed, paymentStatus, orderId}
    G->>Pay: tạo payment (ISSUE-18)
    Pay-->>W: Payment
    Note over W: Không idempotency / không bù trừ (ISSUE-23); địa chỉ không lưu (ISSUE-22)
    W-->>U: {orderId, total}
    U->>S: cart.clear()
    U->>U: Bước 4: hiển thị mã đơn
```

---

## 14. Test

- `frontend/apps/web` **không có** test nào (không có file `*.test.ts`/`*.spec.ts`, không có script `test` trong `package.json`), không có test component hay E2E — ISSUE-46.
- Logic dùng chung có unit test (`bun test`) tại `frontend/packages/lib`:

| File | Phạm vi |
| :--- | :--- |
| `frontend/packages/lib/src/server/session.test.ts` | `toSessionUser`, `isAdmin`, `setSession` (gồm case `user = null` xoá profile — dòng 133), `clearSession`, `session`, `shouldRefresh`, accessor token |
| `frontend/packages/lib/src/server/api.test.ts` | Unwrap envelope, query string, bearer header, JSON body, `ApiError`, exchange ticket không bearer |
| `frontend/packages/lib/src/server/redirect.test.ts` | `safeRedirect` |
| `frontend/packages/lib/src/schemas.test.ts` | `pageQuerySchema`, `productQuerySchema`, `shippingSchema`, `checkoutSchema` (giới hạn 1–99), schema admin |
| `frontend/packages/lib/src/cart.test.ts` | `computeTotals` (ngưỡng miễn phí ship) |
| `frontend/packages/lib/src/format.test.ts` | `formatPrice`, `formatDate` |

Các test dùng mock `fetch`, không kiểm hợp đồng với backend thật → không phát hiện ISSUE-48/49.

---

## 15. Vấn đề đã biết

| ID | Mức độ | Tóm tắt (phạm vi web) | Vị trí |
| :--- | :--- | :--- | :--- |
| ISSUE-04 | Critical | Callback SSO ở auth-service sai `redirect_uri`, fallback mock token → web có thể nhận token không hợp lệ | auth-service |
| ISSUE-05 | Critical | `redirect_uri` của `/auth/login` không có allowlist (web truyền `origin/auth/callback`) | auth-service; `login/+page.server.ts:23-24` |
| ISSUE-100 | Critical | `/api/v1/users/me` trả `fullname` (thường) và không có `roles[{id,name}]` → `toSessionUser` cho `roles: []`, `fullName` rỗng | `frontend/packages/lib/src/server/session.ts:13-21`; `services/auth-service/internal/model/user.go:16-22` |
| ISSUE-101 | High | User Keycloak không có trong DB `authservice` → `/users/me` 404 → kích hoạt ISSUE-67 khi đăng nhập | `auth/callback/+page.server.ts:26-31` |
| ISSUE-11 | Medium | Cookie phiên không đặt `secure` tường minh; trùng tên/path với admin → chia sẻ phiên trên cùng host | `frontend/packages/lib/src/server/session.ts:4-9` |
| ISSUE-20 | High | Order chỉ lưu `productId` dòng đầu, không có số lượng/trạng thái | `checkout/checkout.remote.ts:30-35` |
| ISSUE-21 | High | `orderFee` do client (BFF) gửi; không kiểm/giữ kho | `checkout/checkout.remote.ts` |
| ISSUE-22 | Medium | Địa chỉ giao hàng được validate nhưng không lưu | `checkout/checkout.remote.ts:13-44` |
| ISSUE-23 | Medium | Không idempotency; order + payment không giao dịch | `checkout/checkout.remote.ts:30-41` |
| ISSUE-38 | High | Trang chủ dùng mock; thêm mock vào giỏ làm `placeOrder` lỗi | `src/lib/data/home-mock.ts`, `src/routes/+page.svelte` |
| ISSUE-39 | Low | `product-card` link không `localizeHref` | `product-card.svelte:37` |
| ISSUE-40 | Medium | Lọc giá chỉ trên trang hiện tại, client-side | `products/+page.svelte:57-63` |
| ISSUE-41 | Medium | Trạng thái đơn luôn "Đang xử lý" | `orders/+page.svelte:44` |
| ISSUE-42 | Low | Quick add ghi `quantity: stock ?? 1` vào object Product trong giỏ | `product-card.svelte:27` |
| ISSUE-43 | Low | ~113 key i18n không dùng (59 `auth_*` khai báo, 43 không dùng) | `messages/*.json` |
| ISSUE-44 | Medium | Danh mục header/footer/home hard-code `categoryId=1..10` | `header.svelte:22-33`, `footer.svelte:8-16`, `+page.svelte:34-57` |
| ISSUE-46 | Medium | Không có test component/E2E cho web | `frontend/apps/web` |
| ISSUE-47 | Medium | Phụ thuộc `remoteFunctions` và `experimental.async` | `vite.config.ts:13-19` |
| ISSUE-48 | Critical | `products.list`, `categories.list`, `orders.list` kỳ vọng Page, backend trả mảng → catalog/đơn hàng rỗng hoặc lỗi | `frontend/packages/lib/src/server/api.ts:129,136,143` |
| ISSUE-49 | High | `/orders` hiển thị đơn của mọi người dùng | `orders/orders.remote.ts:5-8` |
| ISSUE-50 | Medium | `Product.description` không tồn tại ở backend | `frontend/packages/lib/src/types.ts:73`; `products/[id]/+page.svelte:151-156` |
| **ISSUE-66** | High | Cookie `session_user` (chứa `roles`) là JSON plaintext **không ký/không xác minh**; hook tin cookie này làm danh tính. Ở web: giả mạo cookie vượt guard `/checkout`, `/orders` (API vẫn cần token). Cùng helper được admin dùng cho guard `isAdmin` → kết hợp ISSUE-01, user thường có token hợp lệ tự đặt `roles:["ADMIN"]` để vào backoffice | `frontend/packages/lib/src/server/session.ts:44-48,70-80`; `frontend/apps/web/src/hooks.server.ts:33,49` |
| **ISSUE-67** | High | Khi `GET /api/v1/users/me` lỗi lúc callback (thường xuyên do ISSUE-101), `setSession(..., null)` xoá `session_user` → `session()` trả null → người dùng có token nhưng bị coi là khách, không được refresh; vào `/checkout` bị đẩy về `/login` → SSO lặp lại | `frontend/apps/web/src/routes/auth/callback/+page.server.ts:23-33`; `frontend/packages/lib/src/server/session.ts:44-45,70-72`; `frontend/apps/web/src/hooks.server.ts:35` |

Ghi chú hành vi khác (chưa đăng ký ID): giỏ không giới hạn số lượng (`cart.svelte.ts:46-71`) trong khi `checkoutSchema` giới hạn 1–99 (`schemas.ts:57`) và không kiểm tồn → đặt hàng thất bại với toast chung; `/products/[id]` không có `<title>` và trả HTTP 200 khi sản phẩm không tồn tại.

---

## 16. Tham chiếu source

| Chủ đề | File |
| :--- | :--- |
| Cấu hình build / dev | `frontend/apps/web/package.json`, `frontend/apps/web/vite.config.ts`, `frontend/apps/web/tsconfig.json`, `frontend/apps/web/Dockerfile`, `frontend/apps/web/.env.example` |
| Compose | `docker-compose.yml:412-431` |
| Hooks | `frontend/apps/web/src/hooks.server.ts`, `frontend/apps/web/src/hooks.ts`, `frontend/apps/web/src/app.html`, `frontend/apps/web/src/app.d.ts` |
| Layout / lỗi / SEO | `frontend/apps/web/src/routes/+layout.server.ts`, `+layout.svelte`, `+error.svelte`, `sitemap.xml/+server.ts`, `frontend/apps/web/static/robots.txt` |
| Trang chủ | `frontend/apps/web/src/routes/+page.svelte`, `frontend/apps/web/src/lib/data/home-mock.ts` |
| Catalog / chi tiết | `frontend/apps/web/src/routes/products/+page.svelte`, `catalog.remote.ts`, `products/[id]/+page.svelte`, `products/[id]/product.remote.ts` |
| Giỏ / checkout | `frontend/apps/web/src/lib/stores/cart.svelte.ts`, `src/routes/cart/+page.svelte`, `src/routes/checkout/{+page.server.ts,+page.svelte,checkout.remote.ts}` |
| Đơn hàng | `frontend/apps/web/src/routes/orders/{+page.server.ts,+page.svelte,orders.remote.ts}` |
| Auth | `frontend/apps/web/src/routes/login/*`, `src/routes/auth/callback/*`, `src/lib/remote/auth.remote.ts`, `src/lib/server/api.ts` |
| Components | `frontend/apps/web/src/lib/components/{header,footer,product-card}.svelte` |
| i18n | `frontend/apps/web/project.inlang/settings.json`, `frontend/apps/web/messages/{vi,en}.json` |
| Shared lib | `frontend/packages/lib/src/server/{api,session,env,redirect,index}.ts`, `frontend/packages/lib/src/{cart,schemas,types,format}.ts` |
| Gateway / auth-service | `deploy/apisix/apisix.yaml:63-82,161-231`, `services/auth-service/internal/handler/auth_handler.go:64-76` |
