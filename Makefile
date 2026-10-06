# Cross-platform builds (pure Go, CGO_ENABLED=0).
# Usage: make help | make all | make windows | make linux

DIST     := dist
LDFLAGS  := -s -w
GO_BUILD := CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)"

.PHONY: help all local windows linux linux-arm64 clean test

help:
	@echo "lexigen build targets:"
	@echo "  make all          - Windows + Linux (amd64) into $(DIST)/"
	@echo "  make local        - binary for this machine ($(DIST)/lexigen)"
	@echo "  make windows      - $(DIST)/lexigen-windows-amd64.exe"
	@echo "  make linux        - $(DIST)/lexigen-linux-amd64"
	@echo "  make linux-arm64  - $(DIST)/lexigen-linux-arm64"
	@echo "  make test         - go test ./..."
	@echo "  make clean        - remove $(DIST)/"
	@echo ""
	@echo "Without make:  ./scripts/build.sh all   or   .\\scripts\\build.ps1 all"

all: windows linux

$(DIST):
	mkdir -p $(DIST)

local: $(DIST)
	go build -ldflags="$(LDFLAGS)" -o $(DIST)/lexigen$(if $(filter Windows_NT,$(OS)),.exe,) .

windows: $(DIST)
	$(GO_BUILD) -o $(DIST)/lexigen-windows-amd64.exe .

linux: $(DIST)
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(DIST)/lexigen-linux-amd64 .

linux-arm64: $(DIST)
	GOOS=linux GOARCH=arm64 $(GO_BUILD) -o $(DIST)/lexigen-linux-arm64 .

test:
	go test ./...

clean:
	rm -rf $(DIST)
