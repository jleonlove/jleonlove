"use client";
import { useState } from "react";
import type { AgentRelease, ToolDefinition, ToolExecution } from "@/lib/types";

type ApiEnvelope<T> = { success: boolean; data?: T; error?: { message: string } };

export function ReleaseConsole({ releases, tools, executions }: { releases: AgentRelease[]; tools: ToolDefinition[]; executions: ToolExecution[] }) {
  const [message, setMessage] = useState("");
  const [busy, setBusy] = useState(false);
  const approved = releases.find((item) => item.status === "APPROVED");

  async function call(url: string, body: unknown) {
    setBusy(true); setMessage("");
    try {
      const response = await fetch(url, { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify(body) });
      const json = await response.json() as ApiEnvelope<unknown>;
      if (!response.ok || !json.success) throw new Error(json.error?.message ?? "Operation failed");
      setMessage("Operation recorded. Refresh to view the immutable release and audit state.");
    } catch (error) { setMessage(error instanceof Error ? error.message : "Operation failed"); }
    finally { setBusy(false); }
  }

  return <div className="stack">
    <section className="grid">
      {releases.map((release) => <article className="card span2" key={release.id}>
        <div className="row"><div><span className={`badge ${release.status === "APPROVED" ? "good" : release.status === "REVOKED" ? "danger" : "warn"}`}>{release.status}</span><h2>{release.agentName} <small>v{release.version}</small></h2></div><strong>{release.evaluationScore}/100</strong></div>
        <p className="muted">{release.digest}</p>
        <p><strong>Tools:</strong> {release.allowedTools.join(", ") || "None"}</p>
        <p><strong>Memory:</strong> {release.memoryScopes.join(", ")}</p>
        <p><strong>Policy suite:</strong> {release.policySuite}</p>
        {release.status === "CANDIDATE" && <div className="actions"><button disabled={busy} className="cta" onClick={() => call("/api/releases", { operation: "approve", id: release.id })}>Approve release</button><button disabled={busy} className="ghost" onClick={() => call("/api/releases", { operation: "revoke", id: release.id })}>Revoke</button></div>}
      </article>)}
    </section>

    <section className="card">
      <p className="eyebrow">Governed Tool Gateway</p><h2>Test a declared capability</h2>
      <p className="muted">Every request checks release status, signed manifest, tool risk, and human approval requirements before execution.</p>
      <div className="tool-grid">{tools.map((tool) => <button className="tool-button" disabled={busy || !approved} key={tool.id} onClick={() => approved && call("/api/tools/execute", { agentReleaseId: approved.id, toolId: tool.id, purpose: `Pilot validation of ${tool.name}` })}><strong>{tool.name}</strong><small>{tool.risk} · {tool.requiresApproval ? "Approval required" : "Auto policy"}</small></button>)}</div>
      {message && <p className="notice">{message}</p>}
    </section>

    <section className="card"><p className="eyebrow">Execution Ledger</p><h2>Recent tool requests</h2>{executions.length ? executions.slice(0, 8).map((item) => <div className="row" key={item.id}><div><strong>{item.toolId}</strong><small>{item.purpose}</small></div><span className={`badge ${item.status === "EXECUTED" ? "good" : item.status === "BLOCKED" ? "danger" : "warn"}`}>{item.status}</span></div>) : <p className="muted">No tool requests recorded yet.</p>}</section>
  </div>;
}
