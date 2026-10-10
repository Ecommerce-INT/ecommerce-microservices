import { describe, expect, it } from "bun:test";
import appObj from "./index";

describe("Media Service API", () => {
  it("should respond to health check", async () => {
    const req = new Request("http://localhost:8083/actuator/health");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = await res.json();
    expect(body.status).toBe("UP");
  });

  it("should respond to /media/actuator/health", async () => {
    const req = new Request("http://localhost:8083/media/actuator/health");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = await res.json();
    expect(body.status).toBe("UP");
  });
});
