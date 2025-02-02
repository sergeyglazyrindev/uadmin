define VERSION_CMD =
eval ' \
	define=""; \
	version=`git describe --abbrev=0 --tags | tr -d "[a-z]"` ; \
	[ -z "$$version" ] && version="unknown"; \
	commit=`git rev-parse --verify HEAD`; \
	tagname=`git show-ref --tags | grep $$commit`; \
	if [ -n "$$tagname" ]; then \
		define=`echo $$tagname | awk -F "/" "{print \\$$NF}" | tr -d "[a-z]"`; \
	else \
		define=`printf "$$version-%.12s" $$commit`; \
	fi; \
	tainted=`git ls-files -m | wc -l` ; \
	if [ "$$tainted" -gt 0 ]; then \
		define="$${define}-tainted"; \
	fi; \
	echo "$$define" \
'
endef

EXTRA_BUILD_TARGET=
VERSION?=$(shell $(VERSION_CMD))
GO?=go
BUILD_ID:=$(shell echo 0x$$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \n'))
VERBOSE_FLAGS?=
VERBOSE?=false
ifeq ($(VERBOSE), true)
  VERBOSE_FLAGS+=-v
endif
UADMIN_GITHUB:=github.com/sergeyglazyrindev/uadmin
UADMIN_GITHUB_VERSION:=$(UADMIN_GITHUB)/version.Version=${VERSION}
BUILD_TAGS?=$(TAGS)

include .mk/check.mk
include .mk/dist.mk
include .mk/proto.mk
include .mk/static.mk
include .mk/tests.mk

define GOCOMPILE
CGO_CFLAGS_ALLOW='.*' CGO_LDFLAGS_ALLOW='.*' $(GO) $1 \
                -ldflags="${LDFLAGS} -B $(BUILD_ID) -X $(UADMIN_GITHUB_VERSION)" \
                ${GOFLAGS} -tags="${BUILD_TAGS}" ${VERBOSE_FLAGS} -o /uadmin/uadmin \
                cmd/uadmin/main.go
endef

.PHONY: .build
.build:
	$(call GOCOMPILE,build)

.PHONY: build
build: gopath moddownload .build

.PHONY: .install
.install:
	$(call GOCOMPILE,install)

.PHONY: uadmin.clean
uadmin.clean:
	go clean -i $(UADMIN_GITHUB)

.PHONY: moddownload
moddownload:
ifneq ($(OFFLINE), true)
	go mod download
endif

.PHONY: genlocalfiles
genlocalfiles: $(EXTRA_BUILD_TARGET)
# .proto

.PHONY: touchlocalfiles
touchlocalfiles: .proto.touch

.PHONY: clean
clean: uadmin.clean .proto.clean \
       go clean -i >/dev/null 2>&1 || true

.PHONY: docker
docker:
	docker build . -t $(DOCKER_IMAGE):$(DOCKER_TAG)