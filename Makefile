build:
	@mkdir -p dist
	@go build -o dist/asset-pricing main.go
	@cp -a dist deployment/azure-functions

run:
	@docker compose up --build

install-azure-tools-apt:
	@sh deployment/azure-functions/azure-scripts/install_azure_function_tools_apt.sh

.PHONY: build install-azure-tools-apt