import { describe, expect, it } from "bun:test";
import appObj from "./index";

describe("Promotion & Tax Service API", () => {
  it("should respond to health checks", async () => {
    const res1 = await appObj.fetch(new Request("http://localhost:8093/actuator/health"));
    expect(res1.status).toBe(200);

    const res2 = await appObj.fetch(new Request("http://localhost:8093/promotion/actuator/health"));
    expect(res2.status).toBe(200);

    const res3 = await appObj.fetch(new Request("http://localhost:8093/tax/actuator/health"));
    expect(res3.status).toBe(200);
  });

  it("should list promotions with paging parameters", async () => {
    const req = new Request("http://localhost:8093/backoffice/promotions?pageNo=0&pageSize=5");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
  });

  it("should reject promotion with invalid ID", async () => {
    const req = new Request("http://localhost:8093/backoffice/promotions/not-a-number");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(400);
  });

  it("should reject promotion creation with invalid payload", async () => {
    const req = new Request("http://localhost:8093/backoffice/promotions", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({}),
    });
    const res = await appObj.fetch(req);
    expect(res.status).toBe(400);
  });

  it("should reject tax rate percent query without required params", async () => {
    const req = new Request("http://localhost:8093/tax/backoffice/tax-rates/tax-percent");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(400);
  });

  it("should reject tax batch query without required params", async () => {
    const req = new Request("http://localhost:8093/tax/backoffice/tax-rates/location-based-batch");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(400);
  });

  it("should list tax classes paging", async () => {
    const req = new Request("http://localhost:8093/tax/backoffice/tax-classes/paging?pageNo=0&pageSize=10");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
  });
});
