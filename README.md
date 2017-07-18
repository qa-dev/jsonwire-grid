# ~~jsonwire-grid~~WebDriverGrid [![Build Status](https://travis-ci.org/qa-dev/jsonwire-grid.svg?branch=master)](https://travis-ci.org/qa-dev/jsonwire-grid)
This is high-performance scalable implementation of Selenium Grid (hub),

## Features
* One session per one node, no more no less😺
* Scaling grid-instances for fault-tolerance
* Support and effective management over 9000 nodes, for parallel testing👹
* Single entry point for all your test apps
* Send metrics to [statsd](https://github.com/etsy/statsd)
* Support on-demand nodes in Kubernetes cluster (Only if grid running in cluster)


## HowTo
### Run grid
1. [Download last release](https://github.com/qa-dev/jsonwire-grid/releases) and unzip
1. cd to `jsonwire-grid_vXXX`
1. Type `export CONFIG_PATH=./config-local-sample.json`
1. Type `./jsonwire-grid`
1. Grid running!

### Run nodes
1. [Download selenium](http://www.seleniumhq.org/download/)
1. `java -jar selenium-server-standalone-3.4.0.jar -role node  -hub http://127.0.0.1:4444/grid/register`
1. Repeat!

### Run test
1. Try create session, such as `curl -X POST http://127.0.0.1:4444/wd/hub/session -d '{"desiredCapabilities":{"browserName": "firefox"}}'`
1. If your something similar with `{"state":null,"sessionId":"515be56a...` all right! If not, submit [issue](https://github.com/qa-dev/jsonwire-grid/issues/new)


## Configuration
### Logger - Configuration of logger.
##### `logger.level` - Support `debug`, `info`, `warning` or `error`.

### DB - Configuration of storage.
##### `db.implementation` - Select your favorite db, or local storage.
Now support's: `mysql`, `local`.
>Note: Local (in-memory) storage not support single session storagen between grid-instances.
##### `db.connection` - DSN for your db
* `mysql` - example `db_user:db_pass@(db_host:3306)/db_name?parseTime=true` (parseTime=true - required option)
* `local` - omit this property, because every instance have own in-memory storage

### Statsd - Configuration of metrics.
##### `statsd.host` - Host of statsd server.
##### `statsd.port` - Host of statsd server `int`.
##### `statsd.protocol` - Network protocol.
##### `statsd.prefix` - Prefix of metrics tag.
##### `statsd.enable` - Enable metrics `true/false`.

### Grid - Configuration of app.
##### `grid.client_type` - Type of used nodes.
* `selenium` - [http://www.seleniumhq.org/]()
* `wda` - [agent](https://github.com/qa-dev/WebDriverAgent) for [WDA](https://github.com/qa-dev/WebDriverAgent)
##### `grid.port` - grit will run on this port.
##### `grid.busy_node_duration` - max session lifetime, when timeout was elapsed grid will kill the session.
##### `grid.reserved_node_duration` - max timeout between send request `POST /session` and opening the browser window, if. (Deprecated will renamed)
##### `grid.strategy_list` - list of strategies, if grid not able create session on first strategy it go to next, until list ends. [Read more about strategies.](#strategy-list)

### Strategy list
##### `type` - [type of strategy.](#types-of-strategies)
##### `limit` - max count active nodes of this strategy. Unlimited if equals `0`. Dependent of strategy [type](#types-of-strategies)
##### `params` - object describes available nodes `ex. docker config, kubernetes config, etc.`. Dependent of strategy [type](#types-of-strategies)
##### `node_list` - list of objects describes available nodes.
##### `node_list.[].params` - list of objects describes available nodes `ex. image_name, etc.`. Dependent of strategy [type](#types-of-strategies)
##### `node_list.[].capabilities_list` - array of available capabilities objects. `ex. [{"foo: "bar"}, {"foo: "baz", "ololo": "trololo"}]`

### Types of strategies
##### `persistent` - working with registered nodes, same as original selenium grid.
* `limit` - omit this property, its always equals `0`.
* `node_list` - omit this property.
##### `kubernetes` - on-demand nodes in kubernetes cluster.
* `params` - omit this property.
* `node_list.[].params.image` - docker image with selenium.
* `node_list.[].params.port` - port of selenium.

