// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/delete_secret.go
// Original timestamp: 2026/06/17 20:01:17

package kv

import (
	"fmt"

	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	vlr "github.com/jeanfrancoisgratton/vaultlib/v2/kv"
)

func (c *Client) Delete(path string) *ce.CustomError {
	// -f is empty, this means we grab the whole secret
	if SecretField != "" {
		return c.deleteFieldFromSecret(path)
	}

	return c.deleteWholeSecret(path)
}

func (c *Client) deleteFieldFromSecret(path string) *ce.CustomError {
	if _, e := c.vc.DeleteSecretField(path, SecretField); e != nil {
		return &ce.CustomError{Title: "Error deleting the field", Message: e.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign("Deleted " + hftx.Green(path+"/"+SecretField)))
	}
	return nil
}

func (c *Client) deleteWholeSecret(path string) *ce.CustomError {
	if e := c.vc.SoftDeleteSecret(path, vlr.DeleteOptions{Versions: []int{SecretVersion}}); e != nil {
		return &ce.CustomError{Title: "Error deleting the secret", Message: e.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign("Deleted " + hftx.Green(path)))
	}
	return nil
}
