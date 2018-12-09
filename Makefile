STACK_NAME ?= aws-lambda-go-turbo
ARTIFACT ?= runtime.zip
DEPLOY_BUCKET ?= sam-deployment

# creates .env with .env.template if it doesn't exist already
.env:
	cp -f .env.template .env

# Entrypoints
test: .env
	docker-compose run --rm build-go make _test
.PHONY: test

build: .env
	docker-compose run --rm build-go make _build
.PHONY: build

deploy: .env
	docker-compose run --rm aws make _deploy
.PHONY: deploy

remove: .env
	docker-compose run --rm aws make _remove
.PHONY: remove

# Helpers
shellGo: .env
	docker-compose run --rm build-go /bin/bash
.PHONY: shellGo

shellAWS: .env
	docker-compose run --rm aws /bin/bash
.PHONY: shellAWS

# Internal targets
_deps:
	dep ensure -v

_test: _deps
	echo $(GOPATH)
	go test -v ./...

_build: _deps
	GOOS=linux go build -ldflags="-w -s" -o bin/bootstrap runtime/main.go
	zip -1j $(ARTIFACT) bin/bootstrap

_validate:
	sam validate

# _package: $(ARTIFACT)
# 	aws cloudformation package --template-file template.yml --s3-bucket $(DEPLOY_BUCKET) --output-template-file packaged.yml

# _deploy: _package
# 	aws cloudformation deploy \
# 		--template-file ./packaged.yml \
# 		--stack-name $(STACK_NAME) \
# 		--capabilities CAPABILITY_IAM \
# 		--no-fail-on-empty-changeset

_deploy: $(ARTIFACT)
	aws lambda create-function \
		--function-name $(STACK_NAME) \
		--zip-file fileb://$(ARTIFACT) \
		--handler handler.Hello \
		--runtime provided \
		--role $(LAMBDA_ROLE)

_describe: _assumeRole
	aws cloudformation describe-stack-events --stack-name $(STACK_NAME)

_remove: _assumeRole
	aws cloudformation delete-stack --stack-name $(STACK_NAME)

_assumeRole:
ifndef AWS_SESSION_TOKEN
	$(eval ROLE = "$(shell aws sts assume-role --role-arn "$(AWS_ROLE)" --role-session-name "sam-assume-role" --query "Credentials.[AccessKeyId, SecretAccessKey, SessionToken]" --output text)")
	$(eval export AWS_ACCESS_KEY_ID = $(shell echo $(ROLE) | cut -f1))
	$(eval export AWS_SECRET_ACCESS_KEY = $(shell echo $(ROLE) | cut -f2))
	$(eval export AWS_SESSION_TOKEN = $(shell echo $(ROLE) | cut -f3))
endif
