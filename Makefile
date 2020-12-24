GOCMD			:=$(shell which go)
GOBUILD			:=$(GOCMD) build

IMPORT_PATH		:=github.com/YouCD/esDump/cmd
BUILD_TIME		:=$(shell date "+%F %T")
COMMIT_ID       :=$(shell git rev-parse HEAD)
GO_VERSION      :=$(shell $(GOCMD) version)
VERSION			:=$(shell git describe --tags)
BUILD_USER		:=$(shell whoami)
PACKAGES_URL	:=$(github.com/YouCD/esDump)
FLAG			:="-X '${IMPORT_PATH}.buildTime=${BUILD_TIME}' -X '${IMPORT_PATH}.commitID=${COMMIT_ID}' -X '${IMPORT_PATH}.goVersion=${GO_VERSION}' -X '${IMPORT_PATH}.goVersion=${GO_VERSION}' -X '${IMPORT_PATH}.version=${VERSION}' -X '${IMPORT_PATH}.buildUser=${BUILD_USER}'"

BINARY_DIR=bin
BINARY_NAME:=esDump

#mac
build:
	CGO_ENABLED=0 $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-mac
# windows
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-win.exe
# linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-linux
# 全平台
build-all:
	make build
	make build-win
	make build-linux
