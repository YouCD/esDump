GOCMD			:=$(shell which go)
GOBUILD			:=$(GOCMD) build

IMPORT_PATH		:=github.com/YouCD/esDump/cmd
BUILD_TIME		:=$(shell date "+%F %T")
COMMIT_ID       :=$(shell git rev-parse HEAD)
GO_VERSION      :=$(shell $(GOCMD) version)
VERSION			:=$(shell git describe --tags)
BUILD_USER		:=$(shell whoami)
PACKAGES_URL	:=$(github.com/YouCD/esDump)
FLAG			:="-X '${IMPORT_PATH}.buildTime=${BUILD_TIME}' -X '${IMPORT_PATH}.commitID=${COMMIT_ID}' -X '${IMPORT_PATH}.goVersion=${GO_VERSION}' -X '${IMPORT_PATH}.goVersion=${GO_VERSION}' -X '${IMPORT_PATH}.Version=${VERSION}' -X '${IMPORT_PATH}.buildUser=${BUILD_USER}'"

BINARY_DIR=bin
BINARY_NAME:=esDump

#mac
build-darwin:
	@CGO_ENABLED=0 GOOS=darwin $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-amd64/$(BINARY_NAME)
# windows
build-win:
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-windows-amd64/$(BINARY_NAME).exe
# linux
build-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64/$(BINARY_NAME)

build:
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)
# 全平台
build-all:
	make build-darwin
	make build-win
	make build-linux
	@cd $(BINARY_DIR)&& tar Jcf $(BINARY_NAME)-darwin-amd64.txz $(BINARY_NAME)-darwin-amd64&&rm -rf $(BINARY_NAME)-darwin-amd64
	@cd $(BINARY_DIR)&& tar Jcf $(BINARY_NAME)-windows-amd64.txz $(BINARY_NAME)-windows-amd64&&rm -rf $(BINARY_NAME)-windows-amd64
	@cd $(BINARY_DIR)&& tar Jcf $(BINARY_NAME)-linux-amd64.txz $(BINARY_NAME)-linux-amd64&&rm -rf $(BINARY_NAME)-linux-amd64

