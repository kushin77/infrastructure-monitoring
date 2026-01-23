# Grafana: fast, secure, and on-brand dashboards

Benefits-first: enable OAuth in minutes, route users to the right dashboards, and keep copy/design aligned with the unified design system.

## OAuth (Google)

Environment variables (set via compose, env file, or Kubernetes manifests):

```bash
GF_AUTH_GOOGLE_ENABLED=true
GF_AUTH_GOOGLE_CLIENT_ID=<from GSM>
GF_AUTH_GOOGLE_CLIENT_SECRET=<from GSM>
GF_AUTH_GOOGLE_ALLOWED_DOMAINS=bioenergystrategies.com
GF_AUTH_DISABLE_LOGIN_FORM=true
GF_AUTH_ANONYMOUS_ENABLED=false
GF_SERVER_ROOT_URL=https://dev.elevatediq.ai/grafana
```bash

Notes

- Client ID/Secret must be sourced from Google Secret Manager (see references)
- Disable anonymous access and the login form; use domain allow-list
- Root URL should match your routed path (e.g., Traefik path prefix)

## Quick start

1) Verify Prometheus targets reachable

```bash
curl -s http://dev.elevatediq.ai:9091/api/v1/targets | jq '.data.activeTargets | length'
```bash

2) Check Grafana API health

```bash
curl -s http://dev.elevatediq.ai:3000/api/health | jq .
```bash

3) Provision Prometheus datasource and System Health dashboard (see provisioning)

4) Set home dashboard to System Health (UID: `system-health`) for faster incident triage.

## Branding & landing copy

- Headline: "Your infrastructure self-heals — dashboards confirm"
- Default landing: System Health dashboard
- Colors: Use semantic status colors consistently (success/info/warning/error)
- Copy: benefits-first; avoid jargon; reassure during errors

## Provisioning (Datasources & Dashboards)

Place provisioning files here:

- Datasource: `provisioning/datasources/prometheus.yaml`
- Dashboard: `provisioning/dashboards/system-health.json`

These will set Prometheus as the default datasource and provide a stub System Health dashboard.

### How to apply (example)

1) Copy env template and fill values from GSM:

```

cp infrastructure/monitoring/grafana/.env.example infrastructure/monitoring/grafana/.env

# Edit .env with CLIENT_ID/SECRET from GSM

```bash

2) Ensure Prometheus URL is reachable (defaults to `http://prometheus:9090`):

```

export PROMETHEUS_URL=<http://prometheus:9090>

```bash

3) Mount provisioning into Grafana (compose or K8s) and restart Grafana:

```bash
# Compose fragment (illustrative):
# volumes:
#   - ./infrastructure/monitoring/grafana/provisioning:/etc/grafana/provisioning:ro

# Restart Grafana (method depends on environment)
```bash

4) Set default landing to System Health (optional):

```bash
# In Grafana UI: set home dashboard to uid `system-health`
```bash

### System Health panels (current)

- Stat: Targets Up — `sum(up)`
- Bar: Targets Up by Job — `sum by (job) (up)`
- Table: Instance Status — `up` (legend `{{instance}} ({{job}})`)
- Timeseries: Up by Job (6h) — `sum by (job) (up)`
- Variable: `job` — `label_values(up, job)` (multi, include All)

## References

- `docs/security/SECRETS_AND_AUTH_CONSOLIDATED.md` — OAuth secrets via GSM
- `docs/design-system/MESSAGING-PATTERNS.md` — Copy patterns
- `docs/design-system/IMPLEMENTATION.md` — Rollout sequence and specifics
- `docs/phase-3/PHASE-3-EXECUTION-PLAN.md` — Unified design system plan (tokens, components)
- `docs/design-system/TOKENS.md` — Design tokens (color, type, spacing)
- `infrastructure/monitoring/grafana/dashboards/` — Dashboard JSON files
- `infrastructure/monitoring/grafana/provisioning/` — Datasource/dashboard provisioning
