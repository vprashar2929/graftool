package e2e

import (
	"fmt"
	"testing"

	e2edb "github.com/efficientgo/e2e/db"
	e2emon "github.com/efficientgo/e2e/monitoring"

	// e2einteractive "github.com/efficientgo/e2e/interactive"
	"github.com/efficientgo/core/testutil"
	"github.com/efficientgo/e2e"
)

func TestE2E(t *testing.T) {

	e, err := e2e.New()
	t.Cleanup(e.Close)
	testutil.Ok(t, err)
	fmt.Println("Starting Node Exporter")
	n := e2emon.AsInstrumented(e.Runnable("node-exporter").WithPorts(map[string]int{"http": 9100}).Init(e2e.StartOptions{Image: "quay.io/prometheus/node-exporter:latest"}), "http")
	testutil.Ok(t, e2e.StartAndWaitReady(n))
	config := fmt.Sprint(`
global:
  scrape_interval: 15s
scrape_configs:
- job_name: 'node-exporter'
  scrape_interval: 5s
  static_configs:
  - targets: ['node-exporter:9100']
`)
	fmt.Println("Starting Prometheus")
	p1 := e2edb.NewPrometheus(e, "prometheus", e2edb.Option(e2edb.WithImage("quay.io/rh_ee_vprashar/custom-prometheus:v1.0")))
	testutil.Ok(t, p1.SetConfigEncoded([]byte(config)))
	testutil.Ok(t, e2e.StartAndWaitReady(p1))
	prometheus_URL := fmt.Sprintf("http://%s", p1.Endpoint("http"))
	fmt.Printf("Prometheus URL: %s\n", prometheus_URL)
	// testutil.Ok(t, e2einteractive.OpenInBrowser(prometheus_URL))
	g1 := e2emon.AsInstrumented(e.Runnable("grafana").WithPorts(map[string]int{"http": 3000}).Init(e2e.StartOptions{Image: "quay.io/rh_ee_vprashar/custom-grafana:v1.0"}), "http")
	testutil.Ok(t, e2e.StartAndWaitReady(g1))
	grafana_URL := fmt.Sprintf("http://%s", g1.Endpoint("http"))

	fmt.Printf("Grafana URL: %s\n", grafana_URL)
	//testutil.Ok(t, e2einteractive.RunUntilEndpointHit())

}
