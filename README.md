# xm-companies-service
REST API microservice to handle Companies


## Dependencies

* Go 1.17
* Postgres 10+
* RabbitMQ

## After clone actions

get Docker
get docker-compose 

prepare local env:
```shell script
 docker-compose -f docker-compose-local.yml up -d
```

run project locally:
```shell script
 go run cmd/companies_service/main.go
```

build binary:

```shell script
 go build cmd/companies_service/main.go
```
