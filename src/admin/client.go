// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/admin/client.go

package admin

import (
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	vadm "github.com/jeanfrancoisgratton/vaultlib/v2/admin"
)

// Client is an admin-scoped Vault client. Construct it with NewClient (for
// authenticated operations such as seal) or NewUnsealClient (for unsealing,
// which does not require an auth token). The resolved config is retained
// alongside the client for operations that take it directly.
type Client struct {
	vc  *vadm.Client
	cfg vadm.AdminConfig
}

// NewClient resolves the Vault auth token and server address, then builds
// the underlying vadm client. Use it for authenticated admin operations.
func NewClient() (*Client, *ce.CustomError) {
	if err := shared.SetVaultToken(); err != nil {
		return nil, err
	}
	return newClient()
}

// NewUnsealClient resolves only the server address, then builds the
// underlying vadm client. Unsealing does not require an auth token — it is
// precisely the operation run when the caller cannot yet authenticate — so
// the token is intentionally not resolved here.
func NewUnsealClient() (*Client, *ce.CustomError) {
	return newClient()
}

// newClient resolves the server address and builds the vadm client. The
// token, if any, is whatever the shared globals already hold.
func newClient() (*Client, *ce.CustomError) {
	if err := shared.SetServerAddress(); err != nil {
		return nil, err
	}

	cfg := vadm.AdminConfig{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken}
	vc, err := vadm.NewClient(cfg)
	if err != nil {
		return nil, &ce.CustomError{Title: "Unable to create Vault client", Message: err.Error()}
	}
	return &Client{vc: vc, cfg: cfg}, nil
}
