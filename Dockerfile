# =========================
# Builder Stage
# =========================
FROM golang:1.26.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o blog-api ./cmd/api


# =========================
# Runtime Stage
# =========================
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/blog-api .

EXPOSE 8080

CMD ["./blog-api"]