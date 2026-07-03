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

func DeleteSecret(kvengine, path string) *ce.CustomError {
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

	// -f is empty, this means we grab the whole secret
	if SecretField != "" {
		return deleteFieldFromSecret(client, path)
	}

	return deleteWholeSecret(client, path)
}

func deleteFieldFromSecret(c *vlr.Client, path string) *ce.CustomError {
	if _, e := c.DeleteSecretField(path, SecretField); e != nil {
		return &ce.CustomError{Title: "Error deleting the field", Message: e.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign("Deleted " + hftx.Green(path+"/"+SecretField)))
	}
	return nil
}

func deleteWholeSecret(c *vlr.Client, path string) *ce.CustomError {
	if e := c.SoftDeleteSecret(path, vlr.DeleteOptions{Versions: []int{SecretVersion}}); e != nil {
		return &ce.CustomError{Title: "Error deleting the secret", Message: e.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign("Deleted " + hftx.Green(path)))
	}
	return nil
}
