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