#
# dev|test|prod
#
stage=dev

-include Makefile.inc
-include Makefile.$(stage)

binaries = go sqlc migrate

checkfor := $(foreach exec,$(binaries), \
	$(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

dburl = "postgres://$(user):$(password)@$(host):$(port)/$(database)?$(options)"

gen:
	@sqlc generate
#
# migrate create -ext sql -dir db/migrations -seq <title>
#
up:
	@migrate -path schema -database $(dburl) up

up1:
	@migrate -path schema -database $(dburl) up 1

down:
	@migrate -path schema -database $(dburl) down

down1:
	@migrate -path schema -database $(dburl) down 1

drop:
	@migrate -path schema -database $(dburl) drop

version:
	@migrate -path schema -database $(dburl) version

.PHONY: gen up up1 down down1 drop version

