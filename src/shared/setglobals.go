// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/shared/setglobals.go
// Original timestamp: 2026/06/18 09:05:07

package shared

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// Sets the server address url in a global (shared) variable

func SetServerAddress() *ce.CustomError {
	if VaultServerAddress == "" {
		VaultServerAddress = os.Getenv("VAULT_ADDR")
	}
	if VaultServerAddress == "" {
		return &ce.CustomError{Title: "Vault address is missing",
			Message: "Neither the $VAULT_ADDR variable or the -a flag were set",
			Code:    ErrVaultServerAddressMissing}
	}
	return nil
}

// Sets the Vault token in a global (shared) variable

func SetVaultToken() *ce.CustomError {
	// We check if we have a valid token value, if not we exit
	if VaultAuthToken == "" {
		VaultAuthToken = os.Getenv("VAULT_TOKEN")
	}
	// the -t flag and VAULT_TOKEN env var are not set, we check if there is a $HOME/.vault-token file
	if VaultAuthToken == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			data, err := os.ReadFile(filepath.Join(homeDir, ".vault-token"))
			if err == nil {
				VaultAuthToken = strings.TrimSpace(string(data))
			}
		}
	}
	// -t and VAULT_TOKEN are not set, $HOME/.vault-token is missing or empty
	if VaultAuthToken == "" {
		errInfo := ErrorMessages[ErrVaultAuthTokenMissing]
		title := fmt.Sprintf("[%s] Vault token is missing", errInfo.Int2StringCode)
		message := fmt.Sprintf("Neither the $VAULT_TOKEN variable, the -t flag or the ~%s/.vault-token file were set.",
			filepath.Base(os.Getenv("HOME")))
		return &ce.CustomError{Title: title, Message: message, Code: ErrVaultAuthTokenMissing}
	}
	return nil
}
