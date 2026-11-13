# Quy ước viết code

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> Quy ước **rút ra từ code hiện có** (không phải bộ quy tắc được ép bằng linter, trừ khi ghi rõ). Khi thêm code mới, hãy theo mẫu đang dùng để giữ nhất quán; chỗ nào code hiện tại tự mâu thuẫn được ghi chú "⚠️ bất nhất". Kiến trúc tham chiếu: [Go](../02-architecture/software/go-service-architecture.md) · [Bun](../02-architecture/software/bun-service-architecture.md) · [Frontend](../02-architecture/software/frontend-architecture.md) · [API](../02-architecture/software/api-conventions.md).

## 1. Công cụ ép buộc tự động

| Phạm vi | Công cụ | Cấu hình | Chạy ở |
| :--- | :--- | :--- | :--- |
| Go | `gofmt` (fail nếu chưa format), `go vet` | mặc định | `make check`, CI `go-check` |
| Bun services | `tsc --noEmit` (`strict`, `verbatimModuleSyntax`) | `services/<svc>/tsconfig.json` | thủ công `bun run typecheck` (CI chỉ chạy `bun test`) |
| Frontend | Prettier 3 + `prettier-plugin-svelte` + `prettier-plugin-tailwindcss` | `frontend/prettier.config.js`, `frontend/apps/*/prettier.config.js` | `bun run lint`, CI |
| Frontend | ESLint 10 flat config: `@eslint/js` recommended, `typescript-eslint` recommended, `eslint-plugin-svelte` recommended + prettier | `frontend/apps/*/eslint.config.js` | `bun run lint`, CI |
| Frontend | `svelte-check`, i18n parity | `package.json` scripts | `bun run check`, `bun run i18n:check`, CI |
| Line endings | LF bắt buộc cho `*.go`, `*.sh`, `Makefile`, `frontend/**` | `.gitattributes` | git |

Không có linter Go bổ sung (golangci-lint), không có Prettier/ESLint cho service Bun, không có pre-commit hook.

## 2. Go

### 2.1 Bố cục package

```text
services/<name>-service/
├── main.go               # chỉ wiring + server; không chứa logic nghiệp vụ
├── main_test.go          # test HTTP qua router
└── internal/
    ├── handler/<x>_handler.go
    ├── service/<x>_service.go
    ├── repository/<x>_repo.go | <x>_repository.go   ⚠️ bất nhất tên file
    └── model/<x>.go
```

- Module path `com.ecommerce/<name>-service`; import shared bằng `com.ecommerce/shared/...` (alias `commonmw` cho middleware shared, `chimiddleware` cho middleware chi — `services/product-service/main.go:20-22`).
- Mỗi tầng một package; một service có thể có nhiều cặp handler/service/repo (product: product + favourite).
- Chỉ đặt vào `services/shared` những gì **mọi** service cần (config, DB pool, middleware, response).

### 2.2 Đặt tên

| Thứ | Quy ước quan sát | Ví dụ |
| :--- | :--- | :--- |
| Constructor | `NewXxx(deps) *Xxx` | `NewProductHandler(svc)` |
| Struct tầng | `XxxHandler`, `XxxService`, `XxxRepository` | `services/product-service/internal/handler/product_handler.go:15` |
| Receiver | 1 chữ: `h` (handler), `s` (service), `r` (repository) | |
| Handler method | động từ + đối tượng: `ListProducts`, `GetProductByID`, `FindAllOrdersPaged`, `UpdateCartById` | ⚠️ `ByID` vs `ById` lẫn lộn |
| DTO | hậu tố `Dto` / `Vm` (di sản Spring): `ProductDto`, `ProductListGetVm` | `services/product-service/internal/model/product.go:10` |
| Field Go | ⚠️ `ProductId`, `CategoryId` (product/order/payment) vs `ID`, `UserID` (auth, favourite, search) | |
| Sentinel error | `ErrXxx` ở package `service` | `services/order-service/internal/service/order_service.go:13-14` |
| Log prefix | `"[Product Service] ..."` | `services/product-service/main.go:60` |

### 2.3 JSON & DB tag

- JSON tag **camelCase**, khớp tên trường trong frontend: `json:"productTitle"`. Ngoại lệ có chủ đích: token Keycloak giữ snake_case (`access_token`) vì là payload chuẩn OAuth. ⚠️ `fullname` (auth) lệch với `fullName` ở request và frontend (ISSUE-100).
- Trường tùy chọn: con trỏ + `omitempty` (order/payment, `services/order-service/internal/model/order.go:19-27`) hoặc giá trị + luôn xuất hiện (product). Chọn theo service đang sửa.
- Struct scan riêng trong repository dùng tag `db:"snake_case"` (ví dụ `productRow` — `services/product-service/internal/repository/product_repo.go:177-187`), sau đó `toDto()`.
- Cột DB snake_case; `COALESCE(col, '')` trong SELECT để tránh `NULL` khi scan vào `string`.

### 2.4 Xử lý lỗi

- Repository: trả `error` của driver nguyên vẹn; `sql.ErrNoRows` → `(nil, nil)`.
- Service: biến "không tìm thấy" thành sentinel, bọc ngữ cảnh bằng `fmt.Errorf("%w: product with id %d not found", ErrNotFound, id)` (`services/product-service/internal/service/product_service.go:34`).
- Handler: `errors.Is` → helper `response.*`; luôn `return` ngay sau khi ghi lỗi.
- Không `panic` trong code ứng dụng (Recoverer chỉ là lưới an toàn).
- ⚠️ Code hiện tại trả `err.Error()` ra client cho lỗi 500 (ISSUE-93) và nuốt lỗi bằng `_ =` ở vài chỗ (`InitSchema`, Kafka send). Code mới **nên** log lỗi chi tiết và trả thông điệp chung.

### 2.5 Context

- Mọi method service/repository nhận `ctx context.Context` là tham số đầu tiên; handler truyền `r.Context()`.
- Đọc danh tính bằng `middleware.UserID(ctx)`, `Username(ctx)`, `Email(ctx)` (không verify — chỉ dùng cho mục đích hiển thị/ghi nhận, không cho phân quyền; [ISSUE-02](../07-plan/known-issues-tech-debt.md)).
- `context.Background()` chỉ dùng ở `main` (khởi động, shutdown, consumer).

### 2.6 Truy vấn DB (bob)

- Dùng builder `psql.Select/Insert/Update/Delete` với `sm.*`, `im.*`, `um.*`, `dm.*`; giá trị luôn qua `psql.Arg(...)` (tham số hóa).
- Tên cột qua `psql.Quote("col")`; cột sort động phải đi qua `switch` whitelist (`services/order-service/internal/repository/order_repository.go:339-348`).
- `bob.All/One` + `scan.StructMapper[T]()`; insert lấy id bằng `im.Returning("id")` + `scan.SingleColumnMapper[int]`.
- Mỗi method kiểm `r.db.Pool == nil` để hỗ trợ unit test không DB.

### 2.7 HTTP

- Route đăng ký qua hàm `register…Routes(r chi.Router, base string)` và lặp qua danh sách base path (prefix + root).
- Thành công: `response.WriteJSON(w, status, v)` (JSON thô) — trừ auth-service dùng `response.OK/Message/Error` (envelope). Không trộn hai kiểu trong cùng service.
- Phân trang kiểu Spring: query `page`, `size`, `sortBy`, `sortOrder` + `PageResponse[T]`.

### 2.8 Test

- `main_test.go`, package `main`; dựng router bằng repo `nil`; helper `performRequest(app, method, target, body)`; assert bằng `testify/assert` (`services/product-service/main_test.go:17-45`).

## 3. TypeScript – Bun services

| Chủ đề | Quy ước |
| :--- | :--- |
| Module | ESM (`"type": "module"`), `import type` cho type (bắt buộc bởi `verbatimModuleSyntax`) |
| Style | 2 space; ⚠️ nháy kép ở `index.ts`/`db.ts` (promotion, rating, media) vs nháy đơn ở `service.ts`, notification — không có formatter |
| File | `index.ts` (app + routes), `service.ts` (class `XxxService`), `db.ts`, `schema.ts`, `types.ts` |
| Route | Khai báo trong `function registerXxxRoutes(router: Hono)`, mount 2 lần (prefix + root) |
| Handler | `async (c) => { try { ... return c.json(x) } catch (err: any) { return c.json({ statusCode: "500", message: err.message }, 500) } }` |
| Validation | Zod schema tên `XxxPostVmSchema`/`XxxPutVmSchema` + `type XxxPostVm = z.infer<...>`; `safeParse` → 400 `{statusCode:"400", message:"Validation error", errors}` |
| Param | `parseInt(c.req.param("id"), 10)` + `isNaN` → 400 |
| Drizzle | Bảng `pgTable('snake_name', { camelProp: type('snake_col') })`; `id: bigserial('id', { mode: 'number' })`; xuất `XxxRow = typeof t.$inferSelect` |
| Response VM | Interface `XxxVm`, `XxxListVm` (`ratingList`, `totalElements`, `totalPages`); hàm `mapRow()` chuyển row → VM, Date → `toISOString()` |
| Phân trang | Query `pageNo`, `pageSize` |
| Log | `console.log("[Rating Service] ...")` |
| Env | `process.env.X \|\| "default"` tại module dùng |
| Test | `src/index.test.ts`, `describe/it/expect` từ `bun:test`, `appObj.fetch(new Request(url))` |

## 4. SvelteKit / frontend

### 4.1 Format & lint

- Prettier: **tab**, `singleQuote: true`, `trailingComma: 'none'`, `printWidth: 100`, sắp xếp class Tailwind tự động (`frontend/prettier.config.js:2-9`).
- ESLint tắt `no-undef` (TS) và `svelte/no-navigation-without-resolve` vì href sinh bởi `localizeHref()` (`frontend/apps/web/eslint.config.js:22-44`).
- Dùng `// eslint-disable-next-line` có lý do khi cần `any` (ví dụ `frontend/packages/ui/src/utils.ts`).

### 4.2 Svelte 5 runes

- Runes bắt buộc cho mọi file của dự án (`compilerOptions.runes` — `frontend/apps/web/vite.config.ts:21-23`): dùng `$props()`, `$state()`, `$derived()`, `$effect()`; **không** dùng `export let`, `$:` hay store `writable`.
- Props: `let { children, data }: LayoutProps = $props();` với type sinh từ `./$types`.
- Render children bằng `{@render children()}`.
- State dùng chung: class với field `$state`/`$derived` trong file `*.svelte.ts`, export singleton (`frontend/apps/web/src/lib/stores/cart.svelte.ts`). Truy cập `localStorage` phải kiểm `typeof localStorage === 'undefined'` và chỉ hydrate trong `onMount`.
- Đọc URL: `page` từ `$app/state` (không phải `$app/stores`); query param thao tác bằng `SvelteURLSearchParams`.

### 4.3 Remote functions & server code

- Đặt cạnh route: `routes/<feature>/<feature>.remote.ts`; remote dùng chung app: `src/lib/remote/*.remote.ts`.
- Luôn truyền Zod schema làm đối số đầu (`query(schema, fn)`, `form(schema, fn)`, `command(schema, fn)`); schema tái sử dụng lấy từ `@ecommerce/lib/schemas`. Với `form`, không dùng `z.coerce` (dùng `formNumber`/`optionalId`).
- Gọi backend **chỉ** qua `api()` của app (`#lib/server/api`) hoặc `createApi({ token })`; không `fetch` gateway trực tiếp trong component.
- `query` cho đọc; `command` cho mutation gọi từ JS, làm mới dữ liệu bằng `.updates(query)` hoặc `requested(query, n).refreshAll()`; `form` cho form có progressive enhancement.
- Bảo vệ trang bằng `+page.server.ts` (web) / guard trong `hooks.server.ts` (admin); redirect đăng nhập kèm `?redirect=` và luôn lọc bằng `safeRedirect()`.
- Không đưa token/secret vào `data` trả cho client; chỉ `locals.user`.

### 4.4 Import

| Cần | Import |
| :--- | :--- |
| Code trong `src/lib` của app | `#lib/...` (ví dụ `#lib/components/header.svelte`, `#lib/server/api`) — không dùng `$lib` |
| Type/schema/format/cart dùng chung | `@ecommerce/lib/types`, `@ecommerce/lib/schemas`, `@ecommerce/lib/format`, `@ecommerce/lib/cart` |
| Helper server (cookie, API) | `@ecommerce/lib/server` — chỉ trong file server |
| Component UI | Sub-path từng component: `import { Button } from '@ecommerce/ui/button'`; nhóm nhiều phần: `import * as Sheet from '@ecommerce/ui/sheet'` |
| Icon | `@lucide/svelte` |
| Toast | `toast` từ `svelte-sonner`, `<Toaster />` từ `@ecommerce/ui/sonner` ở layout |
| Zod | `import * as v from 'zod'` |
| i18n | `import { m } from '#lib/paraglide/messages.js'`, `import { localizeHref } from '#lib/paraglide/runtime.js'` |

### 4.5 i18n (Paraglide)

- Mọi chuỗi hiển thị dùng `m.<key>()`; key dạng `<khu-vực>_<tên>` camelCase phần sau (`products_sortNewest`, `checkout_orderFailed`, `sidebar_dashboard`).
- Thêm key vào **cả** `messages/vi.json` và `messages/en.json`; chạy `bun run i18n:check`.
- Mọi `href`/`goto`/`redirect` nội bộ bọc `localizeHref()` để giữ prefix `/en` ([ISSUE-39](../07-plan/known-issues-tech-debt.md) là ví dụ vi phạm).
- Không sửa thư mục sinh `src/lib/paraglide/`.

### 4.6 UI & style

- Thêm component dùng chung vào `packages/ui/src/components/<name>/` (kiểu shadcn-svelte: file `.svelte` + `index.ts`), rồi khai báo sub-path trong `exports` của `packages/ui/package.json`.
- Ghép class bằng `cn()` từ `@ecommerce/ui/utils`; biến thể bằng `tailwind-variants`.
- Màu/token qua CSS variable trong `packages/ui/src/theme.css`; override theo app trong `routes/layout.css`.
- Tên file component kebab-case (`product-card.svelte`, `coming-soon.svelte`).

## 5. Quy ước API khi thêm endpoint

Theo [`api-conventions.md`](../02-architecture/software/api-conventions.md): JSON camelCase, đăng ký route prefix + root, `/actuator/health` sẵn có, giữ định dạng response/lỗi **của service đang sửa**, phân trang theo kiểu service đó đang dùng. Nhớ cập nhật route APISIX và [`api-catalog.md`](../05-api/api-catalog.md).

## 6. Git, nhánh, commit

Repo không có `CONTRIBUTING.md`, commit template hay commitlint. Dấu hiệu duy nhất:

- CI chạy trên push vào `main`, `develop` và pull request vào `main` (`.github/workflows/ci.yml:3-8`) → mô hình nhánh `main` / `develop` + PR.
- Renovate (`renovate.json`, preset `config:recommended`) tự mở PR nâng dependency.
- `.gitattributes` buộc LF cho Go, shell, Makefile, frontend — cấu hình `core.autocrlf` trên Windows không được làm hỏng điều này.
- Báo lỗ hổng bảo mật theo `SECURITY.md`.
