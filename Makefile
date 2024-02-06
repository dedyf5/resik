GOCMD=go
GOGENERATE=$(GOCMD) generate

## generate: Generate files
generate:
	swag init -g ../../app/http/main.go -d ./entities/response,./app/http -o ./app/http/docs --instanceName http --parseDependency true
	$(GOGENERATE) ./...