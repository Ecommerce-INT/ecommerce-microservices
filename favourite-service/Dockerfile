FROM oven/bun:1-alpine AS base

WORKDIR /app

COPY package.json bun.lock* ./
RUN bun install --frozen-lockfile || bun install

COPY src ./src

EXPOSE 8081

ENTRYPOINT ["bun", "run", "src/index.ts"]
