# ADR-002: Subscription System for Weather Updates

## Status
Accepted

## Context
A secondary requirement was to allow users to **subscribe to periodic weather updates via email**. Users provide their email address and city of interest. They receive weather updates on a scheduled basis (hourly or daily). The architecture had to support scalability, reliability, and asynchronous email delivery without impacting the responsiveness of the main API service.

## Decision

To fulfill this requirement, the following design decisions were made:

- **Subscription Handling in Main Service:**
  The logic for processing `/subscribe`, `/confirm`, and `/unsubscribe` requests is implemented in the **main weather API service**. This includes email validation, token generation, and updating subscription records in the database.

- **Dedicated Microservice for Email Dispatch:**
  A **separate microservice** is responsible for sending periodic weather update emails. This service runs on a scheduler (e.g., `cron`, internal ticker) and queries the database for active subscriptions. It fetches weather data and sends emails via SMTP.

- **Reason for Separation:**
  The email dispatch service is decoupled to:
  - Prevent **resource contention** with the main API, especially when the number of subscriptions grows.
  - Allow **independent scaling** of the email service without affecting the main API.
  - Avoid **race conditions** in a multi-instance setup (e.g., duplicate email sends when scaling horizontally).
  - Ensure the **main service remains stateless and responsive**, serving user requests efficiently.

- **Email Confirmation Flow:**
  Upon a subscription request, the main service sends a confirmation email asynchronously (using a goroutine). The email includes a secure token link for activating the subscription.

- **Asynchronous Email Sending:**
  All one-time transactional emails (confirmation, unsubscribe) are sent from the main service in a non-blocking manner to ensure fast API response.

- **Scheduler for Periodic Emails:**
  The email microservice uses a built-in scheduler to send updates at configurable intervals (hourly, daily). It pulls active subscriptions and fetches weather data using the same logic as the main API, including API fallback and caching.

- **Database:**
  Subscription data (email, city, status, frequency, token) is stored in PostgreSQL and accessed by both services.

- **SMTP Server:**
  An external SMTP server is used by both the main and email services. SMTP credentials and configs are injected via environment variables.

## Consequences
- Separation of concerns improves system maintainability and fault isolation.
- The main service remains responsive and stateless under load.
- The email service can be scaled horizontally or scheduled independently.
- Risk of duplicated emails is avoided by having only one scheduler running in a controlled, isolated environment.

## Alternatives Considered

- **Monolithic Design:**
  Implementing all subscription and email functionality in a single service was rejected to avoid performance bottlenecks and tight coupling.

- **Third-party Email Services (e.g., Mailgun, SendGrid):**
  Considered but rejected to retain full control and avoid reliance on third-party service limitations or costs.

- **Message Queue (e.g., RabbitMQ, NATS) for coordination:**
  Not used in the current MVP for simplicity, but remains a viable future enhancement to decouple communication between services.

## Related ADRs
- ADR-001: Architecture for Weather API Service
