MODULE := github.com/guillaumeparis2000/mansionEyeBot

CMDS := mansionEye

GOFILES := $(shell find . -type f -name '*.go')

RELEASE_BUILDS    := $(patsubst %,build/release/%,$(CMDS))
RELEASE_BUILDS_PI := $(patsubst %,build/release-pi/%,$(CMDS))
DEBUG_BUILDS      := $(patsubst %,build/debug/%,$(CMDS))
DEBUG_BUILDS_PI   := $(patsubst %,build/debug-pi/%,$(CMDS))

GO                 := go
GCFLAGS            :=
GCFLAGS_RELEASE    := $(GCFLAGS)
GCFLAGS_RELEASE_PI := $(GCFLAGS) GOOS=linux GOARCH=arm GOARM=6 go build  -ldflags="-s -w"
GCFLAGS_DEBUG      := $(GCFLAGS) -N -l
GCFLAGS_DEBUG_PI   := $(GCFLAGS_RELEASE_PI) -N -l

default:     release
release:     $(RELEASE_BUILDS)
release-pi:  $(RELEASE_BUILDS_PI)
debug:       $(DEBUG_BUILDS)
debug-pi:    $(DEBUG_BUILDS_PI)

build/debug/%: $(GOFILES)
	@mkdir -p "$(@D)"
	@echo "     \x1b[1;34mgo build \x1b[0;1m(debug)\x1b[0m  $@"
	@$(GO) build -trimpath -i -o "$@" -gcflags "$(GCFLAGS_DEBUG)" "$(MODULE)/cmd/$(@F)"

build/debug-pi/%: $(GOFILES)
	@mkdir -p "$(@D)"
	@echo "     \x1b[1;34mgo build \x1b[0;1m(debug)\x1b[0m  $@"
	@$(GO) build -trimpath -i -o "$@" -gcflags "$(GCFLAGS_DEBUG_PI)" "$(MODULE)/cmd/$(@F)"

build/release/%: $(GOFILES)
	@mkdir -p "$(@D)"
	@echo "   \x1b[1;34mgo build \x1b[0;1m(release)\x1b[0m  $@"
	@$(GO) build -trimpath -i -o "$@" -gcflags "$(GCFLAGS_RELEASE)" "$(MODULE)/cmd/$(@F)"

build/release-pi/%: $(GOFILES)
	@mkdir -p "$(@D)"
	@echo "   \x1b[1;34mgo build \x1b[0;1m(release)\x1b[0m  $@"
	@$(GO) build -trimpath -i -o "$@" -gcflags "$(GCFLAGS_RELEASE_PI)" "$(MODULE)/cmd/$(@F)"

clean:
	@echo "                   \x1b[1;31mrm\x1b[0m  $(CMDS)"
	@rm -f $(CMDS)
	@echo "                   \x1b[1;31mrm\x1b[0m  build"
	@rm -rf build
