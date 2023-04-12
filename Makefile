.DEFAULT_GOAL := dev
.PHONY: clean release dev-install dev build generate check test format vet

PLATFORMS=linux_amd64 darwin_amd64 linux_arm64 darwin_arm64 windows_amd64

BINARIES=$(foreach dir,$(PLATFORMS),$(OUTDIR)/$(dir)/bcn)
ZIPS=$(foreach dir,$(PLATFORMS),$(OUTDIR)/$(VERSION)/bcn_$(dir).zip)

OUTDIR=out
BINNAME=bcn
VERSION=$(shell cat ./VERSION)
PACKAGES=$(shell go list ./... | grep -v /vendor/)

$(BINARIES):
	$(eval TMP := $(subst _, ,$(word 2,$(subst /, ,$@))))
	GOOS=$(firstword $(TMP)) GOARCH=$(lastword $(TMP)) go build -o $@$(if $(findstring windows,$@),.exe)
	touch $@

ZIP := $(wildcard $(OUTDIR)/$(VERSION)/bcn_*.zip)

$(OUTDIR)/$(VERSION)/bcn_%.zip: $(OUTDIR)/%/bcn
	mkdir -p $(OUTDIR)/$(VERSION)
	zip -j $@ $?$(if $(findstring windows,$@),.exe)

generate:
	go generate

build: generate $(BINARIES)

dev: generate
	go build -o bcn

install: dev
	cp ./bcn ~/bin/bcn

check: test format vet

test: generate
	go test -race $(PACKAGES)

format: generate
	go fmt $(PACKAGES)

vet: generate
	go vet $(PACKAGES)

release: build $(ZIPS)

clean:
	go clean -modcache
	go clean -testcache
	rm -rf $(OUTDIR)
	rm -f ./bcn
