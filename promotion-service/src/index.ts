import { Hono } from "hono";
import { cors } from "hono/cors";
import { logger } from "hono/logger";
import { initDb } from "./db";
import { PromotionService } from "./service";
import { PromotionPostVmSchema, PromotionPutVmSchema } from "./types";
import { TaxService } from "./tax-service";
import { TaxClassPostVmSchema, TaxRatePostVmSchema } from "./tax-types";

const app = new Hono();
const promotionService = new PromotionService();
const taxService = new TaxService();

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
app.get("/tax/actuator/health", healthHandler);

// =================== Promotion Routes ===================

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

// =================== Tax Routes ===================

function registerTaxRoutes(router: Hono) {
  // Tax Classes
  router.get("/backoffice/tax-classes/paging", async (c) => {
    const pageNo = parseInt(c.req.query("pageNo") || "0", 10);
    const pageSize = parseInt(c.req.query("pageSize") || "10", 10);
    const res = await taxService.getPageableTaxClasses(pageNo, pageSize);
    return c.json(res);
  });

  router.get("/backoffice/tax-classes/:id", async (c) => {
    const id = parseInt(c.req.param("id"), 10);
    if (isNaN(id)) return c.json({ message: "Invalid tax class ID" }, 400);
    const tc = await taxService.findTaxClassById(id);
    if (!tc) return c.json({ message: "Tax class not found" }, 404);
    return c.json(tc);
  });

  router.get("/backoffice/tax-classes", async (c) => {
    const list = await taxService.findAllTaxClasses();
    return c.json(list);
  });

  router.post("/backoffice/tax-classes", async (c) => {
    const body = await c.req.json();
    const parsed = TaxClassPostVmSchema.safeParse(body);
    if (!parsed.success) return c.json({ message: "Invalid payload" }, 400);
    const created = await taxService.createTaxClass(parsed.data);
    c.header("Location", `/tax-classes/${created.id}`);
    return c.json(created, 201);
  });

  router.put("/backoffice/tax-classes/:id", async (c) => {
    const id = parseInt(c.req.param("id"), 10);
    if (isNaN(id)) return c.json({ message: "Invalid tax class ID" }, 400);
    const body = await c.req.json();
    const parsed = TaxClassPostVmSchema.safeParse(body);
    if (!parsed.success) return c.json({ message: "Invalid payload" }, 400);
    const ok = await taxService.updateTaxClass(id, parsed.data);
    if (!ok) return c.json({ message: "Tax class not found" }, 404);
    return c.body(null, 204);
  });

  router.delete("/backoffice/tax-classes/:id", async (c) => {
    const id = parseInt(c.req.param("id"), 10);
    if (isNaN(id)) return c.json({ message: "Invalid tax class ID" }, 400);
    const ok = await taxService.deleteTaxClass(id);
    if (!ok) return c.json({ message: "Tax class not found" }, 404);
    return c.body(null, 204);
  });

  // Tax Rates
  router.get("/backoffice/tax-rates/paging", async (c) => {
    const pageNo = parseInt(c.req.query("pageNo") || "0", 10);
    const pageSize = parseInt(c.req.query("pageSize") || "10", 10);
    const res = await taxService.getPageableTaxRates(pageNo, pageSize);
    return c.json(res);
  });

  router.get("/backoffice/tax-rates/tax-percent", async (c) => {
    const taxClassId = parseInt(c.req.query("taxClassId") || "", 10);
    const countryId = parseInt(c.req.query("countryId") || "", 10);
    if (isNaN(taxClassId) || isNaN(countryId)) {
      return c.json({ message: "taxClassId and countryId are required" }, 400);
    }
    const stateOrProvinceId = c.req.query("stateOrProvinceId") ? parseInt(c.req.query("stateOrProvinceId")!, 10) : null;
    const zipCode = c.req.query("zipCode") || null;
    const percent = await taxService.getTaxPercent(taxClassId, countryId, stateOrProvinceId, zipCode);
    return c.json(percent);
  });

  router.get("/backoffice/tax-rates/location-based-batch", async (c) => {
    const idsStr = c.req.query("taxClassIds") || "";
    const countryId = parseInt(c.req.query("countryId") || "", 10);
    if (!idsStr || isNaN(countryId)) {
      return c.json({ message: "taxClassIds and countryId are required" }, 400);
    }
    const taxClassIds = idsStr.split(",").map((s) => parseInt(s.trim(), 10)).filter((n) => !isNaN(n));
    const stateOrProvinceId = c.req.query("stateOrProvinceId") ? parseInt(c.req.query("stateOrProvinceId")!, 10) : null;
    const zipCode = c.req.query("zipCode") || null;
    const list = await taxService.getBatchTaxRates(taxClassIds, countryId, stateOrProvinceId, zipCode);
    return c.json(list);
  });

  router.get("/backoffice/tax-rates/:id", async (c) => {
    const id = parseInt(c.req.param("id"), 10);
    if (isNaN(id)) return c.json({ message: "Invalid tax rate ID" }, 400);
    const tr = await taxService.findTaxRateById(id);
    if (!tr) return c.json({ message: "Tax rate not found" }, 404);
    return c.json(tr);
  });

  router.get("/backoffice/tax-rates", async (c) => {
    const list = await taxService.findAllTaxRates();
    return c.json(list);
  });

  router.post("/backoffice/tax-rates", async (c) => {
    const body = await c.req.json();
    const parsed = TaxRatePostVmSchema.safeParse(body);
    if (!parsed.success) return c.json({ message: "Invalid payload" }, 400);
    const created = await taxService.createTaxRate(parsed.data);
    c.header("Location", `/tax-rates/${created.id}`);
    return c.json(created, 201);
  });

  router.put("/backoffice/tax-rates/:id", async (c) => {
    const id = parseInt(c.req.param("id"), 10);
    if (isNaN(id)) return c.json({ message: "Invalid tax rate ID" }, 400);
    const body = await c.req.json();
    const parsed = TaxRatePostVmSchema.safeParse(body);
    if (!parsed.success) return c.json({ message: "Invalid payload" }, 400);
    const ok = await taxService.updateTaxRate(id, parsed.data);
    if (!ok) return c.json({ message: "Tax rate not found" }, 404);
    return c.body(null, 204);
  });

  router.delete("/backoffice/tax-rates/:id", async (c) => {
    const id = parseInt(c.req.param("id"), 10);
    if (isNaN(id)) return c.json({ message: "Invalid tax rate ID" }, 400);
    const ok = await taxService.deleteTaxRate(id);
    if (!ok) return c.json({ message: "Tax rate not found" }, 404);
    return c.body(null, 204);
  });
}

// Register Promotion routes
const promotionGroup = new Hono();
registerPromotionRoutes(promotionGroup);
app.route("/promotion", promotionGroup);
registerPromotionRoutes(app);

// Register Tax routes with /tax prefix and root
const taxGroup = new Hono();
registerTaxRoutes(taxGroup);
app.route("/tax", taxGroup);
registerTaxRoutes(app);

const PORT = parseInt(process.env.PORT || "8093", 10);

console.log(`[Promotion & Tax Service] Starting Bun server on port ${PORT}...`);
await initDb();

export default {
  port: PORT,
  fetch: app.fetch,
};
