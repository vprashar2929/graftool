package client_test

import (
	//"fmt"

	"fmt"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vprashar2929/graftool/pkg/client"
	"github.com/vprashar2929/graftool/pkg/query"
)

var _ = Describe("Client", func() {

	var (
		clientConfig   = client.Config{}
		invalidBaseURL = "localhost/_20_%+off_9000_"
		validBaseURL   = "127.0.0.1:9090"
		// usernmae                 = ""
		// password                 = ""
		token                 = ""
		clientConfigWithToken = client.Config{APIKEY: token}
		validMethod           = "GET"
		validURL              = "/api/v1/query"
		validClient           *client.Client
		res                   query.MetricSearchResponse
	)

	It("Creating new Client", func() {
		By("Check New function in case of invalid base URL")
		parseURL := fmt.Sprintf("http://%s", invalidBaseURL)
		_, err := client.New(parseURL, clientConfig)
		Expect(err).To(HaveOccurred())

	})
	It("When valid base URL is provided", func() {
		By("Check New function in case of valid base URL")
		parseURL := fmt.Sprintf("http://%s", validBaseURL)
		_, err := client.New(parseURL, clientConfig)
		Expect(err).NotTo(HaveOccurred())

	})
	It("When token is provided ", func() {
		By("Check new function in case of API Token")
		_, err := client.New(validBaseURL, clientConfigWithToken)
		Expect(err).NotTo(HaveOccurred())
	})
	It("When valid query is provided Get Request", func() {
		By("Create request with baseURL and config")
		parseURL := fmt.Sprintf("http://%s", validBaseURL)
		validClient, _ = client.New(parseURL, clientConfig)
		params := make(url.Values)
		params.Set("query", "prometheus_http_requests_total")
		err := client.GetRequest(validMethod, validURL, validClient, params, &res)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Status).To(Equal("success"))
	})
})
