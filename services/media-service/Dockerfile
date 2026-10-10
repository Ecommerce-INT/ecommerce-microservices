FROM oven/bun:1.4-alpine AS base

WORKDIR /app

COPY package.json bun.lock* ./
RUN bun install --frozen-lockfile || bun install

COPY src ./src

EXPOSE 8083

ENTRYPOINT ["bun", "run", "src/index.ts"]
