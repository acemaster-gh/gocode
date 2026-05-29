BINARY     := gocode
VERSION    := 2.0.0
BUILD_FLAGS := -ldflags "-X main.Version=$(VERSION) -s -w"

.PHONY: all build install clean tidy fmt vet

all: tidy build

build:
	@echo "  Building gocode v$(VERSION)..."
	@go build $(BUILD_FLAGS) -o $(BINARY) .
	@echo "  Binary: $(PWD)/$(BINARY)"

install: build
	@echo "  Installed: $(PWD)/$(BINARY)"
	@echo "  Alias: alias gocode=\"$$HOME/.gocode/gocode\""

tidy:
	@go mod tidy

fmt:
	@go fmt ./...

vet:
	@go vet ./...

clean:
	@rm -f $(BINARY)
