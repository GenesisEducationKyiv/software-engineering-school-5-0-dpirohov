# syntax=docker/dockerfile:1.3


# Stage 1: Build Scheduler Service
FROM golang:1.24.1-alpine AS backend-builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/scheduler_service ./cmd/scheduler_service

# Stage 2: Final Image
FROM alpine

WORKDIR /app

COPY --from=backend-builder /bin/scheduler_service ./scheduler_service

EXPOSE 8082

CMD ["./scheduler_service"]