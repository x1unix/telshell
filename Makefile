ADDR ?= :1000
APP_SH ?= /bin/bash

.PHONY: run
run:
	@go run ./cmd/telshell -addr=$(ADDR) -shell=$(APP_SH)
