GOCMD:=$(shell which go)
GOBUILD:=$(GOCMD) build

PACKAGES_URL:=$(github.com/YouCD/esDump)

BINARY_DIR=bin
BINARY_NAME:=esDump

#mac
build:
	CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)-mac
# windows
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)-win.exe
# linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)-linux
# 全平台
build-all:
	make build
	make build-win
	make build-linux
