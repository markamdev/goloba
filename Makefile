# Binaries
GO = go

# Params
BUILD_DIR = $(shell pwd)/build

# Default target
all: goloba server

dir:
	@mkdir -p $(BUILD_DIR)

goloba: dir
	@echo BUILD GoLoBa
	@cd cmd/goloba && $(GO) build -o $(BUILD_DIR)/ ./

server: dir
	@echo BUILD dummyserver
	@cd cmd/dummyserver && $(GO) build -o $(BUILD_DIR)/ ./

config:
	@echo COPY config
	@cp ./goloba.conf.sample $(BUILD_DIR)/goloba.conf

clean:
	@rm -rf $(BUILD_DIR)