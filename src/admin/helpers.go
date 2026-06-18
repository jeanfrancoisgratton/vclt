// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/admin/helpers.go
// Original timestamp: 2026/06/18 08:15:05

package admin

import (
	"os"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"vclt/shared"
)

func getAddress() *ce.CustomError {
	if shared.VaultServerAddress == "" {
		shared.VaultServerAddress = os.Getenv("VAULT_ADDR")
	}
	if shared.VaultServerAddress == "" {
		return &ce.CustomError{Title: "Vault address is missing",
			Message: "Neither the $VAULT_ADDR variable or the -a flag were set",
			Code:    shared.ErrVaultServerAddressMissing}
	}
	return nil
}
