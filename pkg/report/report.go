package report

import (
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/vprashar2929/graftool/pkg/dashboard"
	"github.com/vprashar2929/graftool/pkg/parse"
)

func DisplayReport(d *dashboard.DashboardResponseData, grafanaBaseURL *string, startTime int64) {
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
						t.AppendRow(table.Row{row, panel.Title, target.Legends, parse.ParseEpoch(d.FilterResp[uid].Metric[target.Expr][0].Value[0]), d.FilterResp[uid].Metric[target.Expr][0].Value[1]})
					} else {
						t.AppendRow(table.Row{row, panel.Title, target.Legends, d.FilterResp[uid].Metric[target.Expr], d.FilterResp[uid].Metric[target.Expr]})
					}
				}
			}
			t.AppendSeparator()
		}
		t.SetCaption(fmt.Sprint("Dashboard Link: http://%s%s?from=%d&to=%d"), *grafanaBaseURL, d.URL[uid], startTime, time.Now().UnixMilli())
		t.Render()
		t.ResetRows()
	}
}
