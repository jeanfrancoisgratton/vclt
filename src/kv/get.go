package kv

import (
	"github.com/hashicorp/vault/api"
	cerr "github.com/jeanfrancoisgratton/customError"
	vlam "github.com/jeanfrancoisgratton/hcpVaultLib/auth"
	vlkm "github.com/jeanfrancoisgratton/hcpVaultLib/kv"
	"os"
	"vclt/env"
)

// TODO : manage secret version
func Get(path, field string, ver int) (string, *cerr.CustomError) {
	//var serr error

	vtkn := os.Getenv("VAULT_TOKEN")
	if vtkn == "" {
		return "", &cerr.CustomError{Title: "Unable to get secret", Message: "Not logged in"}
	}

	// Now we set the client up :
	// first we load the environment file, to get the server's address and the user's token
	// second we set the config up tru the API
	// third we set the client up, using the config from step 2
	// finally, we create the Auth and KV managers, from the above steps
	cs, err := env.LoadEnvironmentFile()
	if err != nil {
		return "", err
	}
	config := &api.Config{Address: cs.VaultAddress}
	client, se := api.NewClient(config)
	if se != nil {
		return "", &cerr.CustomError{Title: "Unable to create a new Vault client", Message: se.Error()}
	}
	am := vlam.NewAuthManager(client, cs.VaultAddress, "")
	if am == nil {
		return "", &cerr.CustomError{Title: "Cannot create Auth Manager"}
	}
	km := vlkm.NewKVManager(client, am)
	if km == nil {
		return "", &cerr.CustomError{Title: "Cannot create KV Manager"}
	}

	// At long last, we read the secret
	if res, cerr := km.ReadSecret(path, field); cerr != nil {
		return "", cerr
	} else {
		return res.(string), nil
	}
}
