#!/bin/bash
# Grafana Organization & Team Provisioning for Phase 16.4
# Sets up NOC Team, Engineering Team, Finance Team with correct permissions

set -e

GRAFANA_URL="${GRAFANA_URL:-http://localhost:3000}"
ADMIN_USER="${ADMIN_USER:-admin}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-admin}"

echo "========================================================================"
echo "GRAFANA TEAM PROVISIONING (Phase 16.4)"
echo "========================================================================"
echo ""

# Helper function to make API calls
grafana_api() {
    local method=$1
    local endpoint=$2
    local data=$3

    curl -s -X "$method" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $(get_admin_token)" \
        -d "$data" \
        "$GRAFANA_URL/api$endpoint"
}

# Get admin token for authentication
get_admin_token() {
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "{\"user\":\"$ADMIN_USER\",\"password\":\"$ADMIN_PASSWORD\"}" \
        "$GRAFANA_URL/api/auth/login" | jq -r '.token'
}

echo "Step 1: Creating NOC Team..."
noc_team=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $(get_admin_token)" \
    -d '{
        "name": "NOC Team",
        "email": "noc@elevatediq.ai",
        "orgId": 1
    }' \
    "$GRAFANA_URL/api/teams")

noc_team_id=$(echo "$noc_team" | jq -r '.id')
echo "✅ NOC Team created (ID: $noc_team_id)"

echo ""
echo "Step 2: Creating Engineering Team..."
eng_team=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $(get_admin_token)" \
    -d '{
        "name": "Engineering Team",
        "email": "engineering@elevatediq.ai",
        "orgId": 1
    }' \
    "$GRAFANA_URL/api/teams")

eng_team_id=$(echo "$eng_team" | jq -r '.id')
echo "✅ Engineering Team created (ID: $eng_team_id)"

echo ""
echo "Step 3: Creating Finance Team..."
fin_team=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $(get_admin_token)" \
    -d '{
        "name": "Finance Team",
        "email": "finance@elevatediq.ai",
        "orgId": 1
    }' \
    "$GRAFANA_URL/api/teams")

fin_team_id=$(echo "$fin_team" | jq -r '.id')
echo "✅ Finance Team created (ID: $fin_team_id)"

echo ""
echo "Step 4: Granting dashboard permissions to teams..."

# Get dashboard IDs
slo_dash=$(curl -s -H "Authorization: Bearer $(get_admin_token)" \
    "$GRAFANA_URL/api/search?query=slo-dashboard" | jq -r '.[0].id')

perf_dash=$(curl -s -H "Authorization: Bearer $(get_admin_token)" \
    "$GRAFANA_URL/api/search?query=performance-dashboard" | jq -r '.[0].id')

health_dash=$(curl -s -H "Authorization: Bearer $(get_admin_token)" \
    "$GRAFANA_URL/api/search?query=health-dashboard" | jq -r '.[0].id')

# Grant NOC team Edit permissions on all dashboards
for dash_id in $slo_dash $perf_dash $health_dash; do
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $(get_admin_token)" \
        -d '{"teamId": '"$noc_team_id"', "permission": 2}' \
        "$GRAFANA_URL/api/dashboards/id/$dash_id/permissions"
done
echo "✅ NOC Team granted Edit permissions on all dashboards"

# Grant Engineering team View permissions on SLO and Performance dashboards
for dash_id in $slo_dash $perf_dash; do
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $(get_admin_token)" \
        -d '{"teamId": '"$eng_team_id"', "permission": 1}' \
        "$GRAFANA_URL/api/dashboards/id/$dash_id/permissions"
done
echo "✅ Engineering Team granted View permissions on SLO & Performance dashboards"

# Grant Finance team View permissions on Cost dashboard
curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $(get_admin_token)" \
    -d '{"teamId": '"$fin_team_id"', "permission": 1}' \
    "$GRAFANA_URL/api/dashboards/id/$health_dash/permissions"
echo "✅ Finance Team granted View permissions on Health dashboard"

echo ""
echo "Step 5: Setting up team members..."

# Add members to NOC team (mocking member IDs 1, 2, 3)
for user_id in 2 3 4; do
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $(get_admin_token)" \
        -d '{"userId": '"$user_id"'}' \
        "$GRAFANA_URL/api/teams/$noc_team_id/members"
done
echo "✅ NOC Team members assigned"

# Add members to Engineering team
for user_id in 2 3; do
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $(get_admin_token)" \
        -d '{"userId": '"$user_id"'}' \
        "$GRAFANA_URL/api/teams/$eng_team_id/members"
done
echo "✅ Engineering Team members assigned"

# Add member to Finance team
curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $(get_admin_token)" \
    -d '{"userId": 4}' \
    "$GRAFANA_URL/api/teams/$fin_team_id/members"
echo "✅ Finance Team members assigned"

echo ""
echo "========================================================================"
echo "TEAM PROVISIONING COMPLETE"
echo "========================================================================"
echo ""
echo "Team Configuration:"
echo "  NOC Team (ID: $noc_team_id)"
echo "    • Edit access to all dashboards"
echo "    • Members: 3"
echo ""
echo "  Engineering Team (ID: $eng_team_id)"
echo "    • View access to SLO & Performance dashboards"
echo "    • Members: 2"
echo ""
echo "  Finance Team (ID: $fin_team_id)"
echo "    • View access to Health dashboard"
echo "    • Members: 1"
echo ""
