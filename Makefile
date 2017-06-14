export CONFIG_PATH=./config.json

TEST_CONFIG_PATH=./config-test.json
all: get-deps build

.PHONY: help build fmt clean run test coverage check vet lint doc cfpush

help:
	@echo "build    - build application from sources"
	@echo "fmt      - format application sources"
	@echo "gen      - generate files"
	@echo "prepare  - prepare project to build"
	@echo "run      - start application"

build: prepare
	go build -o ${GOPATH}/bin/service-entrypoint

fmt:
	go fmt

gen:
	go-bindata -pkg mysql -o storage/migrations/mysql/bindata.go storage/migrations/mysql

get-deps:
	go get -u github.com/jteeuwen/go-bindata/...

prepare: fmt gen

run: build
	${GOPATH}/bin/service-entrypoint

test:
	go test ./...

concurrency-test-prepare: build
	go install github.com/qa-dev/jsonwire-grid/testing/webdriver-node-mock
	go install github.com/qa-dev/jsonwire-grid/testing/webdriver-mock-creator
	go install github.com/qa-dev/jsonwire-grid/testing/webdriver-concurrency-test
	CONFIG_PATH=$(TEST_CONFIG_PATH) nohup ${GOPATH}/bin/service-entrypoint >/dev/null 2>&1 &

concurrency-test: concurrency-test-prepare
	webdriver-concurrency-test

