import { cookies } from "next/headers";
import { atlasRepository } from "./repository";
import type { AtlasRequestContext } from "./types";

export async function getAtlasContext(): Promise<AtlasRequestContext> {
  const jar = await cookies();
  const organizations = await atlasRepository.organizations();
  const organization = await atlasRepository.organization(jar.get("atlas_org")?.value ?? organizations[0]?.id) ?? organizations[0];
  if (!organization) throw new Error("Atlas requires at least one organization.");
  const workspaces = await atlasRepository.workspaces(organization.id);
  const workspace = await atlasRepository.workspace(jar.get("atlas_workspace")?.value ?? workspaces[0]?.id) ?? workspaces[0];
  if (!workspace) throw new Error("Atlas requires at least one workspace.");
  const memberships = await atlasRepository.memberships(organization.id);
  const membership = memberships.find((m) => m.userId === jar.get("atlas_user")?.value) ?? memberships[0];
  if (!membership) throw new Error("Atlas requires at least one organization member.");
  const user = await atlasRepository.user(membership.userId);
  if (!user) throw new Error("Atlas membership user was not found.");
  return { organization, workspace, user, role: membership.role };
}
