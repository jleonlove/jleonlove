/**
 * Atlas durable workflow core (prototype).
 * Pure transition logic only: persistence, signatures, event authentication and
 * external side effects MUST be supplied by a separately qualified host.
 */
export type WorkflowStatus = "WAITING" | "RUNNING" | "COMPLETED" | "QUARANTINED";
export type WorkflowCheckpoint = Readonly<{
  tenantId: string;
  transactionId: string;
  workflowVersion: number;
  status: WorkflowStatus;
  cursor: number;
  processedEventIds: readonly string[];
  releaseDigest: string;
  policyVersion: string;
  principalId: string;
  agentId: string;
  expiresAtMs: number;
}>;
export type ResumeEvent = Readonly<{
  id: string;
  tenantId: string;
  transactionId: string;
  sequence: number;
  type: "DEPENDENCY_READY" | "COMPLETE";
}>;
export type ResumeAuthority = Readonly<{
  eventAuthenticated: boolean;
  checkpointIntegrityVerified: boolean;
  principalAuthorized: boolean;
  agentAuthorized: boolean;
  releaseAuthorized: boolean;
  policyAuthorized: boolean;
  complianceAuthorized: boolean;
  notRevoked: boolean;
  versionCompatible: boolean;
}>;
export type ResumeResult = Readonly<{
  checkpoint: WorkflowCheckpoint;
  outcome: "DUPLICATE" | "RESUMED" | "COMPLETED" | "QUARANTINED" | "REJECTED";
  reason?: string;
}>;
const valid = (s: string) => s.length > 0 && s.trim() === s;
export function resumeWorkflow(
  checkpoint: WorkflowCheckpoint,
  event: ResumeEvent,
  authority: ResumeAuthority,
  nowMs: number,
): ResumeResult {
  if (!Number.isFinite(nowMs) || !valid(event.id) ||
      !valid(checkpoint.tenantId) || !valid(checkpoint.transactionId) ||
      !valid(checkpoint.principalId) || !valid(checkpoint.agentId) ||
      !valid(checkpoint.releaseDigest) || !valid(checkpoint.policyVersion) ||
      !Number.isSafeInteger(checkpoint.cursor) || checkpoint.cursor < 0 ||
      !Number.isSafeInteger(event.sequence) || event.sequence < 1 ||
      !Number.isSafeInteger(checkpoint.workflowVersion) || checkpoint.workflowVersion < 1) {
    return { checkpoint, outcome: "REJECTED", reason: "Malformed workflow input" };
  }
  if (checkpoint.tenantId !== event.tenantId || checkpoint.transactionId !== event.transactionId) {
    return { checkpoint, outcome: "REJECTED", reason: "Workflow scope mismatch" };
  }
  if (checkpoint.status === "QUARANTINED" || checkpoint.status === "COMPLETED") {
    return { checkpoint, outcome: "REJECTED", reason: "Terminal workflow state" };
  }
  // Authorization checks precede deduplication so revoked actors cannot use
  // a duplicate event to receive an apparent successful response.
  const approved = Object.values(authority).every(value => value === true);
  if (!approved || nowMs >= checkpoint.expiresAtMs) {
    return { checkpoint: { ...checkpoint, status: "QUARANTINED" },
      outcome: "QUARANTINED", reason: "Authorization or expiry gate failed" };
  }
  if (checkpoint.processedEventIds.includes(event.id)) {
    return { checkpoint, outcome: "DUPLICATE" };
  }
  if (checkpoint.status !== "WAITING" || event.sequence !== checkpoint.cursor + 1) {
    return { checkpoint, outcome: "REJECTED", reason: "Invalid state or event order" };
  }
  const completed = event.type === "COMPLETE";
  return {
    checkpoint: {
      ...checkpoint,
      status: completed ? "COMPLETED" : "RUNNING",
      cursor: event.sequence,
      processedEventIds: [...checkpoint.processedEventIds, event.id],
    },
    outcome: completed ? "COMPLETED" : "RESUMED",
  };
}
