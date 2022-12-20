package main

import (
	"fmt"
	"log"
	"math/big"
	"net/url"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	grafanaBaseURL     = kingpin.Flag("grafana.web.listen-address", "Address on which Grafana listen for UI,API").Required().String()
	grafanaUsername    = kingpin.Flag("grafana-username", "Username for Grafana Server. If token is not specified then provide username").String()
	grafanaPassword    = kingpin.Flag("grafana-password", "Password for Grafana Server. If token is not specified then provide password").String()
	prometheusUsername = kingpin.Flag("prometheus-username", "Username for Prometheus Server. If token is not specified then provide username").String()
	prometheusPassword = kingpin.Flag("prometheus-password", "Password for Prometheus Server. If token is not specified then provide password").String()
	prometheusBaseURL  = kingpin.Flag("prometheus.web.listen-address", "Address on which Prometheus listen for UI,API").Required().String()
	grafanaDashboards  = kingpin.Flag("grafana-dashboards", "List of Grafana Dashboards to be monitored").Required().Strings()
	token              = kingpin.Flag("token", "Bearer Token for connecting to Grafana/Prometheus").String()
	//interval          = kingpin.Flag("interval", "Set interval for monitoring").Default("1m").Duration()
)

type MetricOutput struct {
	MetricName  string
	MetricValue []MetricResult
}

// GetGrafanaClient will create a client for Grafana HTTP URL
func GetGrafanaClient(baseURL, username, password, token string) *Client {
	var c *Client
	var err error
	if username != "" && password != "" {
		c, err = New(fmt.Sprintf("http://%s", baseURL), Config{BaseAuth: url.UserPassword(username, password)})
	} else if token != "" {
		c, err = New(fmt.Sprintf("http://%s", baseURL), Config{APIKEY: token})
	} else {
		log.Fatal("Please provide either username/password or Token")
	}
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// GetPrometheusClient will create a client for Prometheus HTTP URL
func GetPrometheusClient(baseURL, token, promusername, prompassword string) *Client {
	var p *Client
	var err error
	if promusername != "" && prompassword != "" {
		p, err = New(fmt.Sprintf("http://%s", baseURL), Config{BaseAuth: url.UserPassword(promusername, prompassword)})
	} else {
		p, err = New(fmt.Sprintf("http://%s", baseURL), Config{})
	}
	if err != nil {
		log.Fatal(fmt.Sprintf("Prometheus Client Error: "), err)
	}
	return p
}

// DisplayReport will display the end report on stdout
func DisplayReport(d *DashboardResponseData) {
	for i := 0; i < len(d.Data); i++ {
		fmt.Printf("Dashboard Name: %s ", d.Data[i].Title)
		fmt.Println()
		for j := 0; j < len(d.Data[i].Panels); j++ {
			fmt.Printf(" Panel Title: %s ", d.Data[i].Panels[j])
			fmt.Println()
			fmt.Printf("  Metric Name: %s ", d.Data[i].Metrics[d.Data[i].Panels[j]].MetricName)
			fmt.Println()
			if len(d.Data[i].Metrics[d.Data[i].Panels[j]].MetricValue) != 0 {
				fmt.Printf("   TimeStamp: %v ", ParseEpoch(d.Data[i].Metrics[d.Data[i].Panels[j]].MetricValue[0].Value[0]))
				fmt.Println()
				fmt.Printf("   Metric Value: %s ", d.Data[i].Metrics[d.Data[i].Panels[j]].MetricValue[0].Value[1])
			} else {
				fmt.Printf("   TimeStamp: %v ", d.Data[i].Metrics[d.Data[i].Panels[j]].MetricValue)
				fmt.Println()
				fmt.Printf("   Metric Value: %s ", d.Data[i].Metrics[d.Data[i].Panels[j]].MetricValue)
			}
			fmt.Println()
		}
		fmt.Println()
		fmt.Println()
	}

}

func ParseEpoch(timestamp interface{}) time.Time {
	flt, _, err := big.ParseFloat(fmt.Sprint("", timestamp), 10, 0, big.ToNearestEven)
	if err != nil {
		log.Fatal(err)

	}
	i, _ := flt.Int64()
	return time.Unix(i, 0)
}

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Vibhu Prashar")
	kingpin.CommandLine.Help = "A tool to monitor results displayed on Grafana Dashboard Panels"
	kingpin.Parse()
	grafanaClient := GetGrafanaClient(*grafanaBaseURL, *grafanaUsername, *grafanaPassword, *token)
	d := new(DashboardResponseData)
	d.GetDashboards(*grafanaClient, *grafanaDashboards)
	d.GetDashboardByUID(*grafanaClient)
	d.GetDashboardModelFromResponse(*grafanaClient)
	d.GetDashboardMetricsFromResponse(*grafanaClient)
	prometheusClient := GetPrometheusClient(*prometheusBaseURL, *token, *prometheusUsername, *prometheusPassword)
	d.GetMetricsValue(prometheusClient)
	DisplayReport(d)

}
