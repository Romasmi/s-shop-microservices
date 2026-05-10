FROM golang:1.25.1-alpine AS builder
WORKDIR /app

COPY services/billing-service/go.mod services/billing-service/go.sum ./services/billing-service/
COPY services/user-service/go.mod services/user-service/go.sum ./services/user-service/
COPY services/order-service/go.mod services/order-service/go.sum ./services/order-service/

WORKDIR /app/services/order-service
RUN go mod download

WORKDIR /app
COPY services/billing-service ./services/billing-service
COPY services/user-service ./services/user-service
COPY services/order-service ./services/order-service

WORKDIR /app/services/order-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api    ./cmd/api
#RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o worker ./cmd/worker
#RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o cli    ./cmd/cli

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/services/order-service/api    .
#COPY --from=builder --chown=appuser:appgroup /app/services/order-service/worker .
#COPY --from=builder --chown=appuser:appgroup /app/services/order-service/cli    .
COPY --from=builder --chown=appuser:appgroup /app/services/order-service/migrations ./migrations
COPY --from=builder --chown=appuser:appgroup /app/services/order-service/config.yaml .
USER appuser

EXPOSE 8000
CMD ["./api"]