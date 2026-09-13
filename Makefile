ENV ?= local

TFVARS := environments/$(ENV)/$(ENV).tfvars

.PHONY: fmt fmt-go validate plan deploy init destroy apply

apply: deploy 

init:
	@echo "Start init"
	terraform init

plan: init
	@echo "Planning $(ENV)..."
	terraform plan  -var "environment=$(ENV)" -var-file="$(TFVARS)" -out="tf-$(ENV)-plan.out"

deploy: init
	@echo "Applying $(ENV)..."
	@plan=$$(find . -maxdepth 1 -type f -name "tf-$(ENV)-plan.out" -print -quit); \
	if [ -n "$$plan" ]; then \
		echo "Using plan: $$plan"; \
		terraform apply -var "environment=$(ENV)" -var-file="$(TFVARS)" -auto-approve "$$plan"; \
	else \
		echo "Plan file not found. Running deploy without plan file"; \
		terraform apply -var "environment=$(ENV)" -var-file="$(TFVARS)" -auto-approve; \
	fi

validate:
	@echo "Validating Terraform..."
	terraform validate

destroy:
	@echo "Destroying environment: $(ENV)"
	terraform destroy -var "environment=$(ENV)" -var-file="$(TFVARS)"

fmt:
	@echo "Start formatting"
	@formatted=$$(terraform fmt -recursive); \
	if [ -n "$$formatted" ]; then \
		echo "$$formatted"; \
		count=$$(echo "$$formatted" | wc -l); \
	else \
		count=0; \
	fi; \
	echo "Formatted files: $$count"

fmt-go:
	@echo "Start formatting Go code"
	@files=$$(find src -name "*.go" -not -path "*/vendor/*"); \
	formatted=$$(gofmt -l -w $$files); \
	if [ -n "$$formatted" ]; then \
		echo "$$formatted"; \
		count=$$(echo "$$formatted" | wc -l); \
	else \
		count=0; \
	fi; \
	echo "Formatted files: $$count"