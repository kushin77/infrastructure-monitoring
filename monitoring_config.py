#!/usr/bin/env python3
"""
Prometheus Metrics & Grafana Dashboard Configuration for Prompt Management System

This module provides:
1. Enhanced Prometheus metrics collection
2. Grafana dashboard JSON
3. Alert rules for uptime monitoring
"""

import json
from typing import Dict, List

# ============================================================================
# GRAFANA DASHBOARD JSON
# ============================================================================

GRAFANA_DASHBOARD = {
    "annotations": {
        "list": [
            {
                "builtIn": 1,
                "datasource": "-- Grafana --",
                "enable": True,
                "hide": True,
                "iconColor": "rgba(0, 211, 255, 1)",
                "name": "Annotations & Alerts",
                "type": "dashboard",
            }
        ]
    },
    "editable": True,
    "gnetId": None,
    "graphTooltip": 0,
    "id": None,
    "links": [],
    "panels": [
        # Panel 1: System Health Overview
        {
            "datasource": "Prometheus",
            "fieldConfig": {
                "defaults": {
                    "color": {"mode": "palette-classic"},
                    "custom": {
                        "axisLabel": "",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 10,
                        "gradientMode": "none",
                        "hideFrom": {
                            "tooltip": False,
                            "viz": False,
                            "legend": False,
                        },
                        "lineInterpolation": "linear",
                        "lineWidth": 1,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear",
                        },
                        "showPoints": "never",
                        "spanNulls": True,
                    },
                    "mappings": [],
                    "max": 100,
                    "min": 0,
                    "thresholds": {
                        "mode": "percentage",
                        "steps": [
                            {"color": "green", "value": None},
                            {"color": "yellow", "value": 80},
                            {"color": "red", "value": 95},
                        ],
                    },
                    "unit": "percent",
                },
                "overrides": [],
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 0,
            },
            "id": 2,
            "options": {
                "legend": {
                    "calcs": ["mean", "last"],
                    "displayMode": "table",
                    "placement": "bottom",
                },
                "tooltip": {
                    "mode": "single",
                },
            },
            "targets": [
                {
                    "expr": "system_cpu_usage_percent",
                    "legendFormat": "CPU Usage %",
                    "refId": "A",
                },
                {
                    "expr": "system_memory_usage_percent",
                    "legendFormat": "Memory Usage %",
                    "refId": "B",
                },
                {
                    "expr": "system_disk_usage_percent",
                    "legendFormat": "Disk Usage %",
                    "refId": "C",
                },
            ],
            "title": "System Resources",
            "type": "timeseries",
        },
        # Panel 2: API Response Latency
        {
            "datasource": "Prometheus",
            "fieldConfig": {
                "defaults": {
                    "color": {"mode": "palette-classic"},
                    "custom": {
                        "axisLabel": "Latency (ms)",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 10,
                        "gradientMode": "none",
                        "hideFrom": {
                            "tooltip": False,
                            "viz": False,
                            "legend": False,
                        },
                        "lineInterpolation": "linear",
                        "lineWidth": 1,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear",
                        },
                        "showPoints": "never",
                        "spanNulls": True,
                    },
                    "mappings": [],
                    "thresholds": {
                        "mode": "absolute",
                        "steps": [
                            {"color": "green", "value": None},
                            {"color": "yellow", "value": 500},
                            {"color": "red", "value": 1000},
                        ],
                    },
                    "unit": "ms",
                },
                "overrides": [],
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 0,
            },
            "id": 3,
            "options": {
                "legend": {
                    "calcs": ["mean", "max", "min"],
                    "displayMode": "table",
                    "placement": "bottom",
                },
                "tooltip": {
                    "mode": "single",
                },
            },
            "targets": [
                {
                    "expr": "histogram_quantile(0.50, rate(prompt_api_request_latency_ms_bucket[5m]))",
                    "legendFormat": "p50",
                    "refId": "A",
                },
                {
                    "expr": "histogram_quantile(0.95, rate(prompt_api_request_latency_ms_bucket[5m]))",
                    "legendFormat": "p95",
                    "refId": "B",
                },
                {
                    "expr": "histogram_quantile(0.99, rate(prompt_api_request_latency_ms_bucket[5m]))",
                    "legendFormat": "p99",
                    "refId": "C",
                },
            ],
            "title": "API Latency Percentiles",
            "type": "timeseries",
        },
        # Panel 3: Request Rate
        {
            "datasource": "Prometheus",
            "fieldConfig": {
                "defaults": {
                    "color": {"mode": "palette-classic"},
                    "custom": {
                        "axisLabel": "Requests/sec",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 10,
                        "gradientMode": "none",
                        "hideFrom": {
                            "tooltip": False,
                            "viz": False,
                            "legend": False,
                        },
                        "lineInterpolation": "linear",
                        "lineWidth": 1,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear",
                        },
                        "showPoints": "never",
                        "spanNulls": True,
                    },
                    "mappings": [],
                    "thresholds": {
                        "mode": "absolute",
                        "steps": [
                            {"color": "green", "value": None},
                        ],
                    },
                    "unit": "reqps",
                },
                "overrides": [],
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 8,
            },
            "id": 4,
            "options": {
                "legend": {
                    "calcs": ["mean", "max"],
                    "displayMode": "table",
                    "placement": "bottom",
                },
                "tooltip": {
                    "mode": "single",
                },
            },
            "targets": [
                {
                    "expr": "rate(prompt_api_requests_total[5m])",
                    "legendFormat": "{{method}} {{endpoint}}",
                    "refId": "A",
                },
            ],
            "title": "Request Rate",
            "type": "timeseries",
        },
        # Panel 4: Error Rate
        {
            "datasource": "Prometheus",
            "fieldConfig": {
                "defaults": {
                    "color": {"mode": "palette-classic"},
                    "custom": {
                        "axisLabel": "Errors/sec",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 10,
                        "gradientMode": "none",
                        "hideFrom": {
                            "tooltip": False,
                            "viz": False,
                            "legend": False,
                        },
                        "lineInterpolation": "linear",
                        "lineWidth": 1,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear",
                        },
                        "showPoints": "never",
                        "spanNulls": True,
                    },
                    "mappings": [],
                    "thresholds": {
                        "mode": "absolute",
                        "steps": [
                            {"color": "green", "value": None},
                            {"color": "yellow", "value": 1},
                            {"color": "red", "value": 5},
                        ],
                    },
                    "unit": "errps",
                },
                "overrides": [],
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 8,
            },
            "id": 5,
            "options": {
                "legend": {
                    "calcs": ["mean", "max"],
                    "displayMode": "table",
                    "placement": "bottom",
                },
                "tooltip": {
                    "mode": "single",
                },
            },
            "targets": [
                {
                    "expr": "rate(prompt_api_errors_total[5m])",
                    "legendFormat": "{{error_type}}",
                    "refId": "A",
                },
            ],
            "title": "Error Rate",
            "type": "timeseries",
        },
        # Panel 5: Token Usage
        {
            "datasource": "Prometheus",
            "fieldConfig": {
                "defaults": {
                    "color": {"mode": "palette-classic"},
                    "custom": {
                        "axisLabel": "Tokens",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 10,
                        "gradientMode": "none",
                        "hideFrom": {
                            "tooltip": False,
                            "viz": False,
                            "legend": False,
                        },
                        "lineInterpolation": "linear",
                        "lineWidth": 1,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear",
                        },
                        "showPoints": "never",
                        "spanNulls": True,
                    },
                    "mappings": [],
                    "thresholds": {
                        "mode": "absolute",
                        "steps": [
                            {"color": "green", "value": None},
                        ],
                    },
                    "unit": "short",
                },
                "overrides": [],
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 16,
            },
            "id": 6,
            "options": {
                "legend": {
                    "calcs": ["sum"],
                    "displayMode": "table",
                    "placement": "bottom",
                },
                "tooltip": {
                    "mode": "single",
                },
            },
            "targets": [
                {
                    "expr": "increase(prompt_tokens_total[5m])",
                    "legendFormat": "{{model}}",
                    "refId": "A",
                },
            ],
            "title": "Token Usage by Model",
            "type": "timeseries",
        },
        # Panel 6: Service Availability
        {
            "datasource": "Prometheus",
            "fieldConfig": {
                "defaults": {
                    "color": {"mode": "palette-classic"},
                    "custom": {
                        "hideFrom": {
                            "tooltip": False,
                            "viz": False,
                            "legend": False,
                        },
                    },
                    "mappings": [],
                    "thresholds": {
                        "mode": "absolute",
                        "steps": [
                            {"color": "red", "value": None},
                            {"color": "yellow", "value": 95},
                            {"color": "green", "value": 99},
                        ],
                    },
                    "unit": "percent",
                },
                "overrides": [],
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 16,
            },
            "id": 7,
            "options": {
                "legend": {
                    "displayMode": "list",
                    "placement": "bottom",
                },
                "pieType": "pie",
                "tooltip": {
                    "mode": "single",
                },
            },
            "targets": [
                {
                    "expr": "(1 - (rate(prompt_api_errors_total[5m]) / rate(prompt_api_requests_total[5m]))) * 100",
                    "legendFormat": "Availability %",
                    "refId": "A",
                },
            ],
            "title": "Service Availability",
            "type": "piechart",
        },
    ],
    "refresh": "10s",
    "schemaVersion": 30,
    "style": "dark",
    "tags": ["prompt-library", "monitoring"],
    "templating": {"list": []},
    "time": {
        "from": "now-6h",
        "to": "now",
    },
    "timepicker": {
        "refresh_intervals": ["10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"],
    },
    "timezone": "browser",
    "title": "Prompt Library System Health",
    "uid": "prompt-library-health",
    "version": 1,
}

# ============================================================================
# PROMETHEUS ALERT RULES
# ============================================================================

PROMETHEUS_ALERT_RULES = """
groups:
  - name: prompt_library_alerts
    interval: 30s
    rules:
      # API Availability
      - alert: PromptAPIDown
        expr: up{job="prompt-api"} == 0
        for: 2m
        labels:
          severity: critical
          component: api
        annotations:
          summary: "Prompt API is down"
          description: "Prompt API server is not responding"

      - alert: PromptAPIHighErrorRate
        expr: rate(prompt_api_errors_total[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
          component: api
        annotations:
          summary: "High error rate in Prompt API"
          description: "Error rate is {{ $value | humanizePercentage }}"

      # Latency Alerts
      - alert: PromptAPIHighLatency
        expr: histogram_quantile(0.95, rate(prompt_api_request_latency_ms_bucket[5m])) > 1000
        for: 5m
        labels:
          severity: warning
          component: api
        annotations:
          summary: "High API latency detected"
          description: "p95 latency is {{ $value | humanize }}ms"

      # System Resources
      - alert: HighCPUUsage
        expr: system_cpu_usage_percent > 80
        for: 5m
        labels:
          severity: warning
          component: system
        annotations:
          summary: "High CPU usage"
          description: "CPU usage is {{ $value | humanize }}%"

      - alert: HighMemoryUsage
        expr: system_memory_usage_percent > 85
        for: 5m
        labels:
          severity: warning
          component: system
        annotations:
          summary: "High memory usage"
          description: "Memory usage is {{ $value | humanize }}%"

      - alert: CriticalDiskSpace
        expr: system_disk_usage_percent > 90
        for: 5m
        labels:
          severity: critical
          component: system
        annotations:
          summary: "Critical disk space"
          description: "Disk usage is {{ $value | humanize }}%"

      # Database
      - alert: DatabaseConnectionPoolExhausted
        expr: db_connection_pool_available < 5
        for: 2m
        labels:
          severity: critical
          component: database
        annotations:
          summary: "Database connection pool nearly exhausted"
          description: "Only {{ $value }} connections available"

      # Cache
      - alert: CacheLowHitRate
        expr: cache_hit_rate < 0.5
        for: 10m
        labels:
          severity: warning
          component: cache
        annotations:
          summary: "Cache hit rate is low"
          description: "Cache hit rate is {{ $value | humanizePercentage }}"

      # Ollama
      - alert: OllamaServiceDown
        expr: up{job="ollama"} == 0
        for: 2m
        labels:
          severity: critical
          component: ollama
        annotations:
          summary: "Ollama service is down"
          description: "Ollama embeddings service is not responding"
"""

# ============================================================================
# DOCKER COMPOSE PROMETHEUS CONFIGURATION
# ============================================================================

PROMETHEUS_CONFIG = """
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'prompt-library'
    environment: 'production'

alerting:
  alertmanagers:
    - static_configs:
        - targets:
            - alertmanager:9093

rule_files:
  - '/etc/prometheus/alert_rules.yml'

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'prompt-api'
    static_configs:
      - targets: ['prompt-api:8080']
    metrics_path: '/metrics'
    scrape_interval: 10s

  - job_name: 'prompt-api-enterprise'
    static_configs:
      - targets: ['prompt-api-enterprise:8081']
    metrics_path: '/metrics'
    scrape_interval: 10s

  - job_name: 'ollama'
    static_configs:
      - targets: ['ollama:11434']
    metrics_path: '/metrics'
    scrape_interval: 30s

  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']
    scrape_interval: 10s

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']
    scrape_interval: 10s
"""


def export_dashboard_json() -> str:
    """Export Grafana dashboard as JSON"""
    return json.dumps(GRAFANA_DASHBOARD, indent=2)


def export_alert_rules() -> str:
    """Export Prometheus alert rules"""
    return PROMETHEUS_ALERT_RULES


def export_prometheus_config() -> str:
    """Export Prometheus configuration"""
    return PROMETHEUS_CONFIG
