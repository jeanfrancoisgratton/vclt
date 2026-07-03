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

// ReadSecrets: Reads a secret from the Vault secret path

func ReadSecrets(kvengine, path string) *ce.CustomError {
	// Check for required globals
	if err := shared.SetVaultToken(); err != nil {
		return err
	}
	if err := shared.SetServerAddress(); err != nil {
		return err
	}

	cfg := vlr.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken, MountPath: kvengine}
	client, cvlrErr := vlr.NewClient(cfg)
	if cvlrErr != nil {
		return &ce.CustomError{Title: "Error creating vault client", Message: cvlrErr.Error()}
	}

	// -f is empty; this means we grab the whole secret
	if SecretField == "" {
		return allSecrets(client, path)
	}

	return singleFieldFromSecret(client, path)
}

// We fetch all the fields of a given secret, optionally rendering it in JSON
func allSecrets(c *vlr.Client, path string) *ce.CustomError {
	var secret *vlr.Secret
	var sErr error

	if secret, sErr = c.ReadSecret(path,
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

func singleFieldFromSecret(c *vlr.Client, path string) *ce.CustomError {
	value, err := c.ReadSecretField(path, SecretField, SecretVersion)
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
