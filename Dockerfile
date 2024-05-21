FROM golang:1.22-alpine AS builder
WORKDIR /wallet-rielta
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o wallet-rielta ./cmd/app/main.go

FROM alpine:latest
WORKDIR /wallet-rielta
COPY --from=builder /wallet-rielta ./
EXPOSE 8080

CMD ["./wallet-rielta"]