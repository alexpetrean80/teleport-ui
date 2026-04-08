setup:
	git config core.hooksPath .githooks

build:
	go build main.go -o /dist/teleport-ui

run:
	go run main.go

test:
	go test ./...

fmt:
	golangci-lint run --fix ./...

lint:
	golangci-lint run ./...
