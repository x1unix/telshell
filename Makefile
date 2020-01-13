ADDR ?= :1000
APP_SH ?= /bin/bash
PKG ?= ./cmd/telshell

.PHONY: run
run:
	@go run $(PKG) -addr=$(ADDR) -shell=$(APP_SH)

.PHONY: build
build:
	GOOS=darwin go build -o ./build/telshell_darwin $(PKG)
	GOOS=windows go build -o ./build/telshell_win-amd64.exe $(PKG)
	GOOS=windows GOARCH=386 go build -o ./build/telshell_win-i386.exe $(PKG)
	GOOS=linux GOARCH=amd64 go build -o ./build/telshell_linux-amd64 $(PKG)
	GOOS=linux GOARCH=386 go build -o ./build/telshell_linux-386 $(PKG)
	GOOS=linux GOARCH=arm64 go build -o ./build/telshell_linux-arm64 $(PKG)
