# Migration Notes: infrastructure-monitoring

This repository was extracted from the elevatediq monorepo during Phase-3 dissection.

## Setup Instructions

### 1. Clone Repository
```bash
git clone https://github.com/kushin77/infrastructure-monitoring
cd infrastructure-monitoring
```

### 2. Install Dependencies

**For YAML**:
```bash
# YAML config files - no installation needed
```

### 3. Validate Configuration
```bash
kubectl apply --dry-run=client -f .
```

### 4. Build Docker Image
```bash
docker build -t kushin77/infrastructure-monitoring:latest .
```

### 5. Deploy

**Terraform**:
```bash
terraform init
terraform plan
terraform apply
```

**Kubernetes**:
```bash
kubectl apply -k .
```

**Helm**:
```bash
helm install infrastructure-monitoring ./
```

**Docker Compose**:
```bash
docker-compose up -d
```

## CI/CD Checklist

- [ ] Update `.github/workflows/ci.yml` with validation commands
- [ ] Configure GitHub Actions secrets if needed
- [ ] Test locally before pushing
- [ ] Update Dockerfile with production settings
- [ ] Test deployment in staging environment
- [ ] Add any missing environment variables to CI config
- [ ] Test in production-like environment

## Source Mapping

Original location: `infrastructure/monitoring`  
Full monorepo: https://github.com/elevatediq-monorepo

## Next Steps

1. Follow the setup instructions above
2. Validate configuration in your environment
3. Customize for your deployment needs
4. Deploy using your CI/CD pipeline

