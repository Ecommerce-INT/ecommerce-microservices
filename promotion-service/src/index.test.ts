import { describe, expect, it } from "bun:test";
import appObj from "./index";

describe("Promotion Service API", () => {
  it("should respond to health check", async () => {
    const req = new Request("http://localhost:8093/actuator/health");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = await res.json();
    expect(body.status).toBe("UP");
  });

  it("should respond to /promotion/actuator/health", async () => {
    const req = new Request("http://localhost:8093/promotion/actuator/health");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = await res.json();
    expect(body.status).toBe("UP");
  });
});
