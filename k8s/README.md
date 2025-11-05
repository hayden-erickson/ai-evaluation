# Kubernetes Deployment Guide for Google Cloud

This directory contains all the Kubernetes configuration files needed to deploy the Habits API to Google Kubernetes Engine (GKE).

## Prerequisites

1. Google Cloud SDK installed
2. kubectl installed
3. Docker installed
4. A Google Cloud project with billing enabled

## Setup Steps

### 1. Set up Google Cloud Project

```bash
# Set your project ID
export PROJECT_ID=your-project-id
gcloud config set project $PROJECT_ID

# Enable required APIs
gcloud services enable container.googleapis.com
gcloud services enable containerregistry.googleapis.com
```

### 2. Create a GKE Cluster

```bash
# Create a GKE cluster
gcloud container clusters create habits-cluster \
  --zone us-central1-a \
  --num-nodes 3 \
  --machine-type n1-standard-1 \
  --disk-size 10

# Get credentials for kubectl
gcloud container clusters get-credentials habits-cluster --zone us-central1-a
```

### 3. Build and Push Docker Image

```bash
# Build the Docker image
docker build -t gcr.io/$PROJECT_ID/habits-api:latest .

# Push to Google Container Registry
docker push gcr.io/$PROJECT_ID/habits-api:latest
```

### 4. Update Kubernetes Configuration

Edit `k8s/app-deployment.yaml` and replace `YOUR_PROJECT_ID` with your actual Google Cloud project ID:

```yaml
image: gcr.io/YOUR_PROJECT_ID/habits-api:latest
```

### 5. Update Secrets

**IMPORTANT**: Before deploying, update the following files with secure values:

- `k8s/mysql-secret.yaml` - Update MySQL passwords
- `k8s/app-secret.yaml` - Update JWT secret

For production, use base64 encoded values instead of `stringData`:

```bash
echo -n "your-secure-password" | base64
```

### 6. Deploy to Kubernetes

Apply the configurations in the following order:

```bash
# Create namespace
kubectl apply -f k8s/namespace.yaml

# Create secrets
kubectl apply -f k8s/mysql-secret.yaml
kubectl apply -f k8s/app-secret.yaml

# Create ConfigMap
kubectl apply -f k8s/app-configmap.yaml

# Deploy MySQL
kubectl apply -f k8s/mysql-pvc.yaml
kubectl apply -f k8s/mysql-deployment.yaml
kubectl apply -f k8s/mysql-service.yaml

# Wait for MySQL to be ready
kubectl wait --for=condition=ready pod -l app=mysql -n habits-app --timeout=300s

# Deploy the application
kubectl apply -f k8s/app-deployment.yaml
kubectl apply -f k8s/app-service.yaml
```

### 7. Verify Deployment

```bash
# Check pod status
kubectl get pods -n habits-app

# Check services
kubectl get services -n habits-app

# Get the external IP (may take a few minutes)
kubectl get service habits-api-service -n habits-app
```

### 8. Test the API

Once the LoadBalancer has an external IP:

```bash
# Get the external IP
EXTERNAL_IP=$(kubectl get service habits-api-service -n habits-app -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

# Test health endpoint
curl http://$EXTERNAL_IP/health
```

## Monitoring and Logs

```bash
# View application logs
kubectl logs -f deployment/habits-api -n habits-app

# View MySQL logs
kubectl logs -f deployment/mysql -n habits-app

# Describe a pod for troubleshooting
kubectl describe pod <pod-name> -n habits-app
```

## Scaling

```bash
# Scale the application
kubectl scale deployment habits-api --replicas=5 -n habits-app
```

## Cleanup

```bash
# Delete all resources
kubectl delete namespace habits-app

# Delete the GKE cluster
gcloud container clusters delete habits-cluster --zone us-central1-a
```

## Configuration Files

- `namespace.yaml` - Creates the habits-app namespace
- `mysql-secret.yaml` - MySQL credentials (update before deploying!)
- `mysql-pvc.yaml` - Persistent storage for MySQL data
- `mysql-deployment.yaml` - MySQL database deployment
- `mysql-service.yaml` - MySQL internal service
- `app-secret.yaml` - Application secrets (JWT, etc.)
- `app-configmap.yaml` - Application configuration
- `app-deployment.yaml` - Main application deployment
- `app-service.yaml` - LoadBalancer service for external access

## Environment Variables

The application uses the following environment variables:

- `PORT` - Application port (default: 8080)
- `DB_HOST` - MySQL host and port
- `DB_NAME` - Database name
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `JWT_SECRET` - Secret key for JWT token generation

## Security Notes

1. **Change all default passwords** in the secret files before deploying to production
2. Consider using Google Secret Manager for sensitive data
3. Enable Cloud Armor for DDoS protection
4. Set up Cloud IAM roles appropriately
5. Enable GKE security features like Binary Authorization
6. Use private GKE clusters for enhanced security
