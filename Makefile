dk-build: ## Build docker image with compiled binary (args VERSION, REPOSITORY, TAG required)
	docker build -t moxy ./