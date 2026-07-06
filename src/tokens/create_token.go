// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/tokens/create_token.go
// Original timestamp: 2026/06/30 11:22:39

package tokens

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	tkn "github.com/jeanfrancoisgratton/vaultlib/v2/tokens"
)

func (c *Client) Create(tknName string, displayOutput bool) *ce.CustomError {
	// Validate the TTL before doing anything else
	if err := validateTTL(TokenTTL); err != nil {
		return err
	}

	boundPolicies := splitPolicies(TokenPolicies)
	if t, e := c.vc.CreateToken(tkn.CreateOptions{DisplayName: tknName, TTL: TokenTTL, Orphan: TokenOrphaned,
		Policies: boundPolicies, Renewable: &TokenRenewable}); e != nil {
		return &ce.CustomError{Title: "Error creating token", Message: e.Error()}
	} else {
		if displayOutput {
			fmt.Println(hftx.EnabledSign("Created token "), hftx.Green(tknName))
			fmt.Println()
			fmt.Println("Token: ", hftx.Green(t.ClientToken))
			fmt.Println("Accessor: ", hftx.Green(t.Accessor))
			fmt.Println("Token policies: ", hftx.Green(strings.Join(t.TokenPolicies, ",")))
			fmt.Println("Effective policies: ", hftx.Green(strings.Join(t.TokenPolicies, ",")))
			fmt.Println("Token TTL: ", hftx.Green(fmt.Sprintf("%d", t.LeaseDuration)))
			fmt.Println("Token is renewable: ", hftx.Green(fmt.Sprintf("%t", t.Renewable)))
			fmt.Println("Token is orphaned: ", hftx.Green(fmt.Sprintf("%t", t.Orphan)))
		}
		if TokenSavefile != "" {
			return saveTokenInfoToFile(TokenSavefile, t)
		}
	}

	return nil
}

// validateTTL ensures that the value passed to -t is a valid, non-negative token TTL.
// A TTL is accepted if it is either a Go duration string (e.g. "1h", "30m", "24h")
// or a bare integer number of seconds (e.g. "3600"), as Vault accepts both forms.
func validateTTL(ttl string) *ce.CustomError {
	ttl = strings.TrimSpace(ttl)
	if ttl == "" {
		return &ce.CustomError{Title: "Invalid TTL", Message: "TTL cannot be empty"}
	}

	// A bare integer is interpreted by Vault as a number of seconds.
	if secs, err := strconv.Atoi(ttl); err == nil {
		if secs < 0 {
			return &ce.CustomError{Title: "Invalid TTL", Message: fmt.Sprintf("TTL cannot be negative: %q", ttl)}
		}
		return nil
	}

	d, err := time.ParseDuration(ttl)
	if err != nil {
		return &ce.CustomError{Title: "Invalid TTL", Message: fmt.Sprintf("%q is not a valid duration (expected values such as \"30m\", \"1h\", \"24h\" or a number of seconds)", ttl)}
	}
	if d < 0 {
		return &ce.CustomError{Title: "Invalid TTL", Message: fmt.Sprintf("TTL cannot be negative: %q", ttl)}
	}

	return nil
}

func splitPolicies(policies string) []string {
	if policies == "" {
		return []string{"default"}
	}
	return strings.Split(policies, ",")
}

func saveTokenInfoToFile(TokenSavefile string, tkninfo *tkn.TokenAuth) *ce.CustomError {
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
