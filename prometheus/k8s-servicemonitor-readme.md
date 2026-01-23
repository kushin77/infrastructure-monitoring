# Prometheus Operator ServiceMonitor: Gateway API

This directory contains the `ServiceMonitor` for scraping the `gateway-api` `/metrics` endpoint.

Apply in Kubernetes:

```
kubectl apply -f infrastructure/kubernetes/observability/servicemonitor-gateway-api.yaml
```bash

Requirements:

- Prometheus Operator installed
- `gateway-api` Service labeled with `app.kubernetes.io/name: gateway-api`
- Metrics exposed at `/metrics` over the `http` port

Scrape interval: 15s
Path: `/metrics`
Namespace: `elevatediq`
