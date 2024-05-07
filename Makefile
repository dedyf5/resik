GOCMD=go
GOGENERATE=$(GOCMD) generate

## generate: Generate files
generate:
	protoc --go_out=. --go-grpc_out=. app/grpc/handler/*/*.proto
	wire gen ./app/grpc/bootstrap
	wire gen ./app/http/bootstrap
	swag init -g ../../app/http/main.go -d ./pkg/response,./app/http -o ./app/http/docs --instanceName http --parseDependency true
	$(GOGENERATE) ./...

## generate-proto: Generate GRPC files
generate-grpc:
	protoc --go_out=. --go-grpc_out=. app/grpc/handler/*/*.proto
	wire gen ./app/grpc/bootstrap

## generate-proto: Generate HTTP files
generate-http:
	wire gen ./app/http/bootstrap
	swag init -g ../../app/http/main.go -d ./pkg/response,./app/http -o ./app/http/docs --instanceName http --parseDependency true
	$(GOGENERATE) ./...

## generate-proto: Generate proto files
generate-proto:
	protoc --go_out=. --go-grpc_out=. app/grpc/handler/*/*.proto

## test: Test all files
test:
	$(GOCMD) test -v -cover ./...
