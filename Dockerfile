FROM golang:1.21-bullseye AS builder

# Set environment variables
ENV CGO_ENABLED=1
ENV GO111MODULE=on

# Install dependencies for robotgo
RUN apt-get update && apt-get install -y --no-install-recommends \
    git gcc libc6-dev \
    libx11-dev xorg-dev libxtst-dev \
    libpng-dev \
    libxcursor-dev \
    libxinerama-dev \
    libxi-dev \
    libx11-xcb-dev \
    libxkbcommon-dev \
    libxkbcommon-x11-dev \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy source code first
COPY *.go ./

# Initialize Go module and add dependencies
# Using a specific older version of robotgo that doesn't have the type mismatch issue
RUN go mod init tracker
RUN go get github.com/go-vgo/robotgo@v0.100.0
RUN go get github.com/robotn/gohook@v0.40.0
RUN go get github.com/gofiber/fiber/v2@v2.50.0
RUN go get github.com/lib/pq@v1.10.9
RUN go mod tidy

# Build application
RUN go build -v -o tracker-app .

# Use Debian slim for runtime
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    libx11-6 libxtst6 libxcursor1 libxinerama1 libxi6 libxcb1 libx11-xcb1 \
    libxkbcommon0 libxkbcommon-x11-0 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/tracker-app .

# Expose port
EXPOSE 3000

# Run
CMD ["./tracker-app"]