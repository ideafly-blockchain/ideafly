# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

# ====variables about cross compile====
BIN_REVISION_STRING := $(shell git show -s --pretty=format:%h)
GIT_TIME := $(shell git show -s --pretty=format:%cI)
GO_PROXY := $(shell go env GOPROXY)
DEST_DIR = ./build/bin
BINARY_NAME_PREFIX = geth
ENTRY_FILE_GETH_DIR = github.com/ethereum/go-ethereum/cmd/geth
# ====end of variables about cross compile====

.PHONY: geth evm all test clean docker-local

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

geth:
	$(GORUN) build/ci.go install ./cmd/geth
	@echo "Done building."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."

all:
	$(GORUN) build/ci.go install

test: all
	$(GORUN) build/ci.go test

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go install golang.org/x/tools/cmd/stringer@latest
	env GOBIN= go install github.com/fjl/gencodec@latest
	env GOBIN= go install github.com/golang/protobuf/protoc-gen-go@latest
	env GOBIN= go install ./cmd/abigen
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

crosstools:
	GO111MODULE=off go install github.com/stars-labs/xgo2

geth-linux-amd64:
	xgo2 --goproxy="${GO_PROXY}" --targets=linux/amd64 -ldflags "-s -w -X 'main.gitCommit=${BIN_REVISION_STRING}' -X 'main.gitDate=${GIT_TIME}'" -dest=${DEST_DIR} -out ${BINARY_NAME_PREFIX} --pkg=${ENTRY_FILE_GETH_DIR} .
	@echo "The output binary file: $(DEST_DIR)/${BINARY_NAME_PREFIX}-linux-amd64"

geth-linux-arm64:
	xgo2 --goproxy="${GO_PROXY}" --targets=linux/arm64 -ldflags "-s -w -X 'main.gitCommit=${BIN_REVISION_STRING}' -X 'main.gitDate=${GIT_TIME}'" -dest=${DEST_DIR} -out ${BINARY_NAME_PREFIX} --pkg=${ENTRY_FILE_GETH_DIR} .
	@echo "The output binary file: $(DEST_DIR)/${BINARY_NAME_PREFIX}-linux-arm64"
