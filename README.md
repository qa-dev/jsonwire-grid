# ~~jsonwire-grid~~WebDriverGrid [![Build Status](https://travis-ci.org/qa-dev/jsonwire-grid.svg?branch=master)](https://travis-ci.org/qa-dev/jsonwire-grid)
Golang implementation of Selenium Grid (hub).

## HowTo
#### Run binary file
1. download binary file and
1. `export CONFIG_PATH=./config.json`
1. `./jsonwire-grid`

### Run From Source
#### Requirements
* Go >= 1.8.1
* [go-bindata](https://github.com/jteeuwen/go-bindata)
1. `git clone https://github.com/qa-dev/jsonwire-grid .`
1. `cd jsonwire-grid`
1. `cp config-sample.json config.json`
1. `make run`

## HowToUse
1. Run app
1. `java -jar selenium-server-standalone-3.4.0.jar -role node  -hub http://127.0.0.1:4444/grid/register`
1. try create session `curl -X POST http://127.0.0.1:4444/wd/hub/session -d '{"desiredCapabilities":{"browserName": "firefox"}}'`
