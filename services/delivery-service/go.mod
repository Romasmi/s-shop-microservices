module github.com/Romasmi/s-shop-microservices/delivery-service

go 1.25.1

require (
	github.com/Romasmi/s-shop v0.0.0
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/prometheus/client_golang v1.23.2
	github.com/spf13/viper v1.21.0
	google.golang.org/grpc v1.81.0
	google.golang.org/protobuf v1.36.11
)

replace github.com/Romasmi/s-shop => ../../api
