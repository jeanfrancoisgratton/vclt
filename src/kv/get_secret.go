// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/get_secret.go
// Original timestamp: 2026/06/14 13:13:52

package kv

import (
	"encoding/json"
	"fmt"

	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hfjson "github.com/jeanfrancoisgratton/helperFunctions/v5/prettyjson"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	vlr "github.com/jeanfrancoisgratton/vaultlib/v2/kv"
)

// Read reads a secret from the Vault secret path.
func (c *Client) Read(path string) *ce.CustomError {
	// -f is empty; this means we grab the whole secret
	if SecretField == "" {
		return c.allSecrets(path)
	}

	return c.singleFieldFromSecret(path)
}

// We fetch all the fields of a given secret, optionally rendering it in JSON
func (c *Client) allSecrets(path string) *ce.CustomError {
	var secret *vlr.Secret
	var sErr error

	if secret, sErr = c.vc.ReadSecret(path,
		vlr.ReadOptions{Version: SecretVersion, FallbackToLatestAvailable: true}); sErr != nil {
		return &ce.CustomError{Title: "Error reading secret", Message: sErr.Error()}
	}

	if shared.OutputFormat == "json" {
		payload, err := json.MarshalIndent(secret.Data, "", "  ")
		if err != nil {
			return &ce.CustomError{Title: "Error serializing secret", Message: err.Error()}
		}
		if e := hfjson.Print(payload); e != nil {
			return &ce.CustomError{Title: "Unable to render secret's payload", Message: e.Error()}
		}
		return nil
	}
	return outputData(secret.Data, shared.QuietOutput)
}

func (c *Client) singleFieldFromSecret(path string) *ce.CustomError {
	value, err := c.vc.ReadSecretField(path, SecretField, SecretVersion)
	if err != nil {
		return &ce.CustomError{Title: "Error reading secret", Message: err.Error()}
	}

	if shared.QuietOutput {
		fmt.Printf("%v\n", value)
	} else {
		fmt.Printf("%s : %v\n", hftx.Green(SecretField), value)
	}
	return nil
}
