include .env

GINS.Name:=ginS
GINS.Gopkg:=github.com/hohice/gin-web
GINS.Version:=$(GINS.Gopkg)/pkg/version

DOCKER_TAGS=latest

define get_build_flags
    $(eval SHORT_VERSION=$(shell git describe --tags --always --dirty="-dev"))
    $(eval SHA1_VERSION=$(shell git show --quiet --pretty=format:%H))
	$(eval DATE=$(shell date -u '%Y-%m-%d %H:%M:%S'))
	$(eval BUILD_FLAG= -X $(GINS.Version).ShortVersion="$(SHORT_VERSION)" \
		-X $(GINS.Version).GitSha1Version="$(SHA1_VERSION)" \
		-X $(GINS.Version).BuildDate="$(DATE)")
endef

default: help

##init-mod: init package module by go mod 
.PHONY: init-mod
init-mod:
	go mod init
	@echo "init-mod done"

##down-mod: download dependce package module by go mod 
.PHONY: down-mod
down-mod:
	go mod download
	@echo "down-mod done"

##update-builder: pull newest builder image
.PHONY: update-builder
update-builder:
	docker pull 172.16.1.99/transwarp/walm-builder:1.0
	@echo "update-builder done"

#all:init-vendor/update-vendor update-builder  build
.PHONY: all
all:swag  build

.PHONY: swag
swag:
	@swag init -g server/routers.go
	@echo "gen-swagger done"


.PHONY: build
build:
	@echo "build" $(GINS.Name):$(DOCKER_TAGS)
	@docker build --rm -t $(GINS.Name):$(DOCKER_TAGS) .
	@docker image prune -f 1>/dev/null 2>&1
	@echo "build done"

.PHONY: install
install:
	$(call get_build_flags)
	time go install -v -ldflags '$(BUILD_FLAG)' $(GINS.Gopkg)/cmd/$(GINS.Name)

.PHONY: test
test:
	@make unit-test
	@make e2e-test
	@echo "test done"

.PHONY: unit-test
unit-test:
	CGO_ENABLED=0 go test -v -a -tags="unit" ./...
	@echo "unit-test done"

.PHONY: e2e-test
e2e-test:
	CGO_ENABLED=0 go test -v -a -tags="e2e db" ./test/...
	@echo "e2e-test done"


.PHONY: clean
clean:
	go clean
	@echo "ok"

.PHONY: help
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

