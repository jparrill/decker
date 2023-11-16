
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
	@sleep 5
	@echo

registry-down:
	@echo "Stopping registry..."
	@cd hack/registry && docker-compose down registry
	@echo

.PHONY: unit
unit: registry
	@echo "Running tests..."
	@go test -race ./...
	@echo

test: registry unit registry-down



