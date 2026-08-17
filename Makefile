# ptp — camera control over USB PTP.
#
#   make              build every package
#   make help         list the targets
#   make tools        build the command-line tools into ./bin
#   make test         go test ./...
#

SHELL := /bin/sh
BIN   := bin

# Command binaries: the two bring-up tools, which need only a camera.
CMDS := fujiprobe sonyprobe

.PHONY: build help tools test race vet fmt cross clean tidy $(CMDS)

build: ## build every package (default)
	@go build ./...

# help: self-documenting — lists every target annotated with a "## " description.
help: ## list the targets
	@echo "ptp — camera control over USB PTP."
	@echo
	@echo "Targets:"
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*## "}{printf "  make %-11s %s\n", $$1, $$2}'
	@echo
	@echo "Tools: $(CMDS)"

tools: $(CMDS) ## build the command-line tools into ./bin

$(CMDS): | $(BIN)
	@echo "building $@"
	@go build -o $(BIN)/$@ ./cmd/$@

$(BIN):
	@mkdir -p $(BIN)

test: ## go test ./...
	@go test ./...

race: ## go test -race ./...
	@go test -race ./...

vet: ## go vet ./...
	@go vet ./...

fmt: ## report anything gofmt would change
	@out=$$(gofmt -l .); \
	if [ -n "$$out" ]; then echo "gofmt would change:"; echo "$$out"; exit 1; fi; \
	echo "gofmt clean"

cross: ## check it still builds for linux and windows
	@GOOS=linux go build ./... && echo "linux ok"
	@GOOS=windows go build ./... && echo "windows ok"

tidy: ## go mod tidy
	@go mod tidy

clean: ## remove ./bin
	@rm -rf $(BIN)
