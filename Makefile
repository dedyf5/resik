GOCMD=go
GOGENERATE=$(GOCMD) generate

## generate: Generate files
generate:
	wire gen ./app/http/bootstrap
	swag init -g ../../app/http/main.go -d ./pkg/response,./app/http -o ./app/http/docs --instanceName http --parseDependency true
	$(GOGENERATE) ./...

## test: Test all files
test:
	$(GOCMD) test -v ./...
