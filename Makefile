.PHONY: fmt lint test build run

fmt:
	go fmt ./...

lint:
	@test -z "$(shell gofmt -l .)" || (echo "run 'make fmt'" && gofmt -l . && exit 1)

test:
	go test ./...

build:
	go build ./cmd/oddjob

run:
	go run ./cmd/oddjob
