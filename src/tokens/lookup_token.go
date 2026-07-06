// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/tokens/lookup_token.go
// Original timestamp: 2026/07/03 21:20:42

package tokens

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hfjson "github.com/jeanfrancoisgratton/helperFunctions/v5/prettyjson"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	tkn "github.com/jeanfrancoisgratton/vaultlib/v2/tokens"
)

// LookupSelf : gives info about the current token (the one the user is currently uses)
func (c *Client) LookupSelf(displayOutput bool) (*tkn.TokenInfo, *ce.CustomError) {
	self, err := c.vc.LookupSelf()
	if err != nil {
		return nil, &ce.CustomError{Title: "Error looking up token", Message: err.Error()}
	}
	if shared.OutputFormat == "json" {
		payload, jerr := json.MarshalIndent(self, "", "  ")
		if jerr != nil {
			return nil, &ce.CustomError{Title: "Error serializing payload", Message: jerr.Error()}
		}
		if e := hfjson.Print(payload); e != nil {
			return nil, &ce.CustomError{Title: "Unable to render token's payload", Message: e.Error()}
		}
	} else if displayOutput {
		displayTokenInformation("self", self)
	}
	if TokenSavefile != "" {
		if err := saveTokenInfo2file(TokenSavefile, self); err != nil {
			return self, err
		}
	}

	return self, nil
}

func (c *Client) LookupToken(tokenName string, displayoutput bool) (*tkn.TokenInfo, *ce.CustomError) {
	tok, err := c.vc.LookupToken(tokenName)
	if err != nil {
		return nil, &ce.CustomError{Title: "Error looking up token", Message: err.Error()}
	}
	if shared.OutputFormat == "json" {
		payload, jerr := json.MarshalIndent(tok, "", "  ")
		if jerr != nil {
			return nil, &ce.CustomError{Title: "Error serializing payload", Message: jerr.Error()}
		}
		if e := hfjson.Print(payload); e != nil {
			return nil, &ce.CustomError{Title: "Unable to render token's payload", Message: e.Error()}
		}
	} else if displayoutput {
		displayTokenInformation(tokenName, tok)
	}
	if TokenSavefile != "" {
		if err := saveTokenInfo2file(TokenSavefile, tok); err != nil {
			return tok, err
		}
	}

	return tok, nil
}

// small helper function to dump the token to the TTY
func displayTokenInformation(tokenName string, tkn *tkn.TokenInfo) {
	fmt.Printf("\nToken: %s\n\n", hftx.Green(tokenName))
	fmt.Println(hftx.Green("ID: "), tkn.ID)
	fmt.Println(hftx.Green("Display name: "), tkn.DisplayName)
	fmt.Println(hftx.Green("Token accessor: "), tkn.Accessor)
	fmt.Println(hftx.Green("Creation time: "), tkn.CreationTime)
	fmt.Println(hftx.Green("Expire time: "), tkn.ExpireTime)
	fmt.Println(hftx.Green("Token TTL: "), tkn.TTL)
	fmt.Println(hftx.Green("Creation TTL: "), tkn.CreationTTL)
	fmt.Println(hftx.Green("Token max TTL: "), tkn.ExplicitMaxTTL)
	fmt.Println(hftx.Green("Number of uses: "), tkn.NumUses)
	fmt.Println(hftx.Green("Orphaned token: "), tkn.Orphan)
	fmt.Println(hftx.Green("Renewable token: "), tkn.Renewable)
	fmt.Println(hftx.Green("Token path: "), tkn.Path)
	fmt.Println(hftx.Green("Token policies: "), tkn.Policies)
	fmt.Println(hftx.Green("Token type: "), tkn.Type)
	fmt.Println(hftx.Green("Token metadata: "), tkn.Meta)
	fmt.Println(hftx.Green("Entity ID: "), tkn.EntityID)
}

func saveTokenInfo2file(TokenSavefile string, tkninfo *tkn.TokenInfo) *ce.CustomError {
	jStream, err := json.MarshalIndent(tkninfo, "", "  ")
	if err != nil {
		return &ce.CustomError{Title: err.Error()}
	}

	if !strings.HasSuffix(strings.ToLower(TokenSavefile), ".json") {
		TokenSavefile += ".json"
	}
	if err := os.WriteFile(TokenSavefile, jStream, 0600); err != nil {
		return &ce.CustomError{Title: "Unable to save the json file", Message: err.Error()}
	}
	return nil
}
