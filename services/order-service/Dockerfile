FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy the shared api module
COPY api/ ./api/

# Copy the service go.mod/sum
COPY services/order-service/go.mod services/order-service/go.sum ./services/order-service/

WORKDIR /app/services/order-service
RUN go mod download

# Copy the rest of the service source
COPY services/order-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api    ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/services/order-service/api    .
COPY --from=builder --chown=appuser:appgroup /app/services/order-service/migrations ./migrations
COPY --from=builder --chown=appuser:appgroup /app/services/order-service/config.yaml .
USER appuser

EXPOSE 8000
CMD ["./api"]