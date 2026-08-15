import { AppShell } from "@/components/AppShell";
import { MemoryWorkspace } from "@/components/memory/MemoryWorkspace";

export default function MemoryPage() {
  return (
    <AppShell>
      <p className="eyebrow">Memory Fabric</p>
      <h1>Governed enterprise knowledge</h1>
      <p className="muted">Ingest knowledge, search it, evaluate access, and preserve evidence in one working flow.</p>
      <MemoryWorkspace />
    </AppShell>
  );
}
