.PHONY: up build deploy restart install-ingress install-db install-grafana apply hosts run migration-up wait-db wait-api clean redeploy status help load-test prometheus-run grafana-run context

up: build deploy restart wait-api


deploy: install-ingress install-db install-grafana apply hosts wait-db
	@echo "Completed"
	@echo "Run 'make run' in a separate terminal to start minikube tunnel"
	@echo "Open: http://arch.homework:8080/health"
	@echo "Open: http://arch.homework:8080/metrics"
	@echo "Prometheus: http://localhost:9090 (run 'make prometheus-run' in a separate terminal)"
	@echo "Grafana: http://localhost:3000 (run 'make grafana-run' in a separate terminal)"
	@echo "Grafana password: make grafana-pass"
	@echo "Hint: run 'make context' to set s-shop-system as your default namespace"

context:
	kubectl config set-context --current --namespace=s-shop-system

grafana-pass:
	@kubectl get secret grafana -o jsonpath="{.data.admin-password}" -n s-shop-system | base64 --decode ; echo ""


restart:
	@echo "Restarting API deployment..."
	kubectl rollout restart deployment/api -n s-shop-system

install-ingress:
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx/ || true
	helm repo update ingress-nginx
	helm upgrade --install nginx-ingress ingress-nginx/ingress-nginx \
		--namespace ingress-nginx \
		--create-namespace \
		--values helm/nginx-ingress.yaml

install-db:
	helm repo add bitnami https://repo.broadcom.com/bitnami-files/ || true
	helm repo update bitnami
	helm upgrade --install postgresql bitnami/postgresql \
		--namespace s-shop-system \
		--create-namespace \
		--values helm/postgresql-values.yaml

install-grafana:
	helm repo add grafana https://grafana.github.io/helm-charts || true
	helm repo update grafana
	helm upgrade --install grafana grafana/grafana \
		--namespace s-shop-system \
		--create-namespace \
		--values helm/grafana-values.yaml

apply:
	kubectl apply -f ./k8s

hosts:
	@if ! grep -q "arch.homework" /etc/hosts; then \
		echo "127.0.0.1 arch.homework" | sudo tee -a /etc/hosts; \
	else \
		echo "arch.homework already exists in /etc/hosts"; \
	fi

run:
	minikube tunnel

prometheus-run:
	kubectl port-forward service/prometheus-server 9090:80 -n s-shop-system

grafana-run:
	kubectl port-forward service/grafana 3000:80 -n s-shop-system

migration-up:
	@echo "Running migrations..."
	@echo "Note: This requires a tunnel or port-forward to localhost:5432"
	docker run --rm -v $(PWD)/migrations:/migrations --add-host=host.docker.internal:host-gateway migrate/migrate \
		-path=/migrations/ -database "postgres://user:password@host.docker.internal:5432/social_network?sslmode=disable" up || \
	docker run --rm -v $(PWD)/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ -database "postgres://user:password@localhost:5432/social_network?sslmode=disable" up

wait-db:
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5
	kubectl wait --namespace s-shop-system --for=condition=ready pod -l app.kubernetes.io/name=postgresql --timeout=120s

wait-api:
	@echo "Waiting for API deployment to be ready..."
	kubectl rollout status deployment/api -n s-shop-system --timeout=120s


redeploy: clean deploy

clean:
	kubectl delete -f ./k8s --ignore-not-found=true
	helm uninstall nginx-ingress -n ingress-nginx --ignore-not-found
	helm uninstall postgresql -n s-shop-system --ignore-not-found
	helm uninstall grafana -n s-shop-system --ignore-not-found
	kubectl delete namespace ingress-nginx --ignore-not-found=true
	kubectl delete namespace s-shop-system --ignore-not-found=true
	@echo "Note: /etc/hosts entry must be removed manually"

status:
	@echo "\n Ingress Controller:"
	@kubectl get pods -n ingress-nginx
	@echo "\n Application:"
	@kubectl get pods -n s-shop-system
	@echo "\n Services (s-shop-system):"
	@kubectl get svc -n s-shop-system
	@echo "\n Services (ingress-nginx):"
	@kubectl get svc -n ingress-nginx
	@echo "\n Ingress:"
	@kubectl get ingress -n s-shop-system

load-test:
	@echo "Running load tests with K6..."
	k6 run --vus 1000 --duration 30s load_testing/users.js

help:
	@echo "make up            - Build image, deploy everything and run migrations"
	@echo "make deploy         - Deploy ingress, database, and application"
	@echo "make build          - Build Docker image locally"
	@echo "make restart        - Force restart of the API deployment"
	@echo "make migration-up   - Run database migrations"
	@echo "make run         - Start minikube tunnel (separate terminal)"
	@echo "make load-test      - Run K6 load tests"
	@echo "make status         - Check deployment status"
	@echo "make clean          - Remove all resources"
	@echo "make redeploy       - Clean and redeploy everything"
	@echo "make install-ingress - Install only ingress controller"
	@echo "make apply - Apply only k8s manifests"
	@echo "make hosts    - Setup /etc/hosts entry"
	@echo ""
	@echo "Quick start:"
	@echo "  1. make deploy"
	@echo "  2. make run"
	@echo "  3. curl http://arch.homework:8080/otusapp/romasmi/health"

docker-build:
	$(MAKE) -C ./services/user-service docker-build

docker-push:
	$(MAKE) -C ./services/user-service docker-push

deploy-helm:
	helm upgrade --install user . --set serviceName=user
	helm upgrade --install auth . --set serviceName=auth