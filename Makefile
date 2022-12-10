BINARY_NAME:=goemu

GO:=$(shell which go)
GOFMT:=$(shell which gofmt)

export GO111MODULE:=on

.PHONY: all
all: fmt build

.PHONY: fmt
fmt:
	@$(GO) mod tidy
	@$(GOFMT) -s -w main.go

.PHONY: build
build:
	@$(GO) build -o $(BINARY_NAME)

.PHONY:
run: fmt
	@$(GO) run main.go

.PHONY: clean
clean:
	@$(GO) clean
	@rm -f $(BINARY_NAME)
