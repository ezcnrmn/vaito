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
	-go run -C services/gateway ./cmd/api/main.go -debug-log

.PHONY: run/storage
run/storage: ## Runs storage service
	-go run -C services/storage ./cmd/api/main.go -debug-log

PROTO_DIR := ./proto/storage/v1
OUT_DIR := ./gen/go/storage
MODULE_NAME := github.com/ezcnrmn/vaito/gen/go/storage

.PHONY: gen/storage
gen/storage: ## Gens go code for storage .proto
	@mkdir $$OUT_DIR -p

	@if [ ! -f $(OUT_DIR)/go.mod ]; then \
		cd $(OUT_DIR) && go mod init $(MODULE_NAME); \
	fi

	@for file in ${wildcard ${PROTO_DIR}/*.proto}; do \
		protoc --proto_path=$(PROTO_DIR) --go_out=$$OUT_DIR --go_opt=paths=source_relative --go-grpc_out=$$OUT_DIR --go-grpc_opt=paths=source_relative $$file; \
	done

	cd $(OUT_DIR) && go mod tidy

