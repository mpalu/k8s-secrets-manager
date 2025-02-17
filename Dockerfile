FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /k8s-secrets-manager ./cmd/k8s-secrets-manager

FROM alpine:3.18
COPY --from=builder /k8s-secrets-manager /k8s-secrets-manager
COPY config/config.yaml /config/config.yaml

EXPOSE 8080
ENTRYPOINT ["/k8s-secrets-manager"] 