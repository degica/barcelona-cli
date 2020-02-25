.DEFAULT_GOAL := dev
.PHONY: clean release dev-install dev build generate check test format vet

OUTDIR=out
BINNAME=bcn
VERSION=$(shell cat ./VERSION)
PACKAGES=$(shell go list ./... | grep -v /vendor/)

$(OUTDIR)/linux_amd64/bcn:
	GOOS=linux GOARCH=amd64 go build -o $(OUTDIR)/linux_amd64/bcn
$(OUTDIR)/darwin_amd64/bcn:
	GOOS=darwin GOARCH=amd64 go build -o $(OUTDIR)/darwin_amd64/bcn
$(OUTDIR)/windows_amd64/bcn.exe:
	GOOS=windows GOARCH=amd64 go build -o $(OUTDIR)/windows_amd64/bcn.exe

$(OUTDIR)/$(VERSION)/bcn_linux_amd64.tar.gz: $(OUTDIR)/linux_amd64/bcn
	mkdir -p $(OUTDIR)/$(VERSION)
	zip -j $(OUTDIR)/$(VERSION)/bcn_linux_amd64.zip $(OUTDIR)/linux_amd64/bcn
$(OUTDIR)/$(VERSION)/bcn_darwin_amd64.tar.gz: $(OUTDIR)/darwin_amd64/bcn
	mkdir -p $(OUTDIR)/$(VERSION)
	zip -j $(OUTDIR)/$(VERSION)/bcn_darwin_amd64.zip $(OUTDIR)/darwin_amd64/bcn
$(OUTDIR)/$(VERSION)/bcn_windows_amd64.zip: $(OUTDIR)/windows_amd64/bcn.exe
	mkdir -p $(OUTDIR)/$(VERSION)
	zip -j $(OUTDIR)/$(VERSION)/bcn_windows_amd64.zip $(OUTDIR)/windows_amd64/bcn.exe

generate:
	go generate

build: generate $(OUTDIR)/linux_amd64/bcn $(OUTDIR)/darwin_amd64/bcn $(OUTDIR)/windows_amd64/bcn.exe

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

release: build $(OUTDIR)/$(VERSION)/bcn_linux_amd64.tar.gz $(OUTDIR)/$(VERSION)/bcn_darwin_amd64.tar.gz $(OUTDIR)/$(VERSION)/bcn_windows_amd64.zip

clean:
	rm -rf $(OUTDIR)
	rm -f ./bcn
