.DEFAULT_GOAL := help

PROJECT_NAME := stdservices

.PHONY: help
help:
	@echo "------------------------------------------------------------------------"
	@echo "${PROJECT_NAME}"
	@echo "------------------------------------------------------------------------"
	@grep -E '^[a-zA-Z0-9_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## run go format, vet and lint code
	@go fmt ./...
	@go vet ./...
	@go install honnef.co/go/tools/cmd/staticcheck@2022.1
	@staticcheck ./...

.PHONY: test-unit
test-unit: ## run unit tests
	@go test -v -short -race -vet=all -count=1 -timeout 240s ./...

.PHONY: test-it
test-it: ## run integration tests
	make db
	@go test -run Integration -v -race -vet=all -count=1 -timeout 240s ./...
	make db-down

.PHONY: test
test: ## Run unit and integration tests
	make test-unit
	make test-it

.PHONY: db-up
db: ## spins up the test database
	@docker-compose -f docker-compose.yaml up db -d
	@sleep 2
	@make migrate

.PHONY: db-down
db-down: ## remove the test database container and its volumes
	@docker-compose -f docker-compose.yaml down -v

.PHONY: migrate
migrate: ## executes the migrations towards the test database
	@docker run -v $(CURDIR)/migrations:/migrations \
	--network host migrate/migrate \-path=/migrations/ -database \
	"postgres://user:password@localhost:5432/testdb?sslmode=disable" up

.PHONY: psql
psql: ## executes a psql command to connect to the test database
	@psql postgres://user:password@localhost:5432/testdb?sslmode=disable
