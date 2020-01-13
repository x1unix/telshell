ADDR ?= :1000
APP_SH ?= /bin/bash
PKG ?= ./cmd/telshell
BUILD_DIR ?= ./build
PROJECT_NAME := telshell

define go_build
$(info - Building for $(1) ($(2)))
$(shell GOOS=$(1) GOARCH=$(2) go build -o $(BUILD_DIR)/telshell_$(3) $(PKG))
endef

.PHONY:build
build: clean build-linux build-darwin build-win32

.PHONY: run
run:
	@go run $(PKG) -addr=$(ADDR) -shell=$(APP_SH)

.PHONY: clean
clean:
	@echo "- Clean build directory" && rm -rf $(BUILD_DIR)

.PHONY:windows
windows:
	$(call go_build,windows,amd64,windows-amd64.exe)
	$(call go_build,windows,386,windows-i386.exe)

.PHONY:linux
linux:
	$(call go_build,linux,amd64,linux-amd64)
	$(call go_build,linux,386,linux-i386)
	$(call go_build,linux,arm64,linux-aarch64)

.PHONY: darwin
darwin:
	$(call go_build,darwin,amd64,darwin-amd64)