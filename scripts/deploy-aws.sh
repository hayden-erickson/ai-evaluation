#!/bin/bash

# Habit Tracker API - AWS EKS Deployment Script
# This script deploys the application to AWS EKS

set -e

# Configuration
AWS_REGION=${AWS_REGION:-us-west-2}
CLUSTER_NAME=${CLUSTER_NAME:-habits-cluster}
ECR_REPO_NAME=${ECR_REPO_NAME:-habits-api}
IMAGE_TAG=${IMAGE_TAG:-latest}

echo "=== Habit Tracker API - AWS Deployment ==="
echo "Region: $AWS_REGION"
echo "Cluster: $CLUSTER_NAME"
echo "ECR Repo: $ECR_REPO_NAME"
echo ""

# Step 1: Get AWS account ID
echo "Step 1: Getting AWS account ID..."
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
echo "AWS Account ID: $AWS_ACCOUNT_ID"

# Step 2: Create ECR repository if it doesn't exist
echo ""
echo "Step 2: Creating ECR repository..."
aws ecr describe-repositories --repository-names $ECR_REPO_NAME --region $AWS_REGION 2>/dev/null || \
    aws ecr create-repository --repository-name $ECR_REPO_NAME --region $AWS_REGION

# Step 3: Authenticate Docker to ECR
echo ""
echo "Step 3: Authenticating Docker to ECR..."
aws ecr get-login-password --region $AWS_REGION | \
    docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com

# Step 4: Build Docker image
echo ""
echo "Step 4: Building Docker image..."
docker build -t $ECR_REPO_NAME:$IMAGE_TAG .

# Step 5: Tag image for ECR
echo ""
echo "Step 5: Tagging image for ECR..."
docker tag $ECR_REPO_NAME:$IMAGE_TAG $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO_NAME:$IMAGE_TAG

# Step 6: Push image to ECR
echo ""
echo "Step 6: Pushing image to ECR..."
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO_NAME:$IMAGE_TAG

# Step 7: Update kubeconfig
echo ""
echo "Step 7: Updating kubeconfig..."
aws eks update-kubeconfig --region $AWS_REGION --name $CLUSTER_NAME

# Step 8: Create namespace
echo ""
echo "Step 8: Creating Kubernetes namespace..."
kubectl apply -f k8s/namespace.yaml

# Step 9: Apply secrets (WARNING: Update secrets.yaml with real values first!)
echo ""
echo "Step 9: Applying secrets..."
echo "WARNING: Ensure you've updated k8s/secrets.yaml with real values!"
read -p "Press enter to continue or Ctrl+C to abort..."
kubectl apply -f k8s/secrets.yaml

# Step 10: Apply ConfigMap
echo ""
echo "Step 10: Applying ConfigMap..."
kubectl apply -f k8s/configmap.yaml

# Step 11: Deploy MySQL
echo ""
echo "Step 11: Deploying MySQL..."
kubectl apply -f k8s/mysql-deployment.yaml

# Wait for MySQL to be ready
echo "Waiting for MySQL to be ready..."
kubectl wait --for=condition=ready pod -l app=mysql -n habits-app --timeout=300s

# Step 12: Run database migrations
echo ""
echo "Step 12: Running database migrations..."
echo "Creating migration job..."
kubectl run mysql-migrations --image=mysql:8.0 --restart=Never -n habits-app -- /bin/sh -c "
  for file in /migrations/*.sql; do
    mysql -h mysql-service -u \$DB_USER -p\$DB_PASSWORD \$DB_NAME < \$file
  done
"

# Step 13: Update API deployment with ECR image
echo ""
echo "Step 13: Updating API deployment configuration..."
sed -i "s|YOUR_ECR_REPO/habits-api:latest|$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO_NAME:$IMAGE_TAG|g" k8s/api-deployment.yaml

# Step 14: Deploy API
echo ""
echo "Step 14: Deploying API..."
kubectl apply -f k8s/api-deployment.yaml

# Step 15: Deploy HPA
echo ""
echo "Step 15: Deploying Horizontal Pod Autoscaler..."
kubectl apply -f k8s/hpa.yaml

# Step 16: Deploy Ingress (optional)
echo ""
echo "Step 16: Deploying Ingress..."
read -p "Do you want to deploy the Ingress? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Update k8s/ingress.yaml with your certificate ARN and domain first!"
    read -p "Press enter when ready..."
    kubectl apply -f k8s/ingress.yaml
fi

# Step 17: Get service URL
echo ""
echo "Step 17: Getting service information..."
echo "Waiting for LoadBalancer to be ready..."
kubectl get service habits-api-service -n habits-app -w &
WATCH_PID=$!
sleep 30
kill $WATCH_PID 2>/dev/null || true

echo ""
echo "=== Deployment Complete ==="
echo ""
echo "To get the API URL, run:"
echo "kubectl get service habits-api-service -n habits-app"
echo ""
echo "To check deployment status:"
echo "kubectl get pods -n habits-app"
echo ""
echo "To view logs:"
echo "kubectl logs -f deployment/habits-api -n habits-app"
