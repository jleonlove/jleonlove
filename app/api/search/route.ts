import { fail, ok } from "@/lib/api";
import { getAtlasContext } from "@/lib/context";
import { searchKnowledge } from "@/lib/search-engine";
export async function GET(request:Request){const context=await getAtlasContext();const query=new URL(request.url).searchParams.get("q")?.trim()??"";if(query.length<2)return fail("QUERY_TOO_SHORT","Query must contain at least 2 characters.",400,request);return ok({hits:await searchKnowledge(query,{organizationId:context.organization.id,workspaceId:context.workspace.id})},request)}
