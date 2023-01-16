package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

// Creating API client
type Client struct {
	config  Config
	baseURL url.URL
	client  *http.Client
}

// Config contains client configuration
type Config struct {
	// API Key is an optional API key
	APIKEY string
	// Baset Auth is optional basic auth credentials
	BaseAuth *url.Userinfo
	// HTTP Headers are optional HTTP Headers
	HTTPHeaders map[string]string
	// Client provides an optional HTTP client, otherwise we will use default one
	Client *http.Client
	// OrgID provides an optional organizational ID, BasicAuth defaults to last used org
	OrgID int64
	// NumRetries contains the number of attempted retries
	NumRetries int
}

// Create a new Client
func New(baseURL string, cfg Config) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if cfg.BaseAuth != nil {
		u.User = cfg.BaseAuth
	}
	cli := cfg.Client
	if cli == nil {
		cli = cleanhttp.DefaultClient()
	}
	return &Client{
		config: cfg, baseURL: *u, client: cli,
	}, nil
}

func (c *Client) Request(method, requestPath string, query url.Values, body io.Reader, responseStruct interface{}) error {

	var (
		req          *http.Request
		resp         *http.Response
		err          error
		bodyContents []byte
	)
	// Now if we want to retry a request that sends data, we will need to stash the request data in memory. Otherwise we lose it since readers cannot be replayed

	var bodyBuffer bytes.Buffer
	if c.config.NumRetries > 0 && body != nil {
		body = io.TeeReader(body, &bodyBuffer)
	}
	// Now retry logic

	for i := 0; i <= c.config.NumRetries; i++ {
		// If it is not the first request, re-use the request body
		if i > 0 {
			body = bytes.NewBuffer(bodyBuffer.Bytes())
		}
		req, err = c.NewRequest(method, requestPath, query, body)
		if err != nil {
			return err
		}
		// Wait if thats not the first time
		if i != 0 {
			time.Sleep(time.Second * 5)
		}
		resp, err = c.client.Do(req)
		// If err is not nil retry again
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		// read the body and HTTP codes as thats what the unit testing expects
		bodyContents, err = ioutil.ReadAll(resp.Body)

		// If there was an error reading the body
		if err != nil {
			continue
		}
		// Exit the loop if we have something final to return. This cant be <500
		if resp.StatusCode < http.StatusInternalServerError && resp.StatusCode != http.StatusTooManyRequests {
			break
		}

	}
	if err != nil {
		return err
	}
	if os.Getenv("GF_LOG") != "" {
		log.Printf("response status %d with body %v", resp.StatusCode, string(bodyContents))
	}
	// check the status code
	if resp.StatusCode >= 400 {
		return fmt.Errorf("status: %d, body: %v", resp.StatusCode, string(bodyContents))
	}
	if responseStruct == nil {
		return nil
	}

	err = json.Unmarshal(bodyContents, responseStruct)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) NewRequest(method, requestPath string, query url.Values, body io.Reader) (*http.Request, error) {
	url := c.baseURL
	url.Path = path.Join(url.Path, requestPath)
	url.RawQuery = query.Encode()
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return req, err
	}
	if c.config.APIKEY != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKEY))
	} else if c.config.OrgID != 0 {
		req.Header.Add("X-Grafana-Org-Id", strconv.FormatInt(c.config.OrgID, 10))
	}
	if c.config.HTTPHeaders != nil {
		for k, v := range c.config.HTTPHeaders {
			req.Header.Add(k, v)
		}
	}

	if os.Getenv("GF_LOG") != "" {
		if body == nil {
			log.Printf("request (%s) to %s with no body data", method, url.String())
		} else {
			log.Printf("request (%s) to %s with body data: %s", method, url.String(), body.(*bytes.Buffer).String())
		}
	}
	req.Header.Add("Content-Type", "application/json")
	return req, err
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
func GetPrometheusClient(baseURL, token, username, password string) *Client {
	var p *Client
	var err error
	if username != "" && password != "" {
		p, err = New(fmt.Sprintf("http://%s", baseURL), Config{BaseAuth: url.UserPassword(username, password)})
	} else if token != "" {
		p, err = New(fmt.Sprintf("http://%s", baseURL), Config{APIKEY: token})
	} else {
		p, err = New(fmt.Sprintf("http://%s", baseURL), Config{})
	}
	if err != nil {
		log.Fatal(fmt.Sprintf("Prometheus Client Error: "), err)
	}
	return p
}

func GetRequest(method, url string, c *Client, params url.Values, resp interface{}) error {
	err := c.Request(method, url, params, nil, &resp)
	return err
}
