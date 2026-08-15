import { describe, expect, it } from "vitest";
import { createRuntimeObservation, emitRuntimeObservation, verifyRuntimeObservation } from "./runtime-observability";

describe("runtime observability", () => {
  it("creates release-bound tamper-evident observations", () => {
    const observation = createRuntimeObservation({ level: "info", event: "runtime.ready", release: "RC-000146", occurredAt: "2026-08-15T00:00:00.000Z" });
    expect(verifyRuntimeObservation(observation)).toBe(true);
    expect(verifyRuntimeObservation({ ...observation, event: "runtime.tampered" })).toBe(false);
  });

  it("redacts nested credentials", () => {
    const observation = createRuntimeObservation({ level: "warn", event: "connector.retry", release: "RC-000146", attributes: { authorization: "Bearer unsafe", nested: { apiKey: "unsafe", safe: "visible" } } });
    expect(observation.attributes).toEqual({ authorization: "[REDACTED]", nested: { apiKey: "[REDACTED]", safe: "visible" } });
  });

  it("emits one structured JSON record", () => {
    const original = console.error;
    const records: unknown[][] = [];
    console.error = (...items: unknown[]) => { records.push(items); };
    const observation = createRuntimeObservation({ level: "error", event: "request.error", release: "RC-000146" });
    emitRuntimeObservation(observation);
    console.error = original;
    expect(JSON.parse(String(records[0][0])).digest).toBe(observation.digest);
  });
});
