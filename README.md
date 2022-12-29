# graftool
An automation tool which monitor the results displayed on Grafana Dashboard Panels.

## Need
Currently we do not have a simple and easy way to interact with the Grafana/Prometheus API from the command line. Manual cURL commands is not an optimal solution and is very confusing for anyone who wants to monitor the metrics in automation fashion.

`graftool` provides a simplified way to monitor the metrics of grafana dashboard panels by allowing users to use the tool in there automation script. If you want to trigger a functional/load data to your application and want to monitor its behaviour for longer duration then `graftool` is the best automation tool to use.

## Features

<TODO>

## Design

<TODO>

## Installation
Requirements:
- Go 1.19+
Install using, 
```
go install <url>
```
Build using,
```
cd cmd/
go build -o graftool
```
## Usage
```
./graftool --help
usage: graftool --grafana.web.listen-address=GRAFANA.WEB.LISTEN-ADDRESS --prometheus.web.listen-address=PROMETHEUS.WEB.LISTEN-ADDRESS --grafana-dashboards=GRAFANA-DASHBOARDS [<flags>]

A tool to monitor results displayed on Grafana Dashboard Panels

Flags:
  --help         Show context-sensitive help (also try --help-long and --help-man).
  --grafana.web.listen-address=GRAFANA.WEB.LISTEN-ADDRESS
                 Address on which Grafana listen for UI,API
  --grafana-username=GRAFANA-USERNAME
                 Username for Grafana Server. If token is not specified then provide username
  --grafana-password=GRAFANA-PASSWORD
                 Password for Grafana Server. If token is not specified then provide password
  --prometheus-username=PROMETHEUS-USERNAME
                 Username for Prometheus Server. If token is not specified then provide username
  --prometheus-password=PROMETHEUS-PASSWORD
                 Password for Prometheus Server. If token is not specified then provide password
  --prometheus.web.listen-address=PROMETHEUS.WEB.LISTEN-ADDRESS
                 Address on which Prometheus listen for UI,API
  --grafana-dashboard=GRAFANA-DASHBOARDS ...
                 Name of Grafana Dashboard to be monitored
  --token=TOKEN  Bearer Token for connecting to Grafana/Prometheus
  --version      Show application version.
```
