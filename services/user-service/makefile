docker-build:
	docker build --platform linux/amd64 -t romasmi/s-shop-system:latest -f Dockerfile .
	@if minikube status >/dev/null 2>&1; then \
		echo "Loading image into minikube..."; \
		minikube image load romasmi/s-shop-system:latest; \
	fi

docker-push:
	docker push romasmi/s-shop-system:latest

kube-restart:
	kubectl rollout restart deployment/user-service -n s-shop-system

kube-status:
	kubectl rollout status deployment/user-service -n s-shop-system