# ADR-001: Architecture for Weather API Service

## Status
Accepted

## Context
The task was to implement a simple Weather API service where a user can input a city name and receive current weather information. The system needed to be responsive, reliable, and reasonably extensible.

## Decision
The architecture was designed with the following key components and technology choices:

- **Backend Language:** Go (Golang)
  Chosen for its performance, strong concurrency support, and easy deployment.

- **Web Framework:** [Gin](https://gin-gonic.com/)
  Gin is a lightweight, high-performance web framework for Go. It offers excellent routing and middleware support, making it ideal for building REST APIs.

- **ORM:** [GORM](https://gorm.io/)
  GORM was used to interact with the database due to its mature ecosystem and compatibility with PostgreSQL.

- **Database:** PostgreSQL
  Selected for its reliability, ACID compliance, and native JSON support, which can be beneficial for storing API responses.

- **Weather API Integration:**
  The application integrates with two external weather APIs:
  - **Primary Provider:** Main weather API for fetching weather data. Weather api is selected as main provider (https://api.weatherapi.com/) for its simple api, low latency, and wide weather information support
  - **Fallback Provider:** Used when the main API fails or returns an error, ensuring high availability and resilience. Open weather map (https://openweathermap.org/) is chosen, as it's as stable as main provider and provides similar api structure.

- **Caching Strategy:**
  To reduce external API calls and improve performance, Redis is used as an in-memory cache layer. Weather data retrieved from external APIs is cached in Redis with a time-to-live (TTL) policy (e.g., N minutes). Subsequent requests for the same city will first check the cache, and only if the data is missing or expired, a fresh API call is made. This approach ensures fast response times and reduces dependency on third-party API rate limits.


## Consequences
- The service is highly maintainable and easily extensible.
- Using Gin and GORM reduces boilerplate and speeds up development.
- Implementing two API providers enhances fault tolerance and service reliability.

## Alternatives Considered
- **Echo or Fiber as a web framework** — Gin was preferred due to better documentation and community support.
- **Raw SQL instead of GORM** — GORM improves productivity and maintainability for this scale of application.

## Related ADRs
None at this time.
