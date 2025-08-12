# Educational Project for Genesis SE-School #5

A simple weather API application that allows you to:
- Fetch current weather for a selected city
- Subscribe to weather updates
- Unsubscribe from weather updates

[![Go](https://img.shields.io/badge/Go-1.22-blue?logo=go&logoColor=white)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18-61dafb?logo=react&logoColor=white)](https://react.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-red?logo=redis&logoColor=white)](https://redis.io/)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-3.13-ff6600?logo=rabbitmq&logoColor=white)](https://www.rabbitmq.com/)
[![Docker](https://img.shields.io/badge/Docker-Compose-blue?logo=docker&logoColor=white)](https://www.docker.com/)
[![Playwright](https://img.shields.io/badge/Playwright-E2E-45ba4b?logo=playwright&logoColor=white)](https://playwright.dev/)
[![Testify](https://img.shields.io/badge/Testify-Unit_Tests-yellow?logo=go&logoColor=white)](https://github.com/stretchr/testify)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
---

## Tech Stack

- **React** with **Material UI** — simple frontend
- **Go (Gin)** — REST API server
- **PostgreSQL** — main database
- **Redis** — caching
- **GORM** — ORM
- **RabbitMQ** — message bus
- **Scheduler** — background tasks
- **Prometheus** & **Grafana** — metrics
- **OpenSearch** — log collection and tracing
> Monitoring is configured with **alerting** for high `HTTP 500` error rates.

---

## Testing

- **Frontend** — covered by **end-to-end (E2E)** tests using [Playwright](https://playwright.dev/).
- **Backend** — covered by **unit tests** using [Testify](https://github.com/stretchr/testify).

> Current coverage:
> - E2E tests: cover main user flows (fetching weather, subscribing/unsubscribing)
> - Unit tests: cover business logic

---

## CI Status

[![E2E Tests](https://github.com/GenesisEducationKyiv/software-engineering-school-5-0-dpirohov/actions/workflows/run_e2e_tests.yaml/badge.svg)](https://github.com/GenesisEducationKyiv/software-engineering-school-5-0-dpirohov/actions/workflows/run_e2e_tests.yaml)
[![Unit Tests](https://github.com/GenesisEducationKyiv/software-engineering-school-5-0-dpirohov/actions/workflows/run_tests.yaml/badge.svg)](https://github.com/GenesisEducationKyiv/software-engineering-school-5-0-dpirohov/actions/workflows/run_tests.yaml)
[![Linter](https://github.com/GenesisEducationKyiv/software-engineering-school-5-0-dpirohov/actions/workflows/linter.yaml/badge.svg)](https://github.com/GenesisEducationKyiv/software-engineering-school-5-0-dpirohov/actions/workflows/linter.yaml)
--- 
## Service Architecture

![Architecture](docs/architecture/microservice-architecture.svg)

---

## Screenshots

<details>
  <summary>Main Page</summary>

  ![Main Page](docs/screenshots/MainPage.jpg)

</details>

<details>
  <summary>Grafana Dashboard</summary>

  ![Grafana](docs/screenshots/grafana.jpg)

</details>

<details>
  <summary>OpenSearch Tracing</summary>

  ![OpenSearch](docs/screenshots/Opensearch.jpg)

</details>

---

## Local Setup

To run locally, create a `.env` file for each service in the repository root.

**`env.api_service`**

```env
ENV=DOCKER
HOST=http://0.0.0.0
PORT=8080
DB_URL=postgres://admin:secret@db:5432/mydb
BROKER_URL=amqp://admin:admin@rabbitmq:5672/
APP_URL=http://localhost:8080

# PROVIDERS KEYS AND ENDPOINTS
OPENWEATHER_API_KEY=<YOUR OPENWEATHER API KEY>
OPENWEATHER_API_ENDPOINT=http://api.openweathermap.org/data/2.5/weather

WEATHER_API_API_KEY=<YOUR WEATHER_API API KEY>
WEATHER_API_API_ENDPOINT=http://api.weatherapi.com/v1/current.json

TOKEN_LIFETIME_MINUTES=15

REDIS_URL=redis:6379
REDIS_PWD="secret"
CACHE_TTL=5m
LOCK_TTL=3s
LOCK_RETRY_DUR=100ms
LOCK_MAX_WAIT=3s
```
---
**`env.notification_service`**
```env
ENV=DOCKER
PORT=8081

BROKER_URL=amqp://admin:admin@rabbitmq:5672/
APP_URL=http://localhost:8080

# SMTP CREDENTIALS
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=<EMAIL>
SMTP_PASS=<APP PASSWORD>
```


Run
```bash
docker compose --profile monitoring --profile tracing up --build
```

This will:
- Build all services from scratch
- Perform initial database migrations
- Serve the UI