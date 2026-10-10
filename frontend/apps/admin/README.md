# `@ecommerce/admin` — Backoffice App

Internal backoffice and administration dashboard built with **SvelteKit 3**, **Svelte 5 (Runes)**, and **Tailwind CSS 4**, running on **Bun** (`:3001`).

For workspace architecture, shared packages (`@ecommerce/lib`, `@ecommerce/ui`), and full monorepo commands, see the **[Frontend Workspace README](../../README.md)**.

---

## Security & Access Control

Every backoffice route (except `/login`, `/auth/callback`, and `/access-denied`) is protected server-side in [`src/hooks.server.ts`](./src/hooks.server.ts):

- Unauthenticated requests redirect to `/login?redirect=<path>`.
- Authenticated sessions without `ADMIN` or `ROLE_ADMIN` in `locals.user.roles` are redirected to `/access-denied` (`403`).

---

## Features & Routes

| Route                                            | Description                                                                                                                |
| :----------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------- |
| `/dashboard`                                     | KPI stat cards, recent orders, and catalog summary (`dashboard.remote.ts`).                                                |
| `/products`                                      | Paginated product table with search, Zod-validated create/edit modal form, and delete confirmation (`products.remote.ts`). |
| `/orders`                                        | Paginated order list with status filter and order line-item detail dialog (`orders.remote.ts`).                            |
| `/categories`                                    | Category card grid with create/edit/delete operations (`categories.remote.ts`).                                            |
| `/users`, `/inventory`, `/shipping`, `/settings` | Module placeholders ready for upcoming backend endpoints.                                                                  |
| `/login`, `/auth/callback`, `/access-denied`     | BFF SSO login handoff, ticket exchange into `HttpOnly` cookies, and non-admin guard page.                                  |

---

## Development

Run from the `frontend/` workspace root (or inside `apps/admin/`):

```bash
# Start dev server on http://localhost:3001
bun run --filter @ecommerce/admin dev

# Typecheck & lint
bun run --filter @ecommerce/admin check
bun run --filter @ecommerce/admin lint

# Build & run production server
bun run --filter @ecommerce/admin build
PORT=3001 bun ./apps/admin/build
```

See [`.env.example`](./.env.example) for runtime environment variables (`API_INTERNAL_URL`, `API_PUBLIC_URL`, `ORIGIN`, `PORT`).
