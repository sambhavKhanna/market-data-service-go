# Stage 1: Builder
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for static binaries in Alpine
ARG SERVICE_NAME
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/${SERVICE_NAME} ./cmd/${SERVICE_NAME}

# Stage 2: Runner
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
ARG SERVICE_NAME
COPY --from=builder /usr/local/bin/${SERVICE_NAME} /usr/local/bin/

# Promote the build-arg SERVICE_NAME to a runtime environment variable
ENV SERVICE_NAME=${SERVICE_NAME}

# Default entrypoint, will be overridden in docker-compose for specific services
ENTRYPOINT /usr/local/bin/${SERVICE_NAME}
