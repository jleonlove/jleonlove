from pathlib import Path
import sys
root=Path(__file__).resolve().parents[1]
checks={'financial_exceptions':root/'lib/financial-exceptions.ts','loss_prevention':root/'lib/financial-loss-prevention.ts','loss_tests':root/'tests/financial-loss-prevention.test.ts','package':root/'package.json'}
missing=[k for k,p in checks.items() if not p.exists()]
text=checks['loss_prevention'].read_text() if not missing else ''
required=['PARTIAL','OVERPAYMENT_REVIEW','OVER_ALLOCATION','HOLD_PROVIDER_DISAGREEMENT','performanceSatisfied']
missing_tokens=[x for x in required if x not in text]
if missing or missing_tokens: print('OFFLINE_QUALIFICATION=FAIL',missing,missing_tokens);sys.exit(1)
print('OFFLINE_QUALIFICATION=PASS')
print('NOTE=structural/offline checks only; Vitest requires installed dependencies')
