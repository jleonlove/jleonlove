"use client";
import { useState } from "react";
import type { AtlasRequestContext, Workspace } from "@/lib/types";

export function ContextSwitcher({ context, workspaces }: { context: AtlasRequestContext; workspaces: Workspace[] }) {
  const [pending, setPending] = useState(false);
  async function switchWorkspace(workspaceId: string) {
    setPending(true);
    await fetch("/api/context", { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify({ organizationId: context.organization.id, workspaceId, userId: context.user.id }) });
    location.reload();
  }
  return <div className="context-switcher"><div><b>{context.organization.name}</b><small>{context.user.name} · {context.role}</small></div><select aria-label="Active workspace" value={context.workspace.id} disabled={pending} onChange={(event) => switchWorkspace(event.target.value)}>{workspaces.map((workspace) => <option key={workspace.id} value={workspace.id}>{workspace.name}</option>)}</select></div>;
}
