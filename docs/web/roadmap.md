# Web Roadmap — SvelteKit 3 + Bun

Sequenced plan for the frontend rewrite. No dates — phases are done when their **exit criteria** hold. Parallel-capable items are marked.

Legend: `[x]` done · `[~]` in progress · `[ ]` todo

---

## Phase 0 — Foundation `[x]`

**Goal:** a buildable, runnable, well-documented SvelteKit workspace with no product features yet.

Deliverables
- [x] `feature/sveltekit-frontend` branch; Next.js frontend removed (history retained)
- [x] Bun workspace root (`bun.lock` committed, no pnpm/Turborepo)
- [x] `apps/web` + `apps/admin` scaffolded (Kit 3, Svelte 5 runes, Vite 8, Tailwind 4, ESLint, Prettier, Paraglide `vi`/`en`)
- [x] `@sveltejs/adapter-bun` wired; `experimental.remoteFunctions` + `compilerOptions.experimental.async` enabled
- [x] Verified: install, production builds, `lint`, `check` (0/0), live `bun ./build` smoke test
- [x] `docs/web/` (this folder): dev guide, roadmap, backlog
- [x] `@ecommerce/lib` — types, Zod schemas, formatters, server API client, session/redirect helpers + 60 unit tests

**Exit criteria:** `bun run build`, `bun run check`, `bun run lint` pass on a clean clone; docs describe the whole workflow; nothing references pnpm/Next in the frontend tree.

---

## Phase 1 — App skeleton `[x]`

**Goal:** the shell every page needs — theme, components, i18n routing, session plumbing.

Deliverables
- [x] `@ecommerce/lib` wired into both apps (workspace dep + Vite SSR handling) with unit tests (`bun test`, 60 tests)
- [x] Brand theme ported into `src/routes/layout.css` (red `#e30019` palette, orange-scale remap, line-clamp/no-scrollbar utilities) + Geist variable font (`@fontsource-variable/geist`, both apps)
- [x] shadcn-svelte base set provided once by `@ecommerce/ui` and consumed by both apps (button, input, label, textarea, badge, card, dialog, alert-dialog, dropdown-menu, select, table, sheet, skeleton, sonner, pagination, breadcrumb, separator, empty-state, table-skeleton)
- [x] `Toaster` mounted in both root layouts; feedback via toasts/dialogs only
- [x] Paraglide i18n routing configured: `vi` unprefixed, `en` at `/en/**` (URL strategy + cookie + base locale), language switcher component
- [x] Message catalogs migrated from the old apps (web 297 × 2, admin 120 × 2 flat keys) — script-assisted flattening, plural-free; `bun run i18n:check` verifies key parity
- [x] Session helpers in `@ecommerce/lib/server`: cookie names, set/clear, refresh-on-expiry, `locals.user`, `safeRedirect`, `isAdmin`
- [x] `hooks.server.ts` extended: Paraglide handle + session bootstrap (auth guard arrives in phase 2/3)
- [x] Storefront cart store (`#lib/stores/cart.svelte.ts`, runes + localStorage key `ecommerce-cart`)

**Exit criteria:** both apps render styled shells in both locales, with a locale switcher and a working toast; session helpers are unit-testable and typed; `check`/`lint` stay green.

---

## Phase 2 — Storefront parity `[x]`

**Goal:** every storefront route of the old app, but server-rendered via remote functions.

Deliverables
- [ ] Layout chrome: announcement bar, header (search, cart badge, profile dropdown, admin link), mobile menu, footer
- [ ] Home: hero carousel, category bar, product rails (quick deals, flash sale, brand strip, USP strip) on ported mock data
- [ ] Catalog `/products`: remote `query`, category filter, price-range filter, sort, **pagination from `totalElements`**, **search wired to the API** (the old app sent nothing)
- [ ] PDP `/products/[id]`: gallery, price/stock, quantity selector, add-to-cart / buy-now, wishlist/share resolved (implement or remove)
- [ ] Cart: line items, quantity update/remove/clear, order summary, free-shipping rule (≥ 500 000 ₫)
- [x] Checkout: 3-step flow (shipping → payment → confirmation) with a Zod schema + remote `form`; order + payment creation; success/error states
- [x] Orders: session-guarded list with status badges
- [x] Auth: `/login` server redirect to gateway SSO, `/auth/callback` server-side ticket exchange → httpOnly cookies, logout clears cookies and lands public

**Exit criteria:** a user can browse → add to cart → log in via SSO → checkout → see their order; cart survives reloads; no token appears in JS-accessible storage; storefront is fully SSR (view-source contains data).

---

## Phase 3 — Backoffice parity `[x]`

**Goal:** the admin app the old one promised, with the role check it never had.

Deliverables
- [x] Admin shell: sidebar, topbar, mobile drawer (sheet); **strict ADMIN guard in `apps/admin/src/hooks.server.ts`**
- [x] Dashboard: count cards from `totalElements` + recent orders/products (revenue and user counts need aggregate endpoints - see backlog)
- [x] Products: table (search + pagination), create/edit dialog on a remote `form` + Zod, delete with confirm dialog
- [x] Orders: table + client-side id filter + row detail dialog (order, product, cart, user, status)
- [x] Categories: card grid + create/edit/delete
- [x] Users / Inventory / Shipping / Settings remain explicit stubs (as in the old app)

**Exit criteria:** non-admin sessions cannot reach `/admin/**` (403/redirect server-side); CRUD operations round-trip through remote forms with validation and toasts; dashboard and list views stay coherent after mutations (single-flight refresh).

---

## Phase 4 — Delivery `[~]`

**Goal:** the new apps ship through the existing pipelines with as little topology change as possible.

Deliverables
- [x] Dockerfiles rewritten on `oven/bun:1` (≥1.4): build stage `bun install --frozen-lockfile && bun run build`; runtime stage runs `bun ./build` as the non-root `bun` user (build output is self-contained); ports/EXPOSE stay 3000/3001
- [x] `NEXT_PUBLIC_API_URL` build plumbing removed from Dockerfiles/CI (stale Dockerfiles deleted with the old frontend; new ones have no build-time API vars)
- [x] compose: legacy `frontend` image reference fixed (`frontend-web`), missing `admin` service added, `API_INTERNAL_URL`/`API_PUBLIC_URL`/`ORIGIN` set
- [x] k8s `frontend`/`admin` manifests: `API_INTERNAL_URL` + `API_PUBLIC_URL` + `PROTOCOL_HEADER=x-forwarded-proto`; probes unchanged
- [x] CI: `frontend-check` job runs `bun install --frozen-lockfile` + `check` + `lint` + `i18n:check` + `bun test` with `oven-sh/setup-bun`; Docker build jobs unchanged and gated in `ci-success`
- [x] Root tooling: `Makefile` and `scripts/dev.ps1` `web-*` targets switched to Bun (plus new `web-check`/`web-lint`/`web-test`); README updated

**Exit criteria:** CI green on the branch; `docker build` for both apps succeeds; `make cluster-up` deploys the SvelteKit apps and the ingress hostnames serve them. *(both images build and run locally with Podman from the `frontend/` context — storefront serves `/` + `/en` with fonts, admin serves `/login` and redirects `/dashboard` server-side; CI run and `make cluster-up` remain)*

---

## Phase 5 — Hardening `[ ]`

- [x] Error boundaries (`+error.svelte`, branded + localized, with retry) in both apps; shared `EmptyState` / `TableSkeleton` adopted on the main list routes (per-route sweep continues opportunistically)
- [ ] Accessibility pass: focus management for dialogs/menus/drawer, keyboard nav, labels (bits-ui primitives already trap focus; no formal audit yet)
- [x] SEO: per-route titles/meta via Paraglide, `hreflang` + canonical alternates in the storefront layout, localized `sitemap.xml` route, robots.txt (storefront: private flows disallowed; admin: `Disallow: /`)
- [ ] Performance: prerender static routes where possible, audit bundle output, image sizing (`loading="lazy"`, `srcset` if a transform pipeline appears)
- [x] i18n completeness sweep: no hardcoded UI copy left (fixed the last four `aria-label`s); the only remaining Vietnamese strings are vendored mock *product data* in `home-mock.ts`, which is data rather than chrome
- [ ] Remove dead code carried from the old apps (unused components, mock leftovers)

**Exit criteria:** no blank/white error screens, Lighthouse a11y ≥ 95 on storefront routes, locale catalogs verified in CI.

---

## Phase 6 — Deferred / future `[ ]`

Explicitly out of scope until asked for:

- [ ] E2E tests (Playwright) — decided against for now
- [ ] Live features with `query.live` (notifications bell, order status streaming)
- [ ] Favourites, ratings/reviews UI, search-suggest dropdown (backend endpoints exist; the old app never wired them)
- [ ] Admin stubs: users, inventory, shipping, settings — product decisions required
- [ ] Migration to whatever replaces the experimental remote-functions/RF API surface once SvelteKit stabilises it
- [ ] Better-Auth / identity changes — depends on the backend phase of the [architecture revamp](../architecture-revamp-plan.md)
- [ ] Server-side carts (the cart stays client-only for now)
