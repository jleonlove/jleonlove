BEGIN;
CREATE SCHEMA IF NOT EXISTS atlas;
CREATE TABLE IF NOT EXISTS atlas.evidence_documents (
 id text PRIMARY KEY,
 tenant_id text NOT NULL,
 workspace_id text NOT NULL,
 source_uri text NOT NULL,
 source_version text NOT NULL,
 content text NOT NULL,
 content_sha256 text NOT NULL CHECK (content_sha256 ~ '^[a-f0-9]{64}$'),
 observed_at timestamptz NOT NULL,
 expires_at timestamptz NOT NULL,
 verified boolean NOT NULL DEFAULT false,
 CHECK (expires_at > observed_at)
);
CREATE INDEX IF NOT EXISTS evidence_scope_idx ON atlas.evidence_documents(tenant_id, workspace_id);
ALTER TABLE atlas.evidence_documents ENABLE ROW LEVEL SECURITY;
ALTER TABLE atlas.evidence_documents FORCE ROW LEVEL SECURITY;
CREATE POLICY evidence_read_scope ON atlas.evidence_documents FOR SELECT
USING (
 tenant_id = nullif(current_setting('atlas.tenant_id', true), '')
 AND workspace_id = nullif(current_setting('atlas.workspace_id', true), '')
);
COMMIT;
