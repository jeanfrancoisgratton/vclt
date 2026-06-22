// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/write_policy.go
// Original timestamp: 2026/06/22 06:08:08

package policies

import (
	ce "github.com/jeanfrancoisgratton/customError/v3"
	vpol "github.com/jeanfrancoisgratton/vaultLib/policies"
	"vclt/shared"
)

// WritePolicy reads a Vault ACL policy from a JSON file, validates it, and
// writes it to Vault under policyName.
//
// policyName is the name under which the policy will be stored in Vault.
// policyFile is the path to a JSON file describing the policy (Vault ACL schema).
func WritePolicy(policyName, policyFile string) *ce.CustomError {
	// Check for required globals
	if err := shared.SetVaultToken(); err != nil {
		return err
	}
	if err := shared.SetServerAddress(); err != nil {
		return err
	}

	// Parse and validate the JSON policy file; ParseFile returns the
	// canonical re-marshaled JSON string ready for Vault's API.
	rules, parseErr := ParseFile(policyFile)
	if parseErr != nil {
		return &ce.CustomError{Title: "Invalid policy file", Message: parseErr.Error()}
	}

	cfg := vpol.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken}
	c, clientErr := vpol.NewClient(cfg)
	if clientErr != nil {
		return &ce.CustomError{Title: "Error creating vault client", Message: clientErr.Error()}
	}

	if err := c.CreatePolicy(policyName, rules); err != nil {
		return &ce.CustomError{Title: "Error writing policy", Message: err.Error()}
	}

	return nil
}
