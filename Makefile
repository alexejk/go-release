include build.properties

.PHONY: all build-all build-linux build-osx build-windows pre-build clean lint test package-all package-linux package-osx package-windows install
.NOTPARALLEL: clean

APP_NAME=go-release
APP_BUILD=`git log --pretty=format:'%h' -n 1`
APP_BUILD_DATE=`date +%Y-%m-%d`
APP_BUILD_TIME=`date +%H:%M`
GO_FLAGS= CGO_ENABLED=0
GO_LDFLAGS= -ldflags="-X main.AppVersion=$(APP_VERSION) -X main.AppName=$(APP_NAME) -X main.AppBuild=$(APP_BUILD) -X main.AppBuildDate=$(APP_BUILD_DATE) -X main.AppBuildTime=$(APP_BUILD_TIME)"
GO_BUILD_CMD=$(GO_FLAGS) go build $(GO_LDFLAGS)
BUILD_DIR=build
MOCK_DIR=mocks
BINARY_NAME=go-release

all: clean generate-all test lint build-all package-all

lint:
	@echo "Linting code..."
	@go vet `go list ./... | grep -v $(MOCK_DIR)`

test:
	@echo "Running tests..."
	@go test `go list ./... | grep -v $(MOCK_DIR)`

mock-gen:
	@echo "Generating mocks..."
#	@mockery -dir api/amazon -all


code-gen:
	@echo "Generating code..."
	@go generate -x ./...

generate-all: code-gen mock-gen

pre-build:
	@mkdir -p $(BUILD_DIR)

build-linux: pre-build
	@echo "Building Linux binary..."
	GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64

build-osx: pre-build
	@echo "Building OSX binary..."
	GOOS=darwin GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64

build-windows: pre-build
	@echo "Building Windows binary..."
	GOOS=windows GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64

build-all: build-linux build-osx build-windows

# For backwards compatibility
build: build-all

package-linux:
	@echo "Packaging Linux binary..."
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BINARY_NAME)-$(APP_VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64

package-osx:
	@echo "Packaging OSX binary..."
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BINARY_NAME)-$(APP_VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64

package-windows:
	@echo "Packaging Windows binary..."
	zip -q -j  $(BUILD_DIR)/$(BINARY_NAME)-$(APP_VERSION)-windows-amd64.zip $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64

package-all: package-linux package-osx package-windows

# For backwards compatibility
package: package-all

clean:
	@echo "Cleaning..."
	@rm -Rf $(BUILD_DIR)
	@rm -Rf $(MOCK_DIR)

install: build-osx
ifeq ($(shell uname -s), Darwin)
	cp $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 /usr/local/bin/$(BINARY_NAME)
else
	$(error Cannot install binary as this is currently only supported for Darwin)
endif

# Docker All task calls the "all" target
docker-all: all
	@echo "Copying compiled artifacts to the output directory"
	cp $(BUILD_DIR)/* /opt/

build-docker-builder-image:
	@echo "Creating builder image"
	docker build -t $(BINARY_NAME)-builder .

# Docker Build triggers build in the builder image
build-in-docker: build-docker-builder-image pre-build
	@echo "Triggering build in docker"
	docker run --rm -v `pwd`/build:/opt $(BINARY_NAME)-builder make docker-all
