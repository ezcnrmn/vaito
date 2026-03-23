include .env.local
export

.PHONY: help 
help: ## Prints help for targets with comments
	@echo 'Usage:'
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

.PHONY: db/up
db/up: ## Runs database in docker
	docker compose up -d postgres_db

.PHONY: db/migrations/user/new
db/migrations/user/new: ## Creates new db migration for user service with name={}
	migrate create -seq -ext=.sql -dir=./migrations/user ${name}

.PHONY: db/migrations/user/up
db/migrations/user/up: ## Apply db migrations for user service
	migrate -path=./migrations/user/ -database=${USER_DB_DSN} up

.PHONY: db/migrations/listing/new
db/migrations/listing/new: ## Creates new db migration for listing service with name={}
	migrate create -seq -ext=.sql -dir=./migrations/listing ${name}

.PHONY: db/migrations/listing/up
db/migrations/listing/up: ## Apply db migrations for listing service
	migrate -path=./migrations/listing/ -database=${LISTING_DB_DSN} up

.PHONY: run/gateway
run/gateway: ## Runs gateway service
	-go run -C services/gateway ./cmd/api/main.go -debug-log

.PHONY: run/user
run/user: ## Runs user service
	-go run -C services/user ./cmd/api/main.go -debug-log

.PHONY: run/listing
run/listing: ## Runs listing service
	-go run -C services/listing ./cmd/api/main.go -debug-log

USER_PROTO_DIR := ./proto/user/v1
USER_OUT_DIR := ./gen/go/user
USER_MODULE_NAME := github.com/ezcnrmn/vaito/gen/go/user

.PHONY: gen/user
gen/user: ## Gens go code for user/v1 .proto
	@mkdir $$USER_OUT_DIR -p

	@if [ ! -f $(USER_OUT_DIR)/go.mod ]; then \
		cd $(USER_OUT_DIR) && go mod init $(USER_MODULE_NAME); \
	fi

	@for file in ${wildcard ${USER_PROTO_DIR}/*.proto}; do \
		echo "Generating for $$file"; \
		protoc --proto_path=$(USER_PROTO_DIR) --go_out=$$USER_OUT_DIR --go_opt=paths=source_relative --go-grpc_out=$$USER_OUT_DIR --go-grpc_opt=paths=source_relative $$file; \
	done

	@cd $(USER_OUT_DIR) && go mod tidy

LISTING_PROTO_DIR := ./proto/listing/v1
LISTING_OUT_DIR := ./gen/go/listing
LISTING_MODULE_NAME := github.com/ezcnrmn/vaito/gen/go/listing

.PHONY: gen/listing
gen/listing: ## Gens go code for listing/v1 .proto
	@mkdir $$LISTING_OUT_DIR -p

	@if [ ! -f $(LISTING_OUT_DIR)/go.mod ]; then \
		cd $(LISTING_OUT_DIR) && go mod init $(LISTING_MODULE_NAME); \
	fi

	@for file in ${wildcard ${LISTING_PROTO_DIR}/*.proto}; do \
		echo "Generating for $$file"; \
		protoc --proto_path=$(LISTING_PROTO_DIR) --go_out=$$LISTING_OUT_DIR --go_opt=paths=source_relative --go-grpc_out=$$LISTING_OUT_DIR --go-grpc_opt=paths=source_relative $$file; \
	done

	@cd $(LISTING_OUT_DIR) && go mod tidy

