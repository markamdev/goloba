# Binaries
GO = go

# Params
BUILD_DIR = $(shell pwd)/build

# Go compiler should always check if it's necessary to re-buld binary
.PHONY: goloba dummyserver docker test

# Default target
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
