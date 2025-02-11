# Binary name
BINARY_NAME=spider

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build clean test run

all: build

## Build the binary
build:
	@echo "Building..."
	go build -o $(GOBIN)/$(BINARY_NAME) ./cmd/spider/main.go

## Clean the binary
clean:
	@echo "Cleaning..."
	rm -rf $(GOBIN)
	go clean

