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
The API service exposes Prometheus metrics at `/metrics`. The `api` service in `k8s/20-service.yaml` is annotated for automatic discovery.

Metrics:
- **Latency**: `http_request_duration_seconds_bucket` (Histogram)
- **RPS (Requests Per Second)**: `rate(http_requests_total[1m])`
- **Error Rate**: `rate(http_requests_total{status=~"4..|5.."}[1m]) / rate(http_requests_total[1m])`

### Prometheus Installation
You can install Prometheus using the provided manifests or Helm.

#### Using Manifests (Recommended)
The Prometheus server is included in the `./k8s` directory and is deployed automatically with `make deploy`.
It is configured to:
- Scrape application metrics from pods/services with `prometheus.io/scrape: "true"` annotations.
- Scrape Kubernetes node metrics (cAdvisor) to provide **CPU and Memory usage by pods**.

To access Prometheus UI:
```shell
kubectl port-forward service/prometheus-service 9090:80 -n s-shop-system
```
Open http://localhost:9090

#### Using Helm
Alternatively, you can use Helm:
```shell
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/prometheus
```

### Nginx Ingress Controller Metrics
To enable Nginx Ingress metrics, we've updated `helm/nginx-ingress.yaml` with the following configuration:
```yaml
controller:
  metrics:
    enabled: true
    service:
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "10254"
```

### Kubernetes Pod Metrics (CPU/Memory)
Recommended Prometheus queries:
- **CPU Usage by Pod**: `sum(rate(container_cpu_usage_seconds_total{container!="", pod!=""}[5m])) by (pod, namespace)`
- **Memory Usage by Pod**: `sum(container_memory_working_set_bytes{container!="", pod!=""}) by (pod, namespace)`

Make sure `metrics-server` is enabled if you are using Minikube:
```shell
minikube addons enable metrics-server
```

## Load Testing

Install [k6](https://k6.io/) for load testing.

### Run tests

To run the user API load test:
```shell
k6 run --vus 1000 --duration 30s load_testing/users.js
```

Or specify a different base URL:
```shell
k6 run --vus 1000 --duration 30s -e BASE_URL=http://arch.homework:8080 load_testing/users.js
```

The test covers:
- Create user
- Get user by ID
- Update user
- Delete user
