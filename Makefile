#
# Makefile
# BrickBang Platform Service
# Copyright (C) 2025 by BrickBang Inc. All Rights Reserved
# --------------------------------------------------------
#
svc := $(or $(BRB_SVC),brb)
env := $(or $(BRB_ENV),dev)
ver := $(or $(BRB_VER),0.1.0)
sys := $(or $(BRB_SYS),brickbang)

dist := $(or $(BRB_DIST),dist)
cert := $(or $(BRB_CERT),./resource/cert)

image := $(or $(BRB_IMAGE),$(sys)/$(svc):$(ver))
#
# Metadata to be embedded
#
version := $(ver)
staging := $(env)
githash := $(shell git rev-parse --short=8 HEAD)
gobuild := $(shell go version | sed -e "s/go version //g;s/ /-/g")
compile := $(shell date "+%FT%T.%N%:z")
#
# Go build: ldflags linker parameters
#
ldflags += -X $(sys)/internal.Version=$(version)
ldflags += -X $(sys)/internal.Staging=$(staging)
ldflags += -X $(sys)/internal.Githash=$(githash)
ldflags += -X $(sys)/internal.Gobuild=$(gobuild)
ldflags += -X $(sys)/internal.Compile=$(compile)
#
# Check the dependent bins
#
bins = go openssl

-include internal/storage/Makefile.inc

checkfor := $(foreach exec,$(bins), \
	$(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))
#
# Main entry point
#
all:
	@echo '*** BrickBang Makefile sections'
	@echo '    ---------------------------'

	@echo '>>> Dbs management section'
	@echo '  - dbs-gen     : Generate sqlc db layer'
	@echo '  - dbs-up      : Install db schema and default data'
	@echo '  - dbs-up1     : Migrate one level of the db schema'
	@echo '  - dbs-down    : Uninstall db schema (all the data purged)'
	@echo '  - dbs-down1   : Migrate down one level of the db schema'
	@echo '  - dbs-drop    : Drop entire db schema (all the data purged)'
	@echo '  - dbs-version : Show the db migration version'
	@echo

	@echo '>>> App management section'
	@echo '  - app-tidy    : Ensure that all imports are satisfied'
	@echo '  - app-build   : Build the application inside the linux container'
	@echo '  - app-up      : Run all the containers from the docker composer yml'
	@echo '  - app-down    : Shut down all the docker compose containers'
	@echo '  - app-clean   : Remove all the docker exited containers'
	@echo '  - app-cert    : Generate app TLS/SSL certificates'
	@echo '  - app-prune   : Prune all in the local docker env'
	@echo

	@exit 0

.PHONY: all
#
# Swagger section
#
docs:
	swag init

.PHONY: docs
#
# App section
#
app-tidy:
	@go mod tidy
app-build:
	@go build -a -ldflags="$(ldflags)" -o $(dist)/$(svc) main.go

.PHONY: vet test ci
vet:
	@go vet ./...

test:
	@go test ./... -v

ci: app-tidy vet test
	@echo "CI checks completed"

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: test-integration test-integration-all
test-integration:
	@docker info > /dev/null || (echo "Docker not available"; exit 1)
	@go test -tags=integration ./internal/repository/dbs -v

test-integration-all:
	@docker info > /dev/null || (echo "Docker not available"; exit 1)
	@go test -tags=integration ./... -v

.PHONY: app-tidy app-build
#
# Composer section
#
app-up:
	@docker compose -f $(svc)-local.yml up --build
app-down:
	@docker compose -f $(svc)-local.yml down  --remove-orphans
app-clean:
	@docker rm -v $(shell docker ps --filter status=exited -q)
	@docker rmi $(image)
app-cert:
	# <dev> mode using self signed certificate
	# <prod> mode using let's encrypt for the real domain
	@openssl req -x509 -newkey rsa:4096 \
		-keyout $(cert)/$(svc).key -out $(cert)/$(svc).crt -days 365 -nodes
app-prune:
	@docker system prune -af

.PHONY: app-up app-down app-clean app-cert app-prune
#
# Dbs section
#
dbs-gen:
	@make -C internal/repository gen
dbs-up:
	@make -C internal/repository up
dbs-up1:
	@make -C internal/repository up1
dbs-down:
	@make -C internal/repository down
dbs-down1:
	@make -C internal/repository down1
dbs-drop:
	@make -C internal/repository drop
dbs-version:
	@make -C internal/repository version

.PHONY: dbs-gen dbs-up dbs-up1 dbs-down dbs-down1 dbs-drop dbs-version

.PHONY: gen-crud

