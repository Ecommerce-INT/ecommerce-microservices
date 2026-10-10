# `@ecommerce/web` — Storefront App

Customer-facing e-commerce storefront built with **SvelteKit 3**, **Svelte 5 (Runes)**, and **Tailwind CSS 4**, running on **Bun** (`:3000`).

For workspace architecture, shared packages (`@ecommerce/lib`, `@ecommerce/ui`), and full monorepo commands, see the **[Frontend Workspace README](../../README.md)**.

---

## Features & Routes

| Route                         | Description                                                                                                                                                              |
| :---------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/` (`/en`)                   | Storefront home — hero banners, category bar, flash sale rail, quick deals, and featured product rails.                                                                  |
| `/products`                   | Server-rendered catalog with search, category & price-range filters, sorting, and pagination (`catalog.remote.ts`).                                                      |
| `/products/[id]`              | Product detail page (PDP) with stock status, quantity selector, add-to-cart, and buy-now (`product.remote.ts`).                                                          |
| `/cart`                       | Client-side runes cart (`#lib/stores/cart.svelte.ts`) persisted under `localStorage["ecommerce-cart"]` with shared free-shipping threshold math (`@ecommerce/lib/cart`). |
| `/checkout`                   | Auth-guarded 3-step checkout wizard (shipping $\rightarrow$ payment $\rightarrow$ review) validated server-side via Zod v4 (`checkout.remote.ts`).                       |
| `/orders`                     | Auth-guarded customer order history with status badges (`orders.remote.ts`).                                                                                             |
| `/login`, `/auth/callback`    | BFF OIDC/SSO login redirect and server-side ticket exchange into `HttpOnly` session cookies.                                                                             |
| `/sitemap.xml`, `/robots.txt` | Localized SEO sitemap and crawler configuration with `vi` / `en` `hreflang` alternates.                                                                                  |

---

## Development

Run from the `frontend/` workspace root (or inside `apps/web/`):

```bash
# Start dev server on http://localhost:3000
bun run --filter @ecommerce/web dev

# Typecheck & lint
bun run --filter @ecommerce/web check
bun run --filter @ecommerce/web lint

# Build & run production server
bun run --filter @ecommerce/web build
bun ./apps/web/build
```

See [`.env.example`](./.env.example) for runtime environment variables (`API_INTERNAL_URL`, `API_PUBLIC_URL`, `ORIGIN`, `PORT`).
