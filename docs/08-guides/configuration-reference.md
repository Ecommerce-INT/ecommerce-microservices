# Configuration Reference — biến môi trường & cấu hình ứng dụng

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> Tham chiếu hợp nhất **mọi biến môi trường được code đọc** (11 service + 2 frontend), biến được set nhưng không dùng, khóa cấu hình APISIX và thiết lập realm Keycloak liên quan tới ứng dụng. **Không ghi giá trị secret** (P-5): default là secret được ghi "(giá trị dev hard-code — không ghi)"; giá trị trong `.env` chỉ ghi tên biến. Hạ tầng k8s (`deploy/k8s/configmap.yaml`, `secrets.yaml`) thuộc bộ tài liệu riêng (P-7).

## 1. Cơ chế đọc cấu hình

| Thành phần | Cơ chế | Biến rỗng `""` | Nguồn |
| :--- | :--- | :--- | :--- |
| Go (7 service) | `config.GetEnv(key, default)`; riêng `DATABASE_URL`, `SPRING_DATASOURCE_URL` đọc bằng `os.Getenv` | coi như **không set** → dùng default | `services/shared/config/config.go:5-10` |
| Bun (4 service) | `process.env.X \|\| "default"` tại module dùng, đọc **một lần lúc import** | coi như không set → default | vd `services/rating-service/src/db.ts:5-13` |
| Frontend (web, admin) | `process.env.X ?? default` trong `@ecommerce/lib/server/env` | **không** coi là thiếu: `API_INTERNAL_URL=""` → base URL rỗng | `frontend/packages/lib/src/server/env.ts:11-27` |
| SvelteKit `adapter-bun` | Runtime adapter đọc `PORT`, `HOST`, `ORIGIN`, `PROTOCOL_HEADER`… (không phải code ứng dụng) | — | `frontend/apps/web/vite.config.ts:2,24`; `frontend/apps/web/.env.example` |
| docker compose | `.env` gốc repo chỉ dùng để **nội suy** `${VAR}` trong `docker-compose.yml`; container chỉ nhận các biến liệt kê ở khối `environment:` (không có `env_file`) | — | `docker-compose.yml` |
| Chạy trên host | Go không tự đọc `.env`; Bun chỉ tự nạp `.env` trong thư mục làm việc → cần `set -a; source .env` | — | [`developer-guide.md`](developer-guide.md) |

**Thứ tự giải quyết kết nối PostgreSQL** (Go và Bun giống nhau): `DATABASE_URL` (dùng nguyên) → nếu không có: dựng từ `POSTGRES_USER/PASSWORD/HOST/PORT/DB`, trong đó `SPRING_DATASOURCE_URL` (dạng `jdbc:postgresql://host:port/db`) **ghi đè** host, port, db; luôn thêm `?sslmode=disable` (`services/auth-service/main.go:25-53`, `services/rating-service/src/db.ts:4-32`).

## 2. Biến theo nhóm

Cột **Bắt buộc?**: "Không" = có default trong code (service vẫn khởi động); "Thực tế" = có default nhưng default sai cho môi trường compose/k8s. Cột **Compose** = có set cho service đó trong `docker-compose.yml` hay không.

### 2.1 Chung

| Biến | Default trong code | Bắt buộc? | Dùng bởi | Tác dụng | Nguồn | Compose |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `PORT` | auth `8088`, product `8086`, order `8084`, payment `8085`, shipping `8087`, inventory `8082`, search `8094`, rating `8089`, media `8083`, notification `8090`, promotion `8093` | Không | 11 service | Cổng HTTP lắng nghe | `services/*/main.go:57` (search `:27`); `services/rating-service/src/index.ts:156`, `media-service/src/index.ts:139`, `notification-service/src/index.ts:194`, `promotion-service/src/index.ts:267` | Không (dùng default) |

### 2.2 PostgreSQL (10 service có DB — search-service không có DB)

| Biến | Default trong code | Bắt buộc? | Dùng bởi | Tác dụng | Nguồn | Compose |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `DATABASE_URL` | — | Không | 10 service | DSN đầy đủ, ưu tiên cao nhất | Go `main.go:27`; Bun `src/db.ts:5-6` | Không |
| `SPRING_DATASOURCE_URL` | — | Thực tế (compose) | 9 service (trừ notification) | Ghi đè host/port/db từ chuỗi JDBC (di sản Spring) | Go `main.go:37-50`; Bun `src/db.ts:15-29` | Có: auth `:186`, product `:207`, order `:226`, payment `:248`, shipping `:270`, inventory `:289`, rating `:308`, media `:327`, promotion `:352` |
| `POSTGRES_USER` | `postgres` | Thực tế | 10 service | User DB | Go `main.go:31`; Bun `src/db.ts:9` | Có (`${POSTGRES_USER}`) cho cả 10 |
| `POSTGRES_PASSWORD` | (giá trị dev hard-code — không ghi) | Thực tế | 10 service | Mật khẩu DB | Go `main.go:32`; Bun `src/db.ts:10` | Có (`${POSTGRES_PASSWORD}`) |
| `POSTGRES_HOST` | `localhost` | Thực tế | 10 service | Host DB (bị `SPRING_DATASOURCE_URL` ghi đè) | Go `main.go:33`; Bun `src/db.ts:11` | Chỉ notification (`:393`) |
| `POSTGRES_PORT` | `5432` | Không | 10 service | Port DB | Go `main.go:34`; Bun `src/db.ts:12` | Chỉ notification (`:394`) |
| `POSTGRES_DB` | auth `authservice`, product `productservice`, order `orderservice`, payment `paymentservice`, shipping `shippingservice`, inventory `inventoryservice`, rating `ratingservice`, media `mediaservice`, promotion `promotionservice`, notification **`ecommerce`** | Không | 10 service | Tên DB (notification dùng DB chung — ISSUE-33) | Go `main.go:35`; Bun `src/db.ts:13` | Chỉ notification (`:395`, `ecommerce`) |

Pool: Go `pgxpool` `MaxConns 25` (`services/shared/database/postgres.go:27-38`); Bun `Bun.SQL` `max 20`, `idleTimeout 30`, `connectionTimeout 5`, `prepare: false` (`src/db.ts:36-42`) — không cấu hình qua biến.

### 2.3 Keycloak / SSO (auth-service)

| Biến | Default trong code | Bắt buộc? | Tác dụng | Nguồn | Compose |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `KEYCLOAK_SERVER_URL` | `http://localhost:8080` | Thực tế | URL nội bộ cho token/logout/admin API | `services/auth-service/main.go:75`; `internal/service/keycloak_client.go:77-98` | Có `http://keycloak:8080` (`:181`) |
| `KEYCLOAK_PUBLIC_SERVER_URL` | `http://keycloak.ecommerce.local` | Thực tế | URL trình duyệt dùng cho trang authorize | `main.go:76`; `keycloak_client.go:81-87` | Có `${KEYCLOAK_PUBLIC_SERVER_URL:-http://localhost:8080}` (`:182`) |
| `KEYCLOAK_REALM` | `ecommerce` | Không | Realm | `main.go:77` | Không |
| `KEYCLOAK_CLIENT_ID` | `ecommerce-client` | Không | Client OIDC | `main.go:78` | Không |
| `KEYCLOAK_CLIENT_SECRET` | `""` | Không | Gửi kèm khi đổi code/refresh/logout nếu khác rỗng (client là public — ISSUE-12) | `main.go:79`; `keycloak_client.go:114-116` | **Không** cho auth-service (chỉ APISIX `:161`) — ISSUE-117 |
| `KEYCLOAK_ADMIN_REALM` | `master` | Không | Realm lấy admin token | `main.go:80` | Không |
| `KEYCLOAK_ADMIN_CLIENT_ID` | `admin-cli` | Không | Client lấy admin token | `main.go:81` | Không |
| `KEYCLOAK_ADMIN_USERNAME` | (giá trị dev hard-code — không ghi) | Thực tế | Tài khoản admin tạo/xóa user Keycloak | `main.go:82` | **Không** — dùng default code (ISSUE-117, ISSUE-135) |
| `KEYCLOAK_ADMIN_PASSWORD` | (giá trị dev hard-code — không ghi) | Thực tế | Như trên | `main.go:83` | **Không** (ISSUE-117, ISSUE-135) |
| `SSO_BACKEND_CALLBACK_URL` | `http://api.ecommerce.local/api/v1/auth/callback` | Thực tế | `redirect_uri` gửi Keycloak lúc authorize | `main.go:84`; `keycloak_client.go:62-64,106` | Có `http://localhost:8088/api/v1/auth/callback` (`:183`) — trỏ thẳng port service, không qua gateway |
| `SSO_DEFAULT_FRONTEND_REDIRECT` | `http://ecommerce.local/auth/callback` | Thực tế | Đích mặc định khi login không có `redirect_uri` / state không hợp lệ | `main.go:85` | Có `http://localhost:3000/auth/callback` (`:184`) |
| `SSO_SCOPE` | `openid profile email` | Không | Scope authorize | `main.go:86`; `keycloak_client.go:59-61` | Không |

### 2.4 Kafka

| Biến | Default trong code | Bắt buộc? | Dùng bởi | Tác dụng | Nguồn | Compose |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `KAFKA_SERVERS` | `localhost:9092` | Thực tế | payment, search, notification | Bootstrap broker (CSV). payment: chuỗi rỗng → tắt producer | `services/payment-service/main.go:74`; `services/search-service/main.go:31`; `services/notification-service/src/kafka.ts:6` | Có `kafka:9092` cho payment `:251`, search `:374`, notification `:398` (và order `:229` — **không dùng**, ISSUE-37) |
| `PRODUCT_TOPIC_NAME` | `dbproduct.public.product` | Không | search | Topic CDC sản phẩm | `services/search-service/main.go:32` | Không |
| `PAYMENT_KAFKA_CONSUMER_GROUP_ID` | `notification-group` | Không | notification | Consumer group | `services/notification-service/src/kafka.ts:8` | Không |

Hằng số không cấu hình được: topic `SUCCESSFUL`, `paymentCompleted` (`services/payment-service/internal/service/payment_service.go:19-20`); topic `profile-onboarding-topic`, `status-payment-successful` và `clientId notification-service` (`kafka.ts:7,23`); group `search-group` (`services/search-service/main.go:59`). Xem [`event-catalog.md`](../05-api/event-catalog.md).

### 2.5 OpenSearch (search-service)

| Biến | Default trong code | Bắt buộc? | Tác dụng | Nguồn | Compose |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `OPENSEARCH_URL` | `http://localhost:9200` | Thực tế | Địa chỉ cluster | `services/search-service/main.go:28` | Có `http://opensearch:9200` (`:371`) |
| `OPENSEARCH_USERNAME` | `""` | Không | Basic auth | `main.go:29` | Có, rỗng (`:372`) — security plugin tắt |
| `OPENSEARCH_PASSWORD` | `""` | Không | Basic auth | `main.go:30` | Có, rỗng (`:373`) |

### 2.6 Object storage (media-service)

| Biến | Default trong code | Bắt buộc? | Tác dụng | Nguồn | Compose |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `STORAGE_ENDPOINT` | `http://localhost:9000` | Thực tế | Endpoint S3 (RustFS), path-style | `services/media-service/src/storage.ts:3,9-16` | Có `http://rustfs:9000` (`:330`) |
| `STORAGE_REGION` | `us-east-1` | Không | Region S3 | `storage.ts:4` | Không |
| `STORAGE_ACCESS_KEY` | (giá trị dev hard-code — không ghi) | Thực tế | Access key | `storage.ts:5` | Có `${STORAGE_ACCESS_KEY}` (`:331`) |
| `STORAGE_SECRET_KEY` | (giá trị dev hard-code — không ghi) | Thực tế | Secret key | `storage.ts:6` | Có `${STORAGE_SECRET_KEY}` (`:332`) |
| `STORAGE_BUCKET` | `ecommerce-media` | Không | Bucket (không tự tạo — ISSUE-32) | `storage.ts:7` | Có `${STORAGE_BUCKET}` (`:333`) |
| `PUBLIC_URL` | `http://localhost:8083` | Thực tế | Tiền tố URL file trong `MediaVm.url` | `services/media-service/src/service.ts:21-25` | **Không** (ISSUE-32, ISSUE-63) |

### 2.7 Mail (notification-service)

| Biến | Default trong code | Bắt buộc? | Tác dụng | Nguồn | Compose |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `MAIL_HOST` | `smtp.gmail.com` | Không | SMTP host | `services/notification-service/src/mailer.ts:4` | Không |
| `MAIL_PORT` | `587` | Không | SMTP port; `465` → `secure: true` | `mailer.ts:5,15` | Không |
| `MAIL_USERNAME` | `""` | Không | Tài khoản SMTP + địa chỉ `from`; rỗng (hoặc thiếu password) → **chế độ mock** | `mailer.ts:6,11-18,37` | Có `${MAIL_USERNAME}` (`:399`) — `.env` có giá trị mẫu → mock không kích hoạt (ISSUE-117) |
| `MAIL_PASSWORD` | `""` | Không | Mật khẩu SMTP | `mailer.ts:7` | Có `${MAIL_PASSWORD}` (`:400`) |

### 2.8 Frontend (web `:3000`, admin `:3001`)

| Biến | Default | Bắt buộc? | Dùng bởi | Tác dụng | Nguồn | Compose |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `API_INTERNAL_URL` | `http://localhost:9080` | Thực tế | web, admin (server) | Base URL gateway cho mọi lời gọi BFF; dấu `/` cuối bị cắt | `frontend/packages/lib/src/server/env.ts:11-14` | Có `http://apisix:9080` (web `:422`, admin `:442`) |
| `API_PUBLIC_URL` | = `API_INTERNAL_URL` | Thực tế (khi khác internal) | web, admin | Origin trình duyệt dùng cho redirect SSO `/api/v1/auth/login` | `env.ts:24-27`; `frontend/apps/web/src/routes/login/+page.server.ts:24` | Có `http://localhost:9080` (`:424`, `:443`) |
| `ORIGIN` | — (adapter suy từ request) | Thực tế khi chạy HTTP sau proxy | adapter-bun | Origin công khai (CSRF check form/remote) | `frontend/apps/web/.env.example`; `frontend/apps/admin/.env.example` | Có `http://localhost:3000` / `:3001` (`:426`, `:444`) |
| `PORT` | `3000` (web Dockerfile) | Không | adapter-bun | Cổng lắng nghe | `frontend/apps/web/Dockerfile` (`ENV PORT=3000`) | Có `3000` / `3001` (`:420`, `:441`) |
| `HOST` | `0.0.0.0` (adapter) | Không | adapter-bun | Địa chỉ bind | `.env.example` | Không |
| `PROTOCOL_HEADER` | — | Không | adapter-bun | Header proto khi sau reverse proxy | `.env.example` | Không |
| `NODE_ENV` | — | Không | (không đọc trong code ứng dụng) | Chỉ ảnh hưởng thư viện | `docker-compose.yml:419,440` | Có `production` |

Không có biến cho tên cookie và thời hạn cookie phiên — hard-code trong `frontend/packages/lib/src/server/session.ts`; locale cấu hình trong `frontend/apps/{web,admin}/project.inlang` (không qua biến môi trường).

## 3. Ma trận biến set cho từng container ứng dụng (compose)

| Container | Biến (`docker-compose.yml`) | Dòng |
| :--- | :--- | :--- |
| auth-service | `JWT_ISSUER_URI`¹, `KEYCLOAK_SERVER_URL`, `KEYCLOAK_PUBLIC_SERVER_URL`, `SSO_BACKEND_CALLBACK_URL`, `SSO_DEFAULT_FRONTEND_REDIRECT`, `SSO_ALLOWED_REDIRECT_URIS`¹, `SPRING_DATASOURCE_URL`, `POSTGRES_USER`, `POSTGRES_PASSWORD` | `:179-188` |
| product / shipping / inventory / rating / promotion | `JWT_ISSUER_URI`¹, `SPRING_DATASOURCE_URL`, `POSTGRES_USER`, `POSTGRES_PASSWORD` | `:205-209`, `:268-272`, `:287-291`, `:306-310`, `:350-354` |
| order-service | như trên + `KAFKA_SERVERS`¹ | `:224-229` |
| payment-service | như trên + `KAFKA_SERVERS` | `:246-251` |
| media-service | như trên + `STORAGE_ENDPOINT`, `STORAGE_ACCESS_KEY`, `STORAGE_SECRET_KEY`, `STORAGE_BUCKET` | `:325-333` |
| search-service | `JWT_ISSUER_URI`¹, `OPENSEARCH_URL`, `OPENSEARCH_USERNAME`, `OPENSEARCH_PASSWORD`, `KAFKA_SERVERS` | `:369-374` |
| notification-service | `JWT_ISSUER_URI`¹, `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `KAFKA_SERVERS`, `MAIL_USERNAME`, `MAIL_PASSWORD` | `:391-400` |
| frontend (web) | `NODE_ENV`, `PORT`, `API_INTERNAL_URL`, `API_PUBLIC_URL`, `ORIGIN` | `:418-426` |
| admin | `NODE_ENV`, `PORT`, `API_INTERNAL_URL`, `API_PUBLIC_URL`, `ORIGIN` | `:439-444` |
| apisix | `KEYCLOAK_CLIENT_SECRET` (có fallback placeholder trong compose — không ghi) | `:160-161` |

¹ Được set nhưng **không có code nào đọc** (ISSUE-37).

## 4. Biến được set nhưng không dùng (ISSUE-37)

| Biến | Set ở đâu | Ghi chú |
| :--- | :--- | :--- |
| `JWT_ISSUER_URI` | compose cho cả 11 service; `.env` | Không service nào verify JWT (ISSUE-02); di sản Spring Resource Server |
| `SSO_ALLOWED_REDIRECT_URIS` | compose auth-service (`:185`) | Code không đọc → open redirect (ISSUE-05) |
| `KAFKA_SERVERS` | compose order-service (`:229`) | order-service không dùng Kafka |
| `KEYCLOAK_JWK_SET_URI` | `.env` | Không đọc (APISIX dùng `discovery`) |
| `REDIS_HOST`, `REDIS_PORT` | `.env` | Không ứng dụng nào dùng Redis; `REDIS_PASSWORD` chỉ nội suy cho container `redis` (`:28`) |
| `EUREKA_URI` | `.env` | Di sản Spring Cloud |
| `APP_VERSION` | `.env` | Di sản Swagger UI |
| `ECOMMERCE_SECURITY_ENABLED`, `ECOMMERCE_SECURITY_CSRF_DISABLED`, `ECOMMERCE_SECURITY_STATELESS_SESSION` | `.env` | Di sản common-spring |
| `ECOMMERCE_CORS_ENABLED`, `CORS_ALLOWED_ORIGIN_PATTERNS`, `CORS_ALLOW_CREDENTIALS`, `CORS_MAX_AGE_SECONDS` | `.env` | CORS hard-code `*` trong service |
| `REST_CLIENT_CONNECT_TIMEOUT`, `REST_CLIENT_READ_TIMEOUT`, `REST_CLIENT_POOL_MAX_TOTAL`, `REST_CLIENT_POOL_MAX_PER_ROUTE`, `REST_CLIENT_POOL_KEEP_ALIVE`, `REST_CLIENT_POOL_EVICT_IDLE_AFTER` | `.env` | Di sản |
| `I18N_DEFAULT_LOCALE`, `I18N_CACHE_DURATION` | `.env` | Backend không có i18n |
| `WEB_CORRELATION_ENABLED`, `WEB_EXCEPTION_HANDLER_ENABLED` | `.env` | Correlation id luôn bật, không cấu hình |
| `HTTP_LOG_REQUEST_ENABLED`, `HTTP_LOG_REQUEST_BODY`, `HTTP_LOG_REQUEST_MAX_BYTES`, `HTTP_LOG_RESPONSE_ENABLED`, `HTTP_LOG_RESPONSE_BODY`, `HTTP_LOG_RESPONSE_MAX_BYTES`, `HTTP_LOG_PERF_ENABLED`, `HTTP_LOG_PERF_THRESHOLD_MS` | `.env` | Di sản |
| `POSTGRES_HOST`, `POSTGRES_PORT`, `KAFKA_SERVERS`, `STORAGE_ENDPOINT`, `STORAGE_REGION`, `OPENSEARCH_*`, `KEYCLOAK_SERVER_URL`, `KEYCLOAK_REALM`, `KEYCLOAK_CLIENT_ID`, `KEYCLOAK_ADMIN_*` trong `.env` | `.env` | **Không** đi vào container (compose không tham chiếu `${…}` các biến này cho service); chỉ có tác dụng khi `source .env` để chạy service trên host |

## 5. Biến nội suy cho container hạ tầng (tham khảo)

| Biến `.env` | Container nhận | Dòng compose |
| :--- | :--- | :--- |
| `POSTGRES_USER`, `POSTGRES_PASSWORD` | postgres; keycloak (`KC_DB_USERNAME`, `KC_DB_PASSWORD`); 10 service | `:10-11`, `:133-134` |
| `REDIS_PASSWORD` | redis (`--requirepass`) | `:28` |
| `STORAGE_ACCESS_KEY`, `STORAGE_SECRET_KEY` | rustfs (`RUSTFS_ACCESS_KEY`, `RUSTFS_SECRET_KEY`); media-service | `:107-108`, `:331-332` |
| `KEYCLOAK_ADMIN`, `KEYCLOAK_ADMIN_PASSWORD` | keycloak (bootstrap admin) — **không** truyền cho auth-service | `:127-128` |
| `KEYCLOAK_CLIENT_SECRET` | apisix | `:161` |
| `KEYCLOAK_PUBLIC_SERVER_URL` | auth-service | `:182` |
| `STORAGE_BUCKET`, `MAIL_USERNAME`, `MAIL_PASSWORD` | media-service, notification-service | `:333`, `:399-400` |

## 6. APISIX — khóa cấu hình liên quan ứng dụng

Chi tiết route: [`gateway-routing.md`](../05-api/gateway-routing.md).

| File | Khóa | Giá trị | Ảnh hưởng tới ứng dụng | Dòng |
| :--- | :--- | :--- | :--- | :--- |
| `deploy/apisix/config.yaml` | `deployment.role` / `role_data_plane.config_provider` | `data_plane` / `yaml` | Standalone, route đọc từ `apisix.yaml` (không Admin API) | `:5-8` |
| | `apisix.node_listen` | `9080` | Cổng BFF gọi (`API_INTERNAL_URL`) | `:10` |
| | `apisix.enable_server_tokens` | `false` | Ẩn version trong header `Server` | `:11` |
| | `plugin_attr.prometheus.export_addr` | `0.0.0.0:9091` | Metric gateway | `:12-16` |
| `deploy/apisix/apisix.yaml` | `upstreams[].nodes` | `<service>:<port>` (id 1–6, 8–12; **thiếu 7, 13**) | Tên container compose; service chạy trên host không nhận traffic (ISSUE-97); ISSUE-13 | `:1-45` |
| | `global_rules[1].plugins.cors` | `allow_origins` 3 origin, methods `GET,POST,PUT,DELETE,PATCH,OPTIONS`, headers `Authorization,Content-Type,X-Requested-With`, `max_age 3600` | Thiếu `localhost:3001` (ISSUE-14) | `:50-54` |
| | `global_rules[1].plugins.limit-count` | `count 1200`, `time_window 60`, `key remote_addr`, `policy local`, `rejected_code 429` | Quota chung cho cả BFF (ISSUE-61) | `:55-61` |
| | `routes[].plugins.openid-connect` | `client_id ecommerce-client`, `client_secret $env://KEYCLOAK_CLIENT_SECRET`, `discovery http://keycloak:8080/realms/ecommerce/.well-known/openid-configuration`, `realm ecommerce`, `bearer_only true`, `use_jwks true`, `token_signing_alg_values_expected RS256` | Chỉ xác thực, không kiểm role (ISSUE-01) | vd `:142-159` |
| | `routes[].plugins.proxy-rewrite.regex_uri` | `^/(.*)` → `/<ctx>/$1` | Route kiểu `/api/...` được thêm context-root | vd `:66-70` |
| | `routes[].uri(s)` | dạng `/<svc>/.../*` cho 9 route | Không khớp path gốc collection (ISSUE-15) | `:268-375` |

## 7. Keycloak realm `ecommerce` (không ghi secret)

Nguồn: `deploy/keycloak/import/ecommerce-realm.json` (import khi khởi động `start-dev --import-realm`, `docker-compose.yml:125`). Keycloak dùng DB `keycloak` trên cùng PostgreSQL (`docker-compose.yml:129-134`).

| Khóa | Giá trị | Ảnh hưởng | Dòng |
| :--- | :--- | :--- | :--- |
| `realm` / `enabled` / `displayName` | `ecommerce` / `true` / `Ecommerce Realm` | `KEYCLOAK_REALM`, `realm` trong APISIX | `:2-4` |
| `registrationAllowed` | `false` | Không tự đăng ký trên trang Keycloak (ISSUE-71) | `:5` |
| `loginWithEmailAllowed` | `true` | Đăng nhập bằng email | `:6` |
| `duplicateEmailsAllowed` | `false` | | `:7` |
| `resetPasswordAllowed` | `true` | Link "quên mật khẩu" (cần SMTP — realm không cấu hình `smtpServer`) | `:8` |
| `roles.realm` | `ADMIN`, `USER`, `MANAGER` | Role trong token (`realm_access.roles`) — **không** được backend kiểm (ISSUE-01); DB `authservice` seed `USER`, `PM`, `ADMIN` khác bộ này | `:10-22` |
| Client `clientId` | `ecommerce-client` (`Ecommerce Public Client`) | `KEYCLOAK_CLIENT_ID`, APISIX `client_id` | `:27-28` |
| `publicClient` | `true` | Không cần secret (ISSUE-12) | `:31` |
| `standardFlowEnabled` | `true` | Authorization code flow (SSO) | `:33` |
| `directAccessGrantsEnabled` | `true` | Cho phép grant `password` (ROPC) (ISSUE-133) | `:32` |
| `implicitFlowEnabled` / `serviceAccountsEnabled` | `false` / `false` | | `:34-35` |
| `fullScopeAllowed` | `true` | Mọi realm role vào token | `:36` |
| `redirectUris` | `http://localhost:*`, `http://127.0.0.1:*`, `http://api.ecommerce.local/*` | Khớp `SSO_BACKEND_CALLBACK_URL` compose (`http://localhost:8088/...`); wildcard port (ISSUE-133) | `:37-41` |
| `webOrigins` | `*` | CORS Keycloak mở mọi origin (ISSUE-133) | `:42-44` |
| `attributes."post.logout.redirect.uris"` | `http://localhost:*##http://127.0.0.1:*##http://ecommerce.local/*##http://admin.ecommerce.local/*` | | `:46` |
| `users` | `testuser` (role `USER`), `adminuser` (`ADMIN`, `MANAGER`), `emailVerified: true`, có credential (không ghi) | Seed dev; không có bản ghi tương ứng trong DB `authservice` (ISSUE-101) | `:50-88` |
| Token lifespan / session timeout | không đặt → mặc định Keycloak | Thời hạn access/refresh token | — |

## 8. Issue liên quan

| ID | Mức độ | Tóm tắt |
| :--- | :--- | :--- |
| ISSUE-37 | Low | Biến set nhưng không dùng (§4) |
| ISSUE-05 | Critical | `SSO_ALLOWED_REDIRECT_URIS` không được đọc |
| ISSUE-32 / ISSUE-63 | Medium | `PUBLIC_URL` không set; URL media sai |
| ISSUE-33 | Medium | notification dùng DB `ecommerce` chung |
| ISSUE-97 | Low | Upstream APISIX/Kafka quảng bá tên container — khó chạy service trên host |
| ISSUE-133 | Medium | (mới) Client Keycloak bật ROPC, `webOrigins: *`, redirect wildcard port |
| ISSUE-117 | Low | (mới) Compose thiếu biến cho auth-service/notification → phụ thuộc default trong code |
| ISSUE-135 | Medium | (mới) Credential mặc định hard-code trong code, không fail-fast |

## 9. Tham chiếu source

- `services/shared/config/config.go:5-10`
- `services/{auth,product,order,payment,shipping,inventory}-service/main.go:25-57`; `services/auth-service/main.go:74-87`; `services/search-service/main.go:27-32,59`; `services/payment-service/main.go:74-76`
- `services/{rating,media,notification,promotion}-service/src/db.ts:4-42`; `services/media-service/src/storage.ts:3-16`, `src/service.ts:21`; `services/notification-service/src/mailer.ts:4-18`, `src/kafka.ts:6-8`
- `frontend/packages/lib/src/server/env.ts:11-27`; `frontend/apps/{web,admin}/.env.example`; `frontend/apps/web/Dockerfile`
- `docker-compose.yml`; `.env` (chỉ tên biến)
- `deploy/apisix/config.yaml:1-16`; `deploy/apisix/apisix.yaml:1-375`
- `deploy/keycloak/import/ecommerce-realm.json:1-89`
