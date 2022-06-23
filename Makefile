.PHONY: fmt test gen clean run help

# command aliases
test := ENV=test go test ./...

fmt: ## Run gofmt
	find . -iname *.go -type f -exec gofmt -w -s {} \;

test: ## Run go vet, then test all files
	go vet ./...
	$(test)

update-snapshots: ## Update snapshots during a go test. Must have cupaloy
	ENV=test UPDATE_SNAPSHOTS=true go test ./...

gen: ## go generate ./...
	go generate ./...

clean: fmt gen ## gofmt, go generate, then go mod tidy, and finally rm -rf bin/
	go mod tidy
	rm -rf ./bin

help: ## Print help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@printf "\033[36m%-30s\033[0m %s\n" "(target)" "Build a target binary: $(targets)"
