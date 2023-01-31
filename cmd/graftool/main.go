package main

import (
	"time"

	"github.com/vprashar2929/graftool/pkg/client"
	"github.com/vprashar2929/graftool/pkg/dashboard"
	"github.com/vprashar2929/graftool/pkg/parse"
	"github.com/vprashar2929/graftool/pkg/report"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	grafanaBaseURL     = kingpin.Flag("grafana.web.listen-address", "Address on which Grafana listen for UI,API").Required().String()
	grafanaUsername    = kingpin.Flag("grafana-username", "Username for Grafana Server. If token is not specified then provide username").String()
	grafanaPassword    = kingpin.Flag("grafana-password", "Password for Grafana Server. If token is not specified then provide password").String()
	prometheusUsername = kingpin.Flag("prometheus-username", "Username for Prometheus Server. If token is not specified then provide username").String()
	prometheusPassword = kingpin.Flag("prometheus-password", "Password for Prometheus Server. If token is not specified then provide password").String()
	prometheusBaseURL  = kingpin.Flag("prometheus.web.listen-address", "Address on which Prometheus listen for UI,API").Required().String()
	grafanaDashboard   = kingpin.Flag("grafana-dashboard", "Name of Grafana Dashboard to be monitored").Required().Strings()
	token              = kingpin.Flag("token", "Bearer Token for connecting to Grafana/Prometheus").String()
	interval           = kingpin.Flag("interval", "Set interval for monitoring").Duration()
	step               = kingpin.Flag("step", "Step to pool the dashboard response. Should be less than interval").Duration()
)
var (
	fromTime int64

	startTime = time.Now()
)

func process(prometheusClient *client.Client, d *dashboard.DashboardResponseData, conf *parse.Config) {
	fromTime = time.Now().UnixMilli()
	dashboard.GetDashboardMetricsFromResponse(prometheusClient, d, conf)
	report.DisplayReport(d, grafanaBaseURL, fromTime)
}

func main() {

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Vibhu Prashar")
	kingpin.CommandLine.Help = "A tool to monitor results displayed on Grafana Dashboard Panels"
	kingpin.Parse()
	conf := parse.ParseConfig()
	grafanaClient := client.GetGrafanaClient(*grafanaBaseURL, *grafanaUsername, *grafanaPassword, *token)
	d := new(dashboard.DashboardResponseData)
	dashboard.GetDashboards(grafanaClient, d, *grafanaDashboard)
	dashboard.GetDashboardByUID(grafanaClient, d)
	dashboard.Filter(d)
	prometheusClient := client.GetPrometheusClient(*prometheusBaseURL, *token, *prometheusUsername, *prometheusPassword)
	dur, _ := time.ParseDuration("0s")
	if *interval != dur && *step != dur {
		for time.Since(startTime).Truncate(*interval) != *interval {
			process(prometheusClient, d, conf)
			time.Sleep(*step)
		}
	} else {
		process(prometheusClient, d, conf)
	}

}
