.PHONY: all build test clean install fmt fix vet lint sec security check help

BINARY_NAME=thefuck
GO=go
GOFLAGS=-v
LDFLAGS=-w -s

all: build check

help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  test       - Run tests with coverage"
	@echo "  coverage   - Generate detailed HTML coverage report"
	@echo "  fmt        - Format code with gofmt"
	@echo "  fix        - Run various fixes on the code base."
	@echo "  vet        - Run go vet"
	@echo "  lint       - Run staticcheck linter"
	@echo "  sec        - Run security checks with gosec"
	@echo "  security   - Alias for sec"
	@echo "  check      - Run all checks (fmt, vet, lint, sec, test)"
	@echo "  clean      - Remove build artifacts"
	@echo "  install    - Install the binary to GOPATH/bin"
	@echo "  all        - Run check and build (default)"

build:
	@echo "Generating code..."
	$(GO) generate ./...
	@echo "Building..."
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) main.go
	@echo "Build complete: bin/$(BINARY_NAME)"

test:
	@echo "Running tests..."
	CGO_ENABLED=1 $(GO) test -v -race -coverprofile=coverage.out ./...
	@echo ""
	@echo "Tests complete"
	@echo "Overall coverage:"
	@$(GO) tool cover -func=coverage.out | tail -n 1 | sed 's/[\t ][\t ]*/ /g'

coverage: test
	@echo "Generating coverage report..."
	@$(GO) tool cover -html=coverage.out -o coverage.html \
	&& $(GO) tool cover -func=coverage.out \
	  | awk '$$1 ~ /\.go:/ { \
	      split($$1, a, ":"); \
	      file = a[1]; \
	      cov[file] += $$3; \
	      cnt[file]++; \
	    } \
	    END { \
	      for (f in cov) { \
		printf "%-50s %6.2f%%\n", f, cov[f]/cnt[f]; \
	      } \
	    } \
	  ' | sort \
	&& $(GO) tool cover -func=coverage.out | tail -n 1 | sed 's/[\t ][\t ]*/ /g'
	@echo "Coverage report generated: coverage.html"

fmt:
	@echo "Checking code formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "The following files are not formatted:"; \
		gofmt -l .; \
		echo "Run 'gofmt -w .' to format them"; \
		exit 1; \
	fi
	@echo "Code formatting check passed"

fix:
	@echo "Formatting code..."
	gofmt -w .
	$(GO) get -u ./...
	$(GO) tool | grep / | xargs -r -iX $(GO) get -tool X
	$(GO) fix ./...
	@echo "Code formatted"

vet:
	@echo "Running go vet..."
	$(GO) vet ./...
	@echo "go vet passed"

lint:
	@echo "Running staticcheck..."
	$(GO) tool staticcheck ./...
	@echo "staticcheck passed"

sec: security

security:
	@echo "Running security checks with gosec..."
	$(GO) tool gosec -exclude-dir=internal/gen -quiet ./...
	@echo "Security checks passed"

check: fmt vet lint sec test
	@echo "All checks passed!"

clean:
	@echo "Cleaning..."
	rm -f bin/$(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

install: build
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install
	@echo "Installation complete"

# Run all checks and build
ci: build check
