import { describe, expect, it } from "vitest";
import { authorizeEvidence, digestEvidence, retrieveAuthorizedEvidence, type EvidenceRecord } from "./trusted-intelligence";

const content = "Wheat shipment inspection certificate #100";
const record: EvidenceRecord = {
  id: "doc-1", tenantId: "tenant-a", workspaceId: "trade-1",
  sourceUri: "s3://evidence/certificate-100", sourceVersion: "v1",
  content, contentSha256: digestEvidence(content),
  observedAt: "2026-11-01T00:00:00Z", expiresAt: "2026-12-01T00:00:00Z", verified: true,
};
const context = { tenantId: "tenant-a", workspaceId: "trade-1", now: "2026-11-05T00:00:00Z" };

describe("trusted evidence retrieval boundary", () => {
  it("allows intact, current, verified evidence in the same tenant and workspace", () => {
    expect(authorizeEvidence(record, context)).toEqual({ allowed: true, record });
  });
  it("denies cross-tenant evidence", () => {
    expect(authorizeEvidence(record, { ...context, tenantId: "tenant-b" })).toEqual({ allowed: false, reason: "TENANT_DENIED" });
  });
  it("denies cross-workspace evidence", () => {
    expect(authorizeEvidence(record, { ...context, workspaceId: "trade-2" })).toEqual({ allowed: false, reason: "WORKSPACE_DENIED" });
  });
  it("denies expired evidence at the exact expiry boundary", () => {
    expect(authorizeEvidence(record, { ...context, now: record.expiresAt })).toEqual({ allowed: false, reason: "STALE" });
  });
  it("denies evidence from the future", () => {
    expect(authorizeEvidence(record, { ...context, now: "2026-10-01T00:00:00Z" })).toEqual({ allowed: false, reason: "STALE" });
  });
  it("denies unverified evidence by default", () => {
    expect(authorizeEvidence({ ...record, verified: false }, context)).toEqual({ allowed: false, reason: "UNVERIFIED" });
  });
  it("denies tampered content", () => {
    expect(authorizeEvidence({ ...record, content: "Altered certificate" }, context)).toEqual({ allowed: false, reason: "DIGEST_MISMATCH" });
  });
  it("denies missing request identity", () => {
    expect(authorizeEvidence(record, { ...context, tenantId: "" })).toEqual({ allowed: false, reason: "INVALID_CONTEXT" });
  });
  it("denies invalid record timestamps", () => {
    expect(authorizeEvidence({ ...record, expiresAt: "tomorrow" }, context)).toEqual({ allowed: false, reason: "INVALID_RECORD" });
  });
  it("filters mixed-tenant records before returning any content", () => {
    expect(retrieveAuthorizedEvidence([record, { ...record, tenantId: "tenant-b", id: "doc-2" }], context)).toEqual([record]);
  });
});
