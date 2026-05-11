FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy the shared api module
COPY api/ ./api/

# Copy the service go.mod/sum
COPY services/warehouse-service/go.mod services/warehouse-service/go.sum ./services/warehouse-service/

WORKDIR /app/services/warehouse-service
RUN go mod download

# Copy the rest of the service source
COPY services/warehouse-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/services/warehouse-service/main .
COPY --from=builder --chown=appuser:appgroup /app/services/warehouse-service/migrations ./migrations
COPY --from=builder --chown=appuser:appgroup /app/services/warehouse-service/config.yaml .

USER appuser

EXPOSE 8000
CMD ["./main"]
