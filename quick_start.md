# Initial setup
```shell
make install-grafana 
make install-prometheus
make install-db
make forward-db # to connect to db via local client
```

## How to add a new DB
```shell
make db-connect
```
then add DB manuall.

# Install traefik
```shell
make install-traefik
make forward-traefik # expose it to get locally http://localhost:9000/
```


# Refresh container in pod
```shell
kubectl rollout restart deployment user-service -n s-shop-system
```