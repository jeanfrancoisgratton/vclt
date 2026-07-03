// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/delete_policies.go
// Original timestamp: 2026/06/21 17:29:34

package policies

import (
	"fmt"

	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftfx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	vpol "github.com/jeanfrancoisgratton/vaultlib/v2/policies"
)

func DeletePolicies(policies []string) *ce.CustomError {
	// Check for required globals
	if err := shared.SetVaultToken(); err != nil {
		return err
	}
	if err := shared.SetServerAddress(); err != nil {
		return err
	}

	cfg := vpol.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken}
	c, cvlrErr := vpol.NewClient(cfg)
	if cvlrErr != nil {
		return &ce.CustomError{Title: "Error creating vault client", Message: cvlrErr.Error()}
	}

	for _, policy := range policies {
		if err := c.DeletePolicy(policy); err != nil {
			return &ce.CustomError{Title: "Error deleting policy", Message: err.Error()}
		}
		if !shared.QuietOutput {
			fmt.Println(hftfx.EnabledSign("Deleted policy " + hftfx.Green(policy)))
		}
	}
	return nil
}
