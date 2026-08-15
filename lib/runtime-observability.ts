import crypto from "node:crypto";

export type ObservationLevel = "info" | "warn" | "error";
export type ObservationAttributes = Record<string, unknown>;

export interface RuntimeObservation {
  schema: "atlas.runtime-observation.v1";
  service: "atlas-web";
  release: string;
  environment: string;
  level: ObservationLevel;
  event: string;
  occurredAt: string;
  requestId?: string;
  attributes: ObservationAttributes;
  digest: string;
}

const SENSITIVE_KEY = /(authorization|cookie|password|secret|token|api[-_]?key|private[-_]?key)/i;

function sanitize(value: unknown, key = "", seen = new WeakSet<object>()): unknown {
  if (SENSITIVE_KEY.test(key)) return "[REDACTED]";
  if (value === null || typeof value !== "object") return value;
  if (seen.has(value)) return "[CIRCULAR]";
  seen.add(value);
  if (Array.isArray(value)) return value.map((item) => sanitize(item, key, seen));
  return Object.fromEntries(Object.entries(value).map(([childKey, child]) => [childKey, sanitize(child, childKey, seen)]));
}

function canonical(value: unknown): string {
  if (value === null || typeof value !== "object") return JSON.stringify(value);
  if (Array.isArray(value)) return `[${value.map(canonical).join(",")}]`;
  return `{${Object.entries(value).sort(([a], [b]) => a.localeCompare(b)).map(([key, item]) => `${JSON.stringify(key)}:${canonical(item)}`).join(",")}}`;
}

export function observationDigest(value: Omit<RuntimeObservation, "digest">) {
  return crypto.createHash("sha256").update(canonical(value)).digest("hex");
}

export function createRuntimeObservation(input: {
  level: ObservationLevel;
  event: string;
  release: string;
  environment?: string;
  occurredAt?: string;
  requestId?: string;
  attributes?: ObservationAttributes;
}): RuntimeObservation {
  if (!input.event.trim()) throw new Error("OBSERVATION_EVENT_REQUIRED");
  if (!/^RC-\d{6}$/.test(input.release)) throw new Error("OBSERVATION_RELEASE_INVALID");
  const unsigned: Omit<RuntimeObservation, "digest"> = {
    schema: "atlas.runtime-observation.v1",
    service: "atlas-web",
    release: input.release,
    environment: input.environment ?? "unknown",
    level: input.level,
    event: input.event,
    occurredAt: input.occurredAt ?? new Date().toISOString(),
    ...(input.requestId ? { requestId: input.requestId } : {}),
    attributes: sanitize(input.attributes ?? {}) as ObservationAttributes,
  };
  return { ...unsigned, digest: observationDigest(unsigned) };
}

export function verifyRuntimeObservation(observation: RuntimeObservation) {
  const { digest, ...unsigned } = observation;
  return /^[a-f0-9]{64}$/.test(digest) && observationDigest(unsigned) === digest;
}

export function emitRuntimeObservation(observation: RuntimeObservation) {
  const line = JSON.stringify(observation);
  if (observation.level === "error") console.error(line);
  else if (observation.level === "warn") console.warn(line);
  else console.info(line);
}
