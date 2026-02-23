# 🛒 E-Commerce Platform

A production-ready microservices-based e-commerce backend built with Go. Designed with distributed systems principles including Saga choreography, Outbox pattern, event-driven architecture, and centralized authentication via Keycloak.

---

## 📋 Table of Contents

- [Architecture Overview](#architecture-overview)
- [Tech Stack](#tech-stack)
- [Services](#services)
- [Getting Started](#getting-started)
- [Infrastructure](#infrastructure)
- [API Documentation](#api-documentation)
- [Design Patterns](#design-patterns)
- [Project Structure](#project-structure)
- [CI/CD](#cicd)
- [Observability](#observability)

---

## 🏗️ Architecture Overview

```
Client
  └── Nginx (Reverse Proxy + Auth)
        ├── auth_request ──► Keycloak (Token Validation)
        ├── /api/v1/users    ──► user-service        :8081
        ├── /api/v1/products ──► product-service     :8082
        ├── /api/v1/orders   ──► order-service       :8083
        ├── /api/v1/payments ──► payment-service     :8084
        ├── /api/v1/search   ──► search-service      :8086
        └── /realms          ──► keycloak            :8180

Async Communication (Kafka)
  product-service ──► Kafka ──► search-service       (index sync)
  order-service   ──► Kafka ──► product-service      (stock saga)
  order-service   ──► Kafka ──► payment-service      (payment saga)
  payment-service ──► Kafka ──► order-service        (saga reply)
  order-service   ──► Kafka ──► notification-service (notify user)
```

---

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22+ |
| HTTP Framework | Fiber v2 |
| ORM | GORM |
| Database | PostgreSQL 16 |
| Search Engine | Elasticsearch 8 |
| Cache | Redis 7 |
| Message Broker | Apache Kafka |
| Authentication | Keycloak 24 (OAuth2 / OpenID Connect) |
| Reverse Proxy | Nginx |
| Containerization | Docker + Docker Compose |
| Orchestration | Kubernetes |
| Monitoring | Prometheus + Grafana |
| Tracing | Jaeger (OpenTelemetry) |
| CI/CD | GitHub Actions + Terraform |

---

## 📦 Services

| Service | Port | Responsibility |
|---|---|---|
| user-service | 8081 | User profile, address management, Keycloak token validation |
| product-service | 8082 | Product catalog, categories, stock management, Kafka events |
| order-service | 8083 | Order lifecycle, Saga choreography, compensating transactions |
| payment-service | 8084 | Payment processing, refunds, idempotency |
| notification-service | 8085 | Email/SMS delivery, Kafka consumer only |
| search-service | 8086 | Full-text search, autocomplete, Elasticsearch index sync |
| Nginx | 80 / 443 | Reverse proxy, Keycloak auth_request, routing |
| Keycloak | 8180 | Identity provider, JWT, OAuth2, user management |

---

## 🚀 Getting Started

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Git](https://git-scm.com/)

### 1. Clone the repository

```bash
git clone https://github.com/yasinbozat/ecommerce-platform.git
cd ecommerce-platform
```

### 2. Start the infrastructure

```bash
cd infra
docker-compose up -d
```

> ⚠️ Wait ~60 seconds for Keycloak to fully initialize before making requests.

### 3. Verify all services are healthy

```bash
docker-compose ps
```

### 4. Access the tooling

| Tool | URL | Credentials |
|---|---|---|
| Keycloak Admin | http://localhost:8180 | admin / admin123 |
| Kafdrop (Kafka UI) | http://localhost:9000 | — |
| Kibana | http://localhost:5601 | — |
| Grafana | http://localhost:3000 | admin / grafana123 |
| Jaeger UI | http://localhost:16686 | — |
| pgAdmin | http://localhost:5050 | admin@ecommerce.com / pgadmin123 |
| Prometheus | http://localhost:9090 | — |

### 5. Run a service locally

```bash
cd services/user-service
cp .env.example .env
go run cmd/main.go
```

---

## 🏛️ Infrastructure

### Databases (PostgreSQL)

Each service has its own isolated database following the Database-per-Service pattern:

| Database | Service |
|---|---|
| users_db | user-service |
| products_db | product-service |
| orders_db | order-service |
| payments_db | payment-service |
| notifications_db | notification-service |
| keycloak_db | Keycloak |

### Kafka Topics

| Topic | Producer | Consumer |
|---|---|---|
| product.created | product-service | search-service |
| product.updated | product-service | search-service |
| product.deleted | product-service | search-service |
| stock.reserve | order-service | product-service |
| stock.reserved | product-service | order-service |
| stock.reserve.failed | product-service | order-service |
| stock.release | order-service | product-service |
| payment.process | order-service | payment-service |
| payment.completed | payment-service | order-service |
| payment.failed | payment-service | order-service |
| order.confirmed | order-service | notification-service |
| order.cancelled | order-service | notification-service |
| notification.send | any service | notification-service |

---

## 📡 API Documentation

### Authentication

All protected endpoints require a valid Keycloak JWT access token:

```bash
# Get access token from Keycloak
curl -X POST http://localhost:8180/realms/ecommerce/protocol/openid-connect/token \
  -d "client_id=ecommerce-app" \
  -d "grant_type=password" \
  -d "username=user@example.com" \
  -d "password=yourpassword"

# Use token in requests
curl http://localhost/api/v1/users/me \
  -H "Authorization: Bearer <access_token>"
```

### Endpoint Overview

| Method | Path | Auth | Service |
|---|---|---|---|
| GET | /api/v1/users/me | JWT | user-service |
| PUT | /api/v1/users/me | JWT | user-service |
| GET | /api/v1/users/me/addresses | JWT | user-service |
| POST | /api/v1/users/me/addresses | JWT | user-service |
| PUT | /api/v1/users/me/addresses/:id | JWT | user-service |
| DELETE | /api/v1/users/me/addresses/:id | JWT | user-service |
| GET | /api/v1/products | Public | product-service |
| GET | /api/v1/products/:id | Public | product-service |
| POST | /api/v1/products | Admin JWT | product-service |
| PUT | /api/v1/products/:id | Admin JWT | product-service |
| DELETE | /api/v1/products/:id | Admin JWT | product-service |
| GET | /api/v1/categories | Public | product-service |
| GET | /api/v1/search/products | Public | search-service |
| GET | /api/v1/search/products/autocomplete | Public | search-service |
| GET | /api/v1/search/products/similar/:id | Public | search-service |
| POST | /api/v1/orders | JWT | order-service |
| GET | /api/v1/orders | JWT | order-service |
| GET | /api/v1/orders/:id | JWT | order-service |
| POST | /api/v1/orders/:id/cancel | JWT | order-service |
| GET | /api/v1/payments/:id | JWT | payment-service |

---

## 🎯 Design Patterns

### Outbox Pattern
Guarantees atomic database write + Kafka event publish. Both the business entity and an `outbox_event` record are written in the same DB transaction. A background poller reads unpublished events and sends them to Kafka, preventing event loss even if Kafka is temporarily unavailable.

### Saga Choreography
Distributed transaction management for the order flow without a central orchestrator. Each service reacts to events and publishes results. Compensating transactions handle failures at every step.

```
order-service     →  stock.reserve
product-service   →  stock.reserved        (or stock.reserve.failed)
order-service     →  payment.process
payment-service   →  payment.completed     (or payment.failed)
order-service     →  order.confirmed
notification-service sends confirmation
```

### Repository Pattern
All database access is behind interfaces. Services depend on abstractions, not implementations. Enables clean unit testing with mock repositories without a real database.

### CQRS (Partial)
Read operations use `search-service` (Elasticsearch). Write operations go through `product-service` (PostgreSQL). Kafka keeps both stores eventually consistent.

### Circuit Breaker
Downstream service failures are isolated. Nginx handles upstream unavailability with proper error responses instead of cascading failures.

---

## 📁 Project Structure

```
ecommerce-platform/
├── services/
│   ├── user-service/
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── domain/         # Entities, DTOs, error types
│   │   │   ├── repository/     # Interfaces + PostgreSQL implementations
│   │   │   ├── service/        # Business logic
│   │   │   ├── handler/        # Fiber HTTP handlers
│   │   │   ├── middleware/      # Keycloak, role-based auth middleware
│   │   │   └── config/         # Config loading, database connection
│   │   ├── migrations/         # golang-migrate SQL files
│   │   ├── .env.example
│   │   └── Dockerfile
│   ├── product-service/        # + kafka/ (outbox poller)
│   ├── order-service/          # + saga/ (state machine, compensator)
│   ├── payment-service/        # + provider/ (payment provider interface)
│   ├── notification-service/   # + sender/ (email, SMS)
│   └── search-service/         # + elasticsearch/ + cache/
├── shared/
│   └── pkg/
│       ├── logger/             # Zap structured logger wrapper
│       ├── errors/             # Common error types
│       ├── validator/          # Input validation helpers
│       └── kafka/              # Kafka producer/consumer wrappers
├── infra/
│   ├── docker-compose.yml
│   ├── nginx/                  # nginx.conf + routing config
│   ├── keycloak/               # Realm export for auto-import
│   ├── prometheus/             # Scrape config
│   ├── grafana/                # Dashboard provisioning
│   ├── k8s/                    # Kubernetes manifests
│   └── terraform/              # Infrastructure as Code
└── docs/
    └── openapi/                # Swagger YAML specs per service
```

---

## 🔄 CI/CD

GitHub Actions pipelines are path-filtered so only the changed service's pipeline triggers on each push.

```
.github/workflows/
├── user-service.yml           # Triggers on: services/user-service/**, shared/**
├── product-service.yml
├── order-service.yml
├── payment-service.yml
├── notification-service.yml
└── search-service.yml
```

**Each pipeline runs:**
1. `go vet` + `golangci-lint` — static analysis
2. `go test ./... -race -cover` — unit + integration tests
3. `docker build` — verify container builds
4. Push image to registry (main branch only)
5. Rolling deploy to Kubernetes

---

## 📊 Observability

**Logs**
Structured JSON logging via `go.uber.org/zap`. Every log entry includes `service_name`, `request_id`, `user_id`, `level`, and `timestamp`. Collected centrally via Loki or ELK.

**Metrics**
Prometheus scrapes `/metrics` on every service. Pre-built Grafana dashboards cover HTTP latency, error rates, Kafka consumer lag, and DB query duration.

**Tracing**
OpenTelemetry traces sent to Jaeger via OTLP HTTP (`http://jaeger:4318/v1/traces`). Full distributed request flow is visible across all services in the Jaeger UI at `http://localhost:16686`.

---

## 📄 License

This project is licensed under the MIT License.
