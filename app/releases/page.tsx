import { AppShell } from "@/components/AppShell";
import { ReleaseConsole } from "@/components/ReleaseConsole";
import { getAtlasContext } from "@/lib/context";
import { atlasRepository } from "@/lib/repository";
import { toolCatalog } from "@/lib/tool-gateway";

export default async function ReleasesPage() {
  const context = await getAtlasContext();
  const [releases, executions] = await Promise.all([
    atlasRepository.agentReleases(context.organization.id, context.workspace.id),
    atlasRepository.toolExecutions(context.organization.id, context.workspace.id),
  ]);
  return <AppShell><p className="eyebrow">Agent Release Registry</p><h1>Signed releases. Governed tools.</h1><p className="muted hero-copy">Atlas treats every agent as immutable software. Only evaluated, approved releases can invoke declared tools through the Trust Fabric.</p><ReleaseConsole releases={releases} tools={toolCatalog} executions={executions}/></AppShell>;
}
