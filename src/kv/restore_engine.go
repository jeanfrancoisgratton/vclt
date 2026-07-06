// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/restore_engine.go
// Original timestamp: 2026/06/27 19:58:26

package kv

import (
	"fmt"
	"os"
	"path/filepath"

	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

func (c *Client) Restore(path string) *ce.CustomError {
	// An encoded backup must be decoded before it can be fed to the restore
	// endpoint. We decode into a temp file and leave the original untouched.
	restorePath := path
	if !Cleartext {
		decoded, derr := decode2temp(path)
		if derr != nil {
			return derr
		}
		defer os.Remove(decoded)
		restorePath = decoded
	}

	if reErr := c.vc.RestoreEngine(restorePath); reErr != nil {
		return &ce.CustomError{Title: "Error restoring secrets", Message: reErr.Error()}
	}

	if !shared.QuietOutput {
		fmt.Printf("%s %s from %s\n", hftx.EnabledSign("Successfully restored"),
			hftx.Green(c.engine), hftx.Green(path))
	}
	return nil
}

func decode2temp(path string) (string, *ce.CustomError) {
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".dec-*")
	if err != nil {
		return "", &ce.CustomError{Title: "Error creating temp file", Message: err.Error()}
	}
	tmpPath := tmp.Name()
	tmp.Close()

	if derr := hf.DecodeFile(path, tmpPath, ""); derr != nil {
		os.Remove(tmpPath)
		return "", &ce.CustomError{Title: "Error decoding backup file", Message: derr.Error()}
	}
	return tmpPath, nil
}
