export type DocumentEvidenceInput = {
  documentId: string;
  declaredHash?: string;
  computedHash?: string;
  signerAuthority?: "VERIFIED" | "UNVERIFIED" | "REVOKED" | "UNKNOWN";
  issuedAt?: string;
  expiresAt?: string;
  extractedText?: string;
  previousVersionHash?: string;
  supersedesHash?: string;
};

export type DocumentEvidenceFinding = {
  code: string;
  severity: "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  message: string;
  failClosed: boolean;
};

const injectionPatterns = [
  /ignore\s+(all|any|the)\s+(previous|prior|system)\s+instructions/i,
  /reveal\s+(the\s+)?(system|developer)\s+prompt/i,
  /bypass\s+(kyc|aml|compliance|verification|approval)/i,
  /send|transfer|release\s+funds\s+without\s+(approval|verification)/i,
];

export function inspectDocumentEvidence(input: DocumentEvidenceInput, now = new Date()): DocumentEvidenceFinding[] {
  const out: DocumentEvidenceFinding[] = [];
  if (!input.declaredHash || !input.computedHash) out.push({ code:"DOC_HASH_MISSING", severity:"HIGH", message:"Document integrity hash is incomplete.", failClosed:true });
  else if (input.declaredHash !== input.computedHash) out.push({ code:"DOC_HASH_MISMATCH", severity:"CRITICAL", message:"Document bytes do not match the declared integrity hash.", failClosed:true });
  if (input.signerAuthority === "REVOKED") out.push({ code:"SIGNER_AUTHORITY_REVOKED", severity:"CRITICAL", message:"Document signer authority is revoked.", failClosed:true });
  else if (!input.signerAuthority || ["UNVERIFIED","UNKNOWN"].includes(input.signerAuthority)) out.push({ code:"SIGNER_AUTHORITY_UNVERIFIED", severity:"HIGH", message:"Document signer authority is not independently verified.", failClosed:true });
  if (input.expiresAt && new Date(input.expiresAt).getTime() <= now.getTime()) out.push({ code:"DOCUMENT_EXPIRED", severity:"HIGH", message:"Document is expired and requires revalidation.", failClosed:true });
  if (input.previousVersionHash && input.computedHash && input.previousVersionHash !== input.computedHash && input.supersedesHash !== input.previousVersionHash) out.push({ code:"UNEXPLAINED_VERSION_CHANGE", severity:"HIGH", message:"Document changed without an explicit supersession chain.", failClosed:true });
  const text = input.extractedText ?? "";
  if (injectionPatterns.some((p) => p.test(text))) out.push({ code:"UNTRUSTED_DOCUMENT_INSTRUCTION", severity:"CRITICAL", message:"Untrusted document content attempts to influence Atlas execution or bypass controls.", failClosed:true });
  return out;
}

export function documentEvidenceMayDriveExecution(input: DocumentEvidenceInput, now = new Date()) {
  return !inspectDocumentEvidence(input, now).some((f) => f.failClosed);
}
