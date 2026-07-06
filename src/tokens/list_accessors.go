// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/tokens/list_accessors.go
// Original timestamp: 2026/07/04 21:21:36

package tokens

import (
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

func (c *Client) ListAccessors() *ce.CustomError {
	accessors, err := c.vc.ListAccessors()
	if err != nil {
		return &ce.CustomError{Title: "Error listing token accessors", Message: err.Error()}
	}

	fmt.Println(hftx.Green("Current accessors:"))
	for _, accessor := range accessors {
		fmt.Println(accessor)
	}
	return nil
}
