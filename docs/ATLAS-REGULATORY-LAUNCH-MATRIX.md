# Atlas Regulatory Launch Matrix

Status: Production governance control
Owner: JLST Holdings / Atlas

> This matrix is a product-control framework, not legal advice. Final classifications and required registrations/licenses must be confirmed by qualified counsel for the actual facts, jurisdictions, counterparties, commodities, and transaction structures.

## Default operating boundary

Atlas V1 is a technology, intelligence, verification, compliance-support, document, and workflow-facilitation platform for physical commodity transactions. Atlas must not silently cross into a regulated activity merely because a model or tool is technically capable of doing so.

Atlas/JLST does not custody customer trade funds by default. Customer settlement occurs through the customer's own bank or an appropriately regulated payment, escrow, exchange, custodian, or financial institution.

## Capability classifications

| Capability | Default status | Runtime rule |
|---|---|---|
| Physical commodity intelligence/research | CLEAR WITH CONTROLS | Permit with provenance, uncertainty, and source controls. |
| Trade-document ingestion/extraction | CLEAR WITH CONTROLS | Permit with privacy, authorization, retention, and provenance controls. |
| Counterparty/KYB/KYC support | CLEAR WITH CONTROLS | Screening is decision support; preserve evidence and escalation. |
| Sanctions/restricted-party screening | CLEAR WITH CONTROLS | Require current authoritative data, audit evidence, and human escalation for ambiguity. |
| Export-control screening/classification support | COUNSEL/COMPLIANCE REVIEW | Do not represent Atlas output as a government license or final legal determination. |
| Physical spot trade workflow facilitation | CLEAR WITH CONTROLS | Keep Atlas on the software/facilitation side unless counsel approves expanded activity. |
| Contract/document workflow | CLEAR WITH CONTROLS | Require authorized principals and independent legal review where appropriate. |
| Customer-to-customer custody of fiat/crypto | PROHIBITED UNTIL CLEARED/LICENSED | Fail closed. Do not accept/transmit customer value for another person absent approved regulatory classification and required licensing. |
| Crypto accepted as payment for JLST's own services | COUNSEL/TAX REVIEW | Keep distinct from transmitting customer value; use approved payment/custody providers. |
| Futures/options/swaps order solicitation or acceptance | PROHIBITED UNTIL CLEARED/REGISTERED | Fail closed unless the relevant CFTC/NFA status or exemption is verified. |
| Personalized compensated advice on regulated commodity interests | PROHIBITED UNTIL CLEARED/REGISTERED | Fail closed pending CTA/exemption analysis. |
| Discretionary trading/customer account control | PROHIBITED UNTIL CLEARED/REGISTERED | No autonomous activation. |
| Autonomous high-value settlement | PROHIBITED BY DEFAULT | Require scoped economic authority, independent approval where designated, regulated rails, and settlement evidence. |
| Export/reexport/transfer requiring government authorization | LICENSE DEPENDENT | Block execution until required license/exception determination and evidence are present. |

## Regulatory activation gate

Before enabling a regulated or potentially regulated capability, Atlas must bind:

1. accountable legal entity;
2. jurisdiction(s);
3. commodity/product and transaction type;
4. customer/counterparty classification;
5. whether Atlas/JLST solicits, advises, accepts orders, exercises discretion, receives compensation tied to regulated activity, or handles value;
6. required registration, license, exemption, or no-action/legal basis;
7. counsel/compliance approval and expiration/review date;
8. applicable KYC/KYB, sanctions, AML, export-control, privacy, recordkeeping, and disclosure controls;
9. permitted settlement/custody providers;
10. signed evidence attached to the release/capability registry.

Missing or expired evidence => DENY / REQUIRE APPROVAL / QUARANTINE. Never infer regulatory permission from technical capability.

## Production invariants

- No customer-fund custody or transmission merely because Atlas can initiate a payment API.
- No derivatives solicitation/order routing/advice merely because Atlas understands the market.
- No claim that screening equals government approval.
- No export-controlled execution without the applicable determination and evidence.
- No regulated capability can self-enable or expand its own authority.
- Jurisdiction, regulatory status, licenses/exemptions, and sanctions/export evidence are freshness-bound and revalidated before consequential execution.
- Every consequential regulatory decision is reconstructable from principal -> intent -> classification -> authority -> evidence -> approval -> effect.

## Launch evidence checklist

- U.S. federal activity classification memorandum from qualified counsel.
- State-law review for actual operating/customer activities, including money-transmission analysis if Atlas ever handles third-party value.
- CFTC/NFA analysis covering CTA, IB, FCM/CPO/AP or other relevant categories and exemptions as applicable.
- FinCEN/BSA analysis if any third-party fiat/crypto acceptance/transmission is contemplated.
- OFAC sanctions compliance policy and escalation procedure.
- BIS/EAR export compliance program where applicable to Atlas/JLST activities.
- Privacy/data-processing terms and jurisdiction-specific privacy review.
- Customer Terms of Service, acceptable-use boundaries, disclosures, and limitation-of-authority language.
- Cyber/E&O insurance review.
- Security assessment, penetration testing, incident response, disaster recovery, and evidence retention controls.

## Product messaging boundary

Approved positioning: Atlas provides commodity intelligence, verification, compliance-support tooling, workflow orchestration, document intelligence, and governed transaction facilitation.

Do not market Atlas as a licensed broker, money transmitter, custodian, exchange, investment adviser/commodity trading adviser, government compliance authority, or guaranteed legal-compliance service unless the relevant status is actually obtained and verified.
