GO_FILES=$(shell find -name '*.go' -print)

help: ## Print usage
	@sed -r '/^(\w[^:]+):[^#]*##/!d;s/^([^:]+):[^#]*##\s*(.*)/\x1b[36m\1\t:\x1b[m \2/g' ${MAKEFILE_LIST} | column -t -s $$'\t'
.PHONY: help

build: ## build package
	go build -v .
.PHONY: build

test: ## run test and generate coverage report
	go test -race -coverprofile=coverage.txt -covermode=atomic -v ./... && \
	go tool cover -html=coverage.txt -o cover.html
.PHONY: test

corpus: ## run corpus tests on files in ./test/corpus
	go run ./internal/cmd/testcorpus
.PHONY: corpus

corpus-update: ## update test files in ./test/corpus with the current kubecolor output
	go run ./internal/cmd/testcorpus -update
.PHONY: corpus-update

testshort: ## run test and generate short report
	go test -timeout 30s -count=1 ./... -test.short
.PHONY: testshort

fmt: ## format code
	go fmt ./...
.PHONY: fmt

lint: ## lint code
	staticcheck ./...
.PHONY: lint

config-schema.json: $(wildcard **/*.go) ## regenerate config-schema.json based on config package
	go run ./internal/cmd/configschema -out config-schema.json

docs: $(patsubst %.txt,%.svg,$(wildcard docs/*.txt)) ## generate docs images
.PHONY: docs

# View available themes in charmbracelet/freeze: https://xyproto.github.io/splash/docs/index.html
docs/%.svg: ./docs/%.txt Makefile ./docs/freeze-config.json ${GO_FILES}
	go run ./internal/cmd/imagegen $<
docs/%-light.svg: ./docs/%-light.txt Makefile ./docs/freeze-config-light.json ${GO_FILES}
	go run ./internal/cmd/imagegen -freeze-config=./docs/freeze-config-light.json -flag-color=blue $<
