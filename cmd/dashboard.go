package main

import (
	"fmt"
	"log"
	"net/url"
)

type DashboardMeta struct {
	IsStarred bool   `json:"isStarred"`
	Slug      string `json:"slug"`
	Folder    int64  `json:"folderId"`
	URL       string `json:"url"`
}
type DashboardTargets struct {
	Expr string `json:"expr"`
}
type DashboardPanel struct {
	DataSource map[string]string  `json:"datasource"`
	Targets    []DashboardTargets `json:"targets"`
	Title      string             `json:"title"`
}
type DashboardModel struct {
	Panels []DashboardPanel `json:"panels"`
	Title  string           `json:"title"`
	UID    string           `json:"uid"`
	//Annotations map[string]interface{} `json:"annotations"`
}
type Dashboard struct {
	Meta      DashboardMeta  `json:"meta"`
	Model     DashboardModel `json:"dashboard"`
	FolderID  int64          `json:"folderId"`
	FolderUID string         `json:"folderUid"`
	OverWrite bool           `json:"overwrite"`
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
type DashboardData struct {
	Title   string
	Panels  []string
	Metrics map[string]*MetricOutput
}

type DashboardResponseData struct {
	UID               []string
	DashboardResponse []*Dashboard
	Model             []*DashboardModel
	Data              []*DashboardData
}

func (c *Client) dashboard(path string) (*Dashboard, error) {
	result := &Dashboard{}
	err := c.request("GET", path, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	result.FolderID = result.Meta.Folder
	return result, err
}
func (c *Client) DashboardByUID(uid string) (*Dashboard, error) {
	return c.dashboard(fmt.Sprintf("/api/dashboards/uid/%s", uid))

}

// FolderDashboardSearch uses the folder and dashboard search endpoint to find dashboards based on the params passed in.
func (c *Client) FolderDashboardSearch(params url.Values) (resp []FolderDashboardSearchResponse, err error) {
	err = c.request("GET", "/api/search", params, nil, &resp)
	return
}

// GetDashboards will fetch the uid of dashboards which matches the search
func (d *DashboardResponseData) GetDashboards(c Client, dashboards []string) {
	for _, dashboard := range dashboards {
		query := make(url.Values)
		query.Set("query", dashboard)
		resp, err := c.FolderDashboardSearch(query)
		if err != nil {
			log.Fatal(err)
		}
		for _, val := range resp {
			d.UID = append(d.UID, val.UID)

		}
	}

}

// GetDashboardByUID will fetch the dashboards by uid from grafana
func (d *DashboardResponseData) GetDashboardByUID(c Client) {
	for _, val := range d.UID {
		res, err := c.DashboardByUID(val)
		if err != nil {
			log.Fatal(err)
		}
		//result.Model.Panels)
		d.DashboardResponse = append(d.DashboardResponse, res)
	}
}

// GetDashboardModelFromResponse will fetch the Dashboard Model from the response
func (d *DashboardResponseData) GetDashboardModelFromResponse(c Client) {
	for _, val := range d.DashboardResponse {
		d.Model = append(d.Model, &val.Model)
	}
}

// GetDashboardMetricsFromResponse will fetch the necessary data from the model object. It contains Panel title, metrics etc
func (d *DashboardResponseData) GetDashboardMetricsFromResponse(c Client) {
	if len(d.Model) == 0 {
		log.Fatal("No Dashboard response returned")
	}
	for i := 0; i < len(d.Model); i++ {
		Panels := d.Model[i].Panels
		dd := new(DashboardData)
		dd.Title = d.Model[i].Title
		dd.Metrics = make(map[string]*MetricOutput)
		if len(Panels) == 0 {
			log.Fatal("No Dashboard panels returned")
		}
		for j := 0; j < len(Panels); j++ {
			Targets := Panels[j].Targets
			if len(Targets) == 0 {
				log.Fatal("No Dashboard targets returned")
			}
			dd.Panels = append(dd.Panels, Panels[j].Title)
			mo := new(MetricOutput)
			for k := 0; k < len(Targets); k++ {
				mo.MetricName = Targets[k].Expr
				dd.Metrics[Panels[j].Title] = mo
			}
		}
		d.Data = append(d.Data, dd)
	}
}
