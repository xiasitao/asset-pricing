
build:
	@mkdir -p dist
	@go build -o dist/asset-pricing main.go

.PHONY: build