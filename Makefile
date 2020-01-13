ADDR ?= :1000
APP_SH ?= /bin/bash
PKG ?= ./cmd/telshell
BUILD_DIR ?= ./build

.PHONY: run
run:
	@go run $(PKG) -addr=$(ADDR) -shell=$(APP_SH)

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: build
build:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/telshell_darwin $(PKG)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/telshell_win-amd64.exe $(PKG)
	GOOS=windows GOARCH=386 go build -o $(BUILD_DIR)/telshell_win-i386.exe $(PKG)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/telshell_linux-amd64 $(PKG)
	GOOS=linux GOARCH=386 go build -o $(BUILD_DIR)/telshell_linux-i386 $(PKG)
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/telshell_linux-arm64 $(PKG)
