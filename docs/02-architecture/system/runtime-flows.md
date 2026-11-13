# Runtime Flows

> **Phiên bản:** 1.0 · **Trạng thái:** Final · **Cập nhật:** 2026-11-12 · **Phạm vi:** Ứng dụng (as-is)

> Sequence diagram cho các luồng có ý nghĩa kiến trúc, theo code hiện tại. Đường dẫn ghi theo **phía gateway** (`/api/...`) và phía service sau rewrite khi cần. Chi tiết nghiệp vụ nằm ở tài liệu FEAT tương ứng.

## RF-01 — Đăng nhập SSO (ticket exchange)

Tài liệu tính năng: [FEAT-01](../../04-features/FEAT-01-sso-login-logout.md).

```mermaid
sequenceDiagram
    autonumber
    actor U as Người dùng
    participant BFF as SvelteKit (web/admin)
    participant GW as APISIX
    participant AUTH as auth-service
    participant KC as Keycloak

    U->>BFF: GET /login?redirect=/checkout
    BFF->>BFF: set cookie auth_redirect (10 phút)<br/>ssoUrl = API_PUBLIC_URL/api/v1/auth/login?redirect_uri=ORIGIN/auth/callback
    U->>GW: GET /api/v1/auth/login?redirect_uri=... (public)
    GW->>AUTH: GET /auth/api/v1/auth/login
    AUTH->>AUTH: CreateLoginState(redirect_uri) — state TTL 5 phút (in-memory)
    AUTH-->>U: 302 → Keycloak authorize (redirect_uri = SSO_BACKEND_CALLBACK_URL)
    U->>KC: Đăng nhập
    KC-->>U: 302 → callback?code&state
    U->>GW: GET /api/v1/auth/callback?code&state
    GW->>AUTH: /auth/api/v1/auth/callback
    AUTH->>AUTH: ConsumeLoginState(state) → frontendRedirect
    AUTH->>KC: POST token (grant=authorization_code, redirect_uri = frontendRedirect ⚠️)
    alt Đổi code thành công
        KC-->>AUTH: access/refresh token
    else Lỗi (vd redirect_uri lệch — ISSUE-04)
        AUTH->>AUTH: dùng mock token ⚠️
    end
    AUTH->>AUTH: StoreTokens → ticket TTL 2 phút
    AUTH-->>U: 302 → ORIGIN/auth/callback?ticket=...
    U->>BFF: GET /auth/callback?ticket=...
    BFF->>GW: GET /api/v1/auth/session?ticket= (không token)
    GW->>AUTH: ConsumeTokens(ticket) — dùng 1 lần
    AUTH-->>BFF: ApiResponse{data: tokens}
    BFF->>GW: GET /api/v1/users/me (Bearer)
    GW->>GW: openid-connect verify JWT (JWKS)
    GW->>AUTH: /auth/api/v1/users/me
    AUTH-->>BFF: UserResponse (không có roles — ISSUE-100) hoặc 404 (ISSUE-101)
    BFF->>BFF: setSession: access_token, refresh_token, session_user, session_expires (httpOnly)
    BFF-->>U: 303 → safeRedirect(auth_redirect)
```

## RF-02 — Làm mới phiên trên mỗi request

Tài liệu: [FEAT-03](../../04-features/FEAT-03-session-refresh.md).

```mermaid
sequenceDiagram
    participant U as Trình duyệt
    participant H as hooks.server.ts (handleSession)
    participant GW as APISIX → auth-service
    participant KC as Keycloak
    U->>H: Bất kỳ request
    H->>H: session(cookies); shouldRefresh? (không có access token hoặc còn < 30s)
    alt Cần refresh và có refresh_token
        H->>GW: POST /api/v1/auth/refresh {refreshToken}
        GW->>KC: grant=refresh_token
        alt Thành công
            KC-->>H: token mới
            H->>H: setSession(tokens) — giữ nguyên session_user
        else Lỗi
            H->>H: clearSession → locals.user = null (guest)
        end
    end
    H->>U: tiếp tục xử lý route với locals.user
```

## RF-03 — Duyệt catalog (SSR + remote query)

Tài liệu: [FEAT-06](../../04-features/FEAT-06-product-catalog.md).

```mermaid
sequenceDiagram
    participant U as Trình duyệt
    participant BFF as web (catalog.remote.ts)
    participant GW as APISIX
    participant P as product-service
    participant DB as productservice DB
    U->>BFF: GET /products?page&sort&categoryId&search
    BFF->>BFF: productQuerySchema.parse
    par
        BFF->>GW: GET /api/products?page&size&sort&categoryId&search (public)
        GW->>P: /product/api/products
        P->>DB: SELECT products JOIN categories (toàn bộ, bỏ qua query)
        P-->>BFF: [Product,...] (mảng thô — ISSUE-48)
    and
        BFF->>GW: GET /api/categories?page=0&size=50
        GW->>P: /product/api/categories
        P-->>BFF: [Category,...] (mảng thô)
    end
    BFF-->>U: HTML SSR (BFF đọc .content — undefined với backend hiện tại)
```

## RF-04 — Checkout / đặt hàng

Tài liệu: [FEAT-09](../../04-features/FEAT-09-checkout-place-order.md).

```mermaid
sequenceDiagram
    autonumber
    actor C as Customer
    participant BFF as web (checkout.remote.ts placeOrder)
    participant GW as APISIX
    participant P as product-service
    participant O as order-service
    participant PAY as payment-service
    participant K as Kafka

    C->>BFF: placeOrder({shipping, paymentMethod, items[]})
    BFF->>BFF: checkoutSchema.parse; có access_token? (không → 401)
    loop mỗi dòng giỏ
        BFF->>GW: GET /api/products/{id}
        GW->>P: lấy giá thật
        P-->>BFF: Product
    end
    BFF->>BFF: computeTotals(lines) → subtotal, shippingFee, total
    BFF->>GW: POST /api/orders {orderDate, orderDesc, orderFee: total, productId: dòng đầu tiên}
    GW->>O: INSERT orders
    O-->>BFF: Order{orderId}
    BFF->>GW: POST /api/payments {orderId, isPayed: method≠COD, paymentStatus}
    GW->>PAY: kiểm tra đã trả? → INSERT payments
    PAY->>K: publish SUCCESSFUL, paymentCompleted (lỗi bị bỏ qua)
    PAY-->>BFF: Payment
    BFF-->>C: {orderId, total} → bước 4 "thành công", xóa giỏ
    Note over BFF,PAY: Không transaction: nếu payment lỗi, order đã tạo vẫn tồn tại (ISSUE-23)
```

## RF-05 — CRUD quản trị (vd sản phẩm)

Tài liệu: [FEAT-23](../../04-features/FEAT-23-admin-product-mgmt.md).

```mermaid
sequenceDiagram
    actor A as Admin
    participant G as admin hooks (handleAdminGuard)
    participant R as products.remote.ts
    participant GW as APISIX
    participant P as product-service
    A->>G: GET /products
    G->>G: có user? isAdmin? (không → /login hoặc /access-denied)
    G->>R: getProducts / getCategoryOptions
    R->>GW: GET /api/products, /api/categories?size=100
    GW->>P: proxy (public)
    A->>R: saveProduct (remote form, productUpsertSchema)
    alt id rỗng
        R->>GW: POST /api/products
    else có id
        R->>GW: PUT /api/products/{id}
    end
    GW->>P: (không kiểm token — ISSUE-102)
    P-->>R: Product
    R-->>A: refresh danh sách, toast
```

## RF-06 — Sự kiện thanh toán → thông báo (thiết kế)

Tài liệu: [FEAT-11](../../04-features/FEAT-11-payment.md), [FEAT-19](../../04-features/FEAT-19-notifications-email.md).

```mermaid
sequenceDiagram
    participant PAY as payment-service
    participant K as Kafka
    participant N as notification-service
    participant DB as ecommerce DB
    participant M as SMTP
    PAY->>K: topic SUCCESSFUL / paymentCompleted<br/>{paymentId, isPayed, paymentStatus, orderId, userId}
    Note over N: subscribe status-payment-successful,<br/>profile-onboarding-topic (group notification-group)
    K--xN: Không khớp topic → không nhận (ISSUE-17)
    opt Nếu topic khớp (hành vi code consumer)
        N->>DB: INSERT payment_notifications
        N->>M: Email tới địa chỉ hard-code (ISSUE-09)
    end
```

## RF-07 — Đồng bộ catalog sang OpenSearch (thiết kế)

Tài liệu: [FEAT-14](../../04-features/FEAT-14-catalog-search-suggest.md).

```mermaid
sequenceDiagram
    participant P as product-service
    participant K as Kafka
    participant S as search-service (consumer, group search-group)
    participant OS as OpenSearch
    Note over P,K: Không có producer / Debezium connector (ISSUE-19)
    K-->>S: (dự kiến) message Debezium topic dbproduct.public.product
    alt op = "d"
        S->>OS: DELETE product/_doc/{before.id}
    else khác
        S->>OS: PUT product/_doc/{after.id} {name, slug, price}
    end
    Note over S: Message lỗi bị bỏ qua im lặng
```

## RF-08 — Upload media

Tài liệu: [FEAT-18](../../04-features/FEAT-18-media-management.md).

```mermaid
sequenceDiagram
    participant C as Client (có token)
    participant GW as APISIX
    participant M as media-service
    participant S3 as RustFS
    participant DB as mediaservice DB
    C->>GW: POST /media/medias/... multipart (multipartFile|file, caption, fileNameOverride)
    GW->>M: verify token → proxy
    M->>M: đọc toàn bộ file vào RAM; key = uuid-tênfile
    M->>S3: PUT object (1 lần)
    M->>DB: INSERT media
    M-->>C: {id, caption, fileName, mediaType, url}
```

## RF-09 — Tìm kiếm

```mermaid
sequenceDiagram
    participant C as Client
    participant GW as APISIX
    participant S as search-service
    participant OS as OpenSearch
    C->>GW: GET /storefront/catalog-search?keyword&brand&category&minPrice&maxPrice&page&size
    GW->>S: /search/storefront/catalog-search (public)
    S->>OS: bool{should: multi_match fuzzy name^3, brand^2, categories; must: filters}
    alt OpenSearch lỗi
        S-->>C: trang rỗng (che lỗi — ISSUE-34)
    else
        OS-->>S: hits
        S-->>C: {products, pageNo, pageSize, totalElements, totalPages, isLast, aggregations: {}}
    end
```

> Storefront hiện **không** gọi search-service; ô tìm kiếm ở header chuyển tới `/products?search=` (product-service).
