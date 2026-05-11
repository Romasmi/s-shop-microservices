

docker-build:
	docker build --platform linux/amd64 -t romasmi/user-service:latest -f Dockerfile .
	@if minikube status >/dev/null 2>&1; then \
		echo "Loading image into minikube..."; \
		minikube image load romasmi/user-service:latest; \
	fi

docker-push:
	docker push romasmi/user-service:latest

kube-restart:
	kubectl rollout restart deployment/user-service -n s-shop-system

kube-status:
	kubectl rollout status deployment/user-service -n s-shop-system

kube-logs:
	kubectl logs deployment/user-service -n s-shop-system

#  forward port outside
kube-forward:
	kubectl port-forward -n s-shop-system svc/user-service 8081:8000