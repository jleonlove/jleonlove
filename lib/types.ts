export type Classification = "Public" | "Internal" | "Confidential" | "Restricted";
export type AtlasRole = "owner" | "admin" | "manager" | "member" | "viewer";

export interface Organization {
  id: string;
  name: string;
  slug: string;
  plan: "Genesis" | "Business" | "Enterprise";
  createdAt: string;
}

export interface Workspace {
  id: string;
  organizationId: string;
  name: string;
  slug: string;
  createdAt: string;
}

export interface AtlasUser {
  id: string;
  name: string;
  email: string;
  createdAt: string;
}

export interface Membership {
  id: string;
  organizationId: string;
  userId: string;
  role: AtlasRole;
  createdAt: string;
}

export interface AtlasRequestContext {
  organization: Organization;
  workspace: Workspace;
  user: AtlasUser;
  role: AtlasRole;
}

export interface AtlasDocument {
  id: string;
  organizationId: string;
  workspaceId: string;
  createdBy: string;
  name: string;
  classification: Classification;
  content: string;
  createdAt: string;
  chunks: number;
}

export interface TrustDecision {
  id: string;
  organizationId: string;
  workspaceId: string;
  actorId: string;
  actorRole: AtlasRole;
  decision: "ALLOW" | "DENY";
  action: string;
  policyId: string;
  riskScore: number;
  evidence: string[];
  createdAt: string;
}

export interface SearchHit {
  documentId: string;
  documentName: string;
  classification: Classification;
  excerpt: string;
  score: number;
}

export type AgentReleaseStatus = "DRAFT" | "CANDIDATE" | "APPROVED" | "REVOKED";
export type ToolRisk = "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";

export interface AgentRelease {
  id: string;
  organizationId: string;
  workspaceId: string;
  agentId: string;
  agentName: string;
  version: string;
  digest: string;
  status: AgentReleaseStatus;
  allowedTools: string[];
  memoryScopes: Classification[];
  policySuite: string;
  evaluationScore: number;
  rollbackReleaseId?: string;
  approvedBy?: string;
  createdAt: string;
  approvedAt?: string;
}

export interface ToolDefinition {
  id: string;
  name: string;
  description: string;
  risk: ToolRisk;
  requiresApproval: boolean;
}

export interface ToolExecution {
  id: string;
  organizationId: string;
  workspaceId: string;
  actorId: string;
  agentReleaseId: string;
  toolId: string;
  purpose: string;
  status: "EXECUTED" | "BLOCKED" | "AWAITING_APPROVAL";
  trustDecisionId: string;
  createdAt: string;
}

export type EvidenceStatus = "CLAIMED" | "DOCUMENT_SUPPORTED" | "INDEPENDENTLY_VERIFIED" | "CONTRADICTED" | "STALE" | "UNKNOWN";
export type DealRequirementStatus = "OPEN" | "IN_PROGRESS" | "SATISFIED" | "BLOCKED" | "WAIVED";
export type DealSeverity = "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
export interface DealEvidence { id:string; claim:string; status:EvidenceStatus; sourceIds:string[]; independentSourceGroups:number; observedAt:string; expiresAt?:string; }
export interface DealRequirement { id:string; title:string; status:DealRequirementStatus; severity:DealSeverity; ownerRole:string; dependsOn:string[]; evidenceIds:string[]; dueAt?:string; }
export interface AtlasDeal { id:string; organizationId:string; workspaceId:string; createdBy:string; title:string; commodity:string; quantity?:number; unit?:string; stage:string; currency?:string; indicativeValue?:number; participants:Array<{id:string;name:string;role:string;authority:EvidenceStatus}>; evidence:DealEvidence[]; requirements:DealRequirement[]; assumptions:Array<{id:string;statement:string;dependencyKeys:string[];valid:boolean}>; createdAt:string; updatedAt:string; }
export interface DealDiagnosis { dealId:string; readiness:number; integrity:"LOW"|"MEDIUM"|"HIGH"; criticalBlockers:DealRequirement[]; nextBestAction?:DealRequirement; contradictions:DealEvidence[]; staleEvidence:DealEvidence[]; unknowns:DealEvidence[]; warnings:string[]; }

export interface DealEvent { id:string; dealId:string; type:string; actorId?:string; occurredAt:string; payload:Record<string,unknown>; idempotencyKey:string; }
export interface ClosingGate { pass:boolean; reasons:string[]; requiredApprovals:string[]; }
