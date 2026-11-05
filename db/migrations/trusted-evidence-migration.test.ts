import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const sql = readFileSync("db/migrations/001_trusted_evidence.sql", "utf8");

describe("proposed evidence migration guardrails", () => {
  it("requires tenant and workspace identifiers", () => {
    expect(sql).toContain("tenant_id text NOT NULL");
    expect(sql).toContain("workspace_id text NOT NULL");
  });
  it("enables and forces row-level security", () => {
    expect(sql).toContain("ENABLE ROW LEVEL SECURITY");
    expect(sql).toContain("FORCE ROW LEVEL SECURITY");
  });
  it("requires both session scopes for reads", () => {
    expect(sql).toContain("current_setting('atlas.tenant_id', true)");
    expect(sql).toContain("current_setting('atlas.workspace_id', true)");
  });
  it("does not define permissive write policies", () => {
    expect(sql).not.toMatch(/CREATE POLICY[^;]+FOR (ALL|INSERT|UPDATE|DELETE)/i);
  });
});
