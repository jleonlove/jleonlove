export interface EconomicInput { quantity:number; unitPrice:number; acquisitionCost?:number; logistics?:number; insurance?:number; inspection?:number; financing?:number; commissions?:number; taxes?:number; other?:number; }
export interface EconomicPlausibility { valid:boolean; grossValue:number; totalCosts:number; expectedMargin:number; marginPct:number|null; findings:string[]; }
export function evaluateEconomics(x:EconomicInput):EconomicPlausibility {
 const findings:string[]=[]; const nums=[x.quantity,x.unitPrice,x.acquisitionCost??0,x.logistics??0,x.insurance??0,x.inspection??0,x.financing??0,x.commissions??0,x.taxes??0,x.other??0];
 if(nums.some(n=>!Number.isFinite(n))) findings.push('All economic inputs must be finite.');
 if(x.quantity<=0) findings.push('Quantity must be greater than zero.'); if(x.unitPrice<0) findings.push('Unit price cannot be negative.');
 if(nums.slice(2).some(n=>n<0)) findings.push('Transaction costs cannot be negative.');
 const grossValue=x.quantity*x.unitPrice,totalCosts=(x.acquisitionCost??0)+(x.logistics??0)+(x.insurance??0)+(x.inspection??0)+(x.financing??0)+(x.commissions??0)+(x.taxes??0)+(x.other??0),expectedMargin=grossValue-totalCosts;
 if(Number.isFinite(grossValue)&&grossValue>0&&totalCosts>grossValue) findings.push('Modeled transaction costs exceed gross commodity value.');
 return {valid:findings.length===0,grossValue,totalCosts,expectedMargin,marginPct:grossValue>0?expectedMargin/grossValue*100:null,findings};
}
