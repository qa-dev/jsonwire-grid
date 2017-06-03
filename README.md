# ~~jsonwire-grid~~WebDriverGrid
Golang implementation of Selenium Grid (hub).

## HowTo
#### Run binary file
1. download binary file
1. `export CONFIG_PATH=./config.json`
1. `./webdriver-{platform_name}`

### Run From Source
#### Requirements
* Go >= 1.8.1
1. `git clone https://github.com/qa-dev/webdriver-grid .`
1. `cd webdriver-grid`
1. `export CONFIG_PATH=./config.json`
1. `go run main.go`

## HowToUse
1. Run app
1. `java -jar selenium-server-standalone-3.4.0.jar -role node  -hub http://127.0.0.1:4444/grid/register`
1. try create session `curl -X POST http://127.0.0.1:4444/wd/hub/session -d '{"desiredCapabilities":{"browserName": "firefox"}}'`