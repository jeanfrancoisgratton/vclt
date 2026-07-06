// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/client.go

package policies

import (
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	vpol "github.com/jeanfrancoisgratton/vaultlib/v2/policies"
)

// Client is a policies-scoped Vault client. Construct it once with
// NewClient, then call its methods; the token/address resolution and the
// underlying vaultlib client creation happen a single time.
type Client struct {
	vc *vpol.Client
}

// NewClient resolves the Vault auth token and server address (flag → env →
// ~/.vault-token file), then builds the underlying vaultlib policies
// client.
func NewClient() (*Client, *ce.CustomError) {
	if err := shared.SetVaultToken(); err != nil {
		return nil, err
	}
	if err := shared.SetServerAddress(); err != nil {
		return nil, err
	}

	cfg := vpol.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken}
	vc, err := vpol.NewClient(cfg)
	if err != nil {
		return nil, &ce.CustomError{Title: "Error creating vault client", Message: err.Error()}
	}
	return &Client{vc: vc}, nil
}
