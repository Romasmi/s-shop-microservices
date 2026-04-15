# S-SHOP system implementation (microservice architecture)

## Quick Start

## Quick Start (from scratch)

Ensure minikube is running else run `minikube start`

To deploy everything in Minikube:

1.  **Build and Deploy**:
    ```shell
    make up
    ```
    This command builds Docker images, installs PostgreSQL, Grafana, and Traefik via Helm, and applies all Kubernetes manifests from `deployment/k8s/`.

2.  **Start API Proxy**:
    ```shell
    make run
    ```
    *Keep this command running in a separate terminal. It enables access via http://arch.homework:8080*

3.  **Verify**:
    ```shell
    make status
    ```

## Accessing the API

The API is exposed at `http://arch.homework:8080` (ensure `make up` added the entry to your `/etc/hosts`).

- **Health Check**: `curl http://arch.homework:8080/health`
- **Auth Endpoint**: `curl http://arch.homework:8080/auth`
- **User API**: `curl http://arch.homework:8080/user`

## Monitoring

- **Prometheus**: Run `make prometheus-run` and open [http://localhost:9090](http://localhost:9090)
- **Grafana**: Run `make grafana-run` and open [http://localhost:3000](http://localhost:3000) (User: `admin`, get password via `make grafana-pass`)

## Project Structure

- `services/`: Source code for microservices.
- `deployment/k8s/`: Kubernetes manifests (ordered 00-70).
- `deployment/helm/`: Helm values for infrastructure components.

## Load tests 
```shell
k6 run load_testing/users.js
```