# Binaries
GO = go

# Params
BUILD_DIR = $(shell pwd)/build

# Go compiler should always check if it's necessary to re-buld binary
.PHONY: goloba dummyserver docker test

# Default target
info:
	@echo -- INFO --
	@echo "Available targets:"
	@echo " all          : Build all binaries and copy config"
	@echo " goloba       : Build goloba binary"
	@echo " dummyserver  : Build dummyserver binary"
	@echo " config       : Copy sample config to build directory"
	@echo " clean        : Clean build directory"
	@echo " test         : Run testbench (requires goloba and dummyserver built)"
	@echo " docker       : Build docker image"
	@echo " publish_goloba : Publish goloba multiarch image to Docker Hub"
	@echo " tools-install : Install development tools (e.g., linters)"
	@echo " lint         : Run code linters"

all: goloba config dummyserver

$(BUILD_DIR):
	@echo -- BUILD DIR --
	@mkdir -p $(BUILD_DIR)

goloba: $(BUILD_DIR)
	@echo -- GOLOBA --
	@cd cmd/goloba && $(GO) build -o $(BUILD_DIR)/ ./

dummyserver: $(BUILD_DIR)
	@echo -- DUMMYSERVER --
	@cd cmd/dummyserver && $(GO) build -o $(BUILD_DIR)/ ./

config:
	@echo -- CONFIG --
	@cp ./goloba.conf.sample $(BUILD_DIR)/goloba.conf

clean:
	@echo -- CLEAN --
	@rm -rf $(BUILD_DIR)

test: goloba dummyserver
	@echo -- TESTING --
	@./scripts/start_testbench.sh

docker:
	@echo -- DOCKER --
	@docker build -t markamdev/goloba -f Dockerfile .

# Targets for publishing images to Docker Hub
# docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
# docker buildx rm builder
# docker buildx create --name builder --driver docker-container --use
# docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t markamdev/goloba:latest -t markamdev/goloba:0.5 --push -f Dockerfile .

.publish_goloba:
	$(eval G_VER := $(shell cat dockerhub/goloba.VERSION | head -n 1))
	@echo -- DOCKERHUB PUBLISHING : GOLOBA v $(G_VER) --
	@echo "INFO: multiarch build requires multiarch/qemu-user-static"
	@echo "Run it using: docker run --rm --privileged multiarch/qemu-user-static --reset -p yes"
	@docker buildx create --name golobabuilder --driver docker-container --use
	@docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 \
		-t markamdev/goloba:latest -t markamdev/goloba:$(G_VER) --push -f Dockerfile .
	@docker buildx rm golobabuilder

tools-install:
	@echo -- TOOLS INSTALL --
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	@echo -- LINTING --
	@`go env GOPATH`/bin/golangci-lint run ./...