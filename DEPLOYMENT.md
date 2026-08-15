# Atlas Runtime Deployment
Required secret: `ATLAS_API_TOKEN`.
Optional: `ATLAS_ADDR` (default `:8080`).
Endpoints: `GET /healthz`, `GET /readyz`, authenticated `POST /v1/execute`.
Deploy staging before production. Never commit runtime secrets.
