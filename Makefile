include .env.local
export

.PHONY: help 
help: ## Prints help for targets with comments
	@echo 'Usage:'
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

.PHONY: db/up
db/up: ## Runs database in docker
	docker compose up -d postgres_db

.PHONY: db/migrations/new
db/migrations/new: ## Creates new db migration with name={}
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up: ## Apply db migrations 
	migrate -path=./migrations -database=${DB_DSN} up

.PHONY: run/gateway
run/gateway: ## Runs gateway service
	-go run -C services/gateway ./cmd/api/main.go

.PHONY: run/storage
run/storage: ## Runs storage service
	go run -C services/storage ./cmd/api/main.go

