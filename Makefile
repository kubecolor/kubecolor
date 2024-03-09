help: ## Print usage
	@sed -r '/^(\w+-*\w+):[^#]*##/!d;s/^([^:]+):[^#]*##\s*(.*)/\x1b[36m\1\t:\x1b[m \2/g' ${MAKEFILE_LIST}
.PHONY: help

build-: ## build package
	go build -v .
.PHONY: build

test: ## run test and generate coverage report
	go test -race -coverprofile=coverage.txt -covermode=atomic -v ./... && \
	go tool cover -html=coverage.txt -o cover.html
.PHONY: test

testshort: ## run test and generate short report
	go test -timeout 30s -count=1 ./... -test.short
.PHONY: testshort

fmt: ## format code
	go fmt ./...
.PHONY: fmt

lint: ## lint code
	staticcheck ./...

.PHONY: lint
