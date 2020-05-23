MODULE := github.com/guillaumeparis2000/mansionEyeBot

CMDS := mansionEye

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= ${GIT_TAG}

GOFILES := $(shell find . -type f -name '*.go')

RELEASE_BUILDS    := $(patsubst %,build/release/%,$(CMDS))
RELEASE_BUILDS_PI := $(patsubst %,build/release-pi/%,$(CMDS))
DEBUG_BUILDS      := $(patsubst %,build/debug/%,$(CMDS))

GO                  := go
GCFLAGS             :=
GCFLAGS_RELEASE     := $(GCFLAGS)
GCFLAGS_RELEASE_PI  := $(GCFLAGS) GOOS=linux GOARCH=arm GOARM=6 go build  -ldflags="-s -w"
GCFLAGS_RELEASE_PI1 := $(GCFLAGS) GOOS=linux GOARCH=arm GOARM=6 go build  -ldflags="-s -w"
GCFLAGS_DEBUG       := $(GCFLAGS) -N -l

default:      release
release:      $(RELEASE_BUILDS)
release-pi:   $(RELEASE_BUILDS_PI)
release-pi1:  $(RELEASE_BUILDS_PI1)
debug:        $(DEBUG_BUILDS)

build/debug/%: $(GOFILES)
	@mkdir -p "$(@D)"
	@echo "     \x1b[1;34mgo build \x1b[0;1m(debug)\x1b[0m  $@"
	@$(GO) build -trimpath -i -o "$@" -gcflags "$(GCFLAGS_DEBUG)" "$(MODULE)/cmd/$(@F)"GS_DEBUG_PI)" "$(MODULE)/cmd/$(@F)"

build/release/%: $(GOFILES)
	@mkdir -p "$(@D)"
	@echo "   \x1b[1;34mgo build \x1b[0;1m(release)\x1b[0m  $@"
	@$(GO) build -trimpath -i -o "$@" -gcflags "$(GCFLAGS_RELEASE)" "$(MODULE)/cmd/$(@F)"

build/release-pi/%: $(GOFILES)
	@mkdir -p "$(@D)"
	@echo "   \x1b[1;34mgo build \x1b[0;1m(release)\x1b[0m  $@"
	@$(GO) build -trimpath -i -o "$@" -gcflags "$(GCFLAGS_RELEASE_PI)" "$(MODULE)/cmd/$(@F)"

build/release-pi1/%: $(GOFILES)
	@mkdir -p "$(@D)"
	@echo "   \x1b[1;34mgo build \x1b[0;1m(release)\x1b[0m  $@"
	@$(GO) build -trimpath -i -o "$@" -gcflags "$(GCFLAGS_RELEASE_PI1)" "$(MODULE)/cmd/$(@F)"

clean:
	@echo "                   \x1b[1;31mrm\x1b[0m  $(CMDS)"
	@rm -f $(CMDS)
	@echo "                   \x1b[1;31mrm\x1b[0m  build"
	@rm -rf build

info:
	@echo "Version:           ${VERSION}"
	@echo "Git Tag:           ${GIT_TAG}"
	@echo "Git Commit:        ${GIT_COMMIT}"
	@echo "Git Tree State:    ${GIT_DIRTY}"