TARG.Name:=pm
TRAG.Gopkg:=transwarp/product-mata-serv
TRAG.Version:=$(TRAG.Gopkg)/pkg/version

DOCKER_TAGS=latest

define get_build_flags
    $(eval SHORT_VERSION=$(shell git describe --tags --always --dirty="-dev"))
    $(eval SHA1_VERSION=$(shell git show --quiet --pretty=format:%H))
	$(eval DATE=$(shell date +'%Y-%m-%dT%H:%M:%S'))
	$(eval BUILD_FLAG= -X $(TRAG.Version).ShortVersion="$(SHORT_VERSION)" \
		-X $(TRAG.Version).GitSha1Version="$(SHA1_VERSION)" \
		-X $(TRAG.Version).BuildDate="$(DATE)")
endef


.PHONY: init-vendor
init-vendor:
	glide init
	@echo "init-vendor done"

.PHONY: update-vendor
update-vendor:
	glide update
	@echo "update-vendor done"

.PHONY: update-builder
update-builder:
	docker pull 172.16.1.99/transwarp/walm-builder:1.0
	@echo "update-builder done"

#all:init-vendor/update-vendor update-builder generate build
.PHONY: all
all:swag generate build

.PHONY: generate
generate:
	go generate ./pkg/version/
	@echo "generate done"

.PHONY: swag
swag:
	@swag init -g router/routers.go
	@echo "gen-swagger done"


.PHONY: build
build:
	@echo "build" $(TARG.Name):$(DOCKER_TAGS)
	@docker build -t $(TARG.Name):$(DOCKER_TAGS) .
	@docker image prune -f 1>/dev/null 2>&1
	@echo "build done"

.PHONY: install
install:
	$(call get_build_flags)
	time go install -v -ldflags '$(BUILD_FLAG)' $(TRAG.Gopkg)/cmd/walm

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
	CGO_ENABLED=0 go test -v -a -tags="e2e" ./test/...
	@echo "e2e-test done"


.PHONY: clean
clean:
	-make -C ./pkg/version clean
	@echo "ok"

