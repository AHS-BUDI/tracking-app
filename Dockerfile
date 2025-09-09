FROM golang:1.25-alpine AS builder

# Install dependensi C yang dibutuhkan oleh robotgo
RUN apk add --no-cache gcc g++ pkgconfig xorg-server-dev libx11-dev libxtst-dev libxkbcommon-dev

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build aplikasi
RUN CGO_ENABLED=1 GOOS=linux go build -o tracker-app .

# Gunakan Alpine untuk image final yang lebih kecil
FROM alpine:3.16

# Install runtime dependencies untuk robotgo
RUN apk add --no-cache libx11 libxtst xorg-server libxkbcommon

WORKDIR /app

# Copy binary dari builder stage
COPY --from=builder /app/tracker-app .

# Expose port untuk API
EXPOSE 3000

# Run aplikasi
CMD ["./tracker-app"]