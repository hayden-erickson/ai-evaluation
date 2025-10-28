# Habit Tracker CronJob

This project contains a Kubernetes CronJob that sends notifications to users of a habit tracking application based on their recent activity.

## Prerequisites

- Docker
- kubectl
- A Kubernetes cluster
- A Docker registry (e.g., Docker Hub, GCR)

## Deployment Steps

### 1. Build and Push the Docker Image

1.  Navigate to the `habit-tracker-cronjob` directory.
2.  Build the Docker image:
    ```sh
    docker build -t your-docker-repo/habit-tracker:latest .
    ```
3.  Push the image to your Docker registry:
    ```sh
    docker push your-docker-repo/habit-tracker:latest
    ```

### 2. Configure Kubernetes Secrets

1.  Encode your secrets in base64. You can use the following command:
    ```sh
    echo -n 'your-secret-value' | base64
    ```
2.  Update `secrets.yaml` with your base64-encoded secrets.

### 3. Deploy to Kubernetes

1.  Apply the ConfigMap, Secret, and CronJob manifests:
    ```sh
    kubectl apply -f configmap.yaml
    kubectl apply -f secrets.yaml
    kubectl apply -f cronjob.yaml
    ```

### 4. Verify the CronJob

1.  Check the status of the CronJob:
    ```sh
    kubectl get cronjob habit-tracker-cronjob
    ```
2.  You can also view the logs of the job's pods to ensure it's running correctly:
    ```sh
    kubectl logs -f <pod-name>
    ```
