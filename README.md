# Spring Petclinic - Kubernetes Microservices & Observability Edition

## Overview
This repository contains a completely refactored, cloud-native iteration of the monolithic Spring Petclinic application. The project has been deliberately rebuilt to demonstrate a highly decoupled **Kubernetes Microservices Architecture** featuring standalone **Headless Browser Observability Clients** that simulate and track real-world OpenTelemetry data.

## Architecture & Boundaries
The monolithic framework has been decomposed into distinctly isolated containers:

### 1. Spring Boot Microservices
- **`api-gateway` (NGINX)**: The unified edge router. This lightning-fast Alpine container binds explicitly to `NodePort 30080` and dynamically proxies traffic headers securely to the internal microservices without exposing them publicly.
- **`customers-service`**: Manages the `Owner` and `Pet` databases, and renders the central Thymeleaf views.
- **`vets-service`**: Manages the independent `Veterinarians` layer.
- **`visits-service`**: Manages the `Visits` tracking and domains.
- **`demo-db`**: Centralized K8s PostgreSQL instance deployed securely as an internal backing service.

### 2. Standalone Observability Synthesizers
Found inside the `/observability-clients/` directory, these robust clients act exactly like human website visitors, triggering massive bursts of synthetic web traffic tracking perfectly into OpenTelemetry platforms like Datadog, Prometheus, or Jaeger.
- **`synthetic-client` (Golang + Chromedp)**: An advanced Golang application equipped with `Playwright`-style capabilities. It spins up random raw `incognito` Headless Chrome browser tabs within the Kubernetes node to perfectly simulate unique HTTP traffic sessions. Accompanied by a sleek UI.

## Prerequisites
To deploy this architecture, your local machine/cluster must run:
- **Kubernetes CLI (`kubectl`)**
- A live K8s cluster (Minikube, AWS EKS, or AliCloud ACK).

*(Optional for manual source compilation: Java 17 JDK and Docker)*

---

## ⚡️ Zero-Compilation Deployment (Recommended)
You do not need to download Java, compile `.jar` files, or build Docker images to test this environment! All microservices have been pre-compiled and securely published to the GitHub Container Registry (`ghcr.io`).

To instantly spin up the entire cluster over the internet, simply apply the pre-configured deployment manifests. Your Kubernetes node will automatically download the correct Docker images from GitHub and assemble the architecture:

```bash
kubectl apply -f k8s/all-in-one.yml
```

Once executed, jump down to the **Interfacing with the Cluster** section below constraints!

---

## 🛠 Manual Source Deployment Guide
If you intend to modify the Java or Golang source code, you must orchestrate the pipeline locally.

### 1. Compile Backend Artifacts
Because Spring Boot requires compiled `.jar` environments, execute the native Gradle build matrix:
```bash
./gradlew bootJar
```

### 2. Docker Image Provisioning
Build and tag the corresponding images onto your local container layer. 
*Hint: If deploying onto an x86 Linux machine from an Apple Silicon (M1/M2) Mac, use `docker buildx --platform linux/amd64` to prevent fatal JVM architecture translation crashes!*
```bash
docker build -t ghcr.io/<YOUR_GITHUB_USER>/spring-petclinic-api-gateway:latest api-gateway/
docker push ghcr.io/<YOUR_GITHUB_USER>/spring-petclinic-api-gateway:latest

# Repeat for customers, visits, vets, and observability clients...
```
*Note: Ensure your YAML files inside `/k8s/` are updated to match your specific `ghcr.io/<YOUR_USERNAME>/...` paths before pushing into Kubernetes.*

---

## 🌐 Interfacing with the Cluster
Execute `kubectl -n petclinic-demo get pods` to verify that your microservices transitioned into the `1/1 RUNNING` state.

Once stable, access your new ecosystem via the NodePort assigned to the NGINX API Gateway:
- **Main Web Interface:** 🔗 `http://<YOUR_CLUSTER_IP>:30080`
- **Go Headless Engine Dashboard:** 🔗 `http://<YOUR_CLUSTER_IP>:30080/synthetic/`
