export CONFIG_PATH=./config.json
all: build

.PHONY: help build fmt clean run test coverage check vet lint doc cfpush

help:
	@echo "build - build application from sources"
	@echo "fmt   - format application sources"
	@echo "run   - start application"
	@echo "gen   - generate files for json rpc 2.0 service"

build: fmt
	go build -o ${GOPATH}/bin/service-entrypoint

fmt:
	go fmt

run: build
	${GOPATH}/bin/service-entrypoint