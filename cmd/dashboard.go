package main

import (
	"fmt"
	"log"
	"net/url"
)

type DashboardTargets struct {
	Legends string `json:"legendFormat"`
	Expr    string `json:"expr"`
}
type DashboardPanel struct {
	DataSource  map[string]string  `json:"datasource"`
	Targets     []DashboardTargets `json:"targets"`
	Description string             `json:"description"`
	Title       string             `json:"title"`
	Type        string             `json:"type"`
}
type FilterData struct {
	FilterPanel map[string][]DashboardPanel
	Metric      map[string][]MetricResult
}
type DashboardResponse struct {
	Panels []DashboardPanel `json:"panels"`
	Title  string           `json:"title"`
	UID    string           `json:"uid"`
	//Annotations map[string]interface{} `json:"annotations"`
}
type Response struct {
	Dashboard DashboardResponse `json:"dashboard"`
}
type FolderDashboardSearchResponse struct {
	ID          uint     `json:"id"`
	UID         string   `json:"uid"`
	Title       string   `json:"title"`
	URI         string   `json:"uri"`
	URL         string   `json:"url"`
	Slug        string   `json:"slug"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags"`
	IsStarred   bool     `json:"isStarred"`
	FolderID    uint     `json:"folderId"`
	FolderUID   string   `json:"folderUid"`
	FolderTitle string   `json:"folderTitle"`
	FolderURL   string   `json:"folderUrl"`
}

type DashboardResponseData struct {
	UID               []string
	URL               map[string]string
	DashboardResponse map[string]*Response
	Rows              map[string][]string
	FilterResp        map[string]*FilterData
}

func (c *Client) dashboard(path string) (*Response, error) {
	result := &Response{}
	err := c.request("GET", path, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}
func (c *Client) DashboardByUID(uid string) (*Response, error) {
	return c.dashboard(fmt.Sprintf("/api/dashboards/uid/%s", uid))

}

// FolderDashboardSearch uses the folder and dashboard search endpoint to find dashboards based on the params passed in.
func (c *Client) FolderDashboardSearch(params url.Values) (resp []FolderDashboardSearchResponse, err error) {
	err = c.request("GET", "/api/search", params, nil, &resp)
	return
}

// GetDashboards will fetch the uid of dashboards which matches the search
func (d *DashboardResponseData) GetDashboards(c Client, dashboards []string) {
	d.URL = make(map[string]string)
	for _, dashboard := range dashboards {
		query := make(url.Values)
		query.Set("query", dashboard)
		resp, err := c.FolderDashboardSearch(query)
		if err != nil {
			log.Fatal(err)
		}
		for _, val := range resp {
			d.UID = append(d.UID, val.UID)
			d.URL[val.UID] = val.URL
		}
	}

}

// GetDashboardByUID will fetch the dashboards by uid from grafana
func (d *DashboardResponseData) GetDashboardByUID(c Client) {
	d.DashboardResponse = make(map[string]*Response)
	if len(d.UID) == 0 {
		log.Fatal("No Dashboard Exists on Grafana")
	}
	for _, uid := range d.UID {
		res, err := c.DashboardByUID(uid)
		if err != nil {
			log.Fatal(err)
		}
		d.DashboardResponse[uid] = res
	}
}

// GetDashboardMetricsFromResponse fetch metric value from prometheus
func (d *DashboardResponseData) GetDashboardMetricsFromResponse(p *Client) {
	for _, uid := range d.UID {
		d.FilterResp[uid].Metric = make(map[string][]MetricResult)
		if len(d.Rows[uid]) == 0 {
			log.Fatal("No Rows found in Grafana Dashboard -> ", d.DashboardResponse[uid].Dashboard.Title)
		}
		for _, title := range d.Rows[uid] {
			if len(d.FilterResp[uid].FilterPanel[title]) == 0 {
				log.Fatal("No Panels found in Grafana Dashboard -> ", d.DashboardResponse[uid].Dashboard.Title)
			}
			for _, panel := range d.FilterResp[uid].FilterPanel[title] {
				if len(panel.Targets) == 0 {
					log.Fatal("No Metrics found in Grafana Dashboard -> ", d.DashboardResponse[uid].Dashboard.Title)
				}
				for _, target := range panel.Targets {
					res := d.GetMetricsValue(p, target.Expr)
					d.FilterResp[uid].Metric[target.Expr] = res

				}
			}

		}

	}
}

// FilterData will filter the dashboard data on the basis of how many rows are present in dashboard and corresponding to rows how many panels are there.
func (d *DashboardResponseData) FilterData() {
	fd := new(FilterData)
	fd.FilterPanel = make(map[string][]DashboardPanel)
	title := ""
	d.FilterResp = make(map[string]*FilterData)
	d.Rows = make(map[string][]string)
	for _, uid := range d.UID {
		if len(d.DashboardResponse[uid].Dashboard.Panels) == 0 {
			log.Fatal("No Panels found in Grafana Dashboard -> ", d.DashboardResponse[uid].Dashboard.Title)
		}
		for _, res := range d.DashboardResponse[uid].Dashboard.Panels {
			if title == "" {
				title = res.Title
			}
			if res.Type == "row" {
				d.Rows[uid] = append(d.Rows[uid], res.Title)
				title = res.Title
			}
			if res.Type != "row" && res.Type != "text" {
				fd.FilterPanel[title] = append(fd.FilterPanel[title], res)
			}
		}
		d.FilterResp[uid] = fd
	}
}
