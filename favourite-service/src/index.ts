import { Hono } from "hono";
import { cors } from "hono/cors";
import { logger } from "hono/logger";
import { initDb } from "./db";
import { FavouriteService } from "./service";
import { FavouriteDtoSchema, FavouriteIdSchema } from "./types";

const app = new Hono();
const favouriteService = new FavouriteService();

// Middlewares
app.use("*", logger());
app.use(
  "*",
  cors({
    origin: "*",
    allowHeaders: ["Authorization", "Content-Type", "X-Correlation-Id", "X-Request-Id"],
    allowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
  })
);

// Correlation ID middleware
app.use("*", async (c, next) => {
  const correlationId = c.req.header("X-Correlation-Id") || crypto.randomUUID();
  c.header("X-Correlation-Id", correlationId);
  await next();
});

// Health checks
const healthHandler = (c: any) =>
  c.json({
    status: "UP",
    components: {
      db: { status: "UP" },
    },
  });

app.get("/actuator/health", healthHandler);
app.get("/favourite/actuator/health", healthHandler);

// Helper to register routes on a specific prefix
function registerFavouriteRoutes(router: Hono) {
  // Find all
  router.get("", async (c) => {
    try {
      const list = await favouriteService.findAll();
      return c.json({ collection: list });
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Find by body
  router.get("/find", async (c) => {
    try {
      const body = await c.req.json();
      const parsed = FavouriteIdSchema.safeParse(body);
      if (!parsed.success) {
        return c.json({ statusCode: "400", message: "Invalid payload" }, 400);
      }
      const item = await favouriteService.findById(parsed.data);
      if (!item) {
        return c.json({ statusCode: "404", message: "Favourite not found" }, 404);
      }
      return c.json(item);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Find by path params
  router.get("/:userId/:productId/:likeDate", async (c) => {
    try {
      const userId = parseInt(c.req.param("userId"), 10);
      const productId = parseInt(c.req.param("productId"), 10);
      const likeDate = decodeURIComponent(c.req.param("likeDate"));

      if (isNaN(userId) || isNaN(productId)) {
        return c.json({ statusCode: "400", message: "Invalid IDs" }, 400);
      }

      const item = await favouriteService.findById({ userId, productId, likeDate });
      if (!item) {
        return c.json({ statusCode: "404", message: "Favourite not found" }, 404);
      }
      return c.json(item);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Save (create)
  router.post("", async (c) => {
    try {
      const body = await c.req.json();
      const parsed = FavouriteDtoSchema.safeParse(body);
      if (!parsed.success) {
        return c.json({ statusCode: "400", message: "Validation error", errors: parsed.error.issues }, 400);
      }
      const saved = await favouriteService.save(parsed.data);
      return c.json(saved);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Update
  router.put("", async (c) => {
    try {
      const body = await c.req.json();
      const parsed = FavouriteDtoSchema.safeParse(body);
      if (!parsed.success) {
        return c.json({ statusCode: "400", message: "Validation error", errors: parsed.error.issues }, 400);
      }
      const updated = await favouriteService.update(parsed.data);
      return c.json(updated);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Delete by body
  router.delete("/delete", async (c) => {
    try {
      const body = await c.req.json();
      const parsed = FavouriteIdSchema.safeParse(body);
      if (!parsed.success) {
        return c.json({ statusCode: "400", message: "Invalid payload" }, 400);
      }
      const result = await favouriteService.deleteById(parsed.data);
      return c.json(result);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Delete by path params
  router.delete("/:userId/:productId/:likeDate", async (c) => {
    try {
      const userId = parseInt(c.req.param("userId"), 10);
      const productId = parseInt(c.req.param("productId"), 10);
      const likeDate = decodeURIComponent(c.req.param("likeDate"));

      if (isNaN(userId) || isNaN(productId)) {
        return c.json({ statusCode: "400", message: "Invalid IDs" }, 400);
      }

      const result = await favouriteService.deleteById({ userId, productId, likeDate });
      return c.json(result);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });
}

// Register for both context path /favourite/api/favourites and direct /api/favourites
const favouriteGroup = new Hono();
registerFavouriteRoutes(favouriteGroup);
app.route("/favourite/api/favourites", favouriteGroup);
app.route("/api/favourites", favouriteGroup);

const PORT = parseInt(process.env.PORT || "8081", 10);

console.log(`[Favourite Service] Starting Bun server on port ${PORT}...`);
await initDb();

export default {
  port: PORT,
  fetch: app.fetch,
};
