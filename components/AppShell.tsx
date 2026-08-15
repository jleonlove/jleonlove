import Link from "next/link";
import { getAtlasContext } from "@/lib/context";
import { atlasRepository } from "@/lib/repository";
import { ContextSwitcher } from "./ContextSwitcher";
const links=[["/dashboard","Overview"],["/memory","Memory"],["/trust","Trust"],["/agents","Agents"],["/releases","Releases"],["/organization","Organization"]];
export async function AppShell({children}:{children:React.ReactNode}){const context=await getAtlasContext();const workspaces=await atlasRepository.workspaces(context.organization.id);return <div className="shell"><aside className="sidebar"><Link href="/" className="brand">ATLAS</Link><p className="muted">Genesis v0.2</p><ContextSwitcher context={context} workspaces={workspaces}/><nav className="nav">{links.map(([href,label])=><Link key={href} href={href}>{label}</Link>)}</nav></aside><main className="main">{children}</main></div>}
