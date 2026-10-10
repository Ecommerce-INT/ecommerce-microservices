# Frontend Workspace

Bun monorepo containing the customer storefront (`@ecommerce/web`) and backoffice admin (`@ecommerce/admin`) built with SvelteKit 3, Svelte 5 (runes), and Tailwind CSS 4.

## Request Flow

```text
Browser ──(SSR + Remote Functions + HttpOnly Cookie)──▶ SvelteKit Server (Bun.serve)
                                                              │
                                                       API_INTERNAL_URL
                                                              ▼
                                                    Apache APISIX (:9080)
```

- All data fetching and mutations run through SvelteKit remote functions (`query`, `form`, `command` in `*.remote.ts`) validated server-side with Zod v4.
- Authentication uses a BFF pattern: `/auth/callback` exchanges the OIDC ticket server-side and stores tokens in `HttpOnly`, `SameSite=Lax` cookies (`hooks.server.ts` refreshes tokens and populates `event.locals.user`).

## Workspace Structure

| Package                                                     | Port   | Description                                                                                                                                                         |
| :---------------------------------------------------------- | :----- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [`apps/web`](./apps/web/README.md) (`@ecommerce/web`)       | `3000` | Customer storefront: home, catalog search/filters/pagination, PDP, runes cart store, 3-step checkout, order history, and localized SEO (`sitemap.xml`, `hreflang`). |
| [`apps/admin`](./apps/admin/README.md) (`@ecommerce/admin`) | `3001` | Backoffice dashboard: server-side `ADMIN` role guard in `hooks.server.ts`, KPI summary, and CRUD views for products, categories, and orders.                        |
| [`packages/lib`](./packages/lib) (`@ecommerce/lib`)         | —      | Shared domain types, Zod v4 schemas, cart totals math (`computeTotals`), formatters, BFF cookie session helpers, and the server APISIX client.                      |
| [`packages/ui`](./packages/ui) (`@ecommerce/ui`)            | —      | Shared `shadcn-svelte` (`bits-ui`) components, `EmptyState`, `TableSkeleton`, Geist font, and Tailwind CSS 4 theme (`theme.css`).                                   |

## Commands

Run from `frontend/` (or via `make web-*` / `.\scripts\dev.ps1 web-*` from the repository root):

```bash
bun install          # Install workspace dependencies
bun run dev          # Start both dev servers (web :3000, admin :3001)
bun run dev:web      # Start storefront only (:3000)
bun run dev:admin    # Start backoffice only (:3001)
bun run build        # Build production bundles for both apps (adapter-bun)
bun run check        # Compile Paraglide messages + svelte-check + tsc
bun run lint         # Prettier check + ESLint
bun run format       # Format with Prettier
bun run i18n:check   # Verify vi.json and en.json key parity in both apps
bun test             # Run unit tests in packages/lib
```

## Environment Variables

See [`apps/web/.env.example`](./apps/web/.env.example) and [`apps/admin/.env.example`](./apps/admin/.env.example):

| Variable           | Default                           | Description                                                                                  |
| :----------------- | :-------------------------------- | :------------------------------------------------------------------------------------------- |
| `API_INTERNAL_URL` | `http://localhost:9080`           | Server-to-gateway URL (`http://apisix:9080` in Compose, `http://apisix-gateway:9080` in K8s) |
| `API_PUBLIC_URL`   | `API_INTERNAL_URL`                | Browser-reachable gateway URL used for the `/api/v1/auth/login` redirect                     |
| `HOST` / `PORT`    | `0.0.0.0` / `3000` (`3001` admin) | Bind address and port for the built `Bun.serve` server                                       |
| `ORIGIN`           | —                                 | Explicit public origin (e.g. `http://localhost:3000`) when serving over plain HTTP           |
| `PROTOCOL_HEADER`  | —                                 | Set to `x-forwarded-proto` behind reverse proxies / Kubernetes Ingress                       |

## Container Builds

Run from `frontend/`:

```bash
docker build -f apps/web/Dockerfile -t ecommerce/frontend-web:latest .
docker build -f apps/admin/Dockerfile -t ecommerce/frontend-admin:latest .
```

## Further Documentation

- [Developer Guide](../docs/web/dev-guide.md)
- [Roadmap](../docs/web/roadmap.md) and [Backlog](../docs/web/backlog.md)
