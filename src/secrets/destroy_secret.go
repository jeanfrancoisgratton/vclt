// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/secrets/destroy_secret.go
// Original timestamp: 2026/06/17 21:32:59

package secrets

import (
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	vlr "github.com/jeanfrancoisgratton/vaultLib/kv"
	"vclt/shared"
)

func DestroySecret(kvengine, path string) *ce.CustomError {
	// Check for required globals
	if err := setGlobals(); err != nil {
		return err
	}

	cfg := vlr.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken, MountPath: kvengine}
	c, cvlrErr := vlr.NewClient(cfg)
	if cvlrErr != nil {
		return &ce.CustomError{Title: "Error creating vault client", Message: cvlrErr.Error()}
	}

	if e := c.DestroySecret(path, vlr.DestroyOptions{Version: SecretVersion}); e != nil {
		return &ce.CustomError{Title: "Error destroying the secret", Message: e.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign("Destroyed " + hftx.Green(path)))
	}
	return nil
}
