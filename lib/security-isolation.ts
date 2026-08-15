export interface TenantResource { organizationId:string; workspaceId:string; }
export interface TenantContext { organizationId:string; workspaceId:string; }
export function assertTenantIsolation(ctx:TenantContext, resource:TenantResource){
 if(resource.organizationId!==ctx.organizationId || resource.workspaceId!==ctx.workspaceId) throw new Error('TENANT_ISOLATION_DENY');
 return true;
}
export function filterTenantResources<T extends TenantResource>(ctx:TenantContext, resources:T[]):T[]{return resources.filter(r=>r.organizationId===ctx.organizationId&&r.workspaceId===ctx.workspaceId);}
