import { ok } from "@/lib/api";
import { getAtlasContext } from "@/lib/context";
import { atlasRepository } from "@/lib/repository";
export async function GET(request:Request){const context=await getAtlasContext();return ok({decisions:(await atlasRepository.decisions(context.organization.id,context.workspace.id)).slice(0,50)},request)}
