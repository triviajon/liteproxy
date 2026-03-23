# Stage 1: Build the Go binary
ARG GO_VERSION=1.26.1
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /app

# Copy the source code into the builder container
COPY processor/go.mod processor/go.sum ./
RUN go mod download
COPY processor .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o liteproxy ./cmd/proxy/main.go

# Stage 2: Create the final image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/liteproxy /liteproxy

# Expose the port that the proxy server will listen on (PROCESSOR_PORT)
EXPOSE ${PROCESSOR_PORT:-8080}

# Run the binary
ENTRYPOINT ["/liteproxy"]