// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/admin/list_mounts.go
// Original timestamp: 2026/06/21 19:08:12

package sys

import (
	"os"

	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"github.com/jeanfrancoisgratton/vaultlib/v2/sys"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func ListMounts(showOutput bool) ([]sys.MountInfo, *ce.CustomError) {
	if err := shared.SetVaultToken(); err != nil {
		return nil, err
	}
	if err := shared.SetServerAddress(); err != nil {
		return nil, err
	}

	cfg := sys.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken}

	mounts, err := sys.ListMounts(cfg)
	if err != nil {
		return nil, &ce.CustomError{Title: "Unable to list mounts", Message: err.Error()}
	}

	if !showOutput {
		return mounts, nil
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Path", "Type", "Version", "Description"})
	for _, mount := range mounts {
		t.AppendRow(table.Row{mount.Path, mount.Type, mount.KVVersion, mount.Description})
	}
	t.SortBy([]table.SortBy{
		{Name: "Path", Mode: table.Asc},
	})
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault

	t.Render()

	return mounts, nil
}
