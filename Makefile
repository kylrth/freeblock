all: fmt-check lint-check test build

BINDIR := bin

LINTER_VERSION := 1.44.0
LINTER := $(BINDIR)/golangci-lint_$(LINTER_VERSION)
DEV_OS := $(shell uname -s | tr A-Z a-z)

$(LINTER):
	mkdir -p $(BINDIR)
	wget "https://github.com/golangci/golangci-lint/releases/download/v$(LINTER_VERSION)/golangci-lint-$(LINTER_VERSION)-$(DEV_OS)-amd64.tar.gz" -O - \
		| tar -xz -C $(BINDIR) --strip-components=1 --exclude=README.md --exclude=LICENSE
	mv $(BINDIR)/golangci-lint $(LINTER)

.PHONY: fmt-check
fmt-check:
	BADFILES=$$(gofmt -l -s -d $$(find . -type f -name '*.go')) && [ -z "$$BADFILES" ] && exit 0

.PHONY: lint-check
lint-check: $(LINTER)
	$(LINTER) run --deadline=2m

# This allows you to run specific tests, for example:
#   - `make test TESTTARGET=./pkg/hosts`
#   - `make test "TESTTARGET=./... -run ^TestReadHosts\$$ -bench=XXX"  # use "\$$" to produce "$" in Make
# If you want to see test coverage, run the following:
# 	make test 'TESTTARGET=./... -coverprofile=/repo/cov.out'
# 	go tool cover -html=cov.out -o cov.html
# and then take a look at cov.html in your browser.
TESTTARGET = ./...

.PHONY: test
test:
	go test -cover -race -bench=. ${TESTTARGET}

## DEBUG BUILDS

GO_FILES = $(shell find . -type f -name '*.go')

$(BINDIR)/freeblock: $(GO_FILES)
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/freeblock ./cmd/freeblock

.PHONY: build
build: $(BINDIR)/freeblock

## RELEASE BUILDS

RELEASEDIR := $(BINDIR)/release
PLATFORMS := darwin linux windows
ARCHES := amd64 arm arm64

# This rule expects targets with the format $(RELEASEDIR)/freeblock-GOOS-GOARCH(.exe)?
$(RELEASEDIR)/freeblock-%: $(GO_FILES)
	mkdir -p $(RELEASEDIR)
	GOOS=$(word 2,$(subst -, ,$@)) GOARCH=$(word 3,$(subst -, ,$(word 1,$(subst ., ,$@)))) \
	  go build -o $@ -ldflags "-s -w" ./cmd/freeblock

windowsExt = $(if $(subst windows,,$(plat)),,.exe)

# darwin-arm and windows-arm64 are not supported by go build
.PHONY: release
release: $(filter-out %-darwin-arm,$(filter-out %-windows-arm64.exe, \
	$(foreach arch,$(ARCHES),\
	$(foreach plat,$(PLATFORMS),$(RELEASEDIR)/freeblock-$(plat)-$(arch)$(windowsExt)))))
