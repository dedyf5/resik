GOCMD = go
GOGENERATE = $(GOCMD) generate

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
	protoc $(PROTOC_FLAGS) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=. --go-grpc_out=. pkg/response/*.proto
	protoc-go-inject-tag -input="pkg/response/*.pb.go" -remove_tag_comment
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

APP_VERSION ?= $(shell git describe --tags --always)
APP_GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
APP_BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS_DEF = -X github.com/dedyf5/resik/buildinfo.AppGitCommit=$(APP_GIT_COMMIT) \
			-X github.com/dedyf5/resik/buildinfo.AppBuildTime=$(APP_BUILD_TIME)
LDFLAGS_TAG = -X github.com/dedyf5/resik/buildinfo.AppVersion=$(APP_VERSION) \
			-X github.com/dedyf5/resik/buildinfo.AppVersionGenerator=tag \
			-X github.com/dedyf5/resik/buildinfo.AppGitCommit=$(APP_GIT_COMMIT) \
			-X github.com/dedyf5/resik/buildinfo.AppBuildTime=$(APP_BUILD_TIME)

run-rest:
	$(GOCMD) run -ldflags "$(LDFLAGS_DEF)" main.go rest

run-rest-with-tag:
	$(GOCMD) run -ldflags "$(LDFLAGS_TAG)" main.go rest

run-grpc:
	$(GOCMD) run -ldflags "$(LDFLAGS_DEF)" main.go grpc

run-grpc-with-tag:
	$(GOCMD) run -ldflags "$(LDFLAGS_TAG)" main.go grpc

build:
	$(GOCMD) build -ldflags "$(LDFLAGS_DEF)" -o resik

build-with-tag:
	$(GOCMD) build -ldflags "$(LDFLAGS_TAG)" -o resik

build-docker:
	CGO_ENABLED=0 GOOS=linux $(GOCMD) build -ldflags "-w -s $(LDFLAGS_DEF)" -o resik

build-docker-with-tag:
	CGO_ENABLED=0 GOOS=linux $(GOCMD) build -ldflags "-w -s $(LDFLAGS_TAG)" -o resik
