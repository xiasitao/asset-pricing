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

.PHONY: build run push-image