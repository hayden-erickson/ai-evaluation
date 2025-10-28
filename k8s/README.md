# Kubernetes Manifests

This directory contains all Kubernetes manifests for the Habit Notifier CronJob.

## Files

- **`namespace.yaml`** - Creates the `habit-tracker` namespace
- **`configmap.yaml`** - Non-sensitive configuration (database host, port, etc.)
- **`secret.yaml`** - Sensitive credentials (database password, Twilio tokens)
- **`serviceaccount.yaml`** - ServiceAccount and RBAC configuration
- **`cronjob.yaml`** - Main CronJob definition
- **`networkpolicy.yaml`** - Network security policies
- **`kustomization.yaml`** - Kustomize configuration for easy deployment

## Quick Start

### Deploy everything

```bash
kubectl apply -k .
```

### Deploy individual resources

```bash
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f serviceaccount.yaml
kubectl apply -f networkpolicy.yaml
kubectl apply -f cronjob.yaml
```

## Configuration

### Before deploying

1. **Update `secret.yaml`** with your actual credentials
2. **Update `configmap.yaml`** with your database host and Twilio number
3. **Update `cronjob.yaml`** with your Docker image reference
4. **Update `kustomization.yaml`** with your image registry

### Using Kustomize overlays (recommended for multiple environments)

Create environment-specific overlays:

```
k8s/
├── base/
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── cronjob.yaml
│   └── kustomization.yaml
└── overlays/
    ├── dev/
    │   ├── kustomization.yaml
    │   └── configmap.yaml
    ├── staging/
    │   ├── kustomization.yaml
    │   └── configmap.yaml
    └── production/
        ├── kustomization.yaml
        └── configmap.yaml
```

Deploy to specific environment:

```bash
kubectl apply -k overlays/production/
```

## Security Notes

⚠️ **NEVER commit actual secrets to version control!**

- Use placeholder values in `secret.yaml`
- Create secrets directly with `kubectl create secret`
- Or use External Secrets Operator for production

## Validation

Validate manifests without applying:

```bash
kubectl apply -k . --dry-run=client
```

## Troubleshooting

View CronJob status:

```bash
kubectl get cronjob -n habit-tracker
```

View job history:

```bash
kubectl get jobs -n habit-tracker
```

View logs:

```bash
kubectl logs -n habit-tracker -l app=habit-notifier
```
