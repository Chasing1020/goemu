BINARY_NAME:=goemu

GO:=$(shell which go)
GOFMT:=$(shell which gofmt)

export GO111MODULE:=on

all: fmt build

fmt:
	@$(GO) mod tidy
	@$(GOFMT) -s -w main.go

build:
	@$(GO) build -o $(BINARY_NAME)

run: fmt
	@$(GO) run main.go

clean:
	@$(GO) clean
	@rm -f $(BINARY_NAME)

.PHONY: all fmt build run clean
