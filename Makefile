.PHONY: dev prod test help

# Default action when you just type 'make'
.DEFAULT_GOAL := help

## dev: Run the application in developer mode (.dev.env)
dev:
	@echo "🚀 Starting in DEVELOPER mode..."
	APPENV=development air

## prod: Run the application in production mode (.prod.env)
prod:
	@echo "🔒 Starting in PRODUCTION mode..."
	APPENV=production go run cmd/app/main.go

## test: Run integration tests (.test.env)
test:
	@echo "🧪 Running INTEGRATION TESTS..."
	APPENV=integration-test go test ./... -v

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed -e 's/## //' | awk 'BEGIN {FS = ": "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
