// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/admin/admin_setkeys.go
// Original timestamp: 2026/06/18 06:27:56

package admin

import (
	"fmt"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	"github.com/jeanfrancoisgratton/vaultLib/admin"
	"vclt/shared"
)

func SetRootKeys(rkfile string) *ce.CustomError {
	if err := shared.SetVaultToken(); err != nil {
		return err
	}
	if err := shared.SetServerAddress(); err != nil {
		return err
	}

	if !strings.HasSuffix(rkfile, ".json") {
		rkfile += ".json"
	}

	vk := VaultRootKeysStruct{}
	cfg := admin.AdminConfig{Address: shared.VaultServerAddress}

	if !OfflineMode {
		status, err := admin.GetSealStatus(cfg)
		if err != nil {
			return &ce.CustomError{Title: "Unable to get the server's status", Message: err.Error()}
		}
		vk.MinimumRequired = status.Threshold
	} else {
		vk.MinimumRequired = 0
	}

	fmt.Println("You will be prompted to enter the root keys. Press ENTER at the prompt when you're done")

	for {
		k := hf.GetPassword("Enter root key part, ENTER to quit : ", shared.DebugMode)
		if k != "" {
			vk.Shards = append(vk.Shards, hf.EncodeString(k, ""))
		} else {
			break
		}
	}

	if len(vk.Shards) < vk.MinimumRequired {
		return &ce.CustomError{Title: "Unable to save the root key shards",
			Message: fmt.Sprintf("You have %d keys parts while the minimal number is %d", len(vk.Shards), vk.MinimumRequired)}
	}

	vk.InitialRootKey = hf.GetPassword("[OPTIONAL, BUT RECOMMENDED] Enter the initial root key : ", shared.DebugMode)
	vk.Comments = hf.GetStringValFromPrompt("[OPTIONAL] Enter comments : ")
	return vk.saveRootKeys(rkfile)
}
