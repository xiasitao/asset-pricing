
build:
	@mkdir -p dist
	@go build -o dist/asset-pricing main.go

install-azure-tools-apt:
	@sh deployment/azure-scripts/install_azure_function_tools_apt.sh

.PHONY: build install-azure-tools-apt