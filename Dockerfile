# Multi-stage build for mediagui
ARG USER_ID=1000
ARG GROUP_ID=1000

# Stage 1: Build UI
FROM node:20-alpine AS ui-builder

WORKDIR /app/ui

# Copy package files first for better caching
COPY ui/package*.json ./
RUN npm ci --only=production

# Copy UI source and build
COPY ui/ ./
RUN npm run build

# Stage 2: Build Go binary
FROM golang:1.22-alpine AS go-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built UI from previous stage
COPY --from=ui-builder /app/ui/dist ./ui/dist

# Build the binary with version information
ARG VERSION
RUN if [ -z "$VERSION" ]; then \
        MB_DATE=$(date '+%Y.%m.%d') && \
        MB_HASH=$(git rev-parse --short HEAD) && \
        VERSION="$MB_DATE-$MB_HASH"; \
    fi && \
    go build \
    -ldflags "-X main.Version=${VERSION}" \
    -o mediagui ./mediagui.go

# Stage 3: Runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create user and group
ARG USER_ID
ARG GROUP_ID
RUN addgroup -g ${GROUP_ID} -S mediagui && \
    adduser -u ${USER_ID} -S mediagui -G mediagui

# Create directories
RUN mkdir -p /app /data && \
    chown -R mediagui:mediagui /app /data

# Copy binary from builder
COPY --from=go-builder /app/mediagui /app/mediagui

# Set ownership
RUN chown -R mediagui:mediagui /app

# Switch to non-root user
USER mediagui

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 7623

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:7623/health || exit 1

# Default command
CMD ["./mediagui"]
