import { Hono } from "hono";
import { cors } from "hono/cors";
import { logger } from "hono/logger";
import { initDb } from "./db";
import { RatingService } from "./service";
import { RatingPostVmSchema } from "./types";

const app = new Hono();
const ratingService = new RatingService();

app.use("*", logger());
app.use(
  "*",
  cors({
    origin: "*",
    allowHeaders: ["Authorization", "Content-Type", "X-Correlation-Id", "X-Request-Id"],
    allowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
  })
);

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
app.get("/rating/actuator/health", healthHandler);

function extractUser(c: any) {
  const authHeader = c.req.header("Authorization") || "";
  let userId = "anonymous";
  let firstName = "Customer";
  let lastName = "";

  if (authHeader.startsWith("Bearer ")) {
    try {
      const parts = authHeader.split(".");
      if (parts.length === 3) {
        const payload = JSON.parse(Buffer.from(parts[1], "base64").toString());
        userId = payload.sub || userId;
        firstName = payload.given_name || payload.name || firstName;
        lastName = payload.family_name || lastName;
      }
    } catch {}
  }

  return { userId, firstName, lastName };
}

function registerRatingRoutes(router: Hono) {
  // Backoffice ratings filter
  router.get("/backoffice/ratings", async (c) => {
    try {
      const proName = c.req.query("productName") || "";
      const cusName = c.req.query("customerName") || "";
      const message = c.req.query("message") || "";
      const createdFrom = c.req.query("createdFrom") || "1970-01-01T00:00:00.000Z";
      const createdTo = c.req.query("createdTo") || new Date().toISOString();
      const pageNo = parseInt(c.req.query("pageNo") || "0", 10);
      const pageSize = parseInt(c.req.query("pageSize") || "5", 10);

      const result = await ratingService.getRatingListWithFilter(
        proName,
        cusName,
        message,
        createdFrom,
        createdTo,
        pageNo,
        pageSize
      );
      return c.json(result);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Backoffice delete rating
  router.delete("/backoffice/ratings/:id", async (c) => {
    try {
      const id = parseInt(c.req.param("id"), 10);
      if (isNaN(id)) {
        return c.json({ statusCode: "400", message: "Invalid ID" }, 400);
      }
      const result = await ratingService.deleteRating(id);
      return c.json(result);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Storefront get ratings by product ID
  router.get("/storefront/ratings/products/:productId", async (c) => {
    try {
      const productId = parseInt(c.req.param("productId"), 10);
      if (isNaN(productId)) {
        return c.json({ statusCode: "400", message: "Invalid productId" }, 400);
      }
      const pageNo = parseInt(c.req.query("pageNo") || "0", 10);
      const pageSize = parseInt(c.req.query("pageSize") || "5", 10);

      const result = await ratingService.getRatingListByProductId(productId, pageNo, pageSize);
      return c.json(result);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Storefront create rating
  router.post("/storefront/ratings", async (c) => {
    try {
      const body = await c.req.json();
      const parsed = RatingPostVmSchema.safeParse(body);
      if (!parsed.success) {
        return c.json({ statusCode: "400", message: "Validation error", errors: parsed.error.issues }, 400);
      }

      const user = extractUser(c);
      const created = await ratingService.createRating(parsed.data, user);
      return c.json(created);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Storefront calculate average star
  router.get("/storefront/ratings/product/:productId/average-star", async (c) => {
    try {
      const productId = parseInt(c.req.param("productId"), 10);
      if (isNaN(productId)) {
        return c.json(0.0);
      }
      const avg = await ratingService.calculateAverageStar(productId);
      return c.json(avg);
    } catch (err: any) {
      return c.json(0.0);
    }
  });
}

// Register for both context path /rating and direct root
const ratingGroup = new Hono();
registerRatingRoutes(ratingGroup);
app.route("/rating", ratingGroup);
registerRatingRoutes(app);

const PORT = parseInt(process.env.PORT || "8089", 10);

console.log(`[Rating Service] Starting Bun server on port ${PORT}...`);
await initDb();

export default {
  port: PORT,
  fetch: app.fetch,
};
