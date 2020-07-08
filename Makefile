.SHELLFLAGS = -c

.SILENT: ;               # no need for @
.ONESHELL: ;             # recipes execute in same shell
.NOTPARALLEL: ;          # wait for this target to finish
.EXPORT_ALL_VARIABLES: ; # send all vars to shell

.PHONY: help

help: ## Show Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

compile: ## Build the app
	packr2
	go build -o moxy
	packr2 clean

clean:
	packr2 clean

release: ## Publish a release
	goreleaser
	packr2 clean

snapshot:
	goreleaser --snapshot --skip-publish --rm-dist
	packr2 clean

dk-build: ## Build docker image with compiled binary (args VERSION, REPOSITORY, TAG required)
	docker build -t moxy ./