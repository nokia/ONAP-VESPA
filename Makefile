ifneq ($(shell uname), Linux)
    ifneq ($(GOPATH),)
export GOPATH := $(GOPATH);$(CURDIR)
    else
export GOPATH := $(CURDIR)
export PATH = $(PATH);$(CURDIR)/bin
    endif
else
    ifneq ($(GOPATH),)
export GOPATH := $(GOPATH):$(CURDIR)
    else
export GOPATH := $(CURDIR)
export PATH := $(PATH):$(CURDIR)/bin
    endif
endif

ifeq ($(REVISION),)
export REVISION := $(shell git rev-parse --short HEAD)
endif
ifeq ($(BRANCH),)
export BRANCH := $(shell git rev-parse --abbrev-ref=loose HEAD)
endif
ifeq ($(VERSION),)
export VERSION := dev.$(BRANCH).$(REVISION)
endif
ifeq ($(BUILDID),)
export BUILDID := 0
endif
export BUILDDATE := $(shell date +%Y-%m-%d-%H:%M:%S)
export LDFLAGS := "-X main.Version=$(VERSION) -X main.Branch=$(BRANCH) -X main.Revision=$(REVISION) -X main.Build=$(BUILDID) -X main.BuildDate=$(BUILDDATE)"
.PHONY: rpm build tools test lint clean sloc

all: build test rpm

tools:
	@echo "GOPATH: $(GOPATH)"
	@echo "PATH: $(PATH)"
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/axw/gocov/gocov/...
	go get -u github.com/AlekSi/gocov-xml
	go get -u github.com/tebeka/go2xunit
	go get -u github.com/hhatto/gocloc/cmd/gocloc
	go get -u github.com/gobuffalo/packr/packr
	gometalinter --install

build:
	@echo "GOPATH: $(GOPATH)"
	cd src/ves-agent && dep ensure -v
	@echo "Version: $(VERSION)"
	mkdir -p build/windows
	mkdir -p build/linux
	packr
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/linux/ves-agent ves-agent
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/windows/ves-agent.exe ves-agent
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/linux/ves-simu ves-agent/ves-simu
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/windows/ves-simu.exe ves-agent/ves-simu

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/linux/gencert utils/gencert
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/windows/gencert.exe utils/gencert
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/linux/testurl utils/testurl
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/windows/testurl.exe utils/testurl
	packr clean

test:
	mkdir -p build
	touch build/.testok
	(go test -parallel 1 -race -v -coverprofile=build/cover.out ves-agent/... || rm -f build/.testok) | tee build/testresults.out 
	go2xunit -input build/testresults.out -output build/testresults.xml
	gocov convert build/cover.out | gocov-xml > build/coverage.xml
	sed -i "s|$(CURDIR)|.|" build/coverage.xml
	gocov convert build/cover.out | gocov report
	test -f build/.testok

lint:
	mkdir -p build
	gometalinter --vendor --enable-all --checkstyle --deadline=120s --line-length=400 src/ves-agent/... > build/checkstyle.xml || true

sloc:
	mkdir -p build
	gocloc --by-file --output-type=sloccount --not-match-d="^pkg|^evel-test-collector|^libevel|vendor|^.vscode" . > build/sloccount.scc

rpm:
	mkdir -p rpmbuild/SOURCES/
	mkdir -p build
	cp -v ves-agent* rpmbuild/SOURCES/
	cp -v build/linux/ves-agent* rpmbuild/SOURCES/
	rpmbuild --define "_topdir $(CURDIR)/rpmbuild" -ba ves-agent.spec
	cp $(CURDIR)/rpmbuild/RPMS/x86_64/*.rpm	build/
	
clean:
	rm -rf build/
	rm -rf rpmbuild/
	go clean
	packr clean || true