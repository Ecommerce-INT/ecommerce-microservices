import { describe, expect, it } from "bun:test";
import appObj from "./index";

describe("Notification Service API", () => {
  it("should respond to health check", async () => {
    const req = new Request("http://localhost:8090/actuator/health");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = (await res.json()) as { status: string };
    expect(body.status).toBe("UP");
  });

  it("should respond to /notification/actuator/health", async () => {
    const req = new Request("http://localhost:8090/notification/actuator/health");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = (await res.json()) as { status: string };
    expect(body.status).toBe("UP");
  });

  it("should handle sendSimpleMail", async () => {
    const req = new Request("http://localhost:8090/api/email/sendSimpleMail", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        recipient: "test@example.com",
        subject: "Hello Test",
        msgBody: "Test body content",
      }),
    });
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const text = await res.text();
    expect(text).toContain("Mail Sent Successfully");
  });

  it("should handle sendMail with query params", async () => {
    const req = new Request("http://localhost:8090/api/email/sendMail?to=customer@example.com&subject=Welcome&body=Hello", {
      method: "POST",
    });
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const text = await res.text();
    expect(text).toContain("Mail Sent Successfully");
  });

  it("should list notifications", async () => {
    const req = new Request("http://localhost:8090/api/notifications");
    const res = await appObj.fetch(req);
    expect(res.status).toBe(200);
    const body = (await res.json()) as { status: string };
    expect(Array.isArray(body)).toBe(true);
  });
});
