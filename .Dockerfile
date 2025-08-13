# Stage 1 — Build
FROM golang:1.22 AS builder

# Enable Go modules and tidy
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

# Stage 2 — Runtime
FROM debian:bookworm-slim

# Install curl & yt-dlp
RUN apt-get update && \
    apt-get install -y curl ca-certificates ffmpeg && \
    curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp \
    -o /usr/local/bin/yt-dlp && \
    chmod +x /usr/local/bin/yt-dlp && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /root

# Copy binary from builder
COPY --from=builder /app/bot .

# Run the bot
CMD ["./bot"]
