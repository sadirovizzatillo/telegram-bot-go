FROM golang:1.22-alpine

# Install Python, pip, ffmpeg, and yt-dlp
RUN apk add --no-cache python3 py3-pip ffmpeg \
    && pip install --no-cache-dir -U yt-dlp

WORKDIR /app

# Copy Go module files first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your code
COPY . .

# Build the bot
RUN go build -ldflags="-w -s" -o bot ./cmd/bot

# Run the bot
CMD ["./bot"]
