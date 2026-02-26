include .env.local
export

.PHONY: help 
help: ## Prints help for targets with comments
	@echo 'Usage:'
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

.PHONY: db/up
db/up: ## Runs database in docker
	docker compose up -d postgres_db

.PHONY: run/gateway
run/gateway: ## Runs gateway service
	go run -C services/gateway ./cmd/api/main.go

.PHONY: run/storage
run/storage: ## Runs storage service
	go run -C services/storage ./cmd/api/main.go

