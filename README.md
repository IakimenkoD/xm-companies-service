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

to get local token for create and delete methods:

```shell script
 curl --request POST \
  --url http://localhost:4000/internal/signin \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxIiwiZXhwIjoxNjUwNTI4NDU1fQ.VFpR-3UUCB-QNyQBKEXxTY3vNBxCxi6-T_MShSh3yig' \
  --header 'Content-Type: application/json' \
  --data '{
	"username" : "user1",
	"password": "password1"
}'
```