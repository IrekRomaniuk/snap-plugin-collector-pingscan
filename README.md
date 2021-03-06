![Alt text](https://img.shields.io/badge/version-development-red.svg)
# snap-plugin-collector-pingcount
This plugin counts IP addresses of responding hosts.

It's used in the [Snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
3. [License](#license-and-authors)
4. [Releases] (#Releases)
5. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.5+](https://golang.org/dl/)  - needed only for building. See also [How to install Go language](http://ask.xmodulo.com/install-go-language-linux.html)

### Operating systems
All OSs currently supported by Snap:
* Linux/amd64

### Installation
#### To build the plugin binary:
```
$ go get -u github.com/IrekRomaniuk/snap-plugin-collector-pingcount
```
### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Load the plugin and create a task, see example in [Examples](https://github.com/IrekRomaniuk/snap-plugin-collector-pingcount/tree/master/examples).

## Documentation

### Collected Metrics

This plugin has the ability to gather the following metric:

Namespace | Description
----------|-----------------------
/niuk/pingcount/total-up | total number of hosts responding


### Example
Example running pingcount collector and writing data to an Influx database.

Load pingcount plugin
```
$ snaptel plugin load $GOPATH/bin/snap-plugin-collector-pingcount
```
List available plugins
```
$ snaptel task watch 4df1ddea-11ef-49e9-867b-6f19658cf16e
Watching Task (4df1ddea-11ef-49e9-867b-6f19658cf16e):
NAMESPACE                        DATA    TIMESTAMP
niuk/pingcount/total-up          1102    2016-12-08 14:58:54.176178073 -0500 EST
```
See available metrics for your system
```
$ snaptel metric list
```

Create a task manifest file (see below) and put full path to the [file](https://github.com/IrekRomaniuk/snap-plugin-collector-pingcount/blob/master/examples/pinglist.txt) listing IP addresses:
```yaml
deadline: "15s"
version: 1
schedule:
  type: "simple"
  interval: "30s"
max-failures: 10
workflow:
  collect:
    metrics:
      /niuk/pingcount/total-up: {}
    config:
      /niuk/pingcount:
        target: "/home/global/path/examples/pinglist.txt"
```
Load influxdb plugin for publishing:
```
$ snaptel plugin load snap-plugin-publisher-influxdb
```

Create a task:
```
$ snaptel task create -t pingcount.yml -n pingcount
Using task manifest to create task
Task created
ID: 4df1ddea-11ef-49e9-867b-6f19658cf16e
Name: pingcount-shields
State: Running
```

List running tasks:
```
$ snaptel task list
ID                                       NAME                                            STATE           HIT     MISS    FAIL    CREATED                 LAST FAILURE
4df1ddea-11ef-49e9-867b-6f19658cf16e     pingcount-shields                                Running         53      0       0       2:50PM 12-08-2016                                         
```
Watch the task
```
$ snaptel task watch 4df1ddea-11ef-49e9-867b-6f19658cf16e
Watching Task (4df1ddea-11ef-49e9-867b-6f19658cf16e):
NAMESPACE                        DATA    TIMESTAMP
niuk/pingcount/total-up          1102    2016-12-08 14:58:54.176178073 -0500 EST
```
Watch metrics in real-time using [Snap plugin for Grafana] (https://blog.raintank.io/using-grafana-with-intels-snap-for-ad-hoc-metric-exploration/)
and use InfluxDB plugin for publishing ![Alt text](examples/grafana-pingcount.JPG "Metrics published to InfluxDB")

## License
This plugin is Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@IrekRomaniuk](https://github.com/IrekRomaniuk/)


