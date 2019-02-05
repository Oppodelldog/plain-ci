BUILD_ARTIFACTS = ".build-artifiacts"
BINARY_NAME = "plainci"

generate-assets: ## generates static assets
	statics -i=webview/assets/templates -o=webview/assets/templates.go -pkg=assets -group=Templates -ignore=\.gitignore -prefix=webview/assets
	statics -i=webview/assets/images    -o=webview/assets/images.go -pkg=assets -group=Images -ignore=\.gitignore -prefix=webview/assets
	
build: generate-assets ##
	go build -o $(BUILD_ARTIFACTS)/$(BINARY_NAME) cmd/main.go

deploy: build ## deploy freshly build binary to server
	deploy-plainci

# Self-Documented Makefile see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help