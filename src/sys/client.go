// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/sys/client.go

package sys

import (
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"github.com/jeanfrancoisgratton/vaultlib/v2/sys"
)

// Client is a sys-scoped Vault client. Construct it once with NewClient,
// then call its methods; the token/address resolution and the underlying
// vaultlib client creation happen a single time. The resolved config is
// retained because some sys operations (e.g. ListMounts) take it directly
// rather than going through the client.
type Client struct {
	vc  *sys.Client
	cfg sys.Config
}

// NewClient resolves the Vault auth token and server address (flag → env →
// ~/.vault-token file), then builds the underlying vaultlib sys client.
func NewClient() (*Client, *ce.CustomError) {
	if err := shared.SetVaultToken(); err != nil {
		return nil, err
	}
	if err := shared.SetServerAddress(); err != nil {
		return nil, err
	}

	cfg := sys.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken}
	vc, err := sys.NewClient(cfg)
	if err != nil {
		return nil, &ce.CustomError{Title: "Error creating vault client", Message: err.Error()}
	}
	return &Client{vc: vc, cfg: cfg}, nil
}
