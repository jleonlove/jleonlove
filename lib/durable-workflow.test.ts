import { describe, expect, it } from "vitest";
import { resumeWorkflow, type ResumeAuthority, type WorkflowCheckpoint } from "./durable-workflow";

const checkpoint: WorkflowCheckpoint = {
  tenantId: "tenant-a", transactionId: "trade-1", workflowVersion: 1,
  status: "WAITING", cursor: 0, processedEventIds: [],
  releaseDigest: "sha256:abc", policyVersion: "p1", principalId: "user-1",
  agentId: "agent-1", expiresAtMs: 2000,
};
const authority: ResumeAuthority = {
  eventAuthenticated: true, checkpointIntegrityVerified: true,
  principalAuthorized: true, agentAuthorized: true, releaseAuthorized: true,
  policyAuthorized: true, complianceAuthorized: true, notRevoked: true,
  versionCompatible: true,
};
const event = { id: "evt-1", tenantId: "tenant-a", transactionId: "trade-1",
  sequence: 1, type: "DEPENDENCY_READY" as const };

describe("durable workflow authorization gates", () => {
  it("resumes an authorized event and records its cursor", () => {
    const result = resumeWorkflow(checkpoint, event, authority, 1000);
    expect(result.outcome).toBe("RESUMED");
    expect(result.checkpoint.cursor).toBe(1);
  });
  it("rejects cross-tenant events", () => {
    expect(resumeWorkflow(checkpoint, { ...event, tenantId: "tenant-b" }, authority, 1000).outcome).toBe("REJECTED");
  });
  it("quarantines revoked agents", () => {
    expect(resumeWorkflow(checkpoint, event, { ...authority, notRevoked: false }, 1000).outcome).toBe("QUARANTINED");
  });
  it("quarantines expired workflows", () => {
    expect(resumeWorkflow(checkpoint, event, authority, 2000).outcome).toBe("QUARANTINED");
  });
  it("rejects out-of-order events", () => {
    expect(resumeWorkflow(checkpoint, { ...event, sequence: 2 }, authority, 1000).outcome).toBe("REJECTED");
  });
  it("deduplicates events without advancing cursor", () => {
    const prior = { ...checkpoint, processedEventIds: ["evt-1"] };
    expect(resumeWorkflow(prior, event, authority, 1000).outcome).toBe("DUPLICATE");
  });
  it("does not accept an already completed workflow", () => {
    expect(resumeWorkflow({ ...checkpoint, status: "COMPLETED" }, event, authority, 1000).outcome).toBe("REJECTED");
  });
  it("fails closed if checkpoint integrity is not verified", () => {
    expect(resumeWorkflow(checkpoint, event, { ...authority, checkpointIntegrityVerified: false }, 1000).outcome).toBe("QUARANTINED");
  });
});
