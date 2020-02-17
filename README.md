# ~~jsonwire-grid~~WebDriverGrid [![Build Status](https://travis-ci.org/qa-dev/jsonwire-grid.svg?branch=master)](https://travis-ci.org/qa-dev/jsonwire-grid) [![Go Report Card](https://goreportcard.com/badge/github.com/qa-dev/jsonwire-grid)](https://goreportcard.com/report/github.com/qa-dev/jsonwire-grid) [![codecov](https://codecov.io/gh/qa-dev/jsonwire-grid/branch/master/graph/badge.svg)](https://codecov.io/gh/qa-dev/jsonwire-grid)
This is high-performance scalable implementation of Selenium Grid (hub),
###### What is Selenium-Grid?
>Selenium-Grid allows you run your tests on different machines against different browsers in parallel. That is, running multiple tests at the same time against different machines running different browsers and operating systems. Essentially, Selenium-Grid support distributed test execution. It allows for running your tests in a distributed test execution environment.

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
1. If you see something similar with `{"state":null,"sessionId":"515be56a...` all right! If not, submit [issue](https://github.com/qa-dev/jsonwire-grid/issues/new)


## Configuration
Configurations are stored in json files. Example:
```
{
  "logger": {
    "level": "debug"
  },
  "db": {
    "implementation": "local"
  },
  "grid": {
    "client_type": "selenium",
    "port": 4444,
    "strategy_list": [
      {
        "type": "persistent"
      }
    ],
    "busy_node_duration": "15m",
    "reserved_node_duration": "5m"
  }
}
```

### Logger - Configuration of logger.
| Option        | Possible values                     | Description            | 
| ------------- | ----------------------------------- | ---------------------- |
| logger.level  | `debug`, `info`, `warning`, `error` | Logging level.         |

### DB - Configuration of storage.
| Option            | Possible values          | Description                                |
| ----------------- | ------------------------ | ------------------------------------------ |
| db.implementation | `mysql`, `local`, `mongo`| Select your favorite db, or local storage. |
| db.connection     | see next table           | DSN for your db.                           |

>Note: Note: Local (in-memory) storage does not support common session storage between grid-instances.

| DB implementation | DSN format                                                                                                                                                              |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| mysql             | [spec](https://github.com/go-sql-driver/mysql#dsn-data-source-name), example `db_user:db_pass@(db_host:3306)/db_name?parseTime=true` (parseTime=true - required option) |
| local             | omit this property, because every instance have its own in-memory storage                                                                                               |
| mongo             | NOTE! Mongo db temporary supports only persistent node strategy, example `mongodb://localhost:27017`

### Statsd - Configuration of metrics(optional).
| Option          | Possible values | Description            |
| --------------- | --------------- | ---------------------- |
| statsd.host     | `string`        | Host of statsd server. |
| statsd.port     | `int`           | Port of statsd server. |
| statsd.protocol | `string`        | Network protocol.      |
| statsd.prefix   | `string`        | Prefix of metrics tag. |
| statsd.enable   | `bool`          | Enable metric.         |

### Grid - Configuration of app.
| Option                      | Possible values          | Description                                                                                                                                                    |
| --------------------------- | ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| grid.client_type            | `selenium`, `wda`        | Type of used nodes.                                                                                                                                            |
| grid.port                   | `int`                    | Grid will run on this port.                                                                                                                                    |
| grid.busy_node_duration     | `string` as `12m`, `60s` | Max session lifetime, when timeout was elapsed grid will kill the session.                                                                                     |
| grid.reserved_node_duration | `string` as `12m`, `60s` | Max timeout between send request `POST /session` and opening the browser window. (Deprecated will renamed)                                                     |
| grid.strategy_list          | `array`                  | List of strategies, if grid not able create session on first strategy it go to next, until list ends. [Read more about strategies.](#element-of-strategy-list) |

> * `selenium` - [http://www.seleniumhq.org/]()
> * `wda` - [agent](https://github.com/qa-dev/WebDriverAgent) for [WDA](https://github.com/qa-dev/WebDriverAgent)

### Element of strategy list
| Option                         | Possible values                                                 | Description                                                                  |
| ------------------------------ | --------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| type                           | `string`, see [type of strategy.](#types-of-strategy)           | Host of statsd server.                                                       |
| limit                          | `int`, unlimited if equals `0`                                  | Max count of active nodes on this strategy.                                  |
| params                         | `object`, dependent on [strategy type](#types-of-strategy)      | Object describes available nodes, ex. docker config, kubernetes config, etc. |
| node_list                      | `array`, dependent on [strategy type](#types-of-strategy)       | Array of objects describing available nodes.                                 |
| node_list.[].params            | `object`, dependent on [strategy type](#types-of-strategy)      | Object of describing node, ex. image_name, etc.                              |
| node_list.[].capabilities_list | `array`, ex. [{"foo: "bar"}, {"foo: "baz", "ololo": "trololo"}] | array of objects describes available capabilities .                          |       

### Types of strategy
##### `persistent` - using externally started nodes, same as original selenium grid.
| Strategy option | Possible values | Description                                          |
|---------------- | --------------- | ---------------------------------------------------- |
| limit           | -               | Omit this property, сount of nodes always unlimited. |
| params          | -               | Omit this property.                                  |
| node_list       | -               | Omit this property.                                  |

##### `kubernetes` - on-demand nodes in kubernetes cluster.
| Strategy option             | Possible values | Description                                  |
|--------------------------   | ---------------------- | ------------------------------------- |
| params.namespace            | string                 | Namespace in k8s for on-demand nodes. |
| params.pod_creation_timeout | string as `12m`, `60s` | Max waiting time for creating a pod.  |
| node_list.[].params.image   | string                 | Docker image with selenium.           |
| node_list.[].params.port    | string                 | Port of selenium.                     |

## API
- `/grid/status` - a method returns a status of a grid
- `/grid/session/info` - a returns a session info by session id. 
 Еxample: `curl -X http://localhost:4444/grid/session/info?sessionid=9fc185d2-7a3d-4660-877f-cd4ca2a2f5c3`