VERSION := $(shell git describe --tag --always --dirty)

PROJECT := covidvax
TESTTIMEOUT := 30s

TESTDB := 'postgres://admin:admin-pwd@localhost:5432/covidvax?sslmode=disable'

MIGRATE_VERSION := v4.15.1
MIGRATE_CLI := .deps/migrate

ifeq ($(UNAME), Darwin)
VENDOR_MIGRATE_URL = https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.darwin-amd64.tar.gz
else
VENDOR_MIGRATE_URL = https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz
endif

.PHONY: version
version: ## display version
	@echo $(VERSION)

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

vendor/migrate:
	mkdir -p vendor
	cd vendor && curl -L $(VENDOR_MIGRATE_URL) | tar xvz migrate

.PHONY: migrations
migrations: vendor/migrate ## run migration on test db
	./vendor/migrate -source file://db/migrations -database $(TESTDB) ${OPTS}

tests: ## Run all the tests
	echo 'mode: atomic' > coverage.txt && go test -covermode=atomic -coverprofile=coverage.txt -race -timeout=$(TESTTIMEOUT) ./...
