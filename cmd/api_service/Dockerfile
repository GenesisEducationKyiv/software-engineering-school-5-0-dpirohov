# syntax=docker/dockerfile:1.3

# Stage 1: Build Frontend
FROM node:18 AS frontend-builder

RUN curl -fsSL https://bun.sh/install | bash && \
    mv /root/.bun/bin/bun /usr/local/bin/

WORKDIR /frontend

COPY ./../../frontend/ ./

RUN bun install

RUN npm run build


# Stage 2: Build Backend
FROM golang:1.24.1-alpine AS backend-builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY ../../ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api_service ./cmd/api_service

# Stage 3: Final Image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=backend-builder /bin/api_service ./api_service
COPY --from=frontend-builder /web ./web

EXPOSE 8080

CMD ["./api_service"]