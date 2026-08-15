import { NextRequest } from "next/server";
import { ok, fail } from "@/lib/api";
import { getAtlasContext } from "@/lib/context";
import { atlasRepository } from "@/lib/repository";
import { approveRelease, registerCandidate, revokeRelease } from "@/lib/agent-release-registry";

export async function GET() {
  const context = await getAtlasContext();
  return ok(await atlasRepository.agentReleases(context.organization.id, context.workspace.id));
}

export async function POST(request: NextRequest) {
  const context = await getAtlasContext();
  try {
    const body = await request.json();
    if (body.operation === "approve") return ok(await approveRelease(String(body.id), context));
    if (body.operation === "revoke") return ok(await revokeRelease(String(body.id), context));
    return ok(await registerCandidate({
      context,
      agentId: String(body.agentId ?? "agent_custom"),
      agentName: String(body.agentName ?? "Custom Agent"),
      version: String(body.version ?? "0.1.0"),
      allowedTools: Array.isArray(body.allowedTools) ? body.allowedTools.map(String) : ["knowledge.search"],
      memoryScopes: Array.isArray(body.memoryScopes) ? body.memoryScopes : ["Internal"],
      policySuite: String(body.policySuite ?? "ATLAS-CORE-1"),
      evaluationScore: Number(body.evaluationScore ?? 90),
      rollbackReleaseId: body.rollbackReleaseId ? String(body.rollbackReleaseId) : undefined,
    }));
  } catch (error) {
    return fail("RELEASE_OPERATION_FAILED", error instanceof Error ? error.message : "Release operation failed", 400);
  }
}
