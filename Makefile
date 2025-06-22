GOCMD=go
GOGENERATE=$(GOCMD) generate

## generate: Generate files
generate: generate-proto generate-doc generate-go

## generate-proto: Generate gRPC files
generate-grpc: generate-proto generate-wire

## generate-proto: Generate REST files
generate-rest: generate-doc generate-wire

## generate-proto: Generate proto files
generate-proto:
	protoc $(PROTOC_FLAGS) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=. --go-grpc_out=. core/*/request/*.proto
	protoc-go-inject-tag -input="core/*/request/*.pb.go" -remove_tag_comment
	protoc $(PROTOC_FLAGS) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=. --go-grpc_out=. core/*/response/*.proto
	protoc-go-inject-tag -input="core/*/response/*.pb.go" -remove_tag_comment
	protoc $(PROTOC_FLAGS) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=. --go-grpc_out=. app/grpc/proto/*/*.proto
	protoc-go-inject-tag -input="app/grpc/proto/*/*.pb.go" -remove_tag_comment
	protoc $(PROTOC_FLAGS) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=. --go-grpc_out=. app/grpc/handler/*/*.proto
	protoc-go-inject-tag -input="app/grpc/handler/*/*.pb.go" -remove_tag_comment

## generate-wire: Generate wire files
generate-wire:
	wire gen ./app/grpc/bootstrap
	wire gen ./app/rest/bootstrap

## generate-go: Run go generate to create mocks, wire and other generated files
generate-go:
	$(GOGENERATE) ./...

## generate-doc: Generate documentation
generate-doc:
	swag init -g ../../app/rest/main.go -d ./pkg/response,./app/rest -o ./app/rest/docs --instanceName rest --parseDependency true

## test: Test all files
test:
	$(GOCMD) test -v -cover ./...
