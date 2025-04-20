.DEFAULT_GOAL := build

.PHONY: fmt vet build
fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

build: vet
	@go build

clean:
	@go clean

test:
	@go test ./... -vet=off

cover:
	@go test ./... -cover

bench:
	@go test ./... -run=^$ -bench=. -benchmem

race:
	@go test ./... -race