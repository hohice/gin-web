GING.Name:=ginS
GING.Gopkg:=github.com/hohice/gin-web
GING.Version:=$(GING.Gopkg)/pkg/version

DOCKER_TAGS=latest

define get_build_flags
    $(eval SHORT_VERSION=$(shell git describe --tags --always --dirty="-dev"))
    $(eval SHA1_VERSION=$(shell git show --quiet --pretty=format:%H))
	$(eval DATE=$(shell date -u '%Y-%m-%d %H:%M:%S'))
	$(eval BUILD_FLAG= -X $(GING.Version).ShortVersion="$(SHORT_VERSION)" \
		-X $(GING.Version).GitSha1Version="$(SHA1_VERSION)" \
		-X $(GING.Version).BuildDate="$(DATE)")
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

#all:init-vendor/update-vendor update-builder  build
.PHONY: all
all:swag  build

.PHONY: swag
swag:
	@swag init -g router/routers.go
	@echo "gen-swagger done"


.PHONY: build
build:
	@echo "build" $(GING.Name):$(DOCKER_TAGS)
	@docker build --rm -t $(GING.Name):$(DOCKER_TAGS) .
	@docker image prune -f 1>/dev/null 2>&1
	@echo "build done"

.PHONY: install
install:
	$(call get_build_flags)
	time go install -v -ldflags '$(BUILD_FLAG)' $(GING.Gopkg)/cmd/walm

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

