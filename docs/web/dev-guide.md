# Frontend Dev Guide — SvelteKit 3 + Bun

> **Branch:** `feature/sveltekit-frontend` (cut from `feature/dev-infra`)
>
> **Status:** scaffold complete — both apps install, build, lint, type-check and run on Bun. Feature porting is in progress.
>
> The previous Next.js 16 frontend (React 19, pnpm, Turborepo) was removed on this branch. It remains in git history and on `main` / `feature/dev-infra` until the rewrite reaches parity.

---

## 1. Stack

| Layer | Choice | Rationale |
| :--- | :--- | :--- |
| Runtime & package manager | **Bun ≥ 1.4 workspaces** (`bun.lock` committed) | Same runtime as the 4 Bun services; one tool for install / dev / build / test. Turborepo and pnpm dropped — `bun run --filter` covers task fan-out. |
| Framework | **SvelteKit 3 + Svelte 5 (runes)** on **Vite 8** | File-based routing + SSR; runes replace React state entirely. TypeScript 6, ESLint 10, Prettier 3.8. |
| Server runtime / adapter | **`@sveltejs/adapter-bun` 1.0.0** (official) | Real `Bun.serve` server, ETags/range requests for assets, graceful shutdown. Requires Bun ≥ 1.4; respects `HOST`/`PORT`. |
| Data layer | **SvelteKit remote functions** — `query` / `form` / `command` / `prerender` from `$app/server` | Type-safe client↔server calls, request dedupe/caching, single-flight mutations (`refresh()` / `set()`), progressive-enhancement forms. Replaces axios + TanStack Query + form actions. |
| Validation | **Zod 4** (Standard Schema) | Every remote-function argument and form payload validates server-side by default; `preflight()` adds optional client-side validation. |
| UI | **`@ecommerce/ui`**: shared shadcn-svelte primitives (bits-ui) + brand theme on **Tailwind CSS 4**, `@lucide/svelte`, **svelte-sonner**, Geist variable font (`@fontsource-variable/geist`) | One copy of the primitives for both apps; toasts/dialogs instead of `alert()` / `confirm()`. |
| i18n | **Paraglide JS** (SvelteKit's official integration) | Tree-shaken, fully typed messages. Locales `vi` (base) and `en`; target URL contract stays vi-unprefixed + `/en/...` like the old app. |
| Auth | **BFF**: httpOnly cookies set in `hooks.server.ts` | Tokens never touch JS/localStorage (helpers in `@ecommerce/lib/server`); server-side route guards — including the admin role check the old app never had. |
| Client state | Svelte runes modules (`*.svelte.ts`) | Only the cart needs client state; its localStorage key (`ecommerce-cart`) is preserved. |
| Tests | `bun test` — pure logic only | No e2e yet. Quality gates: `svelte-check`, ESLint, Prettier. |

### Dependency notes

- All primitives live once in `frontend/packages/ui` (`@ecommerce/ui`); the apps depend on it plus `@lucide/svelte` and `svelte-sonner`.
- `cn` comes from `tailwind-variants` (re-exported by `@ecommerce/ui/utils`); `clsx`/`tailwind-merge` are not used directly.
- `@types/bun` replaces `@types/node` in both apps (the lib already used it).
- Lint/format tooling (eslint, prettier + svelte/tailwind plugins, globals, typescript-eslint) is hoisted to `frontend/package.json` and configured once in `frontend/prettier.config.js`.
- `bun run check` compiles the Paraglide messages first, so typechecking works on a clean clone.
- Shared Zod schemas are also imported by client components (cart/checkout/admin forms); if bundle size matters, `zod/mini` is a drop-in for those (~2 KB vs ~12 KB gzipped).

### Experimental features in use

Remote functions and `await` in components are still flagged experimental. Both apps opt in via `vite.config.ts`:

```ts
sveltekit({
	experimental: { remoteFunctions: true },
	compilerOptions: { experimental: { async: true }, runes: /* ... */ },
	adapter: adapter()
})
```

Accepted risk: the API may change. Mitigation: remote functions stay confined to `*.remote.ts` modules so a future migration is mechanical.

---

## 2. Request flow (what changed architecturally)

```
browser ──(form / query RPC + session cookie)──▶ SvelteKit server (Bun.serve)
                                                    │  httpOnly cookie session (BFF)
                                                    ▼
                                             API_INTERNAL_URL
                                                    ▼
                                             APISIX gateway ──▶ Go / Bun services
```

- The browser **never** calls the API gateway: no CORS allow-list to maintain, no tokens in JS, no build-time `NEXT_PUBLIC_*` API URL.
- Server code uses a single env var, `API_INTERNAL_URL` (internal DNS inside the cluster).
- Remote `query` results resolve **during SSR** and hydrate from the serialized value — real server rendering, unlike the old all-client-side fetching.

---

## 3. Workspace layout

```
frontend/
├─ package.json             # Bun workspace root (apps/*, packages/*)
├─ bun.lock                 # committed
├─ apps/web                 # @ecommerce/web   — storefront (port 3000)
├─ apps/admin               # @ecommerce/admin — backoffice (port 3001)
└─ packages/lib             # @ecommerce/lib   — domain types, Zod schemas, formatters, server-only API client
```

Per-app layout (Kit 3, generated by `sv create`):

```
src/
├─ routes/                  # file-based routes; no [lang] segment — Paraglide reroute handles /en prefix
│  └─ layout.css            # Tailwind 4 entry + brand palette (ported from the old globals.css)
├─ hooks.server.ts          # Paraglide handle → auth/session (BFF) goes here
├─ hooks.ts                 # Paraglide reroute (deLocalizeUrl)
├─ lib/
│  ├─ paraglide/            # generated messages — git-ignored
│  ├─ components/ui/        # shadcn-svelte components (copy-in, per app)
│  └─ server/               # server-only modules (SvelteKit-enforced)
├─ app.html                 # %paraglide.lang% / %paraglide.dir%
messages/{vi,en}.json       # translation source (baseLocale: vi)
project.inlang/settings.json
vite.config.ts              # adapter + plugins + experimental flags
```

### Conventions

- **All data access goes through `*.remote.ts` modules.** Components never call the gateway directly; handlers stay thin: validate → call `@ecommerce/lib` server API client → return domain types.
- **Zod schemas are shared** in `@ecommerce/lib/schemas`; the server-only API client lives in `@ecommerce/lib/server` (never imported by client code).
- **Imports:** app code resolves through the `#lib/*` subpath map (Kit 3 deprecated the custom `$lib` alias). The map lives in each app's `package.json` `imports` field, with a matching `paths` entry in its `tsconfig.json` so svelte-check resolves extensionless targets like `#lib/server/api` → `api.ts`.
- **Server-only code** lives in `src/lib/server/**` — SvelteKit throws if a client component imports it.
- **No hardcoded copy:** every user-facing string goes through `messages/{vi,en}.json` via `m.*` (Paraglide).
- **UI primitives:** add from the shared package (`cd frontend/packages/ui && bunx shadcn-svelte@latest add <component>`) and export via its `package.json`; app-level composition components live in `#lib/components`.
- **Feedback:** svelte-sonner toasts for success/errors, shadcn dialogs for confirms.

---

## 4. Commands

Root (`frontend/`):

```bash
bun install                  # install all workspaces
bun run dev                  # both apps — web :3000, admin :3001
bun run dev:web              # storefront only
bun run dev:admin            # backoffice only
bun run build                # both apps (Vite runs under Bun — required by adapter-bun)
bun run check                # svelte-kit sync + svelte-check per app
bun run lint                 # prettier --check + eslint per app
bun run format               # prettier --write per app
bun test                     # unit tests (bun test)
```

Per app (`apps/web`, `apps/admin`): `bun run dev | build | check | lint | format`. App scripts already force the Bun runtime (`bun run --bun vite …`), so they behave the same however they're invoked.

Local production smoke test:

```bash
cd apps/web
bun run build
bun ./build                  # Bun.serve on $PORT (default 3000)
```

---

## 5. Environment variables

| Variable | Scope | Value |
| :--- | :--- | :--- |
| `API_INTERNAL_URL` | server-side | k8s: `http://apisix-gateway:9080` · compose: `http://apisix:9080` · local dev: `http://localhost:9080` |
| `HOST` / `PORT` | built server | `0.0.0.0` / `3000` (web) or `3001` (admin) — set in k8s/compose; default port 3000 |
| `PROTOCOL_HEADER` | built server behind a proxy | Set `x-forwarded-proto` in the cluster — the adapter otherwise assumes `https` |
| `ORIGIN` (runtime) | built server serving plain HTTP directly | compose sets `ORIGIN=http://localhost:3000`; a build-time `paths.origin` works too |

There are **no public/build-time API variables anymore** — `NEXT_PUBLIC_API_URL` and its Docker/CI plumbing are obsolete.

---

## 6. Remote functions — rules of thumb

- One module per feature area, e.g. `src/routes/products/catalog.remote.ts`, `.../checkout.remote.ts`, `src/routes/admin/products/products.remote.ts`.
- `query` for reads · `form` for `<form>`-bound mutations (progressive enhancement) · `command` for programmatic mutations · `prerender` for build-time data.
- Always pass a Zod schema as the first argument when the function takes input.
- Refresh precisely after mutations: `void getProducts(page).refresh()` inside the handler; use `submit().updates(...)` from the client for filter-dependent instances. Never fall back to "invalidate everything".
- Don't use `url`/`params` inside `query` (throws — queries are cached by argument, not navigation). Pass them as arguments.
- Cookies: `getRequestEvent().cookies`; settable only inside `form`/`command`.

```ts
// src/routes/products/catalog.remote.ts
import { query } from '$app/server';
import * as v from 'zod';
import { productApi } from '@ecommerce/lib/server';

export const getProducts = query(
	v.object({
		page: v.number().int().min(0),
		size: v.number().int().min(1).max(60),
		search: v.string().optional()
	}),
	({ page, size, search }) => productApi.list({ page, size, search })
);
```

---

## 7. i18n (Paraglide)

- Source of truth: `messages/vi.json` + `messages/en.json`; `baseLocale: vi`.
- `%paraglide.lang%` / `%paraglide.dir%` in `app.html` are filled by the handle hook.
- Use `localizeHref()` on links that must keep the locale and `setLocale()` to switch (it reloads by design). Cross-locale links carry `data-sveltekit-reload`.
- The generated `src/lib/paraglide/` directory is git-ignored; never edit it by hand.
- Never hardcode user-facing copy — add keys to both message files.

---

## 8. UI (@ecommerce/ui)

- Primitives live once in `frontend/packages/ui` and are imported as `@ecommerce/ui/<component>`; the brand theme is exported as `@ecommerce/ui/theme.css` and imported by each app layout.css (which only overrides `--background`).
- Add components from the package (its own components.json + tsconfig path aliases): `cd frontend/packages/ui && bunx shadcn-svelte@latest add <component>`, then add the subpath to the package exports map.
- `cn` comes from `tailwind-variants` (re-exported by `@ecommerce/ui/utils` next to the bits-ui prop type helpers) - no clsx/tailwind-merge dependency. Icons: `@lucide/svelte`. Toasts: `svelte-sonner` (light theme; no dark-mode toggle wired, which is why `mode-watcher` was dropped).

---

## 9. Auth flow (BFF)

1. `/login` → stores the intended destination in a short-lived cookie and links to `${API_PUBLIC_URL}/api/v1/auth/login?redirect_uri=<origin>/auth/callback` (set `API_PUBLIC_URL` whenever the browser-reachable gateway origin differs from `API_INTERNAL_URL`).
2. `/auth/callback` exchanges the single-use ticket **server-side** (`/api/v1/auth/session?ticket=`) and sets four httpOnly cookies: `access_token`, `refresh_token`, `session_user` (profile + roles) and `session_expires` (SameSite=Lax).
3. `hooks.server.ts` rotates the access token (persisting the rotated refresh token) and populates `event.locals.user`; layouts and remote handlers read it. The rotated refresh token is persisted (the old app dropped it).
4. `apps/admin/src/hooks.server.ts` guards every route except `/login`, `/auth/callback` and `/access-denied`: anonymous visitors go to `/login?redirect=...`, signed-in non-admins to `/access-denied` (roles come from the profile, `ADMIN` / `ROLE_ADMIN`) — the previous app shipped with none.
5. Logout is a remote `command` that ends the gateway session and clears the cookies, then navigates to `/login` (the old logout bounced straight back into SSO).

---

## 10. Deployment

- **Dockerfiles** (`frontend/apps/{web,admin}/Dockerfile`): `docker.io/oven/bun:1` build stage → `bun install --frozen-lockfile && bun run build`; runtime stage copies only `build/` and runs `bun ./build` as the non-root `bun` user — the SSR bundle is self-contained (no production install needed). `frontend/.dockerignore` keeps `node_modules`, `.svelte-kit` and `build` out of the context. Both images were built and smoke-tested locally with Podman.
- **Unchanged:** ports **3000 / 3001**, Dockerfile paths, CI build context `frontend/`, image names `frontend-web` / `frontend-admin`, k8s probes on `/` (SSR returns 200).
- **k8s** (`k8s/frontend/*.yaml`): `API_INTERNAL_URL=http://apisix-gateway:9080`, `API_PUBLIC_URL`, `PROTOCOL_HEADER=x-forwarded-proto`.
- **compose:** `frontend` (image `frontend-web`) and `admin` services with `API_INTERNAL_URL=http://apisix:9080`, `API_PUBLIC_URL=http://localhost:9080` and `ORIGIN` set for plain-HTTP access.
- **CI (`ci.yml`):** the `frontend-check` job runs `bun install --frozen-lockfile` → `check` → `lint` → `i18n:check` → `bun test` with `oven-sh/setup-bun`; the Docker build/push jobs are unchanged and gated behind it in `ci-success`.

---

## 11. Migration reference (old → new)

Removed on this branch: the Next.js `apps/*`, the old hand-rolled `packages/ui` and `packages/config` (replaced by a new `packages/ui` built on shadcn-svelte), `pnpm-lock.yaml`, `pnpm-workspace.yaml`, `turbo.json`, `.npmrc`.

| Old (Next.js) | New (SvelteKit) |
| :--- | :--- |
| axios client + TanStack Query hooks | `*.remote.ts` (`query` / `form`) + `@ecommerce/lib/server` API client |
| zustand `authStore` + localStorage tokens | httpOnly cookie session (BFF) + `event.locals.user` |
| zustand `cartStore` (persist `ecommerce-cart`) | `#lib/stores/cart.svelte.ts` — same localStorage key |
| next-intl `messages/*.json` (327 lines each app) | Paraglide `messages/{vi,en}.json` |
| hand-rolled `@ecommerce/ui` | shadcn-svelte components per app |
| `alert()` / `window.confirm()` | svelte-sonner toasts + shadcn dialog |
| `next/image` (remote patterns) | plain `<img>` — product images are remote URLs |
| `NEXT_PUBLIC_API_URL` (baked, empty-by-default) | `API_INTERNAL_URL` (server-only, runtime) |

Bugs deliberately fixed while porting: catalog search wired to the API, pagination from `totalElements`, brand filter, admin role guard, dashboard/list cache coherence (single-flight refresh), rotated refresh token persistence, admin `/login` 404, logout loop, dead components/APIs dropped.

---

## 12. Status

- [x] `feature/sveltekit-frontend` branch; old Next.js frontend removed
- [x] Bun workspace root + `apps/web` + `apps/admin` scaffolded (Kit 3, Tailwind 4, ESLint, Prettier, Paraglide vi/en)
- [x] `@sveltejs/adapter-bun` wired in both apps; remote functions + async enabled
- [x] Verified: `bun install`, production builds, `lint`, `check` (0 errors / 0 warnings), and a live `bun ./build` smoke test (`200`, `<html lang="vi">`)
- [x] `packages/lib` — types, Zod schemas, formatters, shared cart math, server-only gateway client + session helpers (`@ecommerce/lib/{types,schemas,format,cart,server}`), wired into both apps, covered by 60 unit tests
- [x] Paraglide i18n-routing config so `vi` stays unprefixed (matches the old `as-needed` behaviour)
- [x] shadcn-svelte components + brand theme, shared via `@ecommerce/ui`
- [x] Auth BFF + storefront/admin guards (admin guard in hooks.server.ts)
- [x] Storefront routes + admin routes (dashboard, products, orders + row detail, categories, stubs)
- [x] `@ecommerce/lib` unit tests (`bun test`, 60 tests) + locale key-parity check (`bun run i18n:check`)
- [x] Geist variable font in both apps (+ fixed the admin app shipping without its stylesheet import)
- [x] Dockerfiles, CI, compose, k8s updates (CI/Docker verification pending a pipeline run)
- [x] Root `Makefile` / `scripts/dev.ps1` `web-*` targets switched to Bun (plus `web-check`/`web-lint`/`web-test`)

Roadmap and granular task list: [roadmap.md](./roadmap.md) · [backlog.md](./backlog.md).
