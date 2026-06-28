// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/backup_engine.go
// Original timestamp: 2026/06/27 19:42:26

package kv

import (
	"fmt"
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	vlr "github.com/jeanfrancoisgratton/vaultLib/kv"
)

func BackupEngine(kvengine, path string) *ce.CustomError {
	// Check for required globals
	if err := shared.SetVaultToken(); err != nil {
		return err
	}
	if err := shared.SetServerAddress(); err != nil {
		return err
	}

	cfg := vlr.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken, MountPath: kvengine}
	c, cvlrErr := vlr.NewClient(cfg)
	if cvlrErr != nil {
		return &ce.CustomError{Title: "Error creating vault client", Message: cvlrErr.Error()}
	}

	if beErr := c.BackupEngine(path); beErr != nil {
		return &ce.CustomError{Title: "Error dumping secrets", Message: beErr.Error()}
	}

	if !shared.QuietOutput {
		fmt.Printf("%s %s to %s", hftx.EnabledSign("Succesfully dumped"),
			hftx.Green(kvengine), hftx.Green(path))
	}
	return nil
}
