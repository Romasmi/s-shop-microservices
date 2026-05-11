

docker-build:
	docker build --platform linux/amd64 -t romasmi/delivery-service:latest -f Dockerfile ../..
	@if minikube status >/dev/null 2>&1; then \
		echo "Loading image into minikube..."; \
		minikube image load romasmi/delivery-service:latest; \
	fi

docker-push:
	docker push romasmi/delivery-service:latest

kube-restart:
	kubectl rollout restart deployment/delivery-service -n s-shop-system

kube-status:
	kubectl rollout status deployment/delivery-service -n s-shop-system

kube-logs:
	kubectl logs deployment/delivery-service -n s-shop-system

#  forward port outside
kube-forward:
	kubectl port-forward -n s-shop-system svc/delivery-service 8083:8000