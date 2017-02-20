.DEFAULT_GOAL := dev
.PHONY: clean release dev-install dev build generate

OUTDIR=out
BINNAME=bcn
VERSION=$(shell cat ./VERSION)

$(OUTDIR)/linux_amd64/bcn:
	GOOS=linux GOARCH=amd64 go build -o $(OUTDIR)/linux_amd64/bcn
$(OUTDIR)/darwin_amd64/bcn:
	GOOS=darwin GOARCH=amd64 go build -o $(OUTDIR)/darwin_amd64/bcn
$(OUTDIR)/windows_amd64/bcn.exe:
	GOOS=windows GOARCH=amd64 go build -o $(OUTDIR)/windows_amd64/bcn.exe

$(OUTDIR)/bcn-$(VERSION)-linux_amd64.tar.gz: $(OUTDIR)/linux_amd64/bcn
	zip -j $(OUTDIR)/bcn-$(VERSION)-linux_amd64.zip $(OUTDIR)/linux_amd64/bcn
$(OUTDIR)/bcn-$(VERSION)-darwin_amd64.tar.gz: $(OUTDIR)/darwin_amd64/bcn
	zip -j $(OUTDIR)/bcn-$(VERSION)-darwin_amd64.zip $(OUTDIR)/darwin_amd64/bcn
$(OUTDIR)/bcn-$(VERSION)-windows_amd64.zip: $(OUTDIR)/windows_amd64/bcn.exe
	zip -j $(OUTDIR)/bcn-$(VERSION)-windows_amd64.zip $(OUTDIR)/windows_amd64/bcn.exe

generate:
	go generate

build: generate $(OUTDIR)/linux_amd64/bcn $(OUTDIR)/darwin_amd64/bcn $(OUTDIR)/windows_amd64/bcn.exe

dev: generate
	go build -o bcn
install: dev
	cp ./bcn ~/bin/bcn

release: build $(OUTDIR)/bcn-$(VERSION)-linux_amd64.tar.gz $(OUTDIR)/bcn-$(VERSION)-darwin_amd64.tar.gz $(OUTDIR)/bcn-$(VERSION)-windows_amd64.zip

clean:
	rm -rf $(OUTDIR)
	rm -f ./bcn
