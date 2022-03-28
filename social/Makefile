tools:
	@go install -v \
		github.com/deepmap/oapi-codegen/cmd/oapi-codegen

gen-api:
	@oapi-codegen --package router \
		-generate chi-server,types api/openapi.yaml  > internal/router/gen.go

.PHONY: tools gen-api
