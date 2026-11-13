# Hướng dẫn cho lập trình viên

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> Hướng dẫn chạy, phát triển và kiểm thử repo `ecommerce-microservices` **theo đúng code hiện tại**. Các lệnh đã đối chiếu với `Makefile`, `scripts/dev.ps1`, `docker-compose.yml` và `frontend/package.json`. Khi README gốc mâu thuẫn với tài liệu này, tài liệu này phản ánh code ([ISSUE-45](../07-plan/known-issues-tech-debt.md)).

## 1. Yêu cầu

| Công cụ | Phiên bản | Ghi chú |
| :--- | :--- | :--- |
| Docker + Compose v2 | — | Chạy hạ tầng và toàn stack |
| Go | 1.27+ | Theo `go.work:1` |
| Bun | ≥ 1.4 (frontend pin `bun@1.4.2`, CI dùng 1.4.3) | Service Bun + frontend |
| GNU make **hoặc** PowerShell 5.1+ | — | `scripts/dev.ps1` mirror 1-1 các target của `Makefile` |
| (tùy chọn) k3d, kubectl | — | Đường k8s cũ (`make cluster-up`) — ngoài phạm vi tài liệu ứng dụng |

RAM khuyến nghị ≥ 8 GB (OpenSearch được cấp heap 512 MB, Keycloak, Kafka, 11 service).

## 2. Bố cục repo

```text
.
├── Makefile, scripts/dev.ps1          # task runner (Unix / Windows)
├── docker-compose.yml, .env           # stack local (hạ tầng + 11 service + 2 frontend image)
├── go.work                            # workspace 8 module Go
├── services/
│   ├── shared/                        # module Go dùng chung (config, database, middleware, response)
│   ├── {auth,product,order,payment,inventory,shipping,search}-service/   # Go
│   └── {promotion,rating,media,notification}-service/                   # Bun / TypeScript
├── frontend/                          # Bun workspace: apps/{web,admin}, packages/{lib,ui}
├── deploy/
│   ├── apisix/{apisix.yaml,config.yaml}   # gateway standalone cho compose
│   ├── keycloak/import/ecommerce-realm.json
│   ├── postgres/init/create-all-databases.sql
│   └── k8s/                           # manifest k8s (hạ tầng – ngoài phạm vi)
├── docs/                              # tài liệu cũ / kiến trúc đích (KHÔNG phản ánh as-is)
└── docs-all/                          # bộ tài liệu này
```

Kiến trúc: [container-view](../02-architecture/system/container-view.md) · [Go](../02-architecture/software/go-service-architecture.md) · [Bun](../02-architecture/software/bun-service-architecture.md) · [Frontend](../02-architecture/software/frontend-architecture.md).

## 3. Chạy toàn bộ stack bằng Docker Compose

```bash
make up            # = docker compose up -d --build
make ps            # trạng thái container
make logs SVC=product-service
make down
```

Windows: `.\scripts\dev.ps1 up`, `.\scripts\dev.ps1 logs product-service`…

Lưu ý quan trọng:

- Compose **build** 11 backend service từ source, nhưng **frontend web/admin là image dựng sẵn** `ghcr.io/ecommerce-int/frontend-{web,admin}:latest` với `pull_policy: always` (`docker-compose.yml:412-449`) — sửa code frontend local **không** có hiệu lực trong compose; chạy frontend bằng `make web-dev` (mục 5).
- Biến `${...}` trong compose được nội suy từ file `.env` ở gốc repo (tên biến xem [`cross-cutting-concerns.md` §2.3](../02-architecture/software/cross-cutting-concerns.md#23-danh-mục-biến-chính); **không** commit secret mới vào file này). README nhắc `.env.example` nhưng file đó không tồn tại.
- Bucket media không tự tạo (không có service `storage-init`) — tạo bucket `STORAGE_BUCKET` thủ công qua console RustFS `:9001` trước khi upload ([ISSUE-32](../07-plan/known-issues-tech-debt.md)).
- Keycloak cần ~90 s để healthy; APISIX chờ Keycloak (`depends_on`).

### 3.1 Bảng port (host)

| Thành phần | Port | URL / ghi chú |
| :--- | :--- | :--- |
| Storefront (web) | 3000 | `http://localhost:3000` |
| Backoffice (admin) | 3001 | `http://localhost:3001` |
| APISIX gateway | 9080 | `http://localhost:9080` (proxy) |
| APISIX Prometheus | 9091 | `/apisix/prometheus/metrics` |
| Keycloak | 8080 | console admin (tài khoản admin lấy từ biến `KEYCLOAK_ADMIN*` trong `.env`) |
| PostgreSQL | 5432 | user/password từ `.env` |
| Redis | 6379 | không được app dùng |
| Kafka | 9092 | advertised `kafka:9092` (xem mục 10) |
| OpenSearch | 9200, 9600 | security tắt |
| RustFS (S3) | 9000 (API), 9001 (console) | |
| auth-service | 8088 | |
| product-service | 8086 | |
| order-service | 8084 | |
| payment-service | 8085 | |
| inventory-service | 8082 | |
| shipping-service | 8087 | |
| search-service | 8094 | |
| promotion-service | 8093 | |
| rating-service | 8089 | |
| media-service | 8083 | |
| notification-service | 8090 | |

Nguồn: `docker-compose.yml`, `PORT` mặc định trong từng `main.go` / `src/index.ts`, `frontend/apps/*/vite.config.ts`.

### 3.2 Tài khoản seed

Realm `ecommerce` import từ `deploy/keycloak/import/ecommerce-realm.json` (mật khẩu xem trực tiếp trong file đó — **không** dùng cho môi trường thật, [ISSUE-12](../07-plan/known-issues-tech-debt.md)):

| Username | Realm role | Dùng cho |
| :--- | :--- | :--- |
| `testuser` | `USER` | khách hàng storefront |
| `adminuser` | `ADMIN`, `MANAGER` | backoffice |

Client: `ecommerce-client` (public, redirect `http://localhost:*`). DB `authservice` seed role `USER`, `PM`, `ADMIN` nhưng **không** seed user.

> ⚠️ Hiện tại đăng nhập bằng user seed qua UI **không tạo được phiên**: BFF gọi `/api/v1/users/me`, user chỉ tồn tại trong Keycloak nên auth-service trả 404 → cookie `session_user` bị xóa; ngay cả user có trong DB thì profile không có `roles` → admin luôn bị chuyển `/access-denied` (ISSUE-100). Ngoài ra luồng đổi code SSO có lỗi [ISSUE-04](../07-plan/known-issues-tech-debt.md). Xem [FEAT-01](../04-features/FEAT-01-sso-login-logout.md).

## 4. Chỉ chạy hạ tầng + một service trên máy

```bash
make infra-up                       # postgres redis kafka opensearch rustfs keycloak apisix
set -a; source .env; set +a         # nạp biến (Go không tự đọc .env; Bun chỉ đọc .env trong thư mục service)
make go-run SVC=product-service     # = cd services/product-service && go run .
make bun-dev SVC=rating-service     # = bun install && bun --watch src/index.ts
```

| Make | PowerShell |
| :--- | :--- |
| `make infra-up` | `.\scripts\dev.ps1 infra-up` |
| `make go-run SVC=<svc>` | `.\scripts\dev.ps1 go-run <svc>` |
| `make bun-dev SVC=<svc>` | `.\scripts\dev.ps1 bun-dev <svc>` |
| `make restart SVC=<svc>` | `.\scripts\dev.ps1 restart <svc>` |

Những điều cần biết (ISSUE-97):

- `.env` đặt `POSTGRES_HOST=localhost`, `KAFKA_SERVERS=localhost:9092`, `OPENSEARCH_URL`, `STORAGE_ENDPOINT`, `KEYCLOAK_SERVER_URL` trỏ `localhost` — phù hợp chạy trên host. DB name lấy theo default của từng service (`productservice`…), đã được script init tạo.
- **APISIX upstream trỏ tên container** (`product-service:8086`… — `deploy/apisix/apisix.yaml:1-45`). Service chạy trên host **không** nhận được request qua gateway; khi chỉ chạy `infra-up`, mọi route gateway trả 5xx. Cách làm việc thực tế:
  - Gọi thẳng service: `curl http://localhost:8086/api/products` (truyền `Authorization: Bearer <token>` nếu handler cần claims), hoặc
  - Chạy `make up` rồi `docker compose stop <svc>` và tạm sửa node upstream trong `apisix.yaml` sang `host.docker.internal:<port>` (chỉ local, không commit).
- Kafka quảng bá `kafka:9092` (`docker-compose.yml:50`) → client trên host sau bước bootstrap sẽ cố kết nối `kafka:9092`; thêm `127.0.0.1 kafka` vào `/etc/hosts` nếu cần thử producer/consumer từ host.

## 5. Phát triển frontend

```bash
make web-install        # cd frontend && bun install
make web-dev            # web :3000 + admin :3001 (bun run --filter '*' dev)
cd frontend && bun run dev:web     # chỉ storefront
cd frontend && bun run dev:admin   # chỉ backoffice
```

- Tạo `frontend/apps/{web,admin}/.env` từ `.env.example` (chỉ cần `API_INTERNAL_URL=http://localhost:9080`).
- Frontend gọi **gateway**, nên cần stack backend chạy qua `make up` (gateway + service).
- File Paraglide sinh vào `src/lib/paraglide` khi chạy dev/check; không sửa tay.
- Đăng nhập trên `localhost:3000` và `localhost:3001` dùng chung cookie ([ISSUE-11](../07-plan/known-issues-tech-debt.md)).

## 6. Thêm endpoint mới

### 6.1 Service Go (ví dụ product-service)

1. **Model** — thêm/sửa DTO trong `internal/model/*.go`, JSON tag **camelCase**.
2. **Repository** — thêm method dùng bob; luôn có nhánh `if r.db.Pool == nil { return <stub> }` để unit test chạy không DB; `sql.ErrNoRows` → `(nil, nil)`. Nếu cần bảng/cột mới, thêm vào chuỗi DDL trong `InitSchema` (chỉ có tác dụng với bảng **chưa tồn tại** — [ISSUE-28](../07-plan/known-issues-tech-debt.md); cột mới cho bảng cũ phải `ALTER` thủ công).
3. **Service** — gọi repo, chuyển `nil` thành sentinel (`fmt.Errorf("%w: ...", ErrNotFound)`).
4. **Handler** — viết `func (h *XHandler) Foo(w, r)`; decode body bằng `json.NewDecoder`, path param bằng `chi.URLParam` + `strconv.Atoi`; map lỗi bằng `errors.Is` → `response.NotFound/BadRequest/InternalError`; thành công → `response.WriteJSON(w, status, v)` (auth-service dùng `response.OK/Error`).
5. **Route** — đăng ký trong hàm `register…Routes(r, base)` để route có cả prefix và root.
6. **Test** — thêm case vào `main_test.go` (router + repo `nil` + `httptest`).
7. **Gateway** — nếu path chưa được route nào bao phủ, thêm route APISIX (mục 6.4).
8. `make go-fmt && make go-vet && make go-test`.

```go
// handler
func (h *ProductHandler) registerProductRoutes(r chi.Router, base string) {
    r.Get(base+"/{productId}/stock", h.GetStock) // mới
}
```

Chi tiết mẫu: [`go-service-architecture.md`](../02-architecture/software/go-service-architecture.md), quy ước: [`coding-conventions.md`](coding-conventions.md).

### 6.2 Service Bun (ví dụ rating-service)

1. `src/schema.ts` — thêm cột/bảng Drizzle **và** cập nhật DDL tương ứng trong `initDb()` của `src/db.ts`.
2. `src/types.ts` — Zod schema cho body (`XxxPostVmSchema`) + interface VM response.
3. `src/service.ts` — method trong class service dùng Drizzle.
4. `src/index.ts` — thêm route trong `registerXxxRoutes(router)` (tự động có ở `/rating/...` và `/...`), bọc `try/catch`, `safeParse` → 400 `{statusCode:"400", message, errors}`.
5. `src/index.test.ts` — thêm case `app.fetch(new Request(...))`.
6. `make bun-test` (hoặc `cd services/<svc> && bun test && bun run typecheck`).

### 6.3 SvelteKit: route + remote function

1. Tạo `frontend/apps/<app>/src/routes/<feature>/+page.svelte`.
2. Nếu trang cần đăng nhập (web): thêm `+page.server.ts` redirect `/login?redirect=` khi `!locals.user` (mẫu `frontend/apps/web/src/routes/orders/+page.server.ts`). Admin đã có guard toàn cục.
3. Thêm schema Zod vào `frontend/packages/lib/src/schemas.ts` (nếu dùng chung) và method gọi gateway vào `createApi` trong `frontend/packages/lib/src/server/api.ts` + type trong `types.ts`.
4. Tạo `routes/<feature>/<feature>.remote.ts`:

   ```ts
   import { query, command } from '$app/server';
   import { pageQuerySchema } from '@ecommerce/lib/schemas';
   import { api } from '#lib/server/api';

   export const getThings = query(pageQuerySchema, ({ page, size }) => api().things.list({ page, size }));
   ```

5. Dùng trong component: `{#await getThings(args)} ... {:then data} ... {:catch} ... {/await}`; mutation `await doThing(input)` / form `{...saveThing.enhance(...)}`.
6. Thêm message vào **cả** `messages/vi.json` và `messages/en.json`, dùng `m.key()`; link nội bộ qua `localizeHref()`.
7. `cd frontend && bun run check && bun run lint && bun run i18n:check && bun test`.

### 6.4 Thêm route APISIX (compose)

Sửa `deploy/apisix/apisix.yaml` (standalone, nạp từ file; **giữ dòng `#END` ở cuối file**):

```yaml
- id: product-stock                 # id duy nhất
  upstream_id: 2                    # phải tồn tại trong upstreams (ISSUE-13)
  uris:
  - /api/products/*/stock           # thêm cả path gốc nếu cần (ISSUE-15)
  plugins:
    openid-connect:                 # bỏ khối này nếu route công khai
      client_id: ecommerce-client
      client_secret: $env://KEYCLOAK_CLIENT_SECRET
      discovery: http://keycloak:8080/realms/ecommerce/.well-known/openid-configuration
      realm: ecommerce
      bearer_only: true
      use_jwks: true
      token_signing_alg_values_expected: RS256
    proxy-rewrite:
      regex_uri: ["^/(.*)", "/product/$1"]
```

Sau đó `make restart SVC=apisix`. Môi trường k8s dùng manifest riêng `deploy/k8s/gateway/02-routes.yaml` — cập nhật song song. Lưu ý không có plugin kiểm role ([ISSUE-01](../07-plan/known-issues-tech-debt.md)). Chi tiết: [`gateway-routing.md`](../05-api/gateway-routing.md).

## 7. Test, lint, format

| Mục đích | Make | PowerShell | Lệnh gốc |
| :--- | :--- | :--- | :--- |
| Build Go | `make build` | `dev.ps1 build` | `go build ./services/<mod>/...` |
| Test Go | `make go-test` | `dev.ps1 go-test` | `go test` 8 module |
| Vet Go | `make go-vet` | `dev.ps1 go-vet` | |
| Format Go | `make go-fmt` / `make go-fmt-check` | tương tự | `gofmt -w` / `gofmt -l` |
| `go mod tidy` | `make go-tidy` | `dev.ps1 go-tidy` | |
| Test Bun | `make bun-test` | `dev.ps1 bun-test` | `bun install && bun test` × 4 |
| Tất cả backend | `make test` (go-test + bun-test) | `dev.ps1 test` | |
| Gate trước commit | `make check` (gofmt-check + vet + go-test + bun-test) | `dev.ps1 check` | |
| Frontend type check | `make web-check` | `dev.ps1 web-check` | `bun run check` |
| Frontend lint | `make web-lint` | `dev.ps1 web-lint` | `prettier --check . && eslint .` |
| Frontend format | — | — | `cd frontend && bun run format` |
| Frontend unit test | `make web-test` | `dev.ps1 web-test` | `cd frontend && bun test` |
| i18n parity | — | — | `cd frontend && bun run i18n:check` |
| Typecheck service Bun | — | — | `cd services/<svc> && bun run typecheck` |

CI (`.github/workflows/ci.yml`) chạy: gofmt check + build/vet/test mọi module Go; `bun install && bun test` cho 4 service Bun; frontend `check`, `lint`, `i18n:check`, `bun test`; sau đó build image. Không có test tích hợp/E2E ([ISSUE-46](../07-plan/known-issues-tech-debt.md)).

## 8. Kiểm tra parity i18n

`bun run i18n:check` chạy `frontend/scripts/locale-parity.ts`: với mỗi app có thư mục `messages/`, so khóa của `vi.json` (base) với các file còn lại, in `missing`/`extra` và exit 1 nếu lệch. Luôn thêm key vào **cả hai** file trong cùng commit.

## 9. Gọi API thủ công

```bash
# public qua gateway
curl http://localhost:9080/api/products
curl "http://localhost:9080/storefront/catalog-search?keyword=ao&page=0&size=12"

# trực tiếp service (bỏ qua gateway)
curl http://localhost:8086/actuator/health
curl http://localhost:8084/api/orders/all?page=0&size=10
```

Lấy token để gọi route bảo vệ: client `ecommerce-client` có `directAccessGrantsEnabled` — có thể xin token qua endpoint token của Keycloak (`/realms/ecommerce/protocol/openid-connect/token`, grant `password`) bằng user seed; chỉ dùng cho local.

## 10. Troubleshooting

| Triệu chứng | Nguyên nhân (as-is) | Cách xử lý / tham chiếu |
| :--- | :--- | :--- |
| Danh sách sản phẩm/danh mục/đơn hàng trên UI trống dù DB có dữ liệu | BFF đọc `.content` nhưng service trả mảng thô | [ISSUE-48](../07-plan/known-issues-tech-debt.md) |
| Đăng nhập xong vẫn hiện "Đăng nhập"; admin luôn `/access-denied` | `/users/me` 404 hoặc không có `roles` | ISSUE-100; [ISSUE-04](../07-plan/known-issues-tech-debt.md) |
| Đăng nhập lỗi ngẫu nhiên sau khi restart auth-service | SSO state/ticket trong bộ nhớ | [ISSUE-06](../07-plan/known-issues-tech-debt.md) |
| `/api/favourites`, tax trả 404/502 qua gateway | upstream 7, 13 không tồn tại | [ISSUE-13](../07-plan/known-issues-tech-debt.md) — gọi trực tiếp port 8086/8093 |
| `/promotion/backoffice/promotions` (không có `/` cuối) 404 | route `/*` không khớp path gốc | [ISSUE-15](../07-plan/known-issues-tech-debt.md) |
| Tìm kiếm luôn rỗng | không có producer CDC cho index `product` | [ISSUE-19](../07-plan/known-issues-tech-debt.md) |
| Không có email/notification sau thanh toán | lệch tên topic | [ISSUE-17](../07-plan/known-issues-tech-debt.md) |
| Upload media lỗi bucket | bucket chưa được tạo | [ISSUE-32](../07-plan/known-issues-tech-debt.md) |
| Đặt hàng từ trang chủ lỗi | sản phẩm mock không tồn tại | [ISSUE-38](../07-plan/known-issues-tech-debt.md) |
| Thêm cột vào `InitSchema` nhưng DB không đổi | `CREATE TABLE IF NOT EXISTS` không sửa bảng cũ | [ISSUE-28](../07-plan/known-issues-tech-debt.md) — `ALTER` thủ công hoặc xóa volume `postgres_data` |
| Health `UP` nhưng mọi API 500 | health không kiểm DB | [ISSUE-36](../07-plan/known-issues-tech-debt.md) — xem log service |
| Sửa code frontend không thấy trong `make up` | compose dùng image ghcr | ISSUE-97 — dùng `make web-dev` |
| Service chạy bằng `go run` không nhận request qua gateway | upstream trỏ tên container | ISSUE-97 — mục 4 |
| 429 Too Many Requests khi test tải qua BFF | `limit-count` theo IP, mọi request đến từ BFF | ISSUE-61 |
| Đăng nhập ở web làm thay đổi phiên admin (cùng localhost) | trùng tên cookie | [ISSUE-11](../07-plan/known-issues-tech-debt.md) |
| Sửa category của sản phẩm trên admin không lưu | payload `categoryId` phẳng | ISSUE-60 |
| `svelte-check` báo thiếu `#lib/paraglide/...` | chưa compile Paraglide | chạy `bun run check` hoặc `bun run dev` một lần |
