
.PHONY: check
check:
	@echo "Checking binaries..."
    EXECUTABLES = docker docker-compose go
    K := $(foreach exec,$(EXECUTABLES),\
 	       $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

.PHONY: registry
registry:
	@echo "Starting registry..."
	@cd hack/registry && docker-compose up -d registry


.PHONY: test
test: registry
	@echo "Running tests..."
	@go test -race ./...