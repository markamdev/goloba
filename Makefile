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
	@docker build -t markamdev/goloba -f Dockerfile.goloba .
	@docker build -t markamdev/dummyserver -f Dockerfile.dummyserver .

# Targets for publishing images to Docker Hub

.publish_goloba:
	@echo -- DOCKERHUB PUBLISHING : GOLOBA --
	@docker build -f dockerhub/Dockerfile.goloba-arm7 -t markamdev/goloba:arm7 --platform linux/arm/v7 .
	@docker image push markamdev/goloba:arm7
	@docker build -f dockerhub/Dockerfile.goloba-arm8 -t markamdev/goloba:arm8 --platform linux/arm/v8 .
	@docker image push markamdev/goloba:arm8
	@docker build -f dockerhub/Dockerfile.goloba-amd64 -t markamdev/goloba:amd64 --platform linux/amd64 .
	@docker image push markamdev/goloba:amd64
	@docker manifest create markamdev/goloba:latest \
		--amend markamdev/goloba:arm7 \
		--amend markamdev/goloba:arm8 \
		--amend markamdev/goloba:amd64
	@docker manifest push markamdev/goloba:latest

.publish_dummyserver:
	@echo -- DOCKERHUB PUBLISHING : DUMMYSERVER --
