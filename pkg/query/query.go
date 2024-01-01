package query

import (
	"log"
	"net/url"

	"github.com/vprashar2929/graftool/pkg/client"
	"github.com/vprashar2929/graftool/pkg/parse"
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
	Data   *MetricsData `json:"data"`
	Status string       `json:"status"`
}

// MetricSearch will query the prometheus HTTP API from the required metrics
func MetricSearch(params url.Values, p *client.Client) (resp MetricSearchResponse, err error) {
	err = client.GetRequest("GET", "/api/v1/query", p, params, &resp)
	return
}

// GetMetricsValue will fetch the metric value from the response
func GetMetricsValue(p *client.Client, query string, conf *parse.Config) []MetricResult {
	pquery := make(url.Values)
	q, err := parse.ParseQuery(query, conf)
	if err != nil {
		return nil // TODO: Return error from here
	}
	pquery.Set("query", q)
	presp, err := MetricSearch(pquery, p)
	if err != nil {
		log.Fatal(err)
	}
	return presp.Data.Result
}
