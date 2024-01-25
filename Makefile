GOCMD=go
GOGENERATE=$(GOCMD) generate

## generate: Generate files
generate:
	$(GOGENERATE) ./...