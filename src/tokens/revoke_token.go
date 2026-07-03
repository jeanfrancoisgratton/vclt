// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/tokens/revoke_token.go
// Original timestamp: 2026/07/02 23:09:18

package tokens

import (
	"fmt"
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	tkn "github.com/jeanfrancoisgratton/vaultlib/v2/tokens"
)

func RevokeToken(tokenName string) *ce.CustomError {

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

	if e := client.RevokeToken(tokenName); e != nil {
		return &ce.CustomError{Title: "Error revoking token", Message: e.Error()}
	}
	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign(hftx.Green("Revoked token")), tokenName)
	}

	return nil
}
