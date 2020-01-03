ADDR ?= :1000
APP_SH ?= /bin/bash
PKG ?= ./cmd/telshell

.PHONY: run
run:
	@go run $(PKG) -addr=$(ADDR) -shell=$(APP_SH)

.PHONY: build
build:
	@go build -o ./build/telshell $(PKG)
