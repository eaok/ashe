.PHONY: default run debug test help

default: debug

run: ## run mode
	@sed -i 's/app_mode = \(debug\|run\)/app_mode = run/g' config/config.ini
	@go build
	@nohup ./ashe > logs/log&

debug: ## debug mode
	@sed -i 's/app_mode = \(debug\|run\)/app_mode = debug/g' config/config.ini
	@go run main.go

test: ## Run package unit testsS
	@go test -v -race -short ./...

help: ## Displays help menu
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
