// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/tokens/tist_accessors.go
// Original timestamp: 2026/07/04 21:21:36

package tokens

import (
	"fmt"
	"log"
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	tkn "github.com/jeanfrancoisgratton/vaultlib/v2/tokens"
)

func ListAccessors() *ce.CustomError {
	// Check for required globals
	if err := shared.SetVaultToken(); err != nil {
		return err
	}
	if err := shared.SetServerAddress(); err != nil {
		return err
	}

	cfg := tkn.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken}
	client, cvlrErr := tkn.NewClient(cfg)
	if cvlrErr != nil {
		return &ce.CustomError{Title: "Error creating vault client", Message: cvlrErr.Error()}
	}

	accessors, err := client.ListAccessors()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hftx.Green("Current accessors:"))
	for _, accessor := range accessors {
		fmt.Println(accessor)
	}
	return nil
}
