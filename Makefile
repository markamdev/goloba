# Binaries
GO = go

# Params
BUILD_DIR = $(shell pwd)/build

# Go compiler should always check if it's necessary to re-buld binary
.PHONY: goloba dummyserver

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
