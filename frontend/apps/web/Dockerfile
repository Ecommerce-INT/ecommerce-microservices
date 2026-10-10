# syntax=docker/dockerfile:1

# Build context is the Bun workspace root (frontend/):
#   docker build -f frontend/apps/web/Dockerfile frontend

FROM docker.io/oven/bun:1 AS build
WORKDIR /app

# Manifests first so the dependency layer survives source-only changes.
COPY package.json bun.lock ./
COPY apps/web/package.json apps/web/package.json
COPY apps/admin/package.json apps/admin/package.json
COPY packages/lib/package.json packages/lib/package.json
COPY packages/ui/package.json packages/ui/package.json
RUN bun install --frozen-lockfile

# Vite must run under Bun — required by @sveltejs/adapter-bun.
COPY . .
RUN bun run --filter @ecommerce/web build

FROM docker.io/oven/bun:1 AS runtime
WORKDIR /app
ENV NODE_ENV=production \
    PORT=3000
COPY --from=build /app/apps/web/build ./build
USER bun
EXPOSE 3000
CMD ["bun", "./build"]
