GOCMD=go
GOGENERATE=$(GOCMD) generate

## generate: Generate files
generate:
	protoc --go_out=. --go-grpc_out=. app/grpc/handler/*/*.proto
	wire gen ./app/grpc/bootstrap
	wire gen ./app/rest/bootstrap
	swag init -g ../../app/rest/main.go -d ./pkg/response,./app/rest -o ./app/rest/docs --instanceName rest --parseDependency true
	$(GOGENERATE) ./...

## generate-proto: Generate GRPC files
generate-grpc:
	protoc --go_out=. --go-grpc_out=. app/grpc/handler/*/*.proto
	wire gen ./app/grpc/bootstrap

## generate-proto: Generate REST files
generate-rest:
	wire gen ./app/rest/bootstrap
	swag init -g ../../app/rest/main.go -d ./pkg/response,./app/rest -o ./app/rest/docs --instanceName rest --parseDependency true
	$(GOGENERATE) ./...

## generate-proto: Generate proto files
generate-proto:
	protoc --go_out=. --go-grpc_out=. app/grpc/handler/*/*.proto

## test: Test all files
test:
	$(GOCMD) test -v -cover ./...
