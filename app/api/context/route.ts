import { fail, ok } from "@/lib/api";
import { atlasRepository } from "@/lib/repository";
export async function POST(request: Request) {
  const body = await request.json() as { organizationId?: string; workspaceId?: string; userId?: string };
  const organization = body.organizationId ? await atlasRepository.organization(body.organizationId) : undefined;
  const workspace = body.workspaceId ? await atlasRepository.workspace(body.workspaceId) : undefined;
  const user = body.userId ? await atlasRepository.user(body.userId) : undefined;
  if (!organization || !workspace || workspace.organizationId !== organization.id || !user) return fail("INVALID_CONTEXT", "The selected Atlas context is invalid.", 400, request);
  const response = ok({ organizationId: organization.id, workspaceId: workspace.id, userId: user.id }, request);
  response.cookies.set("atlas_org", organization.id, { httpOnly: true, sameSite: "lax", path: "/" });
  response.cookies.set("atlas_workspace", workspace.id, { httpOnly: true, sameSite: "lax", path: "/" });
  response.cookies.set("atlas_user", user.id, { httpOnly: true, sameSite: "lax", path: "/" });
  return response;
}
