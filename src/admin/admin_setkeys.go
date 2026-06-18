// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/admin/admin_setkeys.go
// Original timestamp: 2026/06/18 06:27:56

package admin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	"github.com/jeanfrancoisgratton/vaultLib/admin"
	"vclt/shared"
)

func SetRootKeys() *ce.CustomError {
	vk := VaultRootKeysStruct{}

	if err := getAddress(); err != nil {
		return err
	}

	cfg := admin.AdminConfig{Address: shared.VaultServerAddress}

	status, err := admin.GetSealStatus(cfg)
	if err != nil {
		return &ce.CustomError{Title: "Unable to get the server's status", Message: err.Error()}
	}
	vk.MinimumRequired = status.Threshold
	fmt.Println("You will be prompted to enter the root keys. Press ENTER at the prompt when you're done")

	for {
		k := hf.GetPassword("Enter root key part, ENTER to quit : ", shared.DebugMode)
		if k != "" {
			vk.Keys = append(vk.Keys, hf.EncodeString(k, ""))
		} else {
			break
		}
	}

	if len(vk.Keys) < vk.MinimumRequired {
		return &ce.CustomError{Title: "Unable to save the root keys",
			Message: fmt.Sprintf("You have %d keys parts while the minimal number is %d", len(vk.Keys), vk.MinimumRequired)}
	}
	return saveRootKeys(vk)
}

func saveRootKeys(vk VaultRootKeysStruct) *ce.CustomError {
	jStream, err := json.MarshalIndent(vk, "", "  ")
	if err != nil {
		return &ce.CustomError{Title: err.Error()}
	}

	if err := os.WriteFile(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vclt", "rootkeys.json"), jStream, 0600); err != nil {
		return &ce.CustomError{Title: "Unable to write root keys json file", Message: err.Error()}
	}
	return nil
}
