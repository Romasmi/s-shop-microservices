.PHONY: generate build run-api run-worker run-cli compose-up compose-down migrate-up migrate-down



build:
	go build -o bin/api cmd/api/main.go

run-api:
	go run cmd/api/main.go

compose-up:
	docker-compose -f deployments/docker-compose.yaml up -d

compose-down:
	docker-compose -f deployments/docker-compose.yaml down

migrate-up:
	migrate -path migrations -database "postgres://user:password@localhost:5432/order_service?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://user:password@localhost:5432/order_service?sslmode=disable" down

docker-build:
	docker build --platform linux/amd64 -t romasmi/order-service:latest -f Dockerfile ../..
	@if minikube status >/dev/null 2>&1; then \
		echo "Loading image into minikube..."; \
		minikube image load romasmi/order-service:latest; \
	fi

docker-push:
	docker push romasmi/order-service:latest

kube-restart:
	kubectl rollout restart deployment/order-service -n s-shop-system

kube-status:
	kubectl rollout status deployment/order-service -n s-shop-system

kube-logs:
	kubectl logs deployment/order-service -n s-shop-system


lint:
	golangci-lint run ./...