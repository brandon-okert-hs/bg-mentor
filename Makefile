FRONTEND_DIR := frontend
DEPLOY_DIR := deploy
TERRAFORM_DIR := $(DEPLOY_DIR)/terraform
SECRETS_DIR := $(DEPLOY_DIR)/secrets
DECRYPTED_SECRETS_DIR := .secrets-decrypted
PLANS_DIR := .plans
ARTIFACT_DIR := artifact

ifneq ($(env),dev)
ifneq ($(env),production)
$(error env must be set to dev or production, like 'make target env=dev')
endif
endif

WEBPACK_ENV=development
ifeq ($(env),production)
WEBPACK_ENV=production
endif

.PHONY: decrypt-secrets

clean:
	rm -rf $(ARTIFACT_DIR)
	rm -rf $(DECRYPTED_SECRETS_DIR)

deps:
	go get .
	go get -t
	yarn --cwd $(FRONTEND_DIR) install

watch-frontend:
	rm -rf $(ARTIFACT_DIR)/$(env)/static
	mkdir -p $(ARTIFACT_DIR)/$(env)/static
	yarn --cwd frontend webpack --env.$(WEBPACK_ENV) --config webpack.config.js --progress --watch

build: clean build-webserver build-frontend

build-webserver:
	rm -rf $(ARTIFACT_DIR)/$(env)/webserver
	mkdir -p $(ARTIFACT_DIR)/$(env)
	env GOOS=linux GOARCH=amd64 go build -o $(ARTIFACT_DIR)/$(env)/webserver ./cmd/webserver

build-frontend:
	rm -rf $(ARTIFACT_DIR)/$(env)/static
	mkdir -p $(ARTIFACT_DIR)/$(env)/static
	yarn --cwd frontend webpack --env.$(env) --config webpack.config.js --progress

run:
	$(ARTIFACT_DIR)/$(env)/webserver

decrypt-secrets:
	rm -rf $(DECRYPTED_SECRETS_DIR)/$(env)
	mkdir -p $(DECRYPTED_SECRETS_DIR)/$(env)
	ansible-vault decrypt $(SECRETS_DIR)/$(env)/deploy-key.pem.secret --output $(DECRYPTED_SECRETS_DIR)/$(env)/deploy-key.pem
	ansible-vault decrypt $(SECRETS_DIR)/$(env)/terraform-aws-credentials.secret --output $(DECRYPTED_SECRETS_DIR)/$(env)/terraform-aws-credentials
	(sleep 1800 && rm -rf $(DECRYPTED_SECRETS_DIR)/$(env) &)

init-infra: decrypt-secrets
	terraform init $(TERRAFORM_DIR)/$(env)

plan-infra: init-infra
	mkdir -p $(PLANS_DIR)/$(env)
	terraform plan -var-file=$(TERRAFORM_DIR)/$(env)/variables.tfvars -out=$(PLANS_DIR)/$(env)/plan $(TERRAFORM_DIR)/$(env)

apply-infra: init-infra
	terraform apply $(PLANS_DIR)/$(env)/plan

deploy: decrypt-secrets
	./scripts/deploy.sh $(env)
