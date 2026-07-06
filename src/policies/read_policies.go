// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/read_policies.go
// Original timestamp: 2026/06/21 13:05:27

package policies

import (
	"encoding/json"
	"fmt"

	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hfjson "github.com/jeanfrancoisgratton/helperFunctions/v5/prettyjson"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	vpol "github.com/jeanfrancoisgratton/vaultlib/v2/policies"
)

func (c *Client) Read(pname string, showOutput bool) (*vpol.Policy, *ce.CustomError) {
	policy, err := c.vc.ReadPolicy(pname)
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
		fmt.Printf("\nPolicy: %s\n\n", hftx.Green(policy.Name))
		fmt.Println(policy.Rules)
	}

	return policy, nil
}
