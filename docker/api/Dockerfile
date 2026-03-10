FROM golang:1.25.1-alpine AS builder
WORKDIR /app
COPY ../../go.mod ./
RUN go mod download
COPY ../.. .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/main .
COPY --from=builder --chown=appuser:appgroup /app/migrations ./migrations
USER appuser

EXPOSE 8000
CMD ["./main"]