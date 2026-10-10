import { describe, expect, it } from "bun:test";
import appObj from "./index";

describe("Favourite Service API", () => {
  it("should respond to /actuator/health", async () => {
    const req = new Request("http://localhost:8081/actuator/health");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = await res.json();
    expect(body.status).toBe("UP");
  });

  it("should respond to /favourite/actuator/health", async () => {
    const req = new Request("http://localhost:8081/favourite/actuator/health");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = await res.json();
    expect(body.status).toBe("UP");
  });
});
