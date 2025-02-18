build:
	@mkdir -p dist
	@go build -o dist/asset-pricing main.go

run:
	@cd deployment/docker-compose && docker compose up --build

push-image:
	@cd deployment/docker-compose && docker compose build
	@docker tag asset-pricing-api:latest xiasitao/asset-pricing-api:latest
	@docker push xiasitao/asset-pricing-api:latest

.PHONY: build run push-image