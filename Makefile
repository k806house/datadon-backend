
.PHONY: build
build:
	./scripts/build

.PHONY: deploy
deploy:
	cd terraform && terraform apply

.PHONY: init
init:
	cd terraform && terraform init