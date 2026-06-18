// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/admin/admin_seal_unseal.go
// Original timestamp: 2026/06/18 09:01:18

package admin

import (
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	vadm "github.com/jeanfrancoisgratton/vaultLib/admin"
	"vclt/shared"
)

func Seal() *ce.CustomError {
	if err := shared.SetVaultToken(); err != nil {
		return err
	}
	if err := shared.SetServerAddress(); err != nil {
		return err
	}

	cfg := vadm.AdminConfig{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken}

	c, err := vadm.NewClient(cfg)
	if err != nil {
		return &ce.CustomError{Title: "Unable to create Vault client", Message: err.Error()}
	}

	if err := c.Seal(); err != nil {
		return &ce.CustomError{Title: "Unable to seal Vault", Message: err.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign(fmt.Sprintf("Vault %s is %s", shared.VaultServerAddress,
			hftx.Green("sealed"))))
	}
	return nil
}

func Unseal(rkfile string) *ce.CustomError {
	var rootkeys []string
	vk := VaultRootKeysStruct{}
	if err := shared.SetServerAddress(); err != nil {
		return err
	}

	cfg := vadm.AdminConfig{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken}

	if err := vk.loadRootKeys(rkfile); err != nil {
		return err
	}

	for _, key := range vk.Keys {
		rootkeys = append(rootkeys, hf.DecodeString(key, ""))
	}

	if len(rootkeys) < vk.MinimumRequired {
		return &ce.CustomError{Title: "Minimal number parts of root keys not met",
			Message: fmt.Sprintf("Received %d parts, expected %d", len(rootkeys), vk.MinimumRequired)}
	}

	c, err := vadm.NewClient(cfg)
	if err != nil {
		return &ce.CustomError{Title: "Unable to create Vault client", Message: err.Error()}
	}

	// TODO : monitor progress and handle errors more gracefully
	if _, err := c.Unseal(rootkeys[:vk.MinimumRequired]); err != nil {
		return &ce.CustomError{Title: "Unable to unseal Vault", Message: err.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign(fmt.Sprintf("Vault %s is %s", shared.VaultServerAddress, hftx.Green("unsealed"))))
	}
	return nil
}
