FROM golang:1.22

# Install yt-dlp
RUN apt-get update && apt-get install -y python3 python3-pip ffmpeg && \
    pip3 install -U yt-dlp && \
    rm -rf /var/lib/apt/lists/*

# Set workdir
WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build bot
RUN go build -o bot ./cmd/bot

# Run bot
CMD ["./bot"]
