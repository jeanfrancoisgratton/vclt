// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/write_secret.go
// Original timestamp: 2026/06/16 06:49:44

package kv

import (
	ce "github.com/jeanfrancoisgratton/customError/v3"
	vlr "github.com/jeanfrancoisgratton/vaultlib/v2/kv"
)

func (c *Client) Write(path, key, value string) (*vlr.WriteResult, *ce.CustomError) {
	if wRes, kvWriteError := c.vc.WriteSecretField(path, key, value); kvWriteError != nil {
		return wRes, &ce.CustomError{Title: "Error writing secret", Message: kvWriteError.Error()}
	} else {
		return wRes, nil
	}
}
