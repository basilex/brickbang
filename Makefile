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

	@echo '>>> Entities template generator section'
	@echo '  - gen-mapper     : Generate entity mapper layer'
	@echo '  - gen-repository : Generate entity repository layer'
	@echo '  - gen-service    : Generate entity service layer'
	@echo '  - gen-controller : Generate entity controller layer'
	@echo '  - gen-all        : Generate all layers for the entity'
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
	@make -C storage gen
dbs-up:
	@make -C storage up
dbs-up1:
	@make -C storage up1
dbs-down:
	@make -C storage down
dbs-down1:
	@make -C storage down1
dbs-drop:
	@make -C storage drop
dbs-version:
	@make -C storage version

.PHONY: dbs-gen dbs-up dbs-up1 dbs-down dbs-down1 dbs-drop dbs-version
#
# Entities template generator section
#
BLUE := \033[0;34m
GREEN := \033[0;32m
RESET := \033[0m

ENTITY ?= None
MODULE_PATH ?= $(sys)
ENTITY_LOWER := $(shell echo $(ENTITY) | tr '[:upper:]' '[:lower:]')
#
# Generic sed replace (cross-platform)
#
define SED_REPLACE
	sed -i '' "s/{{.Entity}}/$(ENTITY)/g; s/{{.EntityLower}}/$(ENTITY_LOWER)/g; s|{{.ModulePath}}|$(MODULE_PATH)|g" $(1)
endef
#
# Universal entity layer generator
# Usage: $(call GEN_ENTITY_LAYER,<template_name>,<target_dir>)
#
define GEN_ENTITY_LAYER
	@cp internal/template/$(1).tpl internal/$(2)/$(ENTITY_LOWER)_$(subst entity_,,$(1)).go
	@$(call SED_REPLACE,internal/$(2)/$(ENTITY_LOWER)_$(subst entity_,,$(1)).go)
	@echo "$(GREEN)[ok]$(RESET) $(subst entity_,,$(1)) for $(BLUE)$(ENTITY)$(RESET) created → internal/$(2)/$(ENTITY_LOWER)_$(subst entity_,,$(1)).go"
endef

gen-mapper:
	$(call GEN_ENTITY_LAYER,entity_mapper,mapper)

gen-repository:
	$(call GEN_ENTITY_LAYER,entity_repository,repository)

gen-service:
	$(call GEN_ENTITY_LAYER,entity_service,service)

gen-controller:
	$(call GEN_ENTITY_LAYER,entity_controller,controller)

gen-all: gen-mapper gen-repository gen-service gen-controller
	@echo "$(GREEN)[done]$(RESET) All layers scaffolded for entity $(BLUE)$(ENTITY)$(RESET)"

.PHONY: gen-mapper gen-repository gen-service gen-controller gen-all

#
# eof
#
