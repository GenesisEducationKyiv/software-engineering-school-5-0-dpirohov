# Test task for genesis se-school #5

[Deployed service on Render.com](https://weatherapi-d4b3.onrender.com/)


### To run localy setup .env file

```env
PORT=8080
DATABASE_URL=postgres://admin:secret@db:5432/mydb
SMTP_LINK_URL=http://0.0.0.0:8080
OPENWEATHER__API_KEY=<openweathermapApiKey>  # Required API key https://openweathermap.org/
WEATHER_API_API_KEY=<https://www.weatherapi.com/>  # fallback API, type any string(required), service will live on OpenWeatherMap as main provider
TOKEN_LIFETIME_MINUTES=15
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=<your_mail@gmail.com>  # Required gmail login
SMTP_PASS=<app_password>  # Required gmail app password
```

Run
```bash
docker compose up
```

This will build from scratch service, and perform initial migration.
Service hosts UI files


# What should be done better?

1. An event bus or message broker with a separate consumer should be implemented to handle email sending. The current implementation is simplified for the sake of this task.

2. A scheduler should be implemented as a separate service for sending scheduled subscriptions (not implemented, as this was outside the scope of the test task).

3. City names are validated using external APIs. If the API returns a 404, the city is considered invalid. This approach could be made more robust and handled separately

4. The unsubscribe button is not implemented in the UI. It should be included in the email content when sending subscription data. The backend endpoint is implemented and works.

5. Tests need significant improvement, especially for database queries. This can be achieved using a seeded SQLite test database.

6. Consider implementing interfaces for services, at this point they are skipped on purpose