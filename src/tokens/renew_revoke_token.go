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
)

// Revoke: revoke (delete) a token. If the token is a parent token, all of its children will be deleted as well
func (c *Client) Revoke(tokenName string) *ce.CustomError {
	if e := c.vc.RevokeToken(tokenName); e != nil {
		return &ce.CustomError{Title: "Error revoking token", Message: e.Error()}
	}
	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign(hftx.Green("Revoked token")), tokenName)
	}

	return nil
}

// Renew: renew the given token. If the -d flag is unset, a lease will be set to 1h
func (c *Client) Renew(tokenName string) *ce.CustomError {
	if _, e := c.vc.RenewToken(tokenName, TokenDuration); e != nil {
		return &ce.CustomError{Title: "Error renewing token", Message: e.Error()}
	}
	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign(hftx.Green("Renewed token")), tokenName)
	}

	return nil
}
