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
)

func (c *Client) Delete(names []string) *ce.CustomError {
	for _, policy := range names {
		if err := c.vc.DeletePolicy(policy); err != nil {
			return &ce.CustomError{Title: "Error deleting policy", Message: err.Error()}
		}
		if !shared.QuietOutput {
			fmt.Println(hftfx.EnabledSign("Deleted policy " + hftfx.Green(policy)))
		}
	}
	return nil
}
