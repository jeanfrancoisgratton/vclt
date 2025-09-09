// vclt
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename : sys/tokenOps.go
// Original timestamp : 2025/03/21 16:06

package sys

import (
	cerr "github.com/jeanfrancoisgratton/customError"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"os"
	"vclt/env"
)

func StoreTokenInEnv() *cerr.CustomError {
	var cs env.Config_s
	var err *cerr.CustomError

	if cs, err = env.LoadEnvironmentFile(); err != nil {
		return err
	}
	cs.VaultToken = hf.EncodeString(os.Getenv("VAULT_TOKEN"), "")
	return cs.SaveEnvironmentFile(env.ConfigFile)
}
