.PHONY: up build deploy restart install-traefik install-db install-grafana apply hosts run wait-db wait-api clean redeploy status help prometheus-run grafana-run

# Main target to start everything from scratch
up: build deploy wait-api

# Build Docker images and load them into minikube
build:
	$(MAKE) -C ./services/user-service docker-build
	$(MAKE) -C ./services/auth-service docker-build

docker-push:
	$(MAKE) -C ./services/user-service docker-push
	$(MAKE) -C ./services/auth-service docker-push

apply:
	kubectl apply -f ./deployment/k8s/

# Helm installations
install-traefik:
	helm repo add traefik https://traefik.github.io/charts || true
	helm repo update traefik

	helm upgrade --install traefik traefik/traefik \
	  --namespace traefik \
	  --create-namespace \
	  -f ./deployment/helm/traefik-values.yaml

forward-traefik:
	kubectl port-forward -n traefik $$(kubectl get pods -n traefik -o name) 9000:9000

install-db:
	helm repo add bitnami https://repo.broadcom.com/bitnami-files/
	helm repo update bitnami
	helm upgrade --install postgresql bitnami/postgresql \
		--namespace s-shop-system \
		--create-namespace \
		--values deployment/helm/postgresql-values.yaml

db-connect:
	kubectl exec -it postgresql-0 -n s-shop-system -- \
		psql -U user -d postgres

forward-db:
	kubectl port-forward svc/postgresql 5432:5432 -n s-shop-system


install-prometheus:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update prometheus-community
	helm upgrade --install prometheus prometheus-community/prometheus \
		--namespace s-shop-system \
		--create-namespace

install-grafana:
	helm repo add grafana https://grafana.github.io/helm-charts || true
	helm repo update grafana
	helm upgrade --install grafana grafana/grafana \
		--namespace s-shop-system \
		--create-namespace \
		--values deployment/helm/grafana-values.yaml

# Utility targets
hosts:
	@echo "Updating /etc/hosts for arch.homework..."
	@sudo sed -i '' '/arch.homework/d' /etc/hosts || sudo sed -i '/arch.homework/d' /etc/hosts
	@echo "127.0.0.1 arch.homework" | sudo tee -a /etc/hosts

run:
	@echo "Starting port-forwarding for Traefik... (Keep this running)"
	@echo "Access API: http://arch.homework:8080"
	@echo "Access Traefik Dashboard: http://arch.homework:8080/dashboard/"
	@echo "Access Grafana: http://localhost:3000 (after make grafana-run)"
	kubectl port-forward service/traefik 8080:8080 -n traefik

wait-db:
	@echo "Waiting for PostgreSQL to be ready..."
	kubectl wait --namespace s-shop-system --for=condition=ready pod -l app.kubernetes.io/name=postgresql --timeout=120s

wait-api:
	@echo "Waiting for API deployments to be ready..."
	kubectl rollout status deployment/user-service -n s-shop-system --timeout=120s
	kubectl rollout status deployment/auth-service -n s-shop-system --timeout=120s

restart:
	kubectl rollout restart deployment/user-service -n s-shop-system
	kubectl rollout restart deployment/auth-service -n s-shop-system
	$(MAKE) wait-api

clean:
	kubectl delete -f ./deployment/k8s --ignore-not-found=true
	helm uninstall traefik -n traefik --ignore-not-found
	helm uninstall postgresql -n s-shop-system --ignore-not-found
	helm uninstall grafana -n s-shop-system --ignore-not-found
	kubectl delete namespace traefik --ignore-not-found=true
	kubectl delete namespace s-shop-system --ignore-not-found=true

status:
	@echo "\n--- Infrastructure ---"
	@kubectl get pods -n traefik
	@kubectl get pods -n s-shop-system -l app.kubernetes.io/name=postgresql
	@echo "\n--- Application ---"
	@kubectl get pods -n s-shop-system -l app=user-service
	@kubectl get pods -n s-shop-system -l app=auth-service
	@echo "\n--- Services ---"
	@kubectl get svc -n s-shop-system
	@kubectl get svc -n traefik
	@echo "\n--- Routes ---"
	@kubectl get ingressroute -n s-shop-system

prometheus-run:
	kubectl port-forward service/prometheus-server 9090:80 -n s-shop-system

grafana-run:
	kubectl port-forward service/grafana 3000:80 -n s-shop-system

grafana-pass:
	@kubectl get secret grafana -o jsonpath="{.data.admin-password}" -n s-shop-system | base64 --decode ; echo ""

help:
	@echo "Usage:"
	@echo "  make up          - Build images and deploy everything (from scratch)"
	@echo "  make run         - Start minikube tunnel (required for access)"
	@echo "  make status      - Check deployment status"
	@echo "  make clean       - Remove all resources"
	@echo ""
	@echo "Quick Start:"
	@echo "  1. make up"
	@echo "  2. make run (in another terminal)"
	@echo "  3. Access API: http://arch.homework:8080/user"
	@echo "  4. Access Dashboard: http://arch.homework:8080/dashboard/"