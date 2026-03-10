# S-SHOP system implementation (microservice architecture)
Readme will be updated later

It will be implementation of system specified here:
https://github.com/Romasmi/s-shop-decomposition-to-microservices

## How to build docker image
```shell
docker build --platform linux/amd64 -t s-shop-system:latest .
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