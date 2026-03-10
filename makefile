.PHONY: deploy install-ingress apply-manifests setup-hosts tunnel clean help

deploy: install-ingress install-db apply hosts
	@echo "Completed"
	@echo "Run 'make tunnel' in a separate terminal to start minikube tunnel"
	@echo "Open: http://arch.homework:8080/health"

up: deploy migration-up

install-ingress:
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx/ || true
	helm repo update
	helm upgrade --install nginx-ingress ingress-nginx/ingress-nginx \
		--namespace ingress-nginx \
		--create-namespace \
		--values helm/nginx-ingress.yaml

install-db:
	helm repo add bitnami https://charts.bitnami.com/bitnami || true
	helm repo update
	helm upgrade --install postgresql bitnami/postgresql \
		--namespace s-shop-system \
		--create-namespace \
		--values helm/postgresql-values.yaml

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

migration-up:
	@echo "Running migrations..."
	@sleep 10
	docker run --rm -v $(PWD)/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ -database "postgres://user:password@localhost:5432/social_network?sslmode=disable" up

clean:
	kubectl delete -f ./k8s --ignore-not-found=true
	helm uninstall nginx-ingress -n ingress-nginx --ignore-not-found
	helm uninstall postgresql -n s-shop-system --ignore-not-found
	kubectl delete namespace ingress-nginx --ignore-not-found=true
	kubectl delete namespace s-shop-system --ignore-not-found=true
	@echo "Note: /etc/hosts entry must be removed manually"

redeploy: clean deploy

status:
	@echo "\n Ingress Controller:"
	@kubectl get pods -n ingress-nginx
	@echo "\n Application:"
	@kubectl get pods -n s-shop-system
	@echo "\n Services:"
	@kubectl get svc -n ingress-nginx
	@echo "\n Ingress:"
	@kubectl get ingress -n s-shop-system

help:
	@echo "make up            - Deploy everything and run migrations"
	@echo "make deploy         - Deploy ingress, database, and application"
	@echo "make migration-up   - Run database migrations"
	@echo "make run         - Start minikube tunnel (separate terminal)"
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