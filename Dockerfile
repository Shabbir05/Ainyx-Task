# -----------------------------------------------------------------------
# Stage 1: Builder
# -----------------------------------------------------------------------
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /app/server ./cmd/server

# -----------------------------------------------------------------------
# Stage 2: Final minimal image
# -----------------------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot AS final

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/server /app/server

# Copy CA certs (needed for TLS to external services)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 3000

ENTRYPOINT ["/app/server"]
