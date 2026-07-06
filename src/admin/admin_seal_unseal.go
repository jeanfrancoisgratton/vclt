// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/admin/admin_seal_unseal.go
// Original timestamp: 2026/06/18 09:01:18

package admin

import (
	"fmt"

	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

func (c *Client) Seal() *ce.CustomError {
	if err := c.vc.Seal(); err != nil {
		return &ce.CustomError{Title: "Unable to seal Vault", Message: err.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign(fmt.Sprintf("Vault %s is %s", shared.VaultServerAddress,
			hftx.Green("sealed"))))
	}
	return nil
}

func (c *Client) Unseal(rkfile string) *ce.CustomError {
	var keyparts []string
	vk := VaultRootKeysStruct{}

	if err := vk.loadRootKeys(rkfile); err != nil {
		return err
	}

	for _, shard := range vk.Shards {
		keyparts = append(keyparts, hf.DecodeString(shard, ""))
	}

	if len(keyparts) < vk.MinimumRequired {
		return &ce.CustomError{Title: "Minimal number parts of root keys not met",
			Message: fmt.Sprintf("Received %d parts, expected %d", len(keyparts), vk.MinimumRequired)}
	}

	// TODO : monitor progress and handle errors more gracefully
	if _, err := c.vc.Unseal(keyparts[:vk.MinimumRequired]); err != nil {
		return &ce.CustomError{Title: "Unable to unseal Vault", Message: err.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign(fmt.Sprintf("Vault %s is %s", shared.VaultServerAddress, hftx.Green("unsealed"))))
	}
	return nil
}
