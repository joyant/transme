APP_NAME := transme
GO_FILES := ./cmd/main.go

.PHONY: all linux darwin windows help clean

# default help info
default: help

# print help info
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  linux     Compile for Linux"
	@echo "  darwin    Compile for macOS"
	@echo "  windows   Compile for Windows"
	@echo "  clean     Remove compiled files"
	@echo "  help      Show this help message"

# compile for linux
linux: $(GO_FILES)
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) $(GO_FILES)

# compile for mac
darwin: $(GO_FILES)
	GOOS=darwin GOARCH=amd64 go build -o bin/$(APP_NAME) $(GO_FILES)

# compile for windows
windows: $(GO_FILES)
	GOOS=windows GOARCH=amd64 go build -o bin/$(APP_NAME).exe $(GO_FILES)

# clean compile files
clean:
	rm -f bin/$(APP_NAME) bin/$(APP_NAME).exe

# catch unknown target
.DEFAULT:
	@echo "Error: Unknown target '$@'"
	@$(MAKE) help

# default target
all: default