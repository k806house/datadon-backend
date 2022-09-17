
.PHONY: build
build:
	./scripts/build

.PHONY: deploy
deploy:
	./scripts/build && cd terraform && terraform init && terraform apply -auto-approve

.PHONY: init
init:
	cd terraform && terraform init