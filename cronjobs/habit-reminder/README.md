# Habit Reminder CronJob

This directory contains a Kubernetes CronJob that sends reminders to users based on their habit logging activity.

## Functionality

The CronJob performs the following actions:
- Queries a MySQL database for habit logs over the past 2 days.
- Sends a notification via Twilio to users who meet specific criteria (e.g., no logs in the last 2 days).

## Files

- `main.go`: The Go application that contains the core logic for the CronJob.
- `Dockerfile`: Used to build the Docker image for the Go application.
- `deployment.yaml`: The Kubernetes CronJob manifest.
- `secrets.yaml`: A template for creating Kubernetes Secrets to store sensitive data.
- `configmap.yaml`: The Kubernetes ConfigMap for non-sensitive configuration.
- `rbac.yaml`: Kubernetes RBAC configurations for the service account.

## Deployment Steps

### 1. Prerequisites

- A running Kubernetes cluster.
- `kubectl` configured to interact with your cluster.
- A Docker registry (e.g., Docker Hub, GCR) to store your image.

### 2. Build and Push the Docker Image

1.  Build the Docker image:
    ```sh
    docker build -t your-repo/habit-reminder:latest .
    ```
2.  Push the image to your registry:
    ```sh
    docker push your-repo/habit-reminder:latest
    ```

### 3. Configure Secrets

1.  Encode your secrets in Base64.
    ```sh
    echo -n 'your-secret' | base64
    ```
2.  Update `secrets.yaml` with the encoded values.

### 4. Deploy to Kubernetes

1.  Apply the Kubernetes manifests:
    ```sh
    kubectl apply -f rbac.yaml
    kubectl apply -f configmap.yaml
    kubectl apply -f secrets.yaml
    kubectl apply -f deployment.yaml
    ```

## Environment Variables

The application is configured using the following environment variables:

| Variable              | Description                        |
| --------------------- | ---------------------------------- |
| `DB_USER`             | MySQL database user                |
| `DB_PASSWORD`         | MySQL database password            |
| `DB_HOST`             | MySQL database host                |
| `DB_PORT`             | MySQL database port                |
| `DB_NAME`             | MySQL database name                |
| `TWILIO_ACCOUNT_SID`  | Twilio Account SID                 |
| `TWILIO_AUTH_TOKEN`   | Twilio Auth Token                  |
| `TWILIO_PHONE_NUMBER` | Your Twilio phone number           |
