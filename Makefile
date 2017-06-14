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
	@echo "Install jsonwire-grid"
	go install github.com/qa-dev/jsonwire-grid

fmt:
	go fmt

gen:
	go-bindata -pkg mysql -o storage/migrations/mysql/bindata.go storage/migrations/mysql

get-deps:
	go get -u github.com/jteeuwen/go-bindata/...

prepare: fmt gen

run: build
	@echo "Start jsonwire-grid"
	jsonwire-grid

test:
	go test ./...

concurrency-test-prepare: build
	@echo "Install jsonwire-grid"
	go install github.com/qa-dev/jsonwire-grid/testing/webdriver-node-mock
	@echo "Install webdriver-node-mock"
	go install github.com/qa-dev/jsonwire-grid/testing/webdriver-mock-creator
	@echo "Install webdriver-mock-creator"
	go install github.com/qa-dev/jsonwire-grid/testing/webdriver-concurrency-test
	@echo "Kill all running jsonwire-grid"
	killall -9 jsonwire-grid &
	@echo "Wait 1s"
	sleep 1
	@echo "Start jsonwire-grid"
	CONFIG_PATH=$(TEST_CONFIG_PATH) jsonwire-grid &

concurrency-test: concurrency-test-prepare
	@echo "Start webdriver-concurrency-test"
	webdriver-concurrency-test

