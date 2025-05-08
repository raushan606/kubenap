# 📦 KubeNap

KubeNap is a general-purpose Kubernetes controller that automatically suspends and resumes Deployments based on real
HTTP traffic. It's designed to save compute resources by scaling apps to zero when they're idle — and waking them on
demand when a user accesses them.

🚀 Key Features

⏱️ Automatically detects idle workloads via ingress traffic

🔄 Scales deployments to zero and back to full replicas

🧠 Ingress-aware, plugin-based architecture (Traefik, NGINX, supported)

🔁 Clean HTTP 307 redirect after waking the app

📊 Prometheus metrics for full observability

### Normal Request flow

```
Ingress (Traefik, NGINX, etc.)
   └───► App Service ───► App Deployment
```
### When App is suspended

```
Ingress
  └───► App Service → (points to KubeNap)
                    └───► /wake
                             └───► Scale up + readiness wait
                             └───► HTTP 307 to original path
```

## 🔧 Components

| Component             | Role                                                              |
|-----------------------|-------------------------------------------------------------------|
| `kubenap-controller`  | Core logic: watches deployments, monitors traffic, handles resume |
| Proxy Plugin          | Platform-specific plugin for NGINX, Traefik, Istio, etc.          |
| Wake Endpoint Handler | Receives HTTP requests to resume apps                             |

---

## 🔍 Core Features

### 1. Deployment Monitor

- Watches all deployments with `kubenap/enabled: true`
- Uses pluggable `ProxyMetricsSource` to track traffic per service
- Scales to 0 if no traffic after `idleAfter`

### 2. Wake Handler

- HTTP server at `/wake?original=/foo/bar`
- Looks up app by ingress + service annotations
- If suspended:
    - Scales it up to `replicaCount`
    - Waits until deployment is ready
    - Sends HTTP `307 Temporary Redirect` to original path

### 3. Endpoint Substitution

- While app is suspended, the controller:
    - Registers as a Service for the app
    - Creates an EndpointSlice pointing to itself
- Reverse proxy routes requests to controller as if the app were up

## 📌 Annotations

Deployments must opt in using these annotations:
```yaml
metadata:
  labels:
    kubenap/enabled: "true"
annotations:
  kubenap/service: "myapp-svc"
  kubenap/ingress: "myapp-ing"
  kubenap/idleAfter: "10m"
  kubenap/replicaCount: "1"
  kubenap/proxy: "nginx" # or traefik
```

## 📡 Proxy Plugin Responsibilities

Each plugin/middleware should:

- Detect if app is suspended (e.g. via 502)
- Forward request to `/wake?original=...`
- Follow 307 redirect to retry

### Traefik Plugin Example

- Use service error handler middleware to redirect 502s
- Call KubeNap controller with original path

### NGINX Plugin Example

```nginx
error_page 502 = @wake_handler;

location @wake_handler {
    proxy_pass http://kubenap-controller/wake?original=$request_uri;
}
```

---

## 🔐 Security

- RBAC for scaling and reading deployments
- Optional: require signed token in `/wake` query param
- Optional: mTLS between plugin and controller

---

## 📈 Metrics (Pluggable Backends)

KubeNap will support pluggable traffic metric providers to track activity and determine when to suspend deployments.

### ✅ Supported Providers

| Provider   | Description                                          |
|------------|------------------------------------------------------|
| Prometheus | Default via HTTP query to `http_requests_total`      |
| Grafana    | Supported if backed by Prometheus, Loki, or Cortex   |
| InfluxDB   | Via Flux or InfluxQL API                             |
| Datadog    | Via REST API or custom agent integration             |
| Custom     | Tail logs or emit events to REST/Redis-based store   |

For now, It only support **Prometheus**. You can configure the provider via a Kubernetes `ConfigMap` or environment 
variables:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: kubenap-config
data:
  metricsProvider: "prometheus"
  prometheusURL: "http://prometheus.monitoring.svc.cluster.local:9090"
```


## 📊 Metrics (Prometheus)

| Metric Name                       | Description                        |
|-----------------------------------|------------------------------------|
| `kubenap_apps_total`              | Count of suspended and active apps |
| `kubenap_resume_requests_total`   | Number of resume attempts          |
| `kubenap_resume_duration_seconds` | Histogram of resume time           |

---

## 📦 Directory Structure

```
kubenap/
├── cmd/                   # Entry point
├── controller/            # Core logic
├── plugins/               # NGINX, Traefik, etc. integrations
├── internal/
│   ├── kube/              # Kubernetes client helpers
│   ├── metrics/           # Prometheus exporter
│   ├── proxyinterfaces/   # Proxy-specific metrics interfaces
│   └── wakemgr/           # Wake/resume logic
```

---

## ✅ Summary

| Design Principle          | Description                                           |
|---------------------------|-------------------------------------------------------|
| General core controller   | Handles suspend/resume logic for annotated apps       |
| Plugin-based ingress      | Lightweight proxy plugins forward to wake endpoint    |
| Compatible with any proxy | Design can integrate with NGINX, Traefik, Istio, etc. |
| 307-based resume handoff  | Clean redirect after deployment is ready              |

🧪 Development
```cmd
go run ./cmd/kubenap
```


📄 License

This project is licensed under the Apache 2.0 License

