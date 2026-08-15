import { NextRequest } from "next/server";
import { ok, fail } from "@/lib/api";
import { getAtlasContext } from "@/lib/context";
import { executeTool } from "@/lib/tool-gateway";

export async function POST(request: NextRequest) {
  try {
    const context = await getAtlasContext();
    const body = await request.json();
    if (!body.agentReleaseId || !body.toolId || !body.purpose) return fail("INVALID_TOOL_REQUEST", "agentReleaseId, toolId, and purpose are required.", 400);
    return ok(await executeTool({ context, agentReleaseId: String(body.agentReleaseId), toolId: String(body.toolId), purpose: String(body.purpose).slice(0, 500) }));
  } catch (error) {
    return fail("TOOL_EXECUTION_FAILED", error instanceof Error ? error.message : "Tool execution failed", 400);
  }
}
