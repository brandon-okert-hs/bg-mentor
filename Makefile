FRONTEND_DIR := frontend
DEPLOY_DIR := deploy
TERRAFORM_DIR := $(DEPLOY_DIR)/terraform
SECRETS_DIR := $(DEPLOY_DIR)/secrets
CONFIG_DIR := $(DEPLOY_DIR)/files/$(env)/etc/bg-mentor
DECRYPTED_SECRETS_DIR := .secrets-decrypted
PLANS_DIR := .plans
ARTIFACT_DIR := artifact

ifneq ($(env),dev)
ifneq ($(env),production)
ifneq ($(env),local)
$(error env must be set to local, dev, or production, like 'make target env=dev')
endif
endif
endif

WEBPACK_CMD=yarn --cwd frontend webpack --env.SERVER_ENV=$(env) --config webpack.config.js --progress
BUILD_CMD=env GOOS=linux GOARCH=amd64 go build
ifeq ($(env),local)
BUILD_CMD=go build
endif

.PHONY: decrypt-secrets

clean:
	rm -rf $(ARTIFACT_DIR)
	rm -rf $(DECRYPTED_SECRETS_DIR)
	rm -rf coverage.out

deps:
	go get .
	go get -t
	yarn --cwd $(FRONTEND_DIR) install

coverage:
	go test -covermode=count -coverprofile=coverage.out && go tool cover -html=coverage.out

watch-frontend:
	rm -rf $(ARTIFACT_DIR)/$(env)/static
	mkdir -p $(ARTIFACT_DIR)/$(env)/static
	$(WEBPACK_CMD) --watch

build: build-webserver build-frontend

build-webserver:
	rm -rf $(ARTIFACT_DIR)/$(env)/webserver
	mkdir -p $(ARTIFACT_DIR)/$(env)
	$(BUILD_CMD) -o $(ARTIFACT_DIR)/$(env)/webserver ./cmd/webserver

build-frontend:
	rm -rf $(ARTIFACT_DIR)/$(env)/static
	mkdir -p $(ARTIFACT_DIR)/$(env)/static
	$(WEBPACK_CMD)

# Runs locally, but simulates the environment that would be deployed
# Will not work in envs other than local if local machine is not the same os as deploy machines
run:
	$(ARTIFACT_DIR)/$(env)/webserver $(DEPLOY_DIR)/files/$(env)/etc/bg-mentor

decrypt-secrets-deploy:
	@mkdir -p $(DECRYPTED_SECRETS_DIR)/$(env)
ifeq ("","$(wildcard $(CONFIG_DIR)/config.json)")
	@echo "Decrypting Deploy Configuration Secrets. You'll need permission."
	@ansible-vault decrypt $(CONFIG_DIR)/config.json.secret --output $(CONFIG_DIR)/config.json
	@(sleep 1800 && rm -rf $(CONFIG_DIR)/config.json.secret &)
endif
ifeq ("","$(wildcard $(DECRYPTED_SECRETS_DIR)/$(env)/deploy-key.pem)")
	@echo "Decrypting Deploy Secrets. You'll need permission."
	@ansible-vault decrypt $(SECRETS_DIR)/$(env)/deploy-key.pem.secret --output $(DECRYPTED_SECRETS_DIR)/$(env)/deploy-key.pem
	@(sleep 1800 && rm -rf $(DECRYPTED_SECRETS_DIR)/$(env)/deploy-key.pem &)
endif

decrypt-secrets-infra:
	@mkdir -p $(DECRYPTED_SECRETS_DIR)/$(env)
ifeq ("","$(wildcard $(DECRYPTED_SECRETS_DIR)/$(env)/terraform-aws-credentials)")
	@echo "Decrypting Infra Secrets. You'll need permission."
	@ansible-vault decrypt $(SECRETS_DIR)/$(env)/terraform-aws-credentials.secret --output $(DECRYPTED_SECRETS_DIR)/$(env)/terraform-aws-credentials
	@(sleep 1800 && rm -rf $(DECRYPTED_SECRETS_DIR)/$(env)/terraform-aws-credentials &)
endif

init-infra: decrypt-secrets-infra
	terraform init $(TERRAFORM_DIR)/$(env)

plan-infra: init-infra
	mkdir -p $(PLANS_DIR)/$(env)
	terraform plan -var-file=$(TERRAFORM_DIR)/$(env)/variables.tfvars -out=$(PLANS_DIR)/$(env)/plan $(TERRAFORM_DIR)/$(env)

apply-infra: init-infra
	terraform apply $(PLANS_DIR)/$(env)/plan

deploy: decrypt-secrets-deploy
	./scripts/deploy.sh $(env)
