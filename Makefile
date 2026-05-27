BINARY     := gocode
VERSION    := 2.0.0
INSTALL_DIR := $(HOME)/.gocode
BUILD_FLAGS := -ldflags "-X main.Version=$(VERSION) -s -w"
SRC        := .

.PHONY: all build install clean tidy fmt vet test

all: tidy build install

build:
	@echo "  Building gocode v$(VERSION)..."
	@go build $(BUILD_FLAGS) -o $(BINARY) $(SRC)
	@echo "  Binary: ./$(BINARY)"

install: build
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@chmod +x $(INSTALL_DIR)/$(BINARY)
	@echo "  Installed: $(INSTALL_DIR)/$(BINARY)"
	@echo ""
	@echo "  Add to PATH if not already:"
	@echo "    echo 'export PATH=\"$$HOME/.gocode:$$PATH\"' >> ~/.bashrc"
	@echo "    source ~/.bashrc"

tidy:
	@go mod tidy

fmt:
	@go fmt ./...

vet:
	@go vet ./...

test:
	@go test -v ./...

clean:
	@rm -f $(BINARY)

# Fast rebuild (skip tidy)
quick: build install
