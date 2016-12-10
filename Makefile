.PHONE: test test-lint

all: test

test: test-lint
	@go test -v ./...

test-lint:
	@! gofmt -d . 2>&1 | read
