

docker-build:
	docker build --platform linux/amd64 -t romasmi/warehouse-service:latest -f Dockerfile ../..
	@if minikube status >/dev/null 2>&1; then \
		echo "Loading image into minikube..."; \
		minikube image load romasmi/warehouse-service:latest; \
	fi

docker-push:
	docker push romasmi/warehouse-service:latest

kube-restart:
	kubectl rollout restart deployment/warehouse-service -n s-shop-system

kube-status:
	kubectl rollout status deployment/warehouse-service -n s-shop-system

kube-logs:
	kubectl logs deployment/warehouse-service -n s-shop-system

#  forward port outside
kube-forward:
	kubectl port-forward -n s-shop-system svc/warehouse-service 8082:8000