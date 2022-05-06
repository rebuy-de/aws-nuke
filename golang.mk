# Source: https://github.com/rebuy-de/golang-template

TARGETS?="."
PACKAGE=$(shell GOPATH= go list $(TARGET))
NAME=$(notdir $(shell echo $(PACKAGE) | sed 's/\/v2//'))

BUILD_VERSION=$(shell git describe --always --dirty --tags | tr '-' '.' )
BUILD_DATE=$(shell LC_ALL=C date)
BUILD_HASH=$(shell git rev-parse HEAD)
BUILD_MACHINE=$(shell uname -n)
BUILD_USER=$(shell whoami)
BUILD_ENVIRONMENT=$(BUILD_USER)@$(BUILD_MACHINE)

BUILD_XDST=$(PACKAGE)/cmd
BUILD_FLAGS=-ldflags "\
	$(ADDITIONAL_LDFLAGS) -s -w \
	-X '$(BUILD_XDST).BuildVersion=$(BUILD_VERSION)' \
	-X '$(BUILD_XDST).BuildDate=$(BUILD_DATE)' \
	-X '$(BUILD_XDST).BuildHash=$(BUILD_HASH)' \
	-X '$(BUILD_XDST).BuildEnvironment=$(BUILD_ENVIRONMENT)' \
"

GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./.git/*")
GOPKGS=$(shell go list ./...)

OUTPUT_FILE=$(NAME)-$(BUILD_VERSION)-$(shell go env GOOS)-$(shell go env GOARCH)$(shell go env GOARM)$(shell go env GOEXE)
OUTPUT_LINK=$(NAME)$(shell go env GOEXE)
WINDOWS_ZIP=$(shell echo $(OUTPUT_FILE) | sed 's/\.exe/\.zip/')

default: build

vendor: go.mod go.sum
	go mod vendor
	touch vendor

format:
	gofmt -s -w $(GOFILES)

vet: go_generate vendor
	go vet $(GOPKGS)

lint:
	$(foreach pkg,$(GOPKGS),golint $(pkg);)

go_generate:
	rm -rvf mocks
	go generate ./...

test_packages: go_generate vendor
	go test $(GOPKGS)

test_format:
	gofmt -s -l $(GOFILES)

test: test_format vet lint test_packages

cov: go_generate
	gocov test -v $(GOPKGS) \
		| gocov-html > coverage.html

_build: vendor
	mkdir -p dist
	$(foreach TARGET,$(TARGETS),go build \
		$(BUILD_FLAGS) \
		-o dist/$(OUTPUT_FILE) \
		$(TARGET);\
	)

build: go_generate _build
	$(foreach TARGET,$(TARGETS),ln -sf $(OUTPUT_FILE) dist/$(OUTPUT_LINK);)

compress: _build
ifeq ($(GOOS),windows)
	zip -j dist/$(WINDOWS_ZIP) dist/$(OUTPUT_FILE)
else
	tar czf dist/$(OUTPUT_FILE).tar.gz -C dist $(OUTPUT_FILE)
endif
	rm -f dist/$(OUTPUT_FILE)

xc: go_generate
	GOOS=linux GOARCH=amd64 make compress
	GOOS=linux GOARCH=arm64 make compress
	GOOS=linux GOARCH=arm GOARM=7 make compress
	GOOS=darwin GOARCH=amd64 make compress
	GOOS=darwin GOARCH=arm64 make compress
	GOOS=windows GOARCH=amd64 make compress

install: test
	$(foreach TARGET,$(TARGETS),go install \
		$(BUILD_FLAGS);)

clean:
	rm dist/ -rvf
	rm mocks/ -rvf
