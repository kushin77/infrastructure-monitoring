#!/usr/bin/env python3
"""Phase 16.4: Dashboard Annotation Provisioning

Generate Grafana provisioning YAML for annotations and alert integration.
This script creates the datasource configuration for Alertmanager annotation queries.
"""

import json
import sys

sys.path.insert(0, ".")

from api.rca_noc_templates import GrafanaProvisioningConfig


def generate_annotation_provisioning():
    """Generate annotation datasource provisioning config"""
    return {
        "apiVersion": 1,
        "datasources": [
            {
                "name": "Prometheus",
                "type": "prometheus",
                "access": "proxy",
                "url": "http://prometheus:9090",
                "isDefault": True,
                "jsonData": {"timeInterval": "15s"},
            },
            {
                "name": "Alertmanager",
                "type": "alertmanager",
                "access": "proxy",
                "url": "http://alertmanager:9093",
                "jsonData": {},
            },
        ],
    }


def generate_alert_rules_provisioning():
    """Generate alert rules provisioning for Phase 16.4 alerts"""
    return {
        "apiVersion": 1,
        "groups": [
            {
                "name": "slo_alerts",
                "interval": "30s",
                "rules": [
                    {
                        "uid": "slo_violation_99_9",
                        "title": "SLO Violation - Availability 99.9%",
                        "condition": "B",
                        "data": [
                            {
                                "refId": "A",
                                "queryType": "",
                                "model": {
                                    "expr": "increase(http_requests_total[5m]) - (increase(http_requests_success_total[5m]))",
                                    "interval": "",
                                    "refId": "A",
                                },
                                "datasourceUid": "prometheus",
                                "relativeTimeRange": {"from": 600, "to": 0},
                            },
                            {
                                "refId": "B",
                                "queryType": "",
                                "mathExpression": "$A / increase(http_requests_total[5m]) * 100 > 0.1",
                                "datasourceUid": "math",
                                "expression": "$A / increase(http_requests_total[5m]) * 100 > 0.1",
                            },
                        ],
                        "noDataState": "NoData",
                        "execErrState": "Alerting",
                        "for": "5m",
                        "annotations": {
                            "summary": "Service SLO violation: Error rate exceeded threshold",
                            "runbook": "https://wiki.elevatediq.ai/runbooks/slo-violations",
                        },
                        "labels": {"severity": "critical", "slo": "availability"},
                    },
                    {
                        "uid": "budget_burn_rate_fast",
                        "title": "Fast Budget Burn",
                        "condition": "B",
                        "data": [
                            {
                                "refId": "A",
                                "queryType": "",
                                "model": {
                                    "expr": "increase(error_rate_per_minute[5m])",
                                    "interval": "",
                                    "refId": "A",
                                },
                                "datasourceUid": "prometheus",
                                "relativeTimeRange": {"from": 600, "to": 0},
                            },
                            {
                                "refId": "B",
                                "queryType": "",
                                "mathExpression": "$A > 10",
                                "datasourceUid": "math",
                                "expression": "$A > 10",
                            },
                        ],
                        "noDataState": "NoData",
                        "execErrState": "Alerting",
                        "for": "2m",
                        "annotations": {
                            "summary": "Error budget burning at excessive rate",
                            "runbook": "https://wiki.elevatediq.ai/runbooks/budget-burn",
                        },
                        "labels": {"severity": "critical", "alert_type": "burn_rate"},
                    },
                ],
            },
            {
                "name": "infrastructure_alerts",
                "interval": "30s",
                "rules": [
                    {
                        "uid": "high_memory_usage",
                        "title": "High Memory Usage",
                        "condition": "B",
                        "data": [
                            {
                                "refId": "A",
                                "queryType": "",
                                "model": {
                                    "expr": "container_memory_usage_bytes / container_spec_memory_limit_bytes * 100",
                                    "interval": "",
                                    "refId": "A",
                                },
                                "datasourceUid": "prometheus",
                            },
                            {
                                "refId": "B",
                                "queryType": "",
                                "mathExpression": "$A > 85",
                                "datasourceUid": "math",
                            },
                        ],
                        "for": "5m",
                        "annotations": {
                            "summary": "High memory usage detected on {{ $labels.pod }}"
                        },
                        "labels": {"severity": "warning"},
                    }
                ],
            },
        ],
    }


def main():
    print("=" * 70)
    print("PHASE 16.4: DASHBOARD ANNOTATION PROVISIONING")
    print("=" * 70)
    print("")

    # Generate provisioning configs
    print("Step 1: Generating annotation datasource provisioning...")
    annotation_config = generate_annotation_provisioning()

    with open("monitoring/grafana/provisioning/datasources/annotations.yml", "w") as f:
        import yaml

        yaml.dump(annotation_config, f, default_flow_style=False)
    print("✅ Annotation datasource provisioning created")

    # Generate alert rules
    print("")
    print("Step 2: Generating alert rules provisioning...")
    alert_rules = generate_alert_rules_provisioning()

    with open("monitoring/grafana/provisioning/alerting/alert-rules.json", "w") as f:
        json.dump(alert_rules, f, indent=2)
    print("✅ Alert rules provisioning created")

    # Generate complete provisioning config
    print("")
    print("Step 3: Generating complete provisioning configuration...")
    all_configs = GrafanaProvisioningConfig.dashboard_provisioning_config()

    with open("monitoring/grafana/provisioning/config.json", "w") as f:
        json.dump(all_configs, f, indent=2)
    print("✅ Complete provisioning configuration created")

    print("")
    print("=" * 70)
    print("PROVISIONING CONFIGURATION COMPLETE")
    print("=" * 70)
    print("")
    print("Created files:")
    print("  • monitoring/grafana/provisioning/datasources/annotations.yml")
    print("  • monitoring/grafana/provisioning/alerting/alert-rules.json")
    print("  • monitoring/grafana/provisioning/config.json")
    print("")
    print("Next steps:")
    print("  1. Restart Grafana to apply datasource provisioning")
    print("  2. Run: bash monitoring/grafana/provisioning/org-teams-setup.sh")
    print("  3. Verify alerts are flowing to dashboard annotations")
    print("")


if __name__ == "__main__":
    try:
        import yaml
    except ImportError:
        print("Installing PyYAML for provisioning...")
        import subprocess

        subprocess.run(["pip", "install", "-q", "pyyaml"], check=True)
        import yaml

    main()
