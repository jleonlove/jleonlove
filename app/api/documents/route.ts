import { fail, ok } from "@/lib/api";
import { atlasRepository } from "@/lib/repository";
import { getAtlasContext } from "@/lib/context";
import type { AtlasDocument, Classification } from "@/lib/types";

export const runtime = "nodejs";
const allowed = new Set<Classification>(["Public", "Internal", "Confidential", "Restricted"]);

export async function GET(request: Request) {
  const context = await getAtlasContext();
  return ok({ documents: await atlasRepository.documents(context.organization.id, context.workspace.id) }, request);
}

export async function POST(request: Request) {
  const context = await getAtlasContext();
  const form = await request.formData();
  const file = form.get("document");
  const pasted = String(form.get("content") ?? "").trim();
  const raw = String(form.get("classification") ?? "Internal") as Classification;
  const classification = allowed.has(raw) ? raw : "Internal";
  if (!(file instanceof File) && !pasted) return fail("INVALID_INPUT", "Provide a file or text content.", 400, request);
  if (file instanceof File && file.size > 5_000_000) return fail("FILE_TOO_LARGE", "File exceeds the 5 MB demo limit.", 413, request);

  let content = pasted;
  let name = "Pasted-Knowledge.txt";
  if (file instanceof File && file.size > 0) {
    name = file.name;
    if (/^(text\/|application\/(json|xml))/.test(file.type) || /\.(txt|md|csv|json|xml)$/i.test(file.name)) content = await file.text();
    else if (!content) content = `Document metadata indexed for ${file.name}. Add extracted text for governed search.`;
  }
  const document: AtlasDocument = {
    id: `doc_${crypto.randomUUID()}`,
    organizationId: context.organization.id,
    workspaceId: context.workspace.id,
    createdBy: context.user.id,
    name,
    classification,
    content: content.slice(0, 150_000),
    createdAt: new Date().toISOString(),
    chunks: Math.max(1, Math.ceil(content.length / 700)),
  };
  await atlasRepository.addDocument(document);
  return ok({ document }, request, { status: 201 });
}
