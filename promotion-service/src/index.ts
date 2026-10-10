import { Hono } from "hono";
import { cors } from "hono/cors";
import { logger } from "hono/logger";
import { initDb } from "./db";
import { PromotionService } from "./service";
import { PromotionPostVmSchema, PromotionPutVmSchema } from "./types";

const app = new Hono();
const promotionService = new PromotionService();

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
app.get("/promotion/actuator/health", healthHandler);

function registerPromotionRoutes(router: Hono) {
  // List promotions
  router.get("/backoffice/promotions", async (c) => {
    try {
      const pageNo = parseInt(c.req.query("pageNo") || "0", 10);
      const pageSize = parseInt(c.req.query("pageSize") || "5", 10);
      const promotionName = c.req.query("promotionName") || "";
      const couponCode = c.req.query("couponCode") || "";
      const startDate = c.req.query("startDate");
      const endDate = c.req.query("endDate");

      const res = await promotionService.getPromotions(
        pageNo,
        pageSize,
        promotionName,
        couponCode,
        startDate,
        endDate
      );
      return c.json(res);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Get single promotion
  router.get("/backoffice/promotions/:promotionId", async (c) => {
    try {
      const id = parseInt(c.req.param("promotionId"), 10);
      if (isNaN(id)) return c.json({ statusCode: "400", message: "Invalid ID" }, 400);

      const p = await promotionService.getPromotion(id);
      if (!p) return c.json({ statusCode: "404", message: "Promotion not found" }, 404);
      return c.json(p);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Create promotion
  router.post("/backoffice/promotions", async (c) => {
    try {
      const body = await c.req.json();
      const parsed = PromotionPostVmSchema.safeParse(body);
      if (!parsed.success) {
        return c.json({ statusCode: "400", message: "Validation error", errors: parsed.error.issues }, 400);
      }
      const created = await promotionService.createPromotion(parsed.data);
      return c.json(created, 201);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Update promotion
  router.put("/backoffice/promotions", async (c) => {
    try {
      const body = await c.req.json();
      const parsed = PromotionPutVmSchema.safeParse(body);
      if (!parsed.success) {
        return c.json({ statusCode: "400", message: "Validation error", errors: parsed.error.issues }, 400);
      }
      const updated = await promotionService.updatePromotion(parsed.data);
      return c.json(updated);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Delete promotion
  router.delete("/backoffice/promotions/:promotionId", async (c) => {
    try {
      const id = parseInt(c.req.param("promotionId"), 10);
      if (isNaN(id)) return c.json({ statusCode: "400", message: "Invalid ID" }, 400);

      await promotionService.deletePromotion(id);
      return c.body(null, 200);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });
}

// Register with /promotion prefix and root
const promotionGroup = new Hono();
registerPromotionRoutes(promotionGroup);
app.route("/promotion", promotionGroup);
registerPromotionRoutes(app);

const PORT = parseInt(process.env.PORT || "8093", 10);

console.log(`[Promotion Service] Starting Bun server on port ${PORT}...`);
await initDb();

export default {
  port: PORT,
  fetch: app.fetch,
};
