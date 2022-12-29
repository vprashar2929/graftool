package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

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
	//interval          = kingpin.Flag("interval", "Set interval for monitoring").Default("1m").Duration()
)
var (
	startTime int64
)

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

func DisplayReport(d *DashboardResponseData) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.SetAutoIndex(true)
	t.AppendHeader(table.Row{"Row Title", "Panel Title", "Legends", "TimeStamp", "Metric Value"})
	for _, uid := range d.UID {
		t.SetTitle(d.DashboardResponse[uid].Dashboard.Title)
		for _, row := range d.Rows[uid] {
			for _, panel := range d.FilterResp[uid].FilterPanel[row] {
				for _, target := range panel.Targets {
					if len(d.FilterResp[uid].Metric[target.Expr]) > 0 {
						t.AppendRow(table.Row{row, panel.Title, target.Legends, ParseEpoch(d.FilterResp[uid].Metric[target.Expr][0].Value[0]), d.FilterResp[uid].Metric[target.Expr][0].Value[1]})
					} else {
						t.AppendRow(table.Row{row, panel.Title, target.Legends, d.FilterResp[uid].Metric[target.Expr], d.FilterResp[uid].Metric[target.Expr]})
					}
				}

			}
			t.AppendSeparator()

		}
		t.SetCaption(fmt.Sprint("Dashboard Link: http://%s%s?from=%d&to=%d"), *grafanaBaseURL, d.URL[uid], startTime, time.Now().UnixMilli())
		t.Render()
	}

}
func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Vibhu Prashar")
	kingpin.CommandLine.Help = "A tool to monitor results displayed on Grafana Dashboard Panels"
	kingpin.Parse()
	grafanaClient := GetGrafanaClient(*grafanaBaseURL, *grafanaUsername, *grafanaPassword, *token)
	d := new(DashboardResponseData)
	d.GetDashboards(*grafanaClient, *grafanaDashboard)
	d.GetDashboardByUID(*grafanaClient)
	d.FilterData()
	prometheusClient := GetPrometheusClient(*prometheusBaseURL, *token, *prometheusUsername, *prometheusPassword)
	startTime = time.Now().UnixMilli()
	d.GetDashboardMetricsFromResponse(prometheusClient)
	DisplayReport(d)

}
