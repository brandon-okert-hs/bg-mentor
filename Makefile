FRONTEND_DIR := frontend
ASSETS_DIR := assets
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

TUNNEL_HOST=""
DB_NAME=""
MIGRATION_USER=""
DB_HOST=""
ifeq ($(env),dev)
TUNNEL_HOST=dev.borngosugaming.com
DB_NAME=borngosu_dev
MIGRATION_USER=migrations-dev
DB_HOST=borngosudb-production-0.caqyd4okx8bs.us-east-1.rds.amazonaws.com
endif
ifeq ($(env),production)
TUNNEL_HOST=borngosugaming.com
DB_NAME=borngosu
MIGRATION_USER=migrations
DB_HOST=borngosudb-production-0.caqyd4okx8bs.us-east-1.rds.amazonaws.com
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

build: build-webserver build-frontend build-assets

build-webserver:
	rm -rf $(ARTIFACT_DIR)/$(env)/webserver
	mkdir -p $(ARTIFACT_DIR)/$(env)
	$(BUILD_CMD) -o $(ARTIFACT_DIR)/$(env)/webserver ./cmd/webserver

build-frontend:
	rm -rf $(ARTIFACT_DIR)/$(env)/static/js
	rm -f $(ARTIFACT_DIR)/$(env)/static/index.html
	mkdir -p $(ARTIFACT_DIR)/$(env)/static
	$(WEBPACK_CMD)

build-assets:
	mkdir -p $(ARTIFACT_DIR)/$(env)/static
	cp $(ASSETS_DIR)/* $(ARTIFACT_DIR)/$(env)/static

# Runs locally, but simulates the environment that would be deployed
# Will not work in envs other than local if local machine is not the same os as deploy machines
run: decrypt-secrets-config
	$(ARTIFACT_DIR)/$(env)/webserver $(DEPLOY_DIR)/files/$(env)/etc/bg-mentor

decrypt-secrets-config:
	@mkdir -p $(DECRYPTED_SECRETS_DIR)/$(env)
ifeq ("","$(wildcard $(CONFIG_DIR)/config.json)")
	@echo "Decrypting Configuration Secrets. You'll need permission."
	@ansible-vault decrypt $(CONFIG_DIR)/config.json.secret --output $(CONFIG_DIR)/config.json
	@(sleep 1800 && rm -rf $(CONFIG_DIR)/config.json &)
endif

decrypt-secrets-deploy:
	@mkdir -p $(DECRYPTED_SECRETS_DIR)/$(env)
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

tunnel-db: decrypt-secrets-deploy
ifeq ($(env),local)
	@echo "Cannot tunnel to local db. Local test server uses dev db"
	exit 1
endif
	ssh -i ${DECRYPTED_SECRETS_DIR}/$(env)/deploy-key.pem ubuntu@$(TUNNEL_HOST) -L 3307:$(DB_HOST):3306 -N

baseline-db:
	flyway -user="$(MIGRATION_USER)" -locations="filesystem:$(shell pwd)/deploy/db/migrations" -url='jdbc:mysql://127.0.0.1:3307/$(DB_NAME)' -baselineVersion='0' baseline

migrate-db:
	flyway -user="$(MIGRATION_USER)" -locations="filesystem:$(shell pwd)/deploy/db/migrations" -url='jdbc:mysql://127.0.0.1:3307/$(DB_NAME)' migrate

deploy: decrypt-secrets-deploy decrypt-secrets-config
	./scripts/deploy.sh $(env)
