# ADR-003: Frontend Stack and Integration Strategy

## Status
Accepted

## Context
A user-facing interface was required to allow city-based weather queries and to support subscription flow (subscribe, confirm, unsubscribe). The frontend needed to be simple, fast to develop, and easily integrable with the backend API.

## Decision

The following choices were made for the frontend:

- **React + Material UI Stack:**
  - **React** was chosen as the frontend framework due to its wide adoption, component-driven structure, and ease of integration with REST APIs.
  - **Material UI** was selected as the UI library for rapid development, responsive design out of the box, and polished prebuilt components that follow Google's Material Design system.

- **Frontend Served by Main API Service:**
  - For simplicity, the built frontend is statically served by the **main Go API service**.
  - This eliminates the need for a separate web server (like Nginx) or deploying a second container in the MVP stage.
  - It ensures tight coupling between frontend and backend versions during development and testing.

## Consequences
- Using React + MUI allows fast prototyping and a clean user experience with minimal styling effort.
- Serving the frontend from the backend simplifies deployment, especially in early development or small-scale deployments.
- In future iterations, the frontend can be decoupled and deployed independently (e.g., via CDN or static hosting).

## Alternatives Considered

- **Vue / Svelte / Angular**:
  - These frameworks were considered but rejected due to personal familiarity with React and existing ecosystem support.

- **Hosting frontend separately (e.g., Netlify, Vercel)**:
  - Initially considered for cleaner separation, but deprioritized to avoid additional complexity during MVP development.

## Related ADRs
- ADR-001: Architecture for Weather API Service
- ADR-002: Subscription System for Weather Updates
