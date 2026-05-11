FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy the shared api module
COPY api/ ./api/

# Copy the service go.mod/sum
COPY services/delivery-service/go.mod services/delivery-service/go.sum ./services/delivery-service/

WORKDIR /app/services/delivery-service
RUN go mod download

# Copy the rest of the service source
COPY services/delivery-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/services/delivery-service/main .
COPY --from=builder --chown=appuser:appgroup /app/services/delivery-service/migrations ./migrations

USER appuser

EXPOSE 8000
CMD ["./main"]
