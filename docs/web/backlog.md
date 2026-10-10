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
- [~] `@ecommerce/lib` — files written: `types.ts`, `format.ts`, `schemas.ts` (Zod), `server/env.ts`, `server/api.ts`, `server/index.ts`, `index.ts`
- [x] `@ecommerce/lib` wired into both apps (devDependency so the Bun adapter bundles the TS source) and verification passed
- [ ] Unit tests in `packages/lib` (`bun test`): formatters, schema coercion (`optionalId`, `formNumber`), API envelope unwrapping (mock fetch)
- [ ] Root script sanity: `bun run check/lint/format/test` across all workspaces

## Theming & UI kit

- [x] Brand theme ported into `apps/*/src/routes/layout.css`:
  - [x] red palette `primary-50…900` (`#e30019` at 500), orange-scale remap, `--background`/`--foreground`
  - [x] utilities: `no-scrollbar`, `brand-pulse` (line-clamp is native in Tailwind 4)
  - [x] merged with shadcn-svelte's theme variables (don't clobber its `@theme`/`:root` blocks)
- [ ] Fonts: Geist (old app used `next/font`); pick `@fontsource-variable/geist` and set `--font-sans`
- [x] shadcn-svelte set up in both apps (manual init; v1.7 preset prompt is interactive). NOTE: `init` — CSS path `src/routes/layout.css`, aliases `$lib/components`, `$lib/components/ui`, `$lib/utils` (utils re-exports `cn` from `@ecommerce/lib/format`)
- [x] Base components added in both apps:
  - [ ] button, input, label, textarea, badge
  - [ ] card, separator, skeleton
  - [ ] dialog + alert-dialog (replaces `window.confirm`), sheet (mobile menus/drawer)
  - [ ] dropdown-menu (profile, row actions), select (category pickers, payment method)
  - [ ] table (admin lists), pagination, breadcrumb
  - [ ] sonner (`svelte-sonner`) — toasts only, no `alert()`
- [x] Mounted `<Toaster />` in both root layouts (toasts replace alert/confirm)
- [x] Icon set: `@lucide/svelte` (installed in both apps)
- [ ] Empty/loading states as small shared components (`EmptyState`, `TableSkeleton`)

## i18n

- [x] Paraglide i18n routing (verified live): `vi` unprefixed / `en` prefixed, strategy `['url', 'cookie', 'baseLocale']`, cookie persistence, `reroute` sanity
- [ ] Migrate message catalogs (script-assisted flatten of old nested JSON → Paraglide flat keys):
  - [x] web: 297 flat keys x 2 locales migrated (script-assisted flatten)
  - [x] admin: 109 flat keys x 2 locales migrated
- [x] Language switcher in the storefront header (link-based + `data-sveltekit-reload`, SSR-correct)
- [ ] Key-parity check between locales (script, runnable in CI)
- [ ] Sweep: no hardcoded user-facing strings (including mock data titles where feasible)

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

- [x] Shell: sidebar (8 nav items), topbar, mobile nav, logout (drawer later)
- [x] Dashboard: stat cards, recent orders, recent products (query-driven; same query keys/args as lists so mutations refresh both — old cache-key bug)
- [x] Products: table (search, pagination), create/edit dialog with remote `form(productFormSchema)`, delete confirm, toasts
- [~] Orders: table + client-side id filter (row detail still pending)
- [x] Categories: card grid, create/edit/delete
- [x] Stubs: users, inventory, shipping, settings (keep the “under development” pattern)

## Infrastructure & tooling

- [ ] Dockerfile `apps/web` — `oven/bun:1`, `bun install --frozen-lockfile`, `bun run build`; runtime `bun install --production --frozen-lockfile` + `bun ./build`, non-root user, EXPOSE 3000
- [ ] Dockerfile `apps/admin` — same, port 3001
- [ ] Remove `NEXT_PUBLIC_API_URL` arg/ENV from both Dockerfiles and CI
- [ ] compose: add `admin` service, fix `frontend` image ref (`ghcr.io/<owner>/frontend-web`), set `API_INTERNAL_URL=http://apisix:9080`, ports 3000/3001
- [ ] k8s manifests: `API_INTERNAL_URL=http://apisix-gateway:9080`, `PROTOCOL_HEADER=x-forwarded-proto`; verify probes on `/`
- [ ] CI (`ci.yml`): frontend job = `bun install` → `check` → `lint` → `test` with `oven-sh/setup-bun`; keep Docker build/push jobs
- [ ] Root `Makefile` + `scripts/dev.ps1`: `web-install/web-dev/web-build` → Bun workspace commands (drop pnpm)
- [ ] Root `README.md`: frontend stack rows/commands updated (pnpm → Bun, Next.js → SvelteKit)

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
- [ ] Cart totals/free-shipping logic defined once

## Quality

- [ ] `+error.svelte` per app; route-level error states for gateway failures (helpful message when the API is unreachable in dev)
- [ ] A11y: focus trap/restore for dialog/sheet/drawer, keyboard nav for menus, form labels + `aria-invalid`
- [ ] SEO (storefront): per-route `<title>`/meta via Paraglide, `hreflang` for `vi`/`en`, sitemap, robots
- [ ] Performance: prerender marketing/static routes; review bundle output after the port
- [ ] Dev experience: `.env.example` per app documenting `API_INTERNAL_URL`

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
