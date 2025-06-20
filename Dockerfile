FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest

RUN apk add --no-cache \
    python3 \
    py3-pip \
    ffmpeg \
    ca-certificates \
    && pip3 install --break-system-packages --no-cache-dir yt-dlp

RUN mkdir -p /app/downloads

COPY --from=builder /app/main /app/main

WORKDIR /app

RUN chmod 755 /app/downloads

EXPOSE 8080

CMD ["./main"]