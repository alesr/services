.DEFAULT_GOAL := help

PROJECT_NAME := stdservices

.PHONY: help

help:
	@echo "------------------------------------------------------------------------"
	@echo "${PROJECT_NAME}"
	@echo "------------------------------------------------------------------------"
	@grep -E '^[a-zA-Z0-9_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test-unit: ## run unit tests
	go test -v -race -vet=all -count=1 -timeout 240s ./...
