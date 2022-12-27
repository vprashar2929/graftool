package main

import (
	"fmt"
	"log"
	"net/url"
)

type MetricMetadata struct {
	Name     string `json:"__name__"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
}
type MetricResult struct {
	Metric MetricMetadata `json:"metric"`
	Value  []interface{}  `json:"value"`
}
type MetricsData struct {
	ResultType string         `json:"resultType"`
	Result     []MetricResult `json:"result"`
}
type MetricSearchResponse struct {
	Status string       `json:"status"`
	Data   *MetricsData `json:"data"`
}

// MetricSearch will query the prometheus HTTP API from the required metrics
func (c *Client) MetricSearch(params url.Values) (resp MetricSearchResponse, err error) {
	err = c.request("GET", "/api/v1/query", params, nil, &resp)
	return
}

// GetMetricsValue will fetch the metric value from the response
func (d *DashboardResponseData) GetMetricsValue(p *Client, query string) []MetricResult {
	pquery := make(url.Values)
	pquery.Set("query", ParseQuery(query, "", "node-exporter.dashboard-testing.svc.cluster.local:9100", "node-exporter", "5m"))
	presp, err := p.MetricSearch(pquery)
	if err != nil {
		log.Fatal(fmt.Sprintf("Prometheus Error: "), err)
	}
	return presp.Data.Result

}
