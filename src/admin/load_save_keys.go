// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/admin/load_save_keys.go
// Original timestamp: 2026/06/18 09:34:27

package admin

import (
	"encoding/json"
	"os"
	"path/filepath"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

func (vk VaultRootKeysStruct) saveRootKeys(rkfile string) *ce.CustomError {
	jStream, err := json.MarshalIndent(vk, "", "  ")
	if err != nil {
		return &ce.CustomError{Title: err.Error()}
	}

	if err := os.WriteFile(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vclt", rkfile), jStream, 0600); err != nil {
		return &ce.CustomError{Title: "Unable to write root keys json file", Message: err.Error()}
	}
	return nil
}

func (vk *VaultRootKeysStruct) loadRootKeys(rkfile string) *ce.CustomError {
	jFile, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vclt", rkfile))
	if err != nil {
		return &ce.CustomError{Title: err.Error()}
	}
	err = json.Unmarshal(jFile, &vk)
	if err != nil {
		return &ce.CustomError{Title: err.Error()}
	}
	return nil
}
