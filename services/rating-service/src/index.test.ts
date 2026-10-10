import { describe, expect, it } from "bun:test";
import appObj from "./index";

describe("Rating Service API", () => {
  it("should respond to health checks", async () => {
    const req = new Request("http://localhost:8089/actuator/health");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = (await res.json()) as { status: string };
    expect(body.status).toBe("UP");
  });

  it("should respond to context path /rating/actuator/health", async () => {
    const req = new Request("http://localhost:8089/rating/actuator/health");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = (await res.json()) as { status: string };
    expect(body.status).toBe("UP");
  });

  it("should calculate average star for product", async () => {
    const req = new Request("http://localhost:8089/storefront/ratings/product/101/average-star");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
  });

  it("should reject ratings query with invalid productId", async () => {
    const req = new Request("http://localhost:8089/storefront/ratings/products/xyz");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(400);
  });

  it("should reject rating submission with missing fields", async () => {
    const req = new Request("http://localhost:8089/storefront/ratings", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({}),
    });
    const res = await appObj.fetch(req);
    expect(res.status).toBe(400);
  });

  it("should reject delete rating with non-numeric ID", async () => {
    const req = new Request("http://localhost:8089/backoffice/ratings/invalid-id", {
      method: "DELETE",
    });
    const res = await appObj.fetch(req);
    expect(res.status).toBe(400);
  });
});
