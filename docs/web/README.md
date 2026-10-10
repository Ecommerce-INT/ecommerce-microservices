# Web (Storefront + Backoffice) Docs

The two frontends are being rewritten from **Next.js 16** to **SvelteKit 3 + Bun** on the `feature/sveltekit-frontend` branch (cut from `feature/dev-infra`). The old implementation stays on `main` / `feature/dev-infra` until parity is reached.

| Doc | Contents |
| :--- | :--- |
| [dev-guide.md](./dev-guide.md) | Stack decisions, request flow, workspace layout, conventions, commands, env vars, remote-function rules, deployment, old → new migration map |
| [roadmap.md](./roadmap.md) | Phased plan with goals, deliverables and exit criteria |
| [backlog.md](./backlog.md) | Granular, checkbox-level task list — what is done and what is left |

**Related:** [`../architecture-revamp-plan.md`](../architecture-revamp-plan.md) (platform-wide revamp; its frontend section is implemented by this work) · root [`README.md`](../../README.md).

## Where things live

```
frontend/
├─ apps/
│  ├─ web/            # @ecommerce/web   — public storefront SSR (:3000), public routes only
│  └─ admin/          # @ecommerce/admin — internal backoffice SSR (:3001), strict ADMIN guard in hooks.server.ts
└─ packages/
   ├─ lib/            # @ecommerce/lib   — API client, BFF session/cookie helpers, Zod schemas, types
   └─ ui/             # @ecommerce/ui    — shared shadcn-svelte primitives & Tailwind 4 theme
```

- **`@ecommerce/lib`** — framework-free domain types, Zod 4 schemas (also imported by client forms), formatters, the server-only gateway client plus the BFF session/cookie helpers (`setSession`, `session`, `clearSession`, `shouldRefresh`, `safeRedirect`, `isAdmin`).
- **`@ecommerce/ui`** — one copy of the shadcn-svelte component set, exported per component (`@ecommerce/ui/button`, …), the shared brand theme (`@ecommerce/ui/theme.css`) and `cn` (from `tailwind-variants`) with the bits-ui prop type helpers. Apps import these; they never hold their own component copies.
- **Guards** — the storefront guards `/orders` and `/checkout` in their `+page.server.ts`; the backoffice guards *every* route in `apps/admin/src/hooks.server.ts` (anonymous → `/login?redirect=…`, non-admin → `/access-denied`).

## Status at a glance

- Stack scaffolded and verified (Bun workspaces, Kit 3, Tailwind 4, Paraglide vi/en, `@sveltejs/adapter-bun`, remote functions).
- Storefront is feature-complete: shell, home, catalog (filters/search/pagination), PDP, cart, SSO login/callback/logout, checkout (3-step, server-priced), orders.
- Backoffice is feature-complete for its four real screens (dashboard, products, orders, categories) plus stubs, behind the strict admin guard.
- All four packages pass `check` + `lint`; both apps build and are smoke-tested (guards, locales, SSR).
- Not yet verified against a running backend (no cluster on the dev machine) and not yet deployed: Dockerfiles, CI, compose/k8s env, Makefile/dev.ps1 targets.
