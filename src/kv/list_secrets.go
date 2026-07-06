// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/list_secrets.go
// Original timestamp: 2026/06/15 12:43:29

package kv

import (
	"os"
	"strconv"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	vkvlib "github.com/jeanfrancoisgratton/vaultlib/v2/kv"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (c *Client) List(displayOutput bool) ([]vkvlib.SecretInfo, *ce.CustomError) {
	secretslist, err := c.vc.ListSecrets(ExtendedSecretsList)
	if err != nil {
		return nil, &ce.CustomError{Title: "Error listing kv", Message: err.Error()}
	}

	if !displayOutput {
		return secretslist, nil
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Secret path", "Version"})
	for _, secret := range secretslist {
		s := hftx.Yellow("n/a")
		if ExtendedSecretsList {
			s = strconv.Itoa(secret.Version)
		}
		t.AppendRow(table.Row{secret.Path, s})
	}
	t.SortBy([]table.SortBy{
		{Name: "Secret path", Mode: table.Asc},
	})
	//t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
	t.Style().Format.Header = text.FormatDefault

	t.Render()
	return secretslist, nil
}
