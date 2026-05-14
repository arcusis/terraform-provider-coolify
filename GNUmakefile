default: fmt lint install

build:
	go build -v ./...

install: build
	go install .

lint:
	golangci-lint run

fmt:
	gofmt -s -w .
	terraform fmt -recursive examples/

test:
	go test -v -cover ./internal/...

.PHONY: build install lint fmt test
