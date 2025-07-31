# syntax=docker/dockerfile:1.3


# Stage 1: Build Notification Service
FROM golang:1.24.1-alpine AS backend-builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/notification_service ./cmd/notification_service

# Stage 2: Final Image
FROM alpine

WORKDIR /app

COPY --from=backend-builder /bin/notification_service ./notification_service

EXPOSE 8081

CMD ["./notification_service"]