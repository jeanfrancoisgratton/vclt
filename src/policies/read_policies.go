// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/read_policies.go
// Original timestamp: 2026/06/21 13:05:27

package policies

import (
	"encoding/json"
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hfjson "github.com/jeanfrancoisgratton/helperFunctions/v5/prettyjson"
	vpol "github.com/jeanfrancoisgratton/vaultLib/policies"
	"vclt/shared"
)

func ReadPolicy(pname string, showOutput bool) (*vpol.Policy, *ce.CustomError) {
	// Check for required globals
	if err := shared.SetVaultToken(); err != nil {
		return nil, err
	}
	if err := shared.SetServerAddress(); err != nil {
		return nil, err
	}

	cfg := vpol.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken}
	c, cvlrErr := vpol.NewClient(cfg)
	if cvlrErr != nil {
		return nil, &ce.CustomError{Title: "Error creating vault client", Message: cvlrErr.Error()}
	}

	policy, err := c.ReadPolicy(pname)
	if err != nil {
		return nil, &ce.CustomError{Title: "Error reading the policy", Message: err.Error()}
	}

	if shared.OutputFormat == "json" {
		payload, jerr := json.MarshalIndent(policy.RuleLines(), "", "  ")
		if jerr != nil {
			return nil, &ce.CustomError{Title: "Error serializing secret", Message: jerr.Error()}
		}
		if e := hfjson.Print(payload); e != nil {
			return nil, &ce.CustomError{Title: "Unable to render secret's payload", Message: e.Error()}
		}
	} else {
		fmt.Println(policy.Rules)
	}
	//fmt.Printf("\nPolicy: %s\n\n", hftx.Green(policy.Name))
	//fmt.Println(policy.Rules)

	return policy, nil
}
