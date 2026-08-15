import { createHmac, timingSafeEqual } from "crypto";

export type ProviderObservation = {
  provider: string;
  chain: string;
  txHash: string;
  destination: string;
  amountAtomic: bigint;
  confirmations: number;
  finalized: boolean;
  observedAt: string;
};

export function verifyWebhookSignature(rawBody: string, signatureHex: string, secret: string) {
  if (!signatureHex || !secret) return false;
  const expected = createHmac("sha256", secret).update(rawBody).digest();
  let supplied: Buffer;
  try { supplied = Buffer.from(signatureHex, "hex"); } catch { return false; }
  return supplied.length === expected.length && timingSafeEqual(supplied, expected);
}

export function observationKey(o: ProviderObservation) {
  return `${o.chain}:${o.txHash}:${o.destination}:${o.amountAtomic.toString()}`;
}

export function reconcileProviderObservations(observations: ProviderObservation[], minIndependentProviders = 2) {
  if (!observations.length) return { settled: false, reason: "NO_OBSERVATIONS" as const };
  const keys = new Map<string, Set<string>>();
  for (const o of observations) {
    if (!o.finalized) continue;
    const k = observationKey(o);
    if (!keys.has(k)) keys.set(k, new Set());
    keys.get(k)!.add(o.provider);
  }
  const ranked = [...keys.entries()].sort((a,b)=>b[1].size-a[1].size);
  if (!ranked.length) return { settled:false, reason:"NOT_FINAL" as const };
  if (ranked.length > 1) return { settled:false, reason:"PROVIDER_DISAGREEMENT" as const, variants: ranked.map(([key,p])=>({key,providers:[...p]})) };
  const [key, providers] = ranked[0];
  if (providers.size < minIndependentProviders) return { settled:false, reason:"INSUFFICIENT_INDEPENDENT_PROVIDERS" as const, providers:[...providers] };
  return { settled:true, reason:"INDEPENDENTLY_CONFIRMED" as const, key, providers:[...providers] };
}

export function detectReorg(previous: ProviderObservation, current: ProviderObservation) {
  return previous.chain === current.chain && previous.txHash === current.txHash && previous.finalized && !current.finalized;
}
