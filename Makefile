GO_CMD=GO111MODULE=on go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_FORMAT=$(GO_CMD) fmt
GO_IMPORTS=goimports
BUILD_DIR=build
BINARY_NAME=slclogger
BINARY_PATH=$(BUILD_DIR)/$(BINARY_NAME)

all: setup test build

setup:
	mkdir -p $(BUILD_DIR)

build: setup deps
	$(GO_PACKR)
	$(GO_BUILD) -o ./$(BINARY_PATH) -v

test: deps
	$(GO_GET) "github.com/stretchr/testify"
	$(GO_TEST) -v ./...

clean:
	$(GO_CLEAN)
	rm -f ./$(BINARY_PATH)

fmt:
	$(GO_GET) "golang.org/x/tools/cmd/goimports"
	find . -type f -name '*.go' | xargs $(GO_IMPORTS)  -d -e -w

run: build
	./$(BINARY_PATH)

deps:
	$(GO_GET) -v -d ./...

update:
	$(GO_GET) -v -d -u ./...

build-linux: setup deps
	$(GO_PACKR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o ./$(BINARY_PATH)

