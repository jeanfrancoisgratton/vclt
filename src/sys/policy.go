// vclt
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename : policy.go
// Original timestamp : 2024/06/30 19:10

package sys

import (
	vault "github.com/hashicorp/vault/api"
	cerr "github.com/jeanfrancoisgratton/customError"
)

func CreatePolicy(client *vault.Client, policyName string, policy string) *cerr.CustomError {
	path := "sys/policies/acl/" + policyName
	data := map[string]interface{}{
		"policy": policy,
	}

	_, err := client.Logical().Write(path, data)
	if err != nil {
		return &cerr.CustomError{Title: err.Error()}
	}
	return nil
}
