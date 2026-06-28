// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/backup_engine.go
// Original timestamp: 2026/06/27 19:42:26

package kv

import (
	"fmt"
	"os"
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
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

	if !Cleartext {
		if eerr := encodefile(path); eerr != nil {
			return eerr
		}
	}

	if !shared.QuietOutput {
		fmt.Printf("%s %s to %s\n", hftx.EnabledSign("Successfully dumped"),
			hftx.Green(kvengine), hftx.Green(path))
	}
	return nil
}

func encodefile(path string) *ce.CustomError {
	if renameErr := os.Rename(path, path+".enc"); renameErr != nil {
		return &ce.CustomError{Title: "Error renaming file", Message: renameErr.Error()}
	}
	if eerr := hf.EncodeFile(path+".enc", path, ""); eerr != nil {
		return &ce.CustomError{Title: "Error encoding temp file", Message: eerr.Error()}
	}

	if rerr := os.Remove(path + ".enc"); rerr != nil {
		return &ce.CustomError{Title: "Error removing temp file", Message: rerr.Error()}
	}
	return nil
}
