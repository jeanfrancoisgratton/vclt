// vclt
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename : login.go
// Original timestamp : 2024/06/28 15:38

package sys

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	cerr "github.com/jeanfrancoisgratton/customError"
	hpcv "github.com/jeanfrancoisgratton/hcpVaultLib"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"os"
	"vclt/env"
)

func LoginUserPass() *cerr.CustomError {
	var serr error

	cs, err := env.LoadEnvironmentFile()
	if err != nil {
		return err
	}
	config := &api.Config{Address: cs.VaultAddress}

	vc, err := hpcv.NewClient(cs.VaultAddress)
	if err != nil {
		return err
	}
	vc.Client, serr = api.NewClient(config)
	if serr != nil {
		return &cerr.CustomError{Title: "Error creating vault client", Message: serr.Error()}
	}

	tkn, ce := vc.AuthManager.LoginWithUserPass(cs.VaultUsername, hf.DecodeString(cs.VaultPassword, ""))
	if ce == nil {
		if Verbose {
			fmt.Printf("Authentication as %s on %s %s\n", cs.VaultUsername, cs.VaultAddress, hf.Green("Succeeded"))
		}
		if err := os.Setenv("VAULT_TOKEN", tkn); err != nil {
			return &cerr.CustomError{Title: "Unable to set environment variable: VAULT_TOKEN", Message: err.Error()}
		}
		if SysStoreToken {
			return StoreTokenInEnv()
		}
		_setAuthTkn(tkn)
		_setLoggedIn()
	}
	return ce
}

func LoginToken() *cerr.CustomError {
	var serr error
	cs, err := env.LoadEnvironmentFile()
	if err != nil {
		return err
	}
	config := &api.Config{Address: cs.VaultAddress}

	vc, err := hpcv.NewClient(cs.VaultAddress)
	if err != nil {
		return err
	}
	vc.Client, serr = api.NewClient(config)
	if serr != nil {
		return &cerr.CustomError{Title: "Error creating vault client", Message: serr.Error()}
	}

	if cs.VaultToken == "" {
		cs.VaultToken = os.Getenv("VAULT_TOKEN")
		if cs.VaultToken == "" {
			return &cerr.CustomError{Title: "Authentication failure", Message: "No token in environment file and VAULT_TOKEN is unset",
				Fatality: cerr.Warning}
		}

		vc.AuthManager.LoginWithToken(cs.VaultToken)
		_setAuthTkn(cs.VaultToken)
		_setLoggedIn()
		if SysStoreToken {
			return StoreTokenInEnv()
		}
	}
	return nil
}
