import { createHash } from "node:crypto";
import { atlasRepository } from "./repository";
import type { AgentRelease, AtlasRequestContext, Classification } from "./types";

function digest(input: unknown) {
  return `sha256:${createHash("sha256").update(JSON.stringify(input)).digest("hex").slice(0, 16)}`;
}

export async function registerCandidate(input: {
  context: AtlasRequestContext;
  agentId: string;
  agentName: string;
  version: string;
  allowedTools: string[];
  memoryScopes: Classification[];
  policySuite: string;
  evaluationScore: number;
  rollbackReleaseId?: string;
}) {
  const release: AgentRelease = {
    id: `rel_${crypto.randomUUID()}`,
    organizationId: input.context.organization.id,
    workspaceId: input.context.workspace.id,
    agentId: input.agentId,
    agentName: input.agentName,
    version: input.version,
    digest: digest(input),
    status: "CANDIDATE",
    allowedTools: [...new Set(input.allowedTools)].sort(),
    memoryScopes: [...new Set(input.memoryScopes)],
    policySuite: input.policySuite,
    evaluationScore: Math.max(0, Math.min(100, input.evaluationScore)),
    rollbackReleaseId: input.rollbackReleaseId,
    createdAt: new Date().toISOString(),
  };
  return atlasRepository.addAgentRelease(release);
}

export async function approveRelease(id: string, context: AtlasRequestContext) {
  if (!['owner', 'admin'].includes(context.role)) throw new Error("Only owners and administrators can approve agent releases.");
  const release = await atlasRepository.agentRelease(id);
  if (!release || release.organizationId !== context.organization.id) throw new Error("Agent release not found.");
  if (release.evaluationScore < 90) throw new Error("Release evaluation score must be at least 90.");
  return atlasRepository.updateAgentRelease(id, { status: "APPROVED", approvedBy: context.user.id, approvedAt: new Date().toISOString() });
}

export async function revokeRelease(id: string, context: AtlasRequestContext) {
  if (!['owner', 'admin'].includes(context.role)) throw new Error("Only owners and administrators can revoke agent releases.");
  const release = await atlasRepository.agentRelease(id);
  if (!release || release.organizationId !== context.organization.id) throw new Error("Agent release not found.");
  return atlasRepository.updateAgentRelease(id, { status: "REVOKED" });
}
