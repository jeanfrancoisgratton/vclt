// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/env/envHelpers.go
// Original timestamp: 2023/08/19 10:02

package env

import (
	"fmt"
	"os"
	"strings"

	cerr "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	"vclt/shared"
)

// getEnvParams : Prompts for values to fill up the Environment structure
func getEnvParams(cfgFile string) (Config_s, *cerr.CustomError) {
	cs := Config_s{}

	cs.EnvironmentName = hf.GetStringValFromPrompt("Enter the name of this environment : ")
	cs.VaultAddress = hf.GetStringValFromPrompt("Enter the address of the Vault (ex: https://mydomain:1234) : ")
	cs.VaultToken = hf.GetStringValFromPrompt("OPTIONAL: Please enter the user's auth token : ")
	cs.VaultUsername = hf.GetStringValFromPrompt("Please enter the username using the Vault : ")
	cs.VaultPassword = hf.GetPassword("Please enter that user's password (leaving this empty will get you prompted for it, later) : ", shared.QuietOutput)
	if cs.VaultPassword != "" {
		cs.VaultPassword = hf.EncodeString(cs.VaultPassword, "")
	}
	cs.KVEnginePath = hf.GetStringValFromPrompt("Please enter the secrets store mount (no trailing/leading slash) : ")
	cs.KVEnginePath = strings.TrimPrefix(cs.KVEnginePath, "/")
	cs.KVEnginePath = strings.TrimSuffix(cs.KVEnginePath, "/")
	if cs.KVEnginePath == "" {
		fmt.Println("Environment variable KVEnginePath is required")
		os.Exit(0)
	}
	cs.Comments = hf.GetStringValFromPrompt("(OPTIONAL) Please enter a comment : ")
	// Call the SaveEnvironmentFile() method using a pointer to cs
	if err := cs.SaveEnvironmentFile(cfgFile); err != nil {
		return Config_s{}, err
	}

	return cs, nil
}
