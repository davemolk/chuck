FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o server \
    ./cmd/server/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server .

COPY --from=builder /app/internal/migrations ./migrations

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

USER appuser

ENV PORT=8080

EXPOSE $PORT

ENTRYPOINT ["./server"]