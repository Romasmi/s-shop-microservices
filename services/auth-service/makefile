docker-build:
	docker build --platform linux/amd64 -t romasmi/auth-service:latest -f Dockerfile .
	@if minikube status >/dev/null 2>&1; then \
		echo "Loading image into minikube..."; \
		minikube image load romasmi/auth-service:latest; \
	fi

docker-push:
	docker push romasmi/auth-service:latest
