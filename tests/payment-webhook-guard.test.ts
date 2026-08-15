import {describe,it,expect} from "vitest";
import {createHmac} from "crypto";
import {qualifyWebhook,commitWebhook} from "../lib/payment-webhook-guard";
describe("webhook guard",()=>{
 it("rejects replayed events",()=>{const secret="s",body="{}",now=new Date("2026-08-15T00:00:00Z"),seen=new Set<string>(); const e={provider:"p",eventId:"1",timestamp:now.toISOString(),rawBody:body,signatureHex:createHmac("sha256",secret).update(body).digest("hex")}; const q=qualifyWebhook(e,secret,seen,now); expect(q.accepted).toBe(true); commitWebhook(q.key,seen); expect(qualifyWebhook(e,secret,seen,now).accepted).toBe(false)});
});
