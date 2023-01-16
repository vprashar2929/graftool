package parse_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"time"

	"github.com/vprashar2929/graftool/pkg/parse"
)

var _ = Describe("Parse", func() {

	Describe("Parse Query", func() {
		var (
			queryWithParse    = `(sum by(instance) (irate(node_cpu_seconds_total{instance="$node",job="$job", mode!="idle"}[$__rate_interval])) / on(instance) group_left sum by (instance)((irate(node_cpu_seconds_total{instance="$node",job="$job"}[$__rate_interval])))) * 100`
			namespace         = "default-instance"
			node              = "default-node"
			job               = "default-job"
			interval          = "5m"
			queryWithoutParse = `node_memory_MemTotal_bytes`
			resWithParse      = `(sum by(instance) (irate(node_cpu_seconds_total{instance="default-node",job="default-job", mode!="idle"}[5m])) / on(instance) group_left sum by (instance)((irate(node_cpu_seconds_total{instance="default-node",job="default-job"}[5m])))) * 100`
			resWithoutParse   = "node_memory_MemTotal_bytes"
		)

		Context("When all parameters are given", func() {
			It("Parse Query", func() {
				q, err := parse.ParseQuery(queryWithParse, namespace, node, job, interval)
				Expect(q).To(Equal(resWithParse))
				Expect(err).NotTo(HaveOccurred())
			})

		})
		Context("When query doesnt contain the variables", func() {
			It("returns query without parsing", func() {
				q, err := parse.ParseQuery(queryWithoutParse, namespace, node, job, interval)
				Expect(q).To(Equal(resWithoutParse))
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When arguments are not provided", func() {
			It("returns error when no arguments are provided", func() {
				q, err := parse.ParseQuery("", namespace, node, job, interval)
				Expect(err).To(HaveOccurred())
				Expect(q).To(Equal(""))
			})
		})
	})
	Describe("Parse Epoch Timestamp", func() {
		epochTime := "1672986081.002"
		expectedTime := time.Unix(1672986081, 0)
		It("returns the human readable timestamp", func() {
			t := parse.ParseEpoch(epochTime)
			Expect(t).To(Equal(expectedTime))
		})

	})
})
