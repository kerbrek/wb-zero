.DEFAULT_GOAL := help

SHELL := /usr/bin/env bash

ifeq ($(shell uname -s), Linux)
	OPEN := xdg-open
else
	OPEN := open
endif

project := wb-zero

.PHONY: lint ## Run linter
lint:
	go vet ./...
	staticcheck ./...

test-env-vars := POSTGRES_PORT=5433 STAN_PORT=4223

.PHONY: prepare-test-services
prepare-test-services:
	@echo Starting Compose test services...
	@env ${test-env-vars} docker-compose -p ${project}-test -f docker-compose.dev.yml up -d
	@echo

stop-prepared-test-services := echo; \
	echo Stopping Compose test services...; \
	docker-compose -p ${project}-test down

.PHONY: test ## Run tests
test: prepare-test-services
	@trap '${stop-prepared-test-services}' EXIT && \
		echo Starting tests... && \
		cd app/ && \
		dotenv -f ./../.env.example run -- env ${test-env-vars} go test $(or $(f), ) ./...

.PHONY: coverage ## Run tests with coverage report
coverage: prepare-test-services
	@trap '${stop-prepared-test-services}' EXIT && \
		echo Starting tests... && \
		cd app/ && \
		dotenv -f ./../.env.example run -- env ${test-env-vars} \
			go test $(or $(f), ) -coverprofile=./../coverage.out -coverpkg ./... ./...
	@go tool cover -html=coverage.out -o coverage.html
	@${OPEN} coverage.html

.PHONY: services ## Start Compose dev services
services:
	@trap 'docker-compose -p ${project}-dev down' EXIT && \
		echo Starting Compose services... && \
		docker-compose -p ${project}-dev --env-file .env.example -f docker-compose.dev.yml up

.PHONY: stop-db # Stop db container to cause a db connection error
stop-db:
	docker-compose -p ${project}-dev stop db

.PHONY: start-db
start-db:
	docker-compose -p ${project}-dev start db

.PHONY: run ## Start application (requires running Compose dev services)
run:
	@echo Starting application...
	@cd app/ && dotenv -f ./../.env.example run -- env GIN_MODE=debug go run .

.PHONY: produce ## Publish several messages to NATS Streaming server
produce:
	@echo Starting producer...
	go run ./producer/ -n=$(or $(n), 5)
	@echo Done.

.PHONY: up ## Start Compose services
up:
	docker-compose -p ${project} -f docker-compose.yml up --build

.PHONY: down ## Stop Compose services
down:
	docker-compose -p ${project} down

.PHONY: ci-prepare
ci-prepare:
	@echo Building test docker images...
	@docker-compose -p ${project}-ci -f docker-compose.ci.yml build

.PHONY: ci-lint
ci-lint: ci-prepare
	@echo Starting go vet...
	@docker-compose -p ${project}-ci -f docker-compose.ci.yml run --no-deps --rm app go vet ./...
	@echo Starting staticcheck...
	@docker-compose -p ${project}-ci -f docker-compose.ci.yml run --no-deps --rm app staticcheck ./...

.PHONY: ci-test
ci-test: ci-lint
	@trap 'docker-compose -p ${project}-ci down' EXIT && \
		echo Starting tests... && \
		docker-compose -p ${project}-ci -f docker-compose.ci.yml run --rm --workdir="/code/app" app go test -vet=off ./...

# .PHONY: gitlab-prepare-test-services
# gitlab-prepare-test-services:
# 	@echo Starting Compose test services...
# 	@docker-compose -p ${project}-gitlab-test --env-file .env.example -f docker-compose.dev.yml up -d
# 	@echo

# gitlab-stop-prepared-test-services := echo; \
# 	echo Stopping Compose test services...; \
# 	docker-compose -p ${project}-gitlab-test down

# gitlab-test-env-vars := POSTGRES_HOST=docker STAN_HOST=docker

# .PHONY: gitlab-test
# gitlab-test: gitlab-prepare-test-services
# 	@trap '${gitlab-stop-prepared-test-services}' EXIT && \
# 		echo Starting tests... && \
# 		cd app/ && \
# 		dotenv -f ./../.env.example run -- env ${gitlab-test-env-vars} go test -vet=off ./...

.PHONY: help ## Print list of targets with descriptions
help:
	@echo; \
		for mk in $(MAKEFILE_LIST); do \
			echo \# $$mk; \
			grep '^.PHONY: .* ##' $$mk \
			| sed 's/\.PHONY: \(.*\) ## \(.*\)/\1	\2/' \
			| expand -t20; \
			echo; \
		done
