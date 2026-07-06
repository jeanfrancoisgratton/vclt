// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/list_policies.go
// Original timestamp: 2026/06/20 20:14:33

package policies

import (
	"os"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (c *Client) List(showOutput bool) ([]string, *ce.CustomError) {
	polList, err := c.vc.ListPolicies()
	if err != nil {
		return nil, &ce.CustomError{Title: "Error listing policies", Message: err.Error()}
	}

	if !showOutput {
		return polList, nil
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Policy name"})
	for _, pol := range polList {
		t.AppendRow(table.Row{pol})
	}
	t.SortBy([]table.SortBy{
		{Name: "Policy name", Mode: table.Asc},
	})
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault

	t.Render()
	return polList, nil
}
