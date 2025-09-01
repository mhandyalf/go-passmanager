FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# Stage 2: Run — minimal image
FROM alpine:3.18

WORKDIR /app

# copy binary dari builder
COPY --from=builder /app/myapp .

# kalau perlu sertifikat SSL (buat koneksi HTTPS/DB TLS)
RUN apk add --no-cache ca-certificates

EXPOSE 8080

CMD ["./myapp"]
