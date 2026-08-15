import { promises as fs } from "node:fs";
import path from "node:path";
import type { AgentRelease, AtlasDeal, AtlasDocument, AtlasUser, Membership, Organization, ToolExecution, TrustDecision, Workspace } from "./types";

type AtlasState = {
  organizations: Organization[];
  workspaces: Workspace[];
  users: AtlasUser[];
  memberships: Membership[];
  documents: AtlasDocument[];
  decisions: TrustDecision[];
  agentReleases: AgentRelease[];
  toolExecutions: ToolExecution[];
  deals: AtlasDeal[];
};

const now = Date.now();
const seed: AtlasState = {
  organizations: [{ id: "org_jlst", name: "JLST Holdings", slug: "jlst-holdings", plan: "Enterprise", createdAt: new Date(now - 86400000 * 30).toISOString() }],
  workspaces: [
    { id: "ws_executive", organizationId: "org_jlst", name: "Executive", slug: "executive", createdAt: new Date(now - 86400000 * 29).toISOString() },
    { id: "ws_atlas", organizationId: "org_jlst", name: "Atlas Product", slug: "atlas-product", createdAt: new Date(now - 86400000 * 20).toISOString() },
  ],
  users: [
    { id: "usr_jleon", name: "J’Leon Love", email: "founder@jlstholdings.com", createdAt: new Date(now - 86400000 * 30).toISOString() },
    { id: "usr_ops", name: "Atlas Operations", email: "operations@jlstholdings.com", createdAt: new Date(now - 86400000 * 14).toISOString() },
  ],
  memberships: [
    { id: "mem_jleon", organizationId: "org_jlst", userId: "usr_jleon", role: "owner", createdAt: new Date(now - 86400000 * 30).toISOString() },
    { id: "mem_ops", organizationId: "org_jlst", userId: "usr_ops", role: "manager", createdAt: new Date(now - 86400000 * 14).toISOString() },
  ],
  documents: [
    { id: "doc_policy_handbook", organizationId: "org_jlst", workspaceId: "ws_executive", createdBy: "usr_jleon", name: "Policy-Handbook.txt", classification: "Internal", content: "Atlas requires human approval for high-impact financial, legal, security, and personnel actions. Every governed response must preserve source evidence and an audit record.", createdAt: new Date(now - 1000 * 60 * 35).toISOString(), chunks: 2 },
    { id: "doc_genesis", organizationId: "org_jlst", workspaceId: "ws_atlas", createdBy: "usr_jleon", name: "Atlas-Genesis.txt", classification: "Confidential", content: "Atlas is an enterprise intelligence operating system built around Memory Fabric, Trust Fabric, Agent OS, governed workflows, and measurable outcomes.", createdAt: new Date(now - 1000 * 60 * 18).toISOString(), chunks: 2 },
  ],
  agentReleases: [
    { id: "rel_exec_100", organizationId: "org_jlst", workspaceId: "ws_executive", agentId: "agent_executive", agentName: "Executive Agent", version: "1.0.0", digest: "sha256:8be02d9f", status: "APPROVED", allowedTools: ["knowledge.search", "report.generate"], memoryScopes: ["Public", "Internal", "Confidential"], policySuite: "ATLAS-CORE-1", evaluationScore: 98, approvedBy: "usr_jleon", createdAt: new Date(now - 86400000 * 5).toISOString(), approvedAt: new Date(now - 86400000 * 4).toISOString() },
    { id: "rel_ops_110", organizationId: "org_jlst", workspaceId: "ws_atlas", agentId: "agent_operations", agentName: "Operations Agent", version: "1.1.0-rc.1", digest: "sha256:4fa9d102", status: "CANDIDATE", allowedTools: ["knowledge.search", "workflow.create"], memoryScopes: ["Public", "Internal"], policySuite: "ATLAS-CORE-1", evaluationScore: 91, rollbackReleaseId: "rel_ops_100", createdAt: new Date(now - 86400000).toISOString() }
  ],
  toolExecutions: [],
  deals: [],
  decisions: [
    { id: "trust_seed_1", organizationId: "org_jlst", workspaceId: "ws_executive", actorId: "usr_jleon", actorRole: "owner", decision: "ALLOW", action: "Executive Agent generated board briefing", policyId: "POL-EXEC-04", riskScore: 0.02, evidence: ["Owner role verified", "Approved internal sources only"], createdAt: new Date(now - 1000 * 60 * 12).toISOString() },
    { id: "trust_seed_2", organizationId: "org_jlst", workspaceId: "ws_atlas", actorId: "usr_ops", actorRole: "manager", decision: "DENY", action: "External connector requested restricted file", policyId: "POL-DATA-12", riskScore: 0.87, evidence: ["External destination", "Restricted classification"], createdAt: new Date(now - 1000 * 60 * 31).toISOString() },
  ],
};

const file = path.join(process.cwd(), "data", "atlas-state.json");
let cache: AtlasState | null = null;
let writing = Promise.resolve();

function migrate(raw: Partial<AtlasState>): AtlasState {
  const state = { ...structuredClone(seed), ...raw } as AtlasState;
  state.documents = (state.documents ?? []).map((d) => ({ organizationId: "org_jlst", workspaceId: "ws_atlas", createdBy: "usr_jleon", ...d }));
  state.decisions = (state.decisions ?? []).map((d) => ({ organizationId: "org_jlst", workspaceId: "ws_atlas", actorId: "usr_jleon", actorRole: "owner", ...d }));
  state.agentReleases = state.agentReleases ?? structuredClone(seed.agentReleases);
  state.toolExecutions = state.toolExecutions ?? [];
  state.deals = state.deals ?? [];
  return state;
}

async function load(): Promise<AtlasState> {
  if (cache) return cache;
  try { cache = migrate(JSON.parse(await fs.readFile(file, "utf8")) as AtlasState); }
  catch { cache = structuredClone(seed); await persist(cache); }
  return cache;
}

async function persist(state: AtlasState) {
  await fs.mkdir(path.dirname(file), { recursive: true });
  const temp = `${file}.tmp`;
  writing = writing.then(async () => { await fs.writeFile(temp, JSON.stringify(state, null, 2), "utf8"); await fs.rename(temp, file); });
  await writing;
}

export const atlasRepository = {
  async organizations() { return (await load()).organizations; },
  async workspaces(organizationId?: string) { const all = (await load()).workspaces; return organizationId ? all.filter((w) => w.organizationId === organizationId) : all; },
  async users() { return (await load()).users; },
  async memberships(organizationId?: string) { const all = (await load()).memberships; return organizationId ? all.filter((m) => m.organizationId === organizationId) : all; },
  async organization(id: string) { return (await load()).organizations.find((o) => o.id === id); },
  async workspace(id: string) { return (await load()).workspaces.find((w) => w.id === id); },
  async user(id: string) { return (await load()).users.find((u) => u.id === id); },
  async documents(organizationId?: string, workspaceId?: string) { return (await load()).documents.filter((d) => (!organizationId || d.organizationId === organizationId) && (!workspaceId || d.workspaceId === workspaceId)); },
  async decisions(organizationId?: string, workspaceId?: string) { return (await load()).decisions.filter((d) => (!organizationId || d.organizationId === organizationId) && (!workspaceId || d.workspaceId === workspaceId)); },
  async agentReleases(organizationId?: string, workspaceId?: string) { return (await load()).agentReleases.filter((r) => (!organizationId || r.organizationId === organizationId) && (!workspaceId || r.workspaceId === workspaceId)); },
  async agentRelease(id: string) { return (await load()).agentReleases.find((r) => r.id === id); },
  async toolExecutions(organizationId?: string, workspaceId?: string) { return (await load()).toolExecutions.filter((e) => (!organizationId || e.organizationId === organizationId) && (!workspaceId || e.workspaceId === workspaceId)); },
  async deals(organizationId?: string, workspaceId?: string) { return (await load()).deals.filter((d) => (!organizationId || d.organizationId === organizationId) && (!workspaceId || d.workspaceId === workspaceId)); },
  async deal(id: string) { return (await load()).deals.find((d) => d.id === id); },
  async addDeal(deal: AtlasDeal) { const state = await load(); state.deals.unshift(deal); await persist(state); return deal; },
  async updateDeal(id: string, patch: Partial<AtlasDeal>) { const state = await load(); const index = state.deals.findIndex((d) => d.id === id); if (index < 0) return undefined; state.deals[index] = { ...state.deals[index], ...patch, updatedAt: new Date().toISOString() }; await persist(state); return state.deals[index]; },
  async addDocument(document: AtlasDocument) { const state = await load(); state.documents.unshift(document); await persist(state); return document; },
  async addDecision(decision: TrustDecision) { const state = await load(); state.decisions.unshift(decision); await persist(state); return decision; },
  async addAgentRelease(release: AgentRelease) { const state = await load(); state.agentReleases.unshift(release); await persist(state); return release; },
  async updateAgentRelease(id: string, patch: Partial<AgentRelease>) { const state = await load(); const index = state.agentReleases.findIndex((r) => r.id === id); if (index < 0) return undefined; state.agentReleases[index] = { ...state.agentReleases[index], ...patch }; await persist(state); return state.agentReleases[index]; },
  async addToolExecution(execution: ToolExecution) { const state = await load(); state.toolExecutions.unshift(execution); await persist(state); return execution; },
  async addWorkspace(workspace: Workspace) { const state = await load(); state.workspaces.push(workspace); await persist(state); return workspace; },
  async addUser(user: AtlasUser, membership: Membership) { const state = await load(); state.users.push(user); state.memberships.push(membership); await persist(state); return user; },
  async overview(organizationId?: string, workspaceId?: string) {
    const documents = await this.documents(organizationId, workspaceId);
    const decisions = await this.decisions(organizationId, workspaceId);
    const allows = decisions.filter((item) => item.decision === "ALLOW").length;
    return { trustScore: Number(((allows / Math.max(1, decisions.length)) * 100).toFixed(1)), knowledgeSources: documents.length, trustEvents: decisions.length, indexedChunks: documents.reduce((sum, item) => sum + item.chunks, 0) };
  },
};
