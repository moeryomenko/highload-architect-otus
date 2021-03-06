HOST?=localhost:8080
ACCESS_TOKEN?=$(shell curl -X POST http://${HOST}/api/v1/login -d '{"nickname":"test", "password":"test"}' | jq .access_token)
NAME?=Dangelo
LAST?=Cummings

test:
	@go test -v ./...

restore:
	@gunzip < dumps/profiles.sql.gz | docker exec -i mysql-master /usr/bin/mysql -u social -ptest social

test-user:
	@curl -X POST http://${HOST}/api/v1/signup -d '{"nickname":"test", "password":"test"}'

tools:
	@go install -v \
		github.com/deepmap/oapi-codegen/cmd/oapi-codegen

load-tools:
	@go install -v go.k6.io/xk6/cmd/xk6@latest
	@xk6 build --with github.com/szkiba/xk6-dashboard@latest

load-test:
	@./k6 run -e HOSTNAME=${HOST} -e FIRST_NAME=${NAME} -e LAST_NAME=${LAST} -e TOKEN=${ACCESS_TOKEN} -e NEXT_TOKEN=${NEXT} --out dashboard ./loadtests/search.js

replication-test:
	@./k6 run -e HOSTNAME=${HOST} -e FIRST_NAME=${NAME} -e LAST_NAME=${LAST} -e TOKEN=${ACCESS_TOKEN} --out dashboard ./loadtests/replication.js

gen-api:
	@oapi-codegen --package router \
		-generate chi-server,types api/openapi.yaml  > internal/router/gen.go

.PHONY: test tools gen-api
