# Binaries
GO = go

# Params
BUILD_DIR = $(shell pwd)/build

# Default target
all: goloba config server

dir:
	@echo -- DIR --
	@mkdir -p $(BUILD_DIR)

goloba: dir
	@echo -- GOLOBA --
	@cd cmd/goloba && $(GO) build -o $(BUILD_DIR)/ ./

server: dir
	@echo -- DUMMYSERVER --
	@cd cmd/dummyserver && $(GO) build -o $(BUILD_DIR)/ ./

config:
	@echo -- CONFIG --
	@cp ./goloba.conf.sample $(BUILD_DIR)/goloba.conf

clean:
	@echo -- CLEAN --
	@rm -rf $(BUILD_DIR)
