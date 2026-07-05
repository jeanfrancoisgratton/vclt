// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/tokens/renew_revoke_token.go
// Original timestamp: 2026/07/02 23:09:18

package tokens

import (
	"fmt"
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	tkn "github.com/jeanfrancoisgratton/vaultlib/v2/tokens"
)

// RevokeToken: revoke (delete) a token. If the token is a parent token, all of its children will be deleted as well

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

// RenewToken: renew the given token. If the -d flag is unset, a lease will be set to 1h

func RenewToken(tokenName string) *ce.CustomError {
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

	if _, e := client.RenewToken(tokenName, TokenDuration); e != nil {
		return &ce.CustomError{Title: "Error renewing token", Message: e.Error()}
	}
	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign(hftx.Green("Renewed token")), tokenName)
	}

	return nil
}
