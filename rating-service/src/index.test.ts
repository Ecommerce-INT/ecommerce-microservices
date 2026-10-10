import { describe, expect, it } from "bun:test";
import appObj from "./index";

describe("Rating Service API", () => {
  it("should respond to health checks", async () => {
    const req = new Request("http://localhost:8089/actuator/health");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = await res.json();
    expect(body.status).toBe("UP");
  });

  it("should respond to context path /rating/actuator/health", async () => {
    const req = new Request("http://localhost:8089/rating/actuator/health");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = await res.json();
    expect(body.status).toBe("UP");
  });
});
