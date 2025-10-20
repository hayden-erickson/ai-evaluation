#!/bin/bash

# Habit Tracker API - AWS EKS Cluster Setup Script
# This script creates an EKS cluster with all necessary resources

set -e

# Configuration
AWS_REGION=${AWS_REGION:-us-west-2}
CLUSTER_NAME=${CLUSTER_NAME:-habits-cluster}
NODE_TYPE=${NODE_TYPE:-t3.medium}
MIN_NODES=${MIN_NODES:-2}
MAX_NODES=${MAX_NODES:-5}
DESIRED_NODES=${DESIRED_NODES:-3}

echo "=== Creating AWS EKS Cluster ==="
echo "Region: $AWS_REGION"
echo "Cluster Name: $CLUSTER_NAME"
echo "Node Type: $NODE_TYPE"
echo "Min Nodes: $MIN_NODES"
echo "Max Nodes: $MAX_NODES"
echo "Desired Nodes: $DESIRED_NODES"
echo ""

# Check if eksctl is installed
if ! command -v eksctl &> /dev/null; then
    echo "eksctl is not installed. Please install it first:"
    echo "https://eksctl.io/installation/"
    exit 1
fi

# Step 1: Create EKS cluster
echo "Step 1: Creating EKS cluster (this may take 15-20 minutes)..."
eksctl create cluster \
    --name $CLUSTER_NAME \
    --region $AWS_REGION \
    --nodegroup-name standard-workers \
    --node-type $NODE_TYPE \
    --nodes $DESIRED_NODES \
    --nodes-min $MIN_NODES \
    --nodes-max $MAX_NODES \
    --managed \
    --with-oidc \
    --ssh-access \
    --ssh-public-key ~/.ssh/id_rsa.pub \
    --tags Environment=production,Project=habits-api

# Step 2: Update kubeconfig
echo ""
echo "Step 2: Updating kubeconfig..."
aws eks update-kubeconfig --region $AWS_REGION --name $CLUSTER_NAME

# Step 3: Install AWS Load Balancer Controller
echo ""
echo "Step 3: Installing AWS Load Balancer Controller..."

# Create IAM policy for ALB controller
curl -o iam-policy.json https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.4.4/docs/install/iam_policy.json

aws iam create-policy \
    --policy-name AWSLoadBalancerControllerIAMPolicy \
    --policy-document file://iam-policy.json || true

# Create IAM service account
eksctl create iamserviceaccount \
    --cluster=$CLUSTER_NAME \
    --namespace=kube-system \
    --name=aws-load-balancer-controller \
    --attach-policy-arn=arn:aws:iam::$(aws sts get-caller-identity --query Account --output text):policy/AWSLoadBalancerControllerIAMPolicy \
    --approve \
    --region=$AWS_REGION || true

# Install ALB controller using Helm
helm repo add eks https://aws.github.io/eks-charts
helm repo update

helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
    -n kube-system \
    --set clusterName=$CLUSTER_NAME \
    --set serviceAccount.create=false \
    --set serviceAccount.name=aws-load-balancer-controller \
    --set region=$AWS_REGION \
    --set vpcId=$(aws eks describe-cluster --name $CLUSTER_NAME --query "cluster.resourcesVpcConfig.vpcId" --output text --region $AWS_REGION)

# Step 4: Install Metrics Server
echo ""
echo "Step 4: Installing Metrics Server..."
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Step 5: Install CloudWatch Container Insights
echo ""
echo "Step 5: Installing CloudWatch Container Insights..."
ClusterName=$CLUSTER_NAME
RegionName=$AWS_REGION
FluentBitHttpPort='2020'
FluentBitReadFromHead='Off'
[[ ${FluentBitReadFromHead} = 'On' ]] && FluentBitReadFromTail='Off'|| FluentBitReadFromTail='On'
[[ -z ${FluentBitHttpPort} ]] && FluentBitHttpServer='Off' || FluentBitHttpServer='On'

curl https://raw.githubusercontent.com/aws-samples/amazon-cloudwatch-container-insights/latest/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/quickstart/cwagent-fluent-bit-quickstart.yaml | sed "s/{{cluster_name}}/${ClusterName}/;s/{{region_name}}/${RegionName}/;s/{{http_server_toggle}}/${FluentBitHttpServer}/;s/{{http_server_port}}/${FluentBitHttpPort}/;s/{{read_from_head}}/${FluentBitReadFromHead}/;s/{{read_from_tail}}/${FluentBitReadFromTail}/" | kubectl apply -f -

# Step 6: Configure EBS CSI Driver for persistent volumes
echo ""
echo "Step 6: Installing EBS CSI Driver..."
eksctl create iamserviceaccount \
    --name ebs-csi-controller-sa \
    --namespace kube-system \
    --cluster $CLUSTER_NAME \
    --attach-policy-arn arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy \
    --approve \
    --role-only \
    --role-name AmazonEKS_EBS_CSI_DriverRole \
    --region=$AWS_REGION || true

eksctl create addon \
    --name aws-ebs-csi-driver \
    --cluster $CLUSTER_NAME \
    --service-account-role-arn arn:aws:iam::$(aws sts get-caller-identity --query Account --output text):role/AmazonEKS_EBS_CSI_DriverRole \
    --region=$AWS_REGION --force || true

echo ""
echo "=== EKS Cluster Setup Complete ==="
echo ""
echo "Cluster Name: $CLUSTER_NAME"
echo "Region: $AWS_REGION"
echo ""
echo "To verify the cluster:"
echo "kubectl get nodes"
echo ""
echo "Next steps:"
echo "1. Update k8s/secrets.yaml with production values"
echo "2. Run ./scripts/deploy-aws.sh to deploy the application"
