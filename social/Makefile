test:
	@go test -v ./...

restore:
	@gunzip < dumps/profiles.sql.gz | docker exec -i mysql /usr/bin/mysql -u social -ptest social

tools:
	@go install -v \
		github.com/deepmap/oapi-codegen/cmd/oapi-codegen

gen-api:
	@oapi-codegen --package router \
		-generate chi-server,types api/openapi.yaml  > internal/router/gen.go

.PHONY: test tools gen-api
