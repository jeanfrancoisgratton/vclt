// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/sys/enable_disable_kvengine.go
// Original timestamp: 2026/06/25 20:47:29

package sys

import (
	ce "github.com/jeanfrancoisgratton/customError/v3"
	"github.com/jeanfrancoisgratton/vaultLib/sys"

	"vclt/shared"
)

func EnableKVengine(kvEngine string) *ce.CustomError {
	// Check for required globals
	if err := shared.SetVaultToken(); err != nil {
		return err
	}
	if err := shared.SetServerAddress(); err != nil {
		return err
	}

	cfg := sys.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken}
	client, cvlrErr := sys.NewClient(cfg)
	if cvlrErr != nil {
		return &ce.CustomError{Title: "Error creating vault client", Message: cvlrErr.Error()}
	}

	err := client.EnableKVEngine(kvEngine, sys.EnableKVOptions{
		Version: KVEngineVersion, Description: KVEngineDescription})
	if err != nil {
		return &ce.CustomError{Title: "Error listing kv", Message: err.Error()}
	}

	return nil
}
