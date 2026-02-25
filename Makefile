DSN := postgres://vaito:pa55word@localhost:5432/vaito?sslmode=disable
STORAGE := http://localhost:4001

.PHONY: help 
help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $\$1, $\$2}'

.PHONY: db/up
db/up: ## Runs database in docker
	docker compose up -d postgres_db

.PHONY: run/gateway
run/gateway: ## Runs gateway service
	cd services/gateway && go run ./cmd/api/main.go -storage=${STORAGE}

.PHONY: run/storage
run/storage: ## Runs storage service
	cd services/storage && go run ./cmd/api/main.go -db-dsn=${DSN}


