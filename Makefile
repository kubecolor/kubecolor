help: ## Print usage
	@sed -r '/^(\w+):[^#]*##/!d;s/^([^:]+):[^#]*##\s*(.*)/\x1b[36m\1\t:\x1b[m \2/g' ${MAKEFILE_LIST}

build: ## build package
	go build -v .

test: ## run test and generate coverage report
	go test -race -coverprofile=coverage.txt -covermode=atomic -v ./... && \
	go tool cover -html=coverage.txt -o cover.html

testshort: ## run test and generate short report
	go test -timeout 30s -count=1 ./... -test.short

fmt: ## format code
	go fmt .

lint: ## lint code
	staticcheck ./...

.PHONY: test testshort build help fmt lint
