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
	vkvlib "github.com/jeanfrancoisgratton/vaultLib/kv"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"vclt/shared"
)

func ListSecrets(kvEngine string, displayOutput bool) ([]vkvlib.SecretInfo, *ce.CustomError) {
	// Check for required globals
	if err := shared.SetVaultToken(); err != nil {
		return nil, err
	}
	if err := shared.SetServerAddress(); err != nil {
		return nil, err
	}

	cfg := vkvlib.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken, MountPath: kvEngine}
	client, cvlrErr := vkvlib.NewClient(cfg)
	if cvlrErr != nil {
		return nil, &ce.CustomError{Title: "Error creating vault client", Message: cvlrErr.Error()}
	}

	secretslist, err := client.ListSecrets(ExtendedSecretsList)
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
		{Name: "Name", Mode: table.Asc},
	})
	//t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
	t.Style().Format.Header = text.FormatDefault

	t.Render()
	return secretslist, nil
}
