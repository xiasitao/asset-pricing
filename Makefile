build:
	@mkdir -p dist
	@go build -o dist/asset-pricing main.go
	@cp -a dist deployment/azure-functions

run:
	@docker compose up --build

push-image:
	@docker compose build
	@docker tag asset-pricing-api:latest xiasitao/asset-pricing-api:latest
	@docker push xiasitao/asset-pricing-api:latest


install-azure-tools-apt:
	@sh deployment/azure-functions/azure-scripts/install_azure_function_tools_apt.sh

.PHONY: build install-azure-tools-apt push-image