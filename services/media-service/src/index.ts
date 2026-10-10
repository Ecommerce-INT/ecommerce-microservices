import { Hono } from "hono";
import { cors } from "hono/cors";
import { logger } from "hono/logger";
import { initDb } from "./db";
import { initStorage } from "./storage";
import { MediaService } from "./service";

const app = new Hono();
const mediaService = new MediaService();

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
      storage: { status: "UP" },
    },
  });

app.get("/actuator/health", healthHandler);
app.get("/media/actuator/health", healthHandler);

function registerMediaRoutes(router: Hono) {
  // Upload multipart
  router.post("/medias", async (c) => {
    try {
      const body = await c.req.parseBody();
      const file = body["multipartFile"] || body["file"];
      if (!file || !(file instanceof File)) {
        return c.json({ statusCode: "400", message: "File is required" }, 400);
      }
      const caption = typeof body["caption"] === "string" ? body["caption"] : "";
      const fileNameOverride = typeof body["fileNameOverride"] === "string" ? body["fileNameOverride"] : "";

      const result = await mediaService.saveMedia(file, caption, fileNameOverride);
      return c.json(result);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Delete media
  router.delete("/medias/:id", async (c) => {
    try {
      const id = parseInt(c.req.param("id"), 10);
      if (isNaN(id)) {
        return c.json({ statusCode: "400", message: "Invalid ID" }, 400);
      }
      await mediaService.removeMedia(id);
      return c.body(null, 204);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Get single media
  router.get("/medias/:id", async (c) => {
    try {
      const id = parseInt(c.req.param("id"), 10);
      if (isNaN(id)) {
        return c.json({ statusCode: "400", message: "Invalid ID" }, 400);
      }
      const media = await mediaService.getMediaById(id);
      if (!media) {
        return c.json({ statusCode: "404", message: "Media not found" }, 404);
      }
      return c.json(media);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Get multiple medias by IDs
  router.get("/medias", async (c) => {
    try {
      const idsParam = c.req.query("ids") || "";
      if (!idsParam) {
        return c.json({ statusCode: "400", message: "ids query param is required" }, 400);
      }
      const ids = idsParam
        .split(",")
        .map((s) => parseInt(s.trim(), 10))
        .filter((n) => !isNaN(n));

      const list = await mediaService.getMediaByIds(ids);
      if (list.length === 0) {
        return c.json({ statusCode: "404", message: "No medias found" }, 404);
      }
      return c.json(list);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });

  // Download / stream file
  router.get("/medias/:id/file/:fileName", async (c) => {
    try {
      const id = parseInt(c.req.param("id"), 10);
      const fileName = decodeURIComponent(c.req.param("fileName"));
      if (isNaN(id)) {
        return c.json({ statusCode: "400", message: "Invalid ID" }, 400);
      }
      const fileData = await mediaService.getFile(id, fileName);
      if (!fileData || !fileData.stream) {
        return c.json({ statusCode: "404", message: "File not found" }, 404);
      }
      c.header("Content-Disposition", `attachment; filename="${fileName}"`);
      c.header("Content-Type", fileData.mediaType);
      return c.body(fileData.stream);
    } catch (err: any) {
      return c.json({ statusCode: "500", message: err.message }, 500);
    }
  });
}

// Register with /media prefix and directly
const mediaGroup = new Hono();
registerMediaRoutes(mediaGroup);
app.route("/media", mediaGroup);
registerMediaRoutes(app);

const PORT = parseInt(process.env.PORT || "8083", 10);

console.log(`[Media Service] Starting Bun server on port ${PORT}...`);
await initDb();
await initStorage();

export default {
  port: PORT,
  fetch: app.fetch,
};
