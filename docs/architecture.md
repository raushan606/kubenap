# 🧠 KubeNap Architecture

KubeNap is a Kubernetes-native controller that automatically suspends and resumes deployments based on real-time HTTP activity, without introducing latency during normal operation. This document outlines the architectural design, data flow, and integration touchpoints.

---

## 🎯 Design Goals

* Avoid inline proxies or latency during normal operation
* Use only standard Kubernetes APIs (Service, EndpointSlice, Deployment)
* Support multiple reverse proxies (Traefik, NGINX, Istio, etc.) via plugin integration
* Provide clean wake-resume-redirection cycle
* Be opt-in via annotations

---

## 📐 Key Components

| Component          | Description                                                                    |
| ------------------ | ------------------------------------------------------------------------------ |
| Deployment Monitor | Watches annotated deployments, determines idleness via pluggable metric source |
| Wake Handler       | HTTP handler for resume-on-request functionality                               |
| Proxy Plugin       | Optional proxy-specific code to forward traffic to KubeNap when app is asleep  |
| Metrics Source     | Pluggable interface to query traffic data (Prometheus, InfluxDB, etc.)         |

---

## 🔁 Lifecycle Flow (Text-Based)

```text
1. ┌────────────┐
   │  Ingress   │  ←─ receives user HTTP request
   └────┬───────┘
        │
        ▼
2. ┌────────────────┐
   │ Kubernetes     │
   │ Service        │  → attempts to route to app pods
   └────┬───────────┘
        │
        ▼
3. ┌────────────────────┐
   │ No Active Endpoints│  → app is scaled to 0 replicas
   └────┬───────────────┘
        │
        ▼
4. ┌───────────────────────────────────────────┐
   │ Ingress returns 502 or triggers wake path │
   └────┬──────────────────────────────────────┘
        ▼
5. ┌────────────────────────────┐
   │ KubeNap receives wake call │ ← via special path or 502 handler
   └────┬───────────────────────┘
        ▼
6. ┌──────────────────────────────────────┐
   │ KubeNap looks up associated          │
   │ Deployment via annotations           │
   │ (e.g. kubenap/service) │
   └────┬─────────────────────────────────┘
        ▼
7. ┌───────────────────────────────┐
   │ KubeNap scales replicas > 0   │
   │ (from stored replicaCount)    │
   └────┬──────────────────────────┘
        ▼
8. ┌────────────────────────────┐
   │ Wait for Deployment Ready  │ ← poll Pod readiness or use informers
   └────┬───────────────────────┘
        ▼
9. ┌──────────────────────────────┐
   │ KubeNap sends 307 Redirect   │
   │ to original request URL      │
   └────┬─────────────────────────┘
        ▼
10.┌─────────────────────┐
   │ Client reissues     │ ← browser/curl/client follows redirect
   │ request to Ingress  │
   └────┬────────────────┘
        ▼
11.┌────────────────────────────┐
   │ Ingress → App Service → Pod│
   │ (now app has endpoints)    │
   └────────────────────────────┘
```

---

## 🧩 Reverse Proxy Integration

KubeNap uses a minimal plugin interface per ingress type:

### Traefik:

* Use `errorpage` middleware to rewrite 502 to KubeNap's `/wake` endpoint

### NGINX:

* Use `error_page 502 = @handler` to invoke KubeNap handler via internal redirect

### Istio/Envoy:

* Use WASM or Lua filter to intercept failure response and retry via `/wake`

---

## 📌 Annotations

```yaml
metadata:
  labels:
    kubenap/enabled: "true"
  annotations:
    kubenap/service: "myapp-svc"
    kubenap/ingress: "myapp-ing"
    kubenap/idleAfter: "10m"
    kubenap/replicaCount: "1"
    kubenap/proxy: "nginx" # or traefik, istio
```

---

## 📈 Metrics Design

KubeNap supports plug-and-play metric backends:

| Backend    | Method                             |
| ---------- | ---------------------------------- |
| Prometheus | Query `http_requests_total`        |
| InfluxDB   | Query Flux expressions             |
| Datadog    | REST API to check app traffic      |
| Logs       | Tail NGINX/Envoy logs if necessary |

Backends implement a `TrafficMetricsSource` interface returning the last-seen activity timestamp.

---

## ✅ Summary

* ⚡ KubeNap does not sit inline with traffic — zero latency overhead
* 🧠 Relies on annotations to opt-in and configure behavior
* 📊 Traffic-aware using pluggable metric sources
* 🔁 307-based redirect ensures seamless client experience
* 🔌 Designed for pluggable integration with any ingress
