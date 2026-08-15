import { atlasRepository } from "./repository";
import type { AtlasRequestContext, ToolDefinition, ToolExecution, TrustDecision } from "./types";

export const toolCatalog: ToolDefinition[] = [
  { id: "knowledge.search", name: "Knowledge Search", description: "Retrieve approved enterprise knowledge.", risk: "LOW", requiresApproval: false },
  { id: "report.generate", name: "Report Generator", description: "Create an internal report from approved sources.", risk: "MEDIUM", requiresApproval: false },
  { id: "workflow.create", name: "Workflow Creator", description: "Create a governed business workflow.", risk: "HIGH", requiresApproval: true },
  { id: "external.send", name: "External Sender", description: "Transmit information outside the organization.", risk: "CRITICAL", requiresApproval: true },
];

function riskScore(risk: ToolDefinition["risk"]) {
  return { LOW: 0.05, MEDIUM: 0.2, HIGH: 0.58, CRITICAL: 0.9 }[risk];
}

export async function executeTool(input: { context: AtlasRequestContext; agentReleaseId: string; toolId: string; purpose: string }) {
  const [release, tool] = await Promise.all([
    atlasRepository.agentRelease(input.agentReleaseId),
    Promise.resolve(toolCatalog.find((item) => item.id === input.toolId)),
  ]);
  if (!release || release.organizationId !== input.context.organization.id) throw new Error("Agent release not found.");
  if (!tool) throw new Error("Tool not found.");

  const releaseApproved = release.status === "APPROVED";
  const toolAllowed = release.allowedTools.includes(tool.id);
  const humanAuthorized = ['owner', 'admin'].includes(input.context.role);
  const awaitingApproval = releaseApproved && toolAllowed && tool.requiresApproval && !humanAuthorized;
  const allowed = releaseApproved && toolAllowed && (!tool.requiresApproval || humanAuthorized);
  const decision: TrustDecision = {
    id: `trust_${crypto.randomUUID()}`,
    organizationId: input.context.organization.id,
    workspaceId: input.context.workspace.id,
    actorId: input.context.user.id,
    actorRole: input.context.role,
    decision: allowed ? "ALLOW" : "DENY",
    action: `${release.agentName} requested ${tool.name}`,
    policyId: tool.requiresApproval ? "POL-TOOL-HIGH-01" : "POL-TOOL-BASE-01",
    riskScore: riskScore(tool.risk),
    evidence: [
      `Release status: ${release.status}`,
      `Release digest: ${release.digest}`,
      `Tool declared in manifest: ${toolAllowed ? "yes" : "no"}`,
      `Tool risk: ${tool.risk}`,
      `Human approval required: ${tool.requiresApproval ? "yes" : "no"}`,
      `Actor role: ${input.context.role}`,
    ],
    createdAt: new Date().toISOString(),
  };
  await atlasRepository.addDecision(decision);
  const execution: ToolExecution = {
    id: `tool_${crypto.randomUUID()}`,
    organizationId: input.context.organization.id,
    workspaceId: input.context.workspace.id,
    actorId: input.context.user.id,
    agentReleaseId: release.id,
    toolId: tool.id,
    purpose: input.purpose,
    status: allowed ? "EXECUTED" : awaitingApproval ? "AWAITING_APPROVAL" : "BLOCKED",
    trustDecisionId: decision.id,
    createdAt: new Date().toISOString(),
  };
  await atlasRepository.addToolExecution(execution);
  return { execution, decision, tool, release };
}
