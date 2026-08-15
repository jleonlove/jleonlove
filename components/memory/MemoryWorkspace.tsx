"use client";

import { FormEvent, useEffect, useState } from "react";
import type { AtlasDocument, SearchHit, TrustDecision } from "@/lib/types";

type AskPayload = { answer: string | null; sources: SearchHit[]; decision: TrustDecision };
type ApiResponse<T> = { success: boolean; data?: T; error?: { code: string; message: string }; requestId: string };

export function MemoryWorkspace() {
  const [documents, setDocuments] = useState<AtlasDocument[]>([]);
  const [question, setQuestion] = useState("What makes Atlas trustworthy?");
  const [result, setResult] = useState<AskPayload | null>(null);
  const [message, setMessage] = useState("");
  const [busy, setBusy] = useState(false);

  async function loadDocuments() {
    const response = await fetch("/api/documents", { cache: "no-store" });
    const payload = (await response.json()) as ApiResponse<{ documents: AtlasDocument[] }>;
    setDocuments(payload.data?.documents ?? []);
  }

  useEffect(() => {
    void loadDocuments();
  }, []);

  async function upload(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setBusy(true);
    setMessage("");
    const form = new FormData(event.currentTarget);
    const response = await fetch("/api/documents", { method: "POST", body: form });
    const payload = (await response.json()) as ApiResponse<{ document: AtlasDocument }>;
    setMessage(response.ok ? "Knowledge indexed successfully." : payload.error?.message ?? "Upload failed.");
    if (response.ok) {
      event.currentTarget.reset();
      await loadDocuments();
    }
    setBusy(false);
  }

  async function ask(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setBusy(true);
    const response = await fetch("/api/ask", {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ question, actorRole: "admin" }),
    });
    const payload = (await response.json()) as ApiResponse<AskPayload>;
    setResult(payload.data ?? null);
    if (!response.ok) setMessage(payload.error?.message ?? "Atlas could not complete the request.");
    setBusy(false);
  }

  return (
    <section className="grid">
      <article className="card span2">
        <h2>Ingest governed knowledge</h2>
        <form className="form" onSubmit={upload}>
          <input className="input" type="file" name="document" accept=".txt,.md,.csv,.json,.xml,.pdf,.doc,.docx" />
          <textarea className="input" name="content" rows={5} placeholder="Paste extracted text or notes for immediate search." />
          <select className="input" name="classification" defaultValue="Internal">
            <option>Public</option><option>Internal</option><option>Confidential</option><option>Restricted</option>
          </select>
          <button className="cta" type="submit" disabled={busy}>{busy ? "Working…" : "Index knowledge"}</button>
          {message ? <p className="muted" role="status">{message}</p> : null}
        </form>
      </article>

      <article className="card span2">
        <h2>Governed question</h2>
        <form className="form" onSubmit={ask}>
          <textarea className="input" rows={4} value={question} onChange={(event) => setQuestion(event.target.value)} />
          <button className="cta" type="submit" disabled={busy}>{busy ? "Evaluating…" : "Ask Atlas"}</button>
        </form>
        {result ? (
          <div className="answer">
            <span className={`badge ${result.decision.decision === "ALLOW" ? "good" : "warn"}`}>{result.decision.decision} · {result.decision.policyId}</span>
            <p>{result.answer ?? "Access denied by Trust Fabric."}</p>
            {result.sources.map((source) => <div className="source" key={source.documentId}><b>{source.documentName}</b><span className="muted">{source.classification} · relevance {Math.round(source.score * 100)}%</span></div>)}
          </div>
        ) : null}
      </article>

      <article className="card span2">
        <h2>Memory index</h2>
        <div className="stack">
          {documents.map((document) => <div className="source" key={document.id}><b>{document.name}</b><span className="muted">{document.classification} · {document.chunks} chunks</span></div>)}
        </div>
      </article>
    </section>
  );
}
