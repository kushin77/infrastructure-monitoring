#!/bin/bash

# Phase 5D Grafana Dashboard Generation Script
# Creates all 18 dashboards for Week 1-3 services
# Usage: bash generate-dashboards.sh

set -e

DASHBOARDS_DIR="infrastructure/monitoring/grafana/dashboards/phase-5d"
GRAFANA_API_URL="${GRAFANA_URL:-http://localhost:3000}"
GRAFANA_API_KEY="${GRAFANA_API_KEY}"

if [ -z "$GRAFANA_API_KEY" ]; then
    echo "⚠️  GRAFANA_API_KEY not set. Using read-only mode (preview only)."
fi

echo "🚀 Phase 5D Grafana Dashboard Generation"
echo "========================================"

# Week 1 Dashboards (6 total: 2 per service)
create_week1_dashboards() {
    echo ""
    echo "📊 Generating Week 1 Dashboards (6 total)..."

    # EIP Portal Dashboards
    cat > "$DASHBOARDS_DIR/week1-eip-portal-overview.json" << 'EOFDASH1'
{
  "dashboard": {
    "title": "Phase 5D Week 1 - EIP Portal Overview",
    "panels": [
      {
        "title": "Gateway Status",
        "targets": [{"expr": "up{job=\"eip-portal\"}"}],
        "type": "stat"
      },
      {
        "title": "Response Time (P95)",
        "targets": [{"expr": "histogram_quantile(0.95, rate(eip_request_duration_seconds_bucket[5m]))"}],
        "type": "graph"
      },
      {
        "title": "Error Rate",
        "targets": [{"expr": "rate(eip_errors_total[5m])"}],
        "type": "graph"
      },
      {
        "title": "Memory Usage",
        "targets": [{"expr": "container_memory_usage_bytes{pod=~\"eip-portal-.*\"} / container_spec_memory_limit_bytes"}],
        "type": "gauge"
      }
    ],
    "refresh": "30s",
    "time": {"from": "now-1h", "to": "now"},
    "timezone": "UTC"
  }
}
EOFDASH1
    echo "  ✅ EIP Portal Overview"

    cat > "$DASHBOARDS_DIR/week1-eip-portal-detailed.json" << 'EOFDASH1B'
{
  "dashboard": {
    "title": "Phase 5D Week 1 - EIP Portal Detailed",
    "panels": [
      {
        "title": "Redis Connection Status",
        "targets": [{"expr": "redis_connected_clients{job=\"eip-portal\"}"}],
        "type": "stat"
      },
      {
        "title": "Database Connection Pool",
        "targets": [{"expr": "db_connection_pool_available{job=\"eip-portal\"}"}],
        "type": "gauge"
      },
      {
        "title": "TLS Certificate Expiry",
        "targets": [{"expr": "(eip_tls_cert_expiry_seconds - time()) / 86400"}],
        "type": "stat"
      },
      {
        "title": "Appsmith Proxy Errors",
        "targets": [{"expr": "rate(eip_appsmith_proxy_errors_total[5m])"}],
        "type": "graph"
      }
    ],
    "refresh": "30s",
    "time": {"from": "now-1h", "to": "now"}
  }
}
EOFDASH1B
    echo "  ✅ EIP Portal Detailed"

    # RoundRobin Core Dashboards
    cat > "$DASHBOARDS_DIR/week1-roundrobin-core-overview.json" << 'EOFDASH2'
{
  "dashboard": {
    "title": "Phase 5D Week 1 - RoundRobin Core Overview",
    "panels": [
      {
        "title": "API Availability",
        "targets": [{"expr": "up{job=\"roundrobin-core\"}"}],
        "type": "stat"
      },
      {
        "title": "Error Rate (Critical Threshold: >0.05%)",
        "targets": [{"expr": "rate(rr_core_errors_total[5m])"}],
        "type": "graph"
      },
      {
        "title": "Response Time (P95)",
        "targets": [{"expr": "histogram_quantile(0.95, rate(rr_core_request_duration_seconds_bucket[5m]))"}],
        "type": "graph"
      },
      {
        "title": "PostgreSQL Slow Queries",
        "targets": [{"expr": "rate(rr_core_slow_queries_total[5m])"}],
        "type": "graph"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH2
    echo "  ✅ RoundRobin Core Overview"

    cat > "$DASHBOARDS_DIR/week1-roundrobin-core-detailed.json" << 'EOFDASH2B'
{
  "dashboard": {
    "title": "Phase 5D Week 1 - RoundRobin Core Detailed",
    "panels": [
      {
        "title": "Tenant Isolation Violations (SECURITY)",
        "targets": [{"expr": "increase(rr_core_tenant_isolation_violations_total[5m])"}],
        "type": "stat"
      },
      {
        "title": "Connection Pool Status",
        "targets": [{"expr": "rr_core_connection_pool_available / rr_core_connection_pool_total"}],
        "type": "gauge"
      },
      {
        "title": "Cache Hit Rate",
        "targets": [{"expr": "rate(rr_core_cache_hits_total[5m]) / rate(rr_core_cache_requests_total[5m])"}],
        "type": "graph"
      },
      {
        "title": "JWT Validation Failures",
        "targets": [{"expr": "rate(rr_core_jwt_validation_failures_total[5m])"}],
        "type": "graph"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH2B
    echo "  ✅ RoundRobin Core Detailed"

    # AI Chatbot Orchestrator Dashboards
    cat > "$DASHBOARDS_DIR/week1-chatbot-overview.json" << 'EOFDASH3'
{
  "dashboard": {
    "title": "Phase 5D Week 1 - AI Chatbot Orchestrator Overview",
    "panels": [
      {
        "title": "Service Status",
        "targets": [{"expr": "up{job=\"ai-chatbot-orchestrator\"}"}],
        "type": "stat"
      },
      {
        "title": "Model Routing Success Rate",
        "targets": [{"expr": "rate(chatbot_model_routing_success_total[5m]) / rate(chatbot_model_routing_attempts_total[5m])"}],
        "type": "gauge"
      },
      {
        "title": "Response Time by Model",
        "targets": [{"expr": "histogram_quantile(0.95, rate(chatbot_response_duration_seconds_bucket[5m]))"}],
        "type": "graph"
      },
      {
        "title": "Cost Tracking",
        "targets": [{"expr": "chatbot_total_cost_usd"}],
        "type": "stat"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH3
    echo "  ✅ Chatbot Overview"

    cat > "$DASHBOARDS_DIR/week1-chatbot-detailed.json" << 'EOFDASH3B'
{
  "dashboard": {
    "title": "Phase 5D Week 1 - AI Chatbot Orchestrator Detailed",
    "panels": [
      {
        "title": "Ollama vs Cloud Model Usage",
        "targets": [
          {"expr": "rate(chatbot_ollama_requests_total[5m])", "legendFormat": "Ollama"},
          {"expr": "rate(chatbot_claude_requests_total[5m])", "legendFormat": "Claude"},
          {"expr": "rate(chatbot_openai_requests_total[5m])", "legendFormat": "OpenAI"}
        ],
        "type": "graph"
      },
      {
        "title": "Conversation Storage Status",
        "targets": [{"expr": "chatbot_redis_storage_used_bytes / chatbot_redis_storage_limit_bytes"}],
        "type": "gauge"
      },
      {
        "title": "Cost per Conversation",
        "targets": [{"expr": "chatbot_cost_per_conversation_usd"}],
        "type": "graph"
      },
      {
        "title": "Token Limit Violations",
        "targets": [{"expr": "rate(chatbot_token_limit_exceeded_total[5m])"}],
        "type": "graph"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH3B
    echo "  ✅ Chatbot Detailed"
}

# Week 2 Dashboards (6 total)
create_week2_dashboards() {
    echo ""
    echo "📊 Generating Week 2 Dashboards (6 total)..."

    # News Feed Engine
    cat > "$DASHBOARDS_DIR/week2-news-feed-overview.json" << 'EOFDASH4'
{
  "dashboard": {
    "title": "Phase 5D Week 2 - News Feed Engine Overview",
    "panels": [
      {
        "title": "Service Availability",
        "targets": [{"expr": "up{job=\"news-feed-engine\"}"}],
        "type": "stat"
      },
      {
        "title": "Ingestion Lag",
        "targets": [{"expr": "news_feed_ingestion_lag_seconds"}],
        "type": "stat"
      },
      {
        "title": "Error Rate",
        "targets": [{"expr": "rate(news_feed_errors_total[5m])"}],
        "type": "graph"
      },
      {
        "title": "Kafka Consumer Lag",
        "targets": [{"expr": "kafka_consumer_lag{job=\"news-feed-engine\"}"}],
        "type": "graph"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH4
    echo "  ✅ News Feed Overview"

    cat > "$DASHBOARDS_DIR/week2-news-feed-detailed.json" << 'EOFDASH4B'
{
  "dashboard": {
    "title": "Phase 5D Week 2 - News Feed Engine Detailed",
    "panels": [
      {
        "title": "Video Generation Success Rate",
        "targets": [{"expr": "rate(news_feed_video_generation_success_total[5m]) / rate(news_feed_video_generation_attempts_total[5m])"}],
        "type": "gauge"
      },
      {
        "title": "Content Quality Scores",
        "targets": [{"expr": "news_feed_content_quality_score"}],
        "type": "graph"
      },
      {
        "title": "External API Rate Limits",
        "targets": [{"expr": "rate(news_feed_api_rate_limit_exceeded_total[5m])"}],
        "type": "stat"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH4B
    echo "  ✅ News Feed Detailed"

    # Gateway API
    cat > "$DASHBOARDS_DIR/week2-gateway-overview.json" << 'EOFDASH5'
{
  "dashboard": {
    "title": "Phase 5D Week 2 - Gateway API Overview",
    "panels": [
      {
        "title": "Gateway Status",
        "targets": [{"expr": "up{job=\"gateway-api\"}"}],
        "type": "stat"
      },
      {
        "title": "Error Rate",
        "targets": [{"expr": "rate(gateway_errors_total[5m])"}],
        "type": "graph"
      },
      {
        "title": "Response Time (P95)",
        "targets": [{"expr": "histogram_quantile(0.95, rate(gateway_request_duration_seconds_bucket[5m]))"}],
        "type": "graph"
      },
      {
        "title": "Circuit Breaker Status",
        "targets": [{"expr": "increase(gateway_circuit_breaker_open_total[5m])"}],
        "type": "stat"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH5
    echo "  ✅ Gateway Overview"

    cat > "$DASHBOARDS_DIR/week2-gateway-detailed.json" << 'EOFDASH5B'
{
  "dashboard": {
    "title": "Phase 5D Week 2 - Gateway API Detailed",
    "panels": [
      {
        "title": "Downstream Service Health",
        "targets": [{"expr": "increase(gateway_downstream_unavailable_total[5m])"}],
        "type": "graph"
      },
      {
        "title": "Rate Limit Status",
        "targets": [{"expr": "rate(gateway_rate_limit_exceeded_total[5m])"}],
        "type": "graph"
      },
      {
        "title": "Cache Hit Ratio",
        "targets": [{"expr": "rate(gateway_cache_hits_total[5m]) / rate(gateway_cache_requests_total[5m])"}],
        "type": "gauge"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH5B
    echo "  ✅ Gateway Detailed"

    # Super Admin API
    cat > "$DASHBOARDS_DIR/week2-super-admin-overview.json" << 'EOFDASH6'
{
  "dashboard": {
    "title": "Phase 5D Week 2 - Super Admin API Overview",
    "panels": [
      {
        "title": "API Status",
        "targets": [{"expr": "up{job=\"super-admin-api\"}"}],
        "type": "stat"
      },
      {
        "title": "Billing Sync Health",
        "targets": [{"expr": "increase(super_admin_billing_sync_failures_total[5m])"}],
        "type": "stat"
      },
      {
        "title": "Tenant Query Performance",
        "targets": [{"expr": "histogram_quantile(0.99, rate(super_admin_tenant_query_duration_seconds_bucket[5m]))"}],
        "type": "graph"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH6
    echo "  ✅ Super Admin Overview"

    cat > "$DASHBOARDS_DIR/week2-super-admin-detailed.json" << 'EOFDASH6B'
{
  "dashboard": {
    "title": "Phase 5D Week 2 - Super Admin API Detailed",
    "panels": [
      {
        "title": "RBAC Permission Denials",
        "targets": [{"expr": "rate(super_admin_permission_denied_total[5m])"}],
        "type": "stat"
      },
      {
        "title": "Audit Log Lag",
        "targets": [{"expr": "super_admin_audit_log_lag_seconds"}],
        "type": "graph"
      },
      {
        "title": "Feature Flag Cache Staleness",
        "targets": [{"expr": "time() - super_admin_feature_flags_last_update_timestamp"}],
        "type": "stat"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH6B
    echo "  ✅ Super Admin Detailed"
}

# Week 3 Dashboards (6 total)
create_week3_dashboards() {
    echo ""
    echo "📊 Generating Week 3 Dashboards (6 total)..."

    # Analytics Service
    cat > "$DASHBOARDS_DIR/week3-analytics-overview.json" << 'EOFDASH7'
{
  "dashboard": {
    "title": "Phase 5D Week 3 - Analytics Service Overview",
    "panels": [
      {
        "title": "Service Status",
        "targets": [{"expr": "up{job=\"analytics\"}"}],
        "type": "stat"
      },
      {
        "title": "Query Latency (P95)",
        "targets": [{"expr": "histogram_quantile(0.95, rate(analytics_query_duration_seconds_bucket[5m]))"}],
        "type": "graph"
      },
      {
        "title": "Data Freshness",
        "targets": [{"expr": "time() - analytics_data_freshness_timestamp"}],
        "type": "stat"
      },
      {
        "title": "Storage Usage",
        "targets": [{"expr": "analytics_storage_used_bytes / analytics_storage_limit_bytes"}],
        "type": "gauge"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH7
    echo "  ✅ Analytics Overview"

    cat > "$DASHBOARDS_DIR/week3-analytics-detailed.json" << 'EOFDASH7B'
{
  "dashboard": {
    "title": "Phase 5D Week 3 - Analytics Service Detailed",
    "panels": [
      {
        "title": "ClickHouse Replication Lag",
        "targets": [{"expr": "analytics_clickhouse_replication_lag_seconds"}],
        "type": "graph"
      },
      {
        "title": "Aggregation Job Status",
        "targets": [{"expr": "increase(analytics_aggregation_job_failures_total[5m])"}],
        "type": "stat"
      },
      {
        "title": "Cost Calculation Errors",
        "targets": [{"expr": "increase(analytics_cost_calculation_errors_total[5m])"}],
        "type": "stat"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH7B
    echo "  ✅ Analytics Detailed"

    # AIOps Engine
    cat > "$DASHBOARDS_DIR/week3-aiops-overview.json" << 'EOFDASH8'
{
  "dashboard": {
    "title": "Phase 5D Week 3 - AIOps Engine Overview",
    "panels": [
      {
        "title": "Service Status",
        "targets": [{"expr": "up{job=\"aiops-engine\"}"}],
        "type": "stat"
      },
      {
        "title": "Model Accuracy",
        "targets": [{"expr": "aiops_anomaly_detection_accuracy"}],
        "type": "gauge"
      },
      {
        "title": "Inference Latency (P95)",
        "targets": [{"expr": "histogram_quantile(0.95, rate(aiops_model_serving_latency_seconds_bucket[5m]))"}],
        "type": "graph"
      },
      {
        "title": "GPU Memory Usage",
        "targets": [{"expr": "aiops_gpu_memory_used_bytes / aiops_gpu_memory_total_bytes"}],
        "type": "gauge"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH8
    echo "  ✅ AIOps Overview"

    cat > "$DASHBOARDS_DIR/week3-aiops-detailed.json" << 'EOFDASH8B'
{
  "dashboard": {
    "title": "Phase 5D Week 3 - AIOps Engine Detailed",
    "panels": [
      {
        "title": "Forecast Error (MAE)",
        "targets": [{"expr": "aiops_forecast_mean_absolute_error"}],
        "type": "graph"
      },
      {
        "title": "Model Training Status",
        "targets": [{"expr": "increase(aiops_model_training_failures_total[5m])"}],
        "type": "stat"
      },
      {
        "title": "Prediction Confidence",
        "targets": [{"expr": "avg(aiops_prediction_confidence)"}],
        "type": "gauge"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH8B
    echo "  ✅ AIOps Detailed"

    # Defrag Portal
    cat > "$DASHBOARDS_DIR/week3-defrag-overview.json" << 'EOFDASH9'
{
  "dashboard": {
    "title": "Phase 5D Week 3 - Defrag Portal Overview",
    "panels": [
      {
        "title": "Portal Status",
        "targets": [{"expr": "up{job=\"defrag-portal\"}"}],
        "type": "stat"
      },
      {
        "title": "Error Rate",
        "targets": [{"expr": "rate(defrag_errors_total[5m])"}],
        "type": "graph"
      },
      {
        "title": "Intervention Queue Depth",
        "targets": [{"expr": "defrag_intervention_queue_depth"}],
        "type": "stat"
      },
      {
        "title": "Feedback Accuracy",
        "targets": [{"expr": "defrag_feedback_accuracy_score"}],
        "type": "gauge"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH9
    echo "  ✅ Defrag Overview"

    cat > "$DASHBOARDS_DIR/week3-defrag-detailed.json" << 'EOFDASH9B'
{
  "dashboard": {
    "title": "Phase 5D Week 3 - Defrag Portal Detailed",
    "panels": [
      {
        "title": "Redis Pub/Sub Lag",
        "targets": [{"expr": "defrag_pubsub_lag_seconds"}],
        "type": "graph"
      },
      {
        "title": "Policy Enforcement Lag",
        "targets": [{"expr": "time() - defrag_last_enforcement_timestamp"}],
        "type": "stat"
      },
      {
        "title": "Feedback Loss Events",
        "targets": [{"expr": "increase(defrag_feedback_write_errors_total[5m])"}],
        "type": "stat"
      }
    ],
    "refresh": "30s"
  }
}
EOFDASH9B
    echo "  ✅ Defrag Detailed"
}

# Deploy to Grafana (if API key provided)
deploy_to_grafana() {
    if [ -z "$GRAFANA_API_KEY" ]; then
        echo ""
        echo "⚠️  Skipping Grafana deployment (GRAFANA_API_KEY not set)"
        echo "To deploy: export GRAFANA_API_KEY=... && bash $0"
        return
    fi

    echo ""
    echo "🚀 Deploying dashboards to Grafana..."

    for dashboard in "$DASHBOARDS_DIR"/*.json; do
        name=$(basename "$dashboard")
        echo "  Uploading: $name"

        curl -s -X POST \
            -H "Authorization: Bearer $GRAFANA_API_KEY" \
            -H "Content-Type: application/json" \
            -d @"$dashboard" \
            "$GRAFANA_API_URL/api/dashboards/db" || echo "    ⚠️  Failed"
    done

    echo "✅ Deployment complete"
}

# Main execution
main() {
    mkdir -p "$DASHBOARDS_DIR"
    create_week1_dashboards
    create_week2_dashboards
    create_week3_dashboards
    deploy_to_grafana

    echo ""
    echo "✅ All 18 dashboards generated"
    echo "�� Location: $DASHBOARDS_DIR/"
    echo ""
    echo "📋 Dashboard Summary:"
    echo "  Week 1: 6 dashboards (EIP Portal, RoundRobin Core, ChatBot)"
    echo "  Week 2: 6 dashboards (News Feed, Gateway, SuperAdmin)"
    echo "  Week 3: 6 dashboards (Analytics, AIOps, Defrag)"
    echo ""
    echo "📝 Next steps:"
    echo "  1. Review dashboard JSON files"
    echo "  2. Deploy to Grafana: export GRAFANA_API_KEY=... && bash $0"
    echo "  3. Add to Kubernetes ConfigMap for persistent storage"
}

main
