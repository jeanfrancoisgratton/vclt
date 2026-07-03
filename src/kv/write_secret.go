// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/write_secret.go
// Original timestamp: 2026/06/16 06:49:44

package kv

import (
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	vlr "github.com/jeanfrancoisgratton/vaultlib/v2/kv"
)

func WriteSecrets(kvengine, path, key, value string) (*vlr.WriteResult, *ce.CustomError) {
	// Check for required globals
	if err := shared.SetVaultToken(); err != nil {
		return nil, err
	}
	if err := shared.SetServerAddress(); err != nil {
		return nil, err
	}

	cfg := vlr.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken, MountPath: kvengine}
	client, cvlrErr := vlr.NewClient(cfg)
	if cvlrErr != nil {
		return nil, &ce.CustomError{Title: "Error creating vault client", Message: cvlrErr.Error()}
	}

	if wRes, kvWriteError := client.WriteSecretField(path, key, value); kvWriteError != nil {
		return wRes, &ce.CustomError{Title: "Error writing secret", Message: kvWriteError.Error()}
	} else {
		return wRes, nil
	}
}
