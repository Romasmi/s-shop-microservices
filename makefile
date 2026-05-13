.PHONY: up build deploy restart install-traefik install-db install-kafka install-grafana hosts run wait-db wait-kafka wait-api clean redeploy status help prometheus-run grafana-run forward-kafka forward-db forward-traefik proto-gen

# Main target to start everything from scratch
up: proto-gen build deploy wait-db wait-kafka wait-api

# Proto generation
proto-gen:
	$(MAKE) -C ./api generate

# Build Docker images and load them into minikube
build:
	$(MAKE) -j5 -C ./services/user-service docker-build & \
	$(MAKE) -j5 -C ./services/auth-service docker-build & \
	$(MAKE) -j5 -C ./services/notification-service docker-build & \
	$(MAKE) -j5 -C ./services/order-service docker-build & \
	$(MAKE) -j5 -C ./services/billing-service docker-build & \
	wait

docker-push:
	$(MAKE) -j5 -C ./services/user-service docker-push & \
	$(MAKE) -j5 -C ./services/auth-service docker-push & \
	$(MAKE) -j5 -C ./services/notification-service docker-push & \
	$(MAKE) -j5 -C ./services/order-service docker-push & \
	$(MAKE) -j5 -C ./services/billing-service docker-push & \
	wait

deploy: install-traefik install-db install-kafka install-prometheus install-grafana install-app

install-app:
	helm upgrade --install s-shop-system ./deployment/helm/s-shop-system \
		--namespace s-shop-system \
		--create-namespace

uninstall-app:
	helm uninstall s-shop-system -n s-shop-system --ignore-not-found

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

install-kafka:
	helm repo add redpanda https://charts.redpanda.com
	helm repo update redpanda
	helm upgrade --install redpanda redpanda/redpanda \
		--namespace s-shop-system \
		--create-namespace \
		--values deployment/helm/redpanda-values.yaml

forward-kafka:
	kubectl port-forward svc/kafka 9092:9092 -n s-shop-system

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

wait-kafka:
	@echo "Waiting for Redpanda to be ready..."
	kubectl wait --namespace s-shop-system --for=condition=ready pod -l app.kubernetes.io/name=redpanda --timeout=300s

wait-api:
	@echo "Waiting for API deployments to be ready..."
	kubectl rollout status deployment/user-service -n s-shop-system --timeout=120s
	kubectl rollout status deployment/auth-service -n s-shop-system --timeout=120s
	kubectl rollout status deployment/notification-service-api -n s-shop-system --timeout=120s
	kubectl rollout status deployment/notification-service-worker -n s-shop-system --timeout=120s
	kubectl rollout status deployment/order-service -n s-shop-system --timeout=120s
	kubectl rollout status deployment/billing-service-api -n s-shop-system --timeout=120s
	kubectl rollout status deployment/billing-service-worker -n s-shop-system --timeout=120s

restart:
	kubectl rollout restart deployment/user-service -n s-shop-system
	kubectl rollout restart deployment/auth-service -n s-shop-system
	kubectl rollout restart deployment/notification-service-api -n s-shop-system
	kubectl rollout restart deployment/notification-service-worker -n s-shop-system
	kubectl rollout restart deployment/order-service -n s-shop-system
	kubectl rollout restart deployment/billing-service-api -n s-shop-system
	kubectl rollout restart deployment/billing-service-worker -n s-shop-system
	$(MAKE) wait-api

clean:
	helm uninstall s-shop-system -n s-shop-system --ignore-not-found
	helm uninstall traefik -n traefik --ignore-not-found
	helm uninstall postgresql -n s-shop-system --ignore-not-found
	helm uninstall redpanda -n s-shop-system --ignore-not-found
	helm uninstall prometheus -n s-shop-system --ignore-not-found
	helm uninstall grafana -n s-shop-system --ignore-not-found
	kubectl delete namespace traefik --ignore-not-found=true
	kubectl delete namespace s-shop-system --ignore-not-found=true

status:
	@echo "\n--- Infrastructure ---"
	@echo "Traefik:"
	@kubectl get pods -n traefik
	@echo "PostgreSQL:"
	@kubectl get pods -n s-shop-system -l app.kubernetes.io/name=postgresql
	@echo "Redpanda:"
	@kubectl get pods -n s-shop-system -l app.kubernetes.io/name=redpanda
	@echo "\n--- Application ---"
	@kubectl get pods -n s-shop-system -l app=user-service
	@kubectl get pods -n s-shop-system -l app=auth-service
	@kubectl get pods -n s-shop-system -l app=notification-service
	@kubectl get pods -n s-shop-system -l app=order-service
	@kubectl get pods -n s-shop-system -l app=billing-service
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

redeploy: build docker-push restart

help:
	@echo "Usage:"
	@echo "  make up          - Build images and deploy everything (from scratch)"
	@echo "  make redeploy    - Build, push and restart all services"
	@echo "  make run         - Start minikube tunnel (required for access)"
	@echo "  make status      - Check deployment status"
	@echo "  make clean       - Remove all resources"
	@echo "  make install-app - Install application using Helm"
	@echo "  make forward-kafka - Port-forward Kafka to localhost:9092"
	@echo "  make forward-db    - Port-forward PostgreSQL to localhost:5432"
	@echo ""
	@echo "Quick Start:"
	@echo "  1. make up"
	@echo "  2. make run (in another terminal)"
	@echo "  3. Access API: http://arch.homework:8080/user"
	@echo "  4. Access Dashboard: http://arch.homework:8080/dashboard/"

draw-puml:
	plantuml -tsvg ./docs/puml/*.puml

test-postman:
	newman run docs/postman.json --verbose