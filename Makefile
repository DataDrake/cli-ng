include Makefile.waterlog

GOPATH   = $(shell pwd)/build
GOCC     = GOPATH=$(GOPATH) go

GOBIN    = build/bin
GOSRC    = build/src
PROJROOT = $(GOSRC)/github.com/DataDrake
PKGNAME  = cli-ng
SUBPKGS  = cmd translate

DESTDIR ?=
PREFIX  ?= /usr
BINDIR   = $(PREFIX)/bin

all: build

build: setup
	@$(call stage,BUILD)
	@$(GOCC) install -v github.com/DataDrake/$(PKGNAME)
	@$(call pass,BUILD)

setup:
	@$(call stage,SETUP)
	@$(call task,Setting up GOPATH...)
	@mkdir -p $(GOPATH)
	@$(call task,Setting up src/...)
	@mkdir -p $(GOSRC)
	@$(call task,Setting up project root...)
	@mkdir -p $(PROJROOT)
	@$(call task,Setting up symlinks...)
	@if [ ! -d $(PROJROOT)/$(PKGNAME) ]; then ln -s $(shell pwd) $(PROJROOT)/$(PKGNAME); fi
	@$(call task,Getting dependencies...)
	@$(GOCC) get github.com/DataDrake/waterlog
	@$(GOCC) get github.com/leonelquinteros/gotext
	@cd $(GOPATH)/src/github.com/leonelquinteros/gotext && git checkout -q v1.2.0
	@$(call pass,SETUP)

validate: golint-setup
	@$(call stage,FORMAT)
	@$(GOCC) fmt -x ./...
	@$(call pass,FORMAT)
	@$(call stage,VET)
	@$(GOCC) vet -x ./$(PROJROOT)/...
	@$(call pass,VET)
	@$(call stage,LINT)
	@for d in $(SUBPKGS); do $(GOBIN)/golint -set_exit_status ./$$d; done
	@$(call pass,LINT)

golint-setup:
	@if [ ! -e $(GOBIN)/golint ]; then \
	    printf "Installing golint..."; \
	    $(GOCC) get -u github.com/golang/lint/golint; \
	    printf "DONE\n\n"; \
	    rm -rf $(GOPATH)/src/golang.org $(GOPATH)/src/github.com/golang $(GOPATH)/pkg; \
	fi

install:
	@$(call stage,INSTALL)
	install -D -m 00755 $(GOBIN)/$(PKGNAME) $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,INSTALL)

uninstall:
	@$(call stage,UNINSTALL)
	rm -f $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,UNINSTALL)

clean:
	@$(call stage,CLEAN)
	@$(call task,Removing symlinks...)
	@unlink $(PROJROOT)/$(PKGNAME)
	@$(call task,Removing build directory...)
	@rm -rf build
	@$(call pass,CLEAN)
