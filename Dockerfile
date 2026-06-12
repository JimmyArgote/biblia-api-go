# Multi-stage: compila em Go e roda em Alpine mínimo

# ===== Stage 1: Build =====
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Cache das dependências
COPY go.mod go.sum ./
RUN go mod download

# Compilar
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o biblia-api-go .

# ===== Stage 2: Runtime =====
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/biblia-api-go .

EXPOSE 8081

ENTRYPOINT ["./biblia-api-go"]
