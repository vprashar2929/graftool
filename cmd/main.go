package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/fatih/color"
	"github.com/rodaine/table"

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
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Dashboard Name", "Row Title", "Panel Title", "Labels", "TimeStamp", "Metric Value")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for _, uid := range d.UID {
		for _, row := range d.Rows[uid] {
			for _, panel := range d.FilterResp[uid].FilterPanel[row] {
				for _, target := range panel.Targets {
					if len(d.FilterResp[uid].Metric[target.Expr]) > 0 {
						tbl.AddRow(d.DashboardResponse[uid].Dashboard.Title, row, panel.Title, target.Legends, ParseEpoch(d.FilterResp[uid].Metric[target.Expr][0].Value[0]), d.FilterResp[uid].Metric[target.Expr][0].Value[1])
					} else {
						tbl.AddRow(d.DashboardResponse[uid].Dashboard.Title, row, panel.Title, target.Legends, d.FilterResp[uid].Metric[target.Expr], d.FilterResp[uid].Metric[target.Expr])
					}

				}

			}

		}

	}
	tbl.Print()

}
func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Vibhu Prashar")
	kingpin.CommandLine.Help = "A tool to monitor results displayed on Grafana Dashboard Panels"
	kingpin.Parse()
	grafanaClient := GetGrafanaClient(*grafanaBaseURL, *grafanaUsername, *grafanaPassword, *token)
	d := new(DashboardResponseData)
	d.GetDashboards(*grafanaClient, *grafanaDashboards)
	d.GetDashboardByUID(*grafanaClient)
	d.FilterData()
	prometheusClient := GetPrometheusClient(*prometheusBaseURL, *token, *prometheusUsername, *prometheusPassword)
	d.GetDashboardMetricsFromResponse(prometheusClient)
	DisplayReport(d)

}
