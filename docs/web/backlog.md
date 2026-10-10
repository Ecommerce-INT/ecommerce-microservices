# Web Backlog

Granular task list for the SvelteKit rewrite. Grouped by area, roughly in execution order within each group. `[x]` done · `[~]` in progress · `[ ]` todo.

---

## Foundation

- [x] Branch `feature/sveltekit-frontend`; old Next.js frontend removed (staged in git)
- [x] Bun workspace root (`frontend/package.json`, `bun.lock`, engines `bun >= 1.4`) — pnpm/Turborepo gone
- [x] `apps/web` + `apps/admin` scaffolded via `sv create` (minimal, TS, Tailwind, ESLint, Prettier, Paraglide `vi`/`en`)
- [x] `@sveltejs/adapter-bun` in both apps; `experimental.remoteFunctions` + `compilerOptions.experimental.async`
- [x] `bun install` (205 packages), production builds, `lint`, `check` (0 errors / 0 warnings), live `bun ./build` smoke test (200, `<html lang="vi">`)
- [x] `docs/web/` — dev-guide, roadmap, backlog
- [x] `@ecommerce/lib` — `types.ts`, `format.ts`, `schemas.ts` (Zod), `cart.ts` (shared totals/free-shipping), `server/` (env, api, session, redirect), `index.ts`
- [x] `@ecommerce/lib` wired into both apps (devDependency so the Bun adapter bundles the TS source) and verification passed
- [x] Unit tests in `packages/lib` (`bun test`, 60 tests): formatters, cart totals, schema coercion (`optionalId`, `formNumber`), API envelope unwrapping (mock fetch), session cookies (`setSession`/refresh), `safeRedirect`
- [x] Root script sanity: `bun run check` / `lint` / `format` / `test` / `i18n:check` across all workspaces — green

## Theming & UI kit

- [x] Brand theme ported into `apps/*/src/routes/layout.css`:
  - [x] red palette `primary-50…900` (`#e30019` at 500), orange-scale remap, `--background`/`--foreground`
  - [x] utilities: `no-scrollbar`, `brand-pulse` (line-clamp is native in Tailwind 4)
  - [x] merged with shadcn-svelte's theme variables (don't clobber its `@theme`/`:root` blocks)
- [x] Fonts: Geist via `@fontsource-variable/geist`; `--font-sans` set in the shared theme (both apps; verified woff2 emission and serving)
- [x] shadcn-svelte set up in both apps (manual init; v1.7 preset prompt is interactive). NOTE: `init` — CSS path `src/routes/layout.css`, aliases `$lib/components`, `$lib/components/ui`, `$lib/utils` (utils re-exports `cn` from `@ecommerce/lib/format`). Imports now resolve through the `#lib/*` subpath map — the deprecated `$lib` alias was removed from both options configs.
- [x] Base components added in both apps:
  - [ ] button, input, label, textarea, badge
  - [ ] card, separator, skeleton
  - [ ] dialog + alert-dialog (replaces `window.confirm`), sheet (mobile menus/drawer)
  - [ ] dropdown-menu (profile, row actions), select (category pickers, payment method)
  - [ ] table (admin lists), pagination, breadcrumb
  - [ ] sonner (`svelte-sonner`) — toasts only, no `alert()`
- [x] Mounted `<Toaster />` in both root layouts (toasts replace alert/confirm)
- [x] Icon set: `@lucide/svelte` (installed in both apps)
- [x] Empty/loading states as shared components (`EmptyState`, `TableSkeleton`) in `@ecommerce/ui` — adopted on storefront orders and the admin orders/products tables

## i18n

- [x] Paraglide i18n routing (verified live): `vi` unprefixed / `en` prefixed, strategy `['url', 'cookie', 'baseLocale']`, cookie persistence, `reroute` sanity
- [x] Migrate message catalogs (script-assisted flatten of old nested JSON → Paraglide flat keys):
  - [x] web: 297 flat keys x 2 locales migrated (script-assisted flatten)
  - [x] admin: 109 flat keys x 2 locales migrated
- [x] Language switcher in the storefront header (link-based + `data-sveltekit-reload`, SSR-correct)
- [x] Key-parity check between locales (`bun run i18n:check`, wired into the CI frontend job)
- [x] Sweep: no hardcoded user-facing strings — the last four `aria-label`s now use `common_decrease`/`common_increase`; remaining Vietnamese strings are vendored mock product *data* (`home-mock.ts`)

## Auth & session

- [x] `@ecommerce/lib/server` session helpers — cookie constants, `setSession(cookies, tokens)`, `clearSession(cookies)`, `getSession(cookies)` with expiry + refresh
- [x] `hooks.server.ts` — Paraglide handle + session bootstrap into `locals.user` (+ refresh rotated tokens; persist the rotated refresh token — old bug)
- [x] `/login` route — server redirect to `${API}/api/v1/auth/login?redirect_uri=<origin>/auth/callback`, honours `?redirect=`
- [x] `/auth/callback` — server-side ticket exchange → httpOnly cookies → redirect (locale-aware return path via cookie, not `sessionStorage`)
- [x] Logout — remote `form`/`command`: call gateway logout, clear cookies, land on a public page (no SSO bounce loop)
- [x] Admin guard implemented in `apps/admin/src/hooks.server.ts`: — server-side layout check for `ADMIN`/`ROLE_ADMIN`; non-admin → 403/redirect (old app had **no** role check)
- [x] Guard storefront `/checkout`, `/orders` server-side

## Storefront

- [x] Layout: top bar, header (logo, search form, cart badge with live count, locale switcher, category nav), footer. Account dropdown lands with the auth BFF
- [x] Home: hero slides + side banners + sidebar, quick deals, category grid, flash sale rail with countdown/progress, hot products, phone/laptop sections with brand chips, brand strip, USP strip — data from ported `home-mock`
- [ ] Catalog `/products`:
  - [x] `catalog.remote.ts` — `query(productQuerySchema, …)` → `PaginatedResponse<Product>`
  - [x] filter sidebar (categories query + price ranges), sort select, pagination links from `totalElements` (bits-ui Pagination needs page objects - revisit)
  - [x] **search wired** to the API (`search` param / catalog-search endpoint — verify which the gateway exposes)
- [x] PDP: `product.remote.ts` `query`, gallery, price/stock, `QuantitySelector`, add-to-cart (store), buy-now (add + go to checkout)
- [x] Cart: line items, qty update/remove/clear, summary, free-shipping rule defined once in the runes store, ≥ 500 000 ₫ (single util, not duplicated like the old app)
- [x] Checkout: 3-step remote `form` (shipping Zod schema → payment method → review), order + payment creation, success page state, auth guard
- [x] Orders: guarded list, status badges, empty state
- [x] Remove dead routes: `/favourites`, `/profile` links → only render if implemented (currently 404 in old app)

## Admin

- [x] Shell: sidebar (8 nav items), topbar, mobile drawer (sheet), logout
- [x] Dashboard: stat cards, recent orders, recent products (query-driven; same query keys/args as lists so mutations refresh both — old cache-key bug)
- [x] Products: table (search, pagination), create/edit dialog with remote `form(productFormSchema)`, delete confirm, toasts
- [x] Orders: table + client-side id filter + row detail dialog (order/product/cart/user fields)
- [x] Categories: card grid, create/edit/delete
- [x] Stubs: users, inventory, shipping, settings (keep the “under development” pattern)

## Infrastructure & tooling

- [x] Dockerfile `apps/web` — `docker.io/oven/bun:1`; build `bun install --frozen-lockfile` + `bun run build`; runtime copies `build/` and runs `bun ./build` as the `bun` user (the bundle is self-contained — no production install needed), EXPOSE 3000; `frontend/.dockerignore` added
- [x] Dockerfile `apps/admin` — same, port 3001. Both images built and run locally with Podman: storefront serves `/` + `/en` (fonts verified), admin serves `/login` and guards `/dashboard` server-side
- [x] No `NEXT_PUBLIC_API_URL` anywhere in the new Dockerfiles/CI (stale image references gone)
- [x] compose: `admin` service added, `frontend` image ref fixed (`frontend-web`), `API_INTERNAL_URL`/`API_PUBLIC_URL`/`ORIGIN` set, ports 3000/3001
- [x] k8s manifests: `API_INTERNAL_URL=http://apisix-gateway:9080` + `API_PUBLIC_URL` + `PROTOCOL_HEADER=x-forwarded-proto`; probes unchanged
- [x] CI (`ci.yml`): `frontend-check` job = `bun install --frozen-lockfile` → `check` → `lint` → `i18n:check` → `test` with `oven-sh/setup-bun`; Docker build/push jobs kept and gated in `ci-success`
- [x] Root `Makefile` + `scripts/dev.ps1`: `web-install`/`web-dev`/`web-build` → Bun, plus `web-check`/`web-lint`/`web-test`
- [x] Root `README.md`: frontend stack rows/commands updated (pnpm → Bun, Next.js → SvelteKit)

## Fixes carried from the old app (do **not** re-introduce)

- [x] Catalog search parameter actually sent (old: query-key only)
- [x] Pagination derived from `totalElements` (old: computed from page length → page > 0 broken)
- [x] Brand filter dropped (old: accepted `?brand=` then ignored it); brand chips link to category filters
- [x] Admin role enforcement (old: any authenticated user = full admin)
- [x] Dashboard ↔ list data coherence (old: different query keys → stale dashboard)
- [x] Rotated refresh token persisted (old: dropped it)
- [x] Admin `/login` no longer 404s on refresh failure (old: hard redirect to a storefront-only route)
- [x] Logout no longer bounces straight back into Keycloak SSO
- [x] No `alert()` / `window.confirm()` as UI (svelte-sonner toasts + shadcn dialog/alert-dialog)
- [x] Cart totals/free-shipping logic defined once — `@ecommerce/lib/cart` `computeTotals`; the cart store, cart page and checkout command all consume it

## Quality

- [x] `+error.svelte` per app — branded, localized (404 vs generic), with home/retry actions; gateway failures surface as the generic error state
- [ ] A11y: focus trap/restore for dialog/sheet/drawer, keyboard nav for menus, form labels + `aria-invalid`
- [x] SEO (storefront): per-route `<title>`/meta via Paraglide, `hreflang` + canonical in the root layout, localized `sitemap.xml` route, robots (private flows disallowed; admin app `Disallow: /`)
- [ ] Performance: prerender marketing/static routes; review bundle output after the port
- [x] Dev experience: `.env.example` per app documenting `API_INTERNAL_URL`, `API_PUBLIC_URL`, `PORT`, `ORIGIN`, `PROTOCOL_HEADER`

## Deferred (do not start without a decision)

- [ ] Multi-line orders - the backend `Order` model carries a single `productId`; checkout records the first line plus the grand total (matching the old app contract). Real multi-line order history needs a backend change.

- [ ] Playwright e2e suite
- [ ] `query.live` notifications / live order status
- [ ] Favourites, ratings, search suggestions UI
- [ ] Admin users/inventory/shipping/settings beyond stubs
- [ ] Server-side carts
- [ ] Better-Auth migration (backend-dependent, see revamp plan)
- [ ] Post-stabilisation migration of experimental remote functions / async Svelte API

---

## Decisions log

| Decision | Choice | Notes |
| :--- | :--- | :--- |
| Package manager / task runner | Bun workspaces | No pnpm, no Turborepo; `bun run --filter` instead |
| Server adapter | `@sveltejs/adapter-bun` (official) | Bun ≥ 1.4; ports 3000/3001 preserved; Docker on `oven/bun` |
| Data layer | Remote functions + Zod | Experimental flags accepted; confined to `*.remote.ts` |
| Forms | Remote `form` (no Superforms) | Progressive enhancement + `preflight()` for client validation |
| UI kit | shadcn-svelte in a shared `@ecommerce/ui` | Components, `theme.css` and `cn` live in `frontend/packages/ui`; components are added from that package, not per app |
| `cn` helper | `tailwind-variants` built-in | v3 ships `cn` with class merging; `clsx`/`tailwind-merge` dropped (verified: `cn('p-2','p-4')` -> `p-4`) |
| Node types | `@types/bun` in both apps | Everything runs on Bun; the lib already used `@types/bun` |
| Lint/format tooling | Hoisted to `frontend/package.json` | eslint/prettier plugins + `prettier.config.js` configured once for all four packages |
| BFF session helpers | `@ecommerce/lib/server` | Both apps share cookie names, set/clear/refresh, `safeRedirect`, `isAdmin` |
| i18n | Paraglide JS | `vi` base/unprefixed, `en` prefixed; catalogs as flat JSON |
| Auth | BFF httpOnly cookies | Tokens never in JS; server-side guards incl. admin role |
| API access | Server-only via `API_INTERNAL_URL` | Browser never calls APISIX; no CORS, no build-time public env |
| Cart | Client-side runes store | localStorage key `ecommerce-cart` kept for continuity |
| Tests | `bun test` (logic only) | No e2e for now (explicit decision) |
