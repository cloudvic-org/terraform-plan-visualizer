# Build stage
FROM golang:1.25.3-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o terraform-plan-visualizer .

# Final stage
FROM alpine:3.18

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /workspace

# Copy binary from builder stage
COPY --from=builder /app/terraform-plan-visualizer /usr/local/bin/terraform-plan-visualizer

# Make binary executable
RUN chmod +x /usr/local/bin/terraform-plan-visualizer

# Switch to non-root user
USER appuser

# Set default command
ENTRYPOINT ["terraform-plan-visualizer"]

# Default arguments (can be overridden)
CMD ["-h"]
