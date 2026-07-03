// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/write_policy.go
// Original timestamp: 2026/06/22 06:08:08

package policies

import (
	"fmt"
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	vpol "github.com/jeanfrancoisgratton/vaultlib/v2/policies"
)

// WritePolicy reads a Vault ACL policy from a JSON or HCL file, validates
// it, and writes it to Vault under policyName.  The correct parser is
// selected automatically based on the file extension (.hcl → HCL, anything
// else → JSON).
//
// policyName is the name under which the policy will be stored in Vault.
// policyFile is the path to a .json or .hcl file describing the policy.
func WritePolicy(policyName, policyFile string) *ce.CustomError {
	// Check for required globals
	if err := shared.SetVaultToken(); err != nil {
		return err
	}
	if err := shared.SetServerAddress(); err != nil {
		return err
	}

	// Dispatch to JSON or HCL parser based on file extension; the returned
	// string is already validated and ready for Vault's API.
	rules, parseErr := ParsePolicyFile(policyFile)
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

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign("Created policy "+hftx.Green(policyName)) + " from " + policyFile)
	}
	return nil
}
