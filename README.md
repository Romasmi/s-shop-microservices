# S-SHOP system implementation (microservice architecture)

## Quick Start

To start the API locally using Minikube:

1.  **Start the tunnel** (in a separate terminal):
    ```shell
    make run
    ```

2.  **Start the environment**:
    ```shell
    make up
    ```
    This command builds the Docker image, installs the Ingress controller, PostgreSQL database via Helm, applies application manifests, configures `/etc/hosts`, and waits for the API to be ready (it runs migrations automatically on startup).
    *Note: The tunnel must be running to access the API via hostnames.*

3.  **Access the API**:
    The API is available at `http://arch.homework:8080`.

    Example:
    ```shell
    curl -X POST http://arch.homework:8080/user \
    --header 'Content-Type: application/json' \
    --data-raw '{
      "username": "johndoe",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john@doe.com",
      "phone": "+71002003040"
    }'
    ```

## Development

### Database Migrations
To run migrations manually:
```shell
make migration-up
```

### Build Docker Image
```shell
make build
```
Or manually:
```shell
docker build --platform linux/amd64 -t romasmi/s-shop-system:latest -f docker/api/Dockerfile .
```

https://hub.docker.com/r/romasmi/s-shop-system


## How to deploy kubernates
Run make command `make deploy` or install manually

### Install ingress

```shell
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx/ && \
helm repo update && \
helm install nginx-ingress ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace \
  --values nginx-ingress.yaml
```

### Apply k8s manifests
```shell
kubectl apply -f ./k8s
```

### If you use minikube then start tunnel
```shell
minikube tunnel
```
and add host to /etc/hosts

```shell
echo "127.0.0.1 arch.homework" | sudo tee -a /etc/hosts
```

### Test
Start tunnel on MacOS/Windows `minikube tunnel`

```shell
curl http://arch.homework:8080/health
curl http://arch.homework:8080/otusapp/romasmi/health
```

## How to add Grafana
```shell
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm install grafana grafana/grafana
-- expose 
kubectl port-forward service/grafana 3000:80
-- open
http://localhost:3000
-- login admin/admin

-- get password
kubectl get secret grafana \
  -o jsonpath="{.data.admin-password}" | base64 --decode
```

## Monitoring and Metrics

### Application Metrics
Get metrics endpoint `/metrics`.

Metrics:
- **Latency**: `http_request_duration_seconds_bucket` (Histogram)
- **RPS (Requests Per Second)**: Can be calculated with `rate(http_requests_total[1m])`
- **Error Rate**: Can be calculated with `rate(http_requests_total{status=~"4..|5.."}[1m]) / rate(http_requests_total[1m])`
