.PHONY: test test-lint

all: test

test: test-lint
	@go test ./...

test-lint:
	@! gofmt -d . 2>&1 | read
