import { createHash } from "node:crypto";

/**
 * Phase 1: fail-closed validation of retrieved commodity evidence.
 * This is an application boundary, not a substitute for database RLS.
 */
export type EvidenceRecord = Readonly<{
  id: string;
  tenantId: string;
  workspaceId: string;
  sourceUri: string;
  sourceVersion: string;
  content: string;
  contentSha256: string;
  observedAt: string;
  expiresAt: string;
  verified: boolean;
}>;

export type EvidenceRequest = Readonly<{
  tenantId: string;
  workspaceId: string;
  now: string;
  requireVerified?: boolean;
}>;

export type EvidenceDecision =
  | { allowed: true; record: EvidenceRecord }
  | { allowed: false; reason: "INVALID_CONTEXT" | "INVALID_RECORD" | "TENANT_DENIED" | "WORKSPACE_DENIED" | "STALE" | "UNVERIFIED" | "DIGEST_MISMATCH" };

const sha256 = (value: string): string => createHash("sha256").update(value, "utf8").digest("hex");
const validTime = (value: string): number | null => {
  if (typeof value !== "string" || !/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{1,3})?Z$/.test(value)) return null;
  const parsed = Date.parse(value);
  return Number.isFinite(parsed) ? parsed : null;
};
const validId = (value: unknown): value is string => typeof value === "string" && value.trim().length > 0;

export function digestEvidence(content: string): string {
  return sha256(content);
}

export function authorizeEvidence(record: EvidenceRecord, request: EvidenceRequest): EvidenceDecision {
  if (!validId(request?.tenantId) || !validId(request?.workspaceId) || validTime(request?.now) === null)
    return { allowed: false, reason: "INVALID_CONTEXT" };

  if (!record || !validId(record.id) || !validId(record.tenantId) ||
      !validId(record.workspaceId) || !validId(record.sourceUri) ||
      !validId(record.sourceVersion) || typeof record.content !== "string" ||
      !/^[a-f0-9]{64}$/.test(record.contentSha256 || "") ||
      typeof record.verified !== "boolean" ||
      validTime(record.observedAt) === null || validTime(record.expiresAt) === null ||
      validTime(record.observedAt)! > validTime(record.expiresAt)!)
    return { allowed: false, reason: "INVALID_RECORD" };

  if (record.tenantId !== request.tenantId) return { allowed: false, reason: "TENANT_DENIED" };
  if (record.workspaceId !== request.workspaceId) return { allowed: false, reason: "WORKSPACE_DENIED" };
  if (validTime(request.now)! < validTime(record.observedAt)! || validTime(request.now)! >= validTime(record.expiresAt)!)
    return { allowed: false, reason: "STALE" };
  if (request.requireVerified !== false && !record.verified) return { allowed: false, reason: "UNVERIFIED" };
  if (sha256(record.content) !== record.contentSha256) return { allowed: false, reason: "DIGEST_MISMATCH" };
  return { allowed: true, record };
}

export function retrieveAuthorizedEvidence(records: readonly EvidenceRecord[], request: EvidenceRequest): EvidenceRecord[] {
  return records.flatMap((record) => {
    const decision = authorizeEvidence(record, request);
    return decision.allowed ? [decision.record] : [];
  });
}
