// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/restore_engine.go
// Original timestamp: 2026/06/27 19:58:26

package kv

import (
	"fmt"
	"os"
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	vlr "github.com/jeanfrancoisgratton/vaultlib/v2/kv"
)

func RestoreEngine(kvengine, path string) *ce.CustomError {
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
		return &ce.CustomError{Title: "Error restoring secrets", Message: beErr.Error()}
	}

	if !Cleartext {
		if eerr := decodefile(path); eerr != nil {
			return eerr
		}
	}

	if !shared.QuietOutput {
		fmt.Printf("%s %s to %s\n", hftx.EnabledSign("Successfully restored"),
			hftx.Green(kvengine), hftx.Green(path))
	}
	return nil
}

func decodefile(path string) *ce.CustomError {
	if renameErr := os.Rename(path, path+".enc"); renameErr != nil {
		return &ce.CustomError{Title: "Error renaming file", Message: renameErr.Error()}
	}
	if derr := hf.DecodeFile(path+".enc", path, ""); derr != nil {
		return &ce.CustomError{Title: "Error encoding file", Message: derr.Error()}
	}

	if rerr := os.Remove(path + ".enc"); rerr != nil {
		return &ce.CustomError{Title: "Error removing temp file", Message: rerr.Error()}
	}
	return nil
}
