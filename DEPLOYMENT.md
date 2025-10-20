# Habit Tracker API - Complete Documentation

A production-ready REST API built with Go for tracking personal habits. Features JWT authentication, MySQL database, and complete AWS Kubernetes deployment.

## 📋 Table of Contents
- [Features](#features)
- [Architecture](#architecture)
- [Data Models](#data-models)
- [API Endpoints](#api-endpoints)
- [Local Development](#local-development)
- [AWS Deployment](#aws-deployment)
- [Security](#security)
- [Maintenance](#maintenance)
- [Troubleshooting](#troubleshooting)

## ✨ Features

- ✅ RESTful API with full CRUD operations
- 🔐 JWT-based authentication
- 🏗️ Clean architecture (handlers → services → repositories)
- 🔒 Security best practices (RBAC, validation, secure headers)
- 📊 MySQL with indexing and foreign keys
- 🐳 Docker support
- ☸️ Kubernetes deployment with auto-scaling
- 📈 AWS EKS ready with monitoring
- 📝 Comprehensive logging
- ✨ Idiomatic Go using standard library

## 🏛️ Architecture

### Project Structure
```
.
├── cmd/api/main.go          # Application entry point
├── internal/
│   ├── config/              # Configuration
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # Auth, logging, security
│   ├── models/              # Data structures
│   ├── repository/          # Database layer
│   └── service/             # Business logic
├── migrations/              # SQL migrations
├── k8s/                    # Kubernetes manifests
├── scripts/                # Deployment scripts
└── docker-compose.yml      # Local development
```

### Layers
1. **Handlers** - HTTP request/response, validation
2. **Services** - Business logic, authorization
3. **Repositories** - Database operations
4. **Models** - Data structures

## 📊 Data Models

### User
- ID, Email, Password Hash, Name
- Profile Image URL, Time Zone, Phone
- Created At

### Habit
- ID, User ID, Name, Description
- Created At

### Log
- ID, Habit ID, Created At, Notes

## 🔌 API Endpoints

### Authentication

**Register**
```bash
POST /api/v1/auth/register
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepass123",
  "time_zone": "America/New_York"
}
```

**Login**
```bash
POST /api/v1/auth/login
{
  "email": "john@example.com",
  "password": "securepass123"
}
# Returns: { "token": "...", "user": {...} }
```

### Users (Protected)
- `GET /api/v1/users/{id}` - Get user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Habits (Protected)
- `POST /api/v1/habits` - Create habit
- `GET /api/v1/habits` - List user's habits
- `GET /api/v1/habits/{id}` - Get habit
- `PUT /api/v1/habits/{id}` - Update habit
- `DELETE /api/v1/habits/{id}` - Delete habit

### Logs (Protected)
- `POST /api/v1/logs` - Create log
- `GET /api/v1/logs?habit_id={id}` - List logs
- `GET /api/v1/logs/{id}` - Get log
- `PUT /api/v1/logs/{id}` - Update log
- `DELETE /api/v1/logs/{id}` - Delete log

All protected endpoints require:
```
Authorization: Bearer {token}
```

## 🚀 Local Development

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- MySQL client (optional)

### Quick Start

```bash
# 1. Clone repository
git clone https://github.com/hayden-erickson/ai-evaluation.git
cd ai-evaluation

# 2. Install dependencies
go mod tidy

# 3. Start services
docker-compose up -d

# 4. Run migrations (Linux/Mac)
chmod +x scripts/run-migrations.sh
DB_PASSWORD=habits_password ./scripts/run-migrations.sh

# Windows PowerShell
$env:DB_PASSWORD="habits_password"
bash scripts/run-migrations.sh

# 5. Test API
curl http://localhost:8080/health
```

### Manual Setup (No Docker)

```bash
# Start MySQL
mysql -u root -p
CREATE DATABASE habits_db;
CREATE USER 'habits_user'@'localhost' IDENTIFIED BY 'habits_password';
GRANT ALL PRIVILEGES ON habits_db.* TO 'habits_user'@'localhost';

# Set environment variables
export PORT=8080
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=habits_user
export DB_PASSWORD=habits_password
export DB_NAME=habits_db
export JWT_SECRET=your-secret-key

# Run migrations
DB_PASSWORD=habits_password ./scripts/run-migrations.sh

# Run application
go run cmd/api/main.go
```

### Testing

```bash
# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123",
    "time_zone": "America/New_York"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}'

# Save the token from response and use it:
TOKEN="your-token-here"

# Create habit
curl -X POST http://localhost:8080/api/v1/habits \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Exercise", "description": "Daily workout"}'
```

## ☁️ AWS Deployment

### Prerequisites
- AWS CLI configured
- kubectl installed
- eksctl installed
- Docker installed
- Helm installed

### Step 1: Create EKS Cluster

```bash
chmod +x scripts/setup-eks.sh
./scripts/setup-eks.sh
```

This creates:
- EKS cluster with managed nodes
- AWS Load Balancer Controller
- Metrics Server for autoscaling
- CloudWatch Container Insights
- EBS CSI Driver

**Time:** ~20-25 minutes

### Step 2: Update Secrets

⚠️ **CRITICAL**: Update production secrets!

```bash
# Generate secrets
openssl rand -base64 32  # JWT secret
openssl rand -base64 24  # DB password

# Edit secrets file
nano k8s/secrets.yaml
# Replace all CHANGE_ME values
```

### Step 3: Deploy Application

```bash
chmod +x scripts/deploy-aws.sh

# Optional: set custom values
export AWS_REGION=us-west-2
export CLUSTER_NAME=habits-cluster
export ECR_REPO_NAME=habits-api

# Deploy
./scripts/deploy-aws.sh
```

This will:
1. Create ECR repository
2. Build & push Docker image
3. Deploy MySQL with persistent storage
4. Run database migrations
5. Deploy API (3 replicas)
6. Configure auto-scaling (3-10 pods)
7. Create LoadBalancer

### Step 4: Configure SSL (Optional)

```bash
# Request certificate
aws acm request-certificate \
  --domain-name api.yourdomain.com \
  --validation-method DNS \
  --region us-west-2

# Update k8s/ingress.yaml with:
# - Certificate ARN
# - Your domain name

# Deploy ingress
kubectl apply -f k8s/ingress.yaml

# Update DNS with ALB hostname
kubectl get ingress -n habits-app
```

### Step 5: Verify Deployment

```bash
# Check pods
kubectl get pods -n habits-app

# Check services
kubectl get svc -n habits-app

# Get LoadBalancer URL
kubectl get service habits-api-service -n habits-app

# View logs
kubectl logs -f deployment/habits-api -n habits-app

# Test API
LOAD_BALANCER_URL=$(kubectl get svc habits-api-service -n habits-app -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
curl http://$LOAD_BALANCER_URL/health
```

## 🔒 Security

### Application Security
- ✅ JWT auth (24h expiration)
- ✅ bcrypt password hashing
- ✅ Input validation
- ✅ SQL injection prevention
- ✅ RBAC (users access only their data)
- ✅ Secure headers (HSTS, CSP, X-Frame-Options)
- ✅ HTTPS enforcement

### Kubernetes Security
- ✅ Secrets for sensitive data
- ✅ Non-root containers
- ✅ Resource limits
- ✅ Health checks

### AWS Security
- ✅ VPC isolation
- ✅ IAM roles
- ✅ Security groups
- ✅ EBS encryption
- ✅ ACM certificates

### Production Checklist
- [ ] Update all secrets in k8s/secrets.yaml
- [ ] Enable database backups
- [ ] Configure proper CORS
- [ ] Set up CloudWatch Alarms
- [ ] Enable AWS WAF
- [ ] Implement rate limiting
- [ ] Configure log retention
- [ ] Set up CI/CD pipeline
- [ ] Enable audit logging
- [ ] Plan disaster recovery

## 🔧 Maintenance

### Update Application

```bash
# Build new version
docker build -t habits-api:v2 .

# Tag and push to ECR
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.us-west-2.amazonaws.com
docker tag habits-api:v2 $AWS_ACCOUNT_ID.dkr.ecr.us-west-2.amazonaws.com/habits-api:v2
docker push $AWS_ACCOUNT_ID.dkr.ecr.us-west-2.amazonaws.com/habits-api:v2

# Update deployment
kubectl set image deployment/habits-api habits-api=$AWS_ACCOUNT_ID.dkr.ecr.us-west-2.amazonaws.com/habits-api:v2 -n habits-app

# Check rollout
kubectl rollout status deployment/habits-api -n habits-app
```

### Scale Application

```bash
# Manual scaling
kubectl scale deployment habits-api --replicas=5 -n habits-app

# Update autoscaling
kubectl edit hpa habits-api-hpa -n habits-app
```

### Database Backup

```bash
# Backup
kubectl exec -n habits-app deployment/mysql -- mysqldump -u habits_user -p$DB_PASSWORD habits_db > backup.sql

# Restore
kubectl exec -i -n habits-app deployment/mysql -- mysql -u habits_user -p$DB_PASSWORD habits_db < backup.sql
```

### View Logs

```bash
# Application logs
kubectl logs -f deployment/habits-api -n habits-app

# All pods
kubectl logs -f -l app=habits-api -n habits-app

# CloudWatch
aws logs tail /aws/containerinsights/habits-cluster/application --follow
```

## 🐛 Troubleshooting

### API Not Responding

```bash
kubectl get pods -n habits-app
kubectl describe pod <pod-name> -n habits-app
kubectl logs <pod-name> -n habits-app
```

### Database Connection Issues

```bash
kubectl get pods -l app=mysql -n habits-app
kubectl exec -it -n habits-app <mysql-pod> -- mysql -u habits_user -p
kubectl get svc mysql-service -n habits-app
```

### Authentication Issues
- Verify JWT_SECRET matches in secrets
- Check token hasn't expired (24h)
- Ensure header format: `Authorization: Bearer <token>`

### Performance Issues
- Check HPA: `kubectl get hpa -n habits-app`
- View metrics: `kubectl top pods -n habits-app`
- Check resource limits: `kubectl describe pod <pod> -n habits-app`

## 📦 Dependencies

```go
go 1.21

require (
    github.com/go-sql-driver/mysql v1.7.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/google/uuid v1.5.0
    github.com/go-playground/validator/v10 v10.16.0
    golang.org/x/crypto v0.17.0
)
```

## 📝 License

MIT License

## 👥 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open Pull Request

## 📧 Support

- GitHub Issues: https://github.com/hayden-erickson/ai-evaluation/issues
- Email: hayden.erickson@example.com

---

Built with ❤️ using Go and Kubernetes
