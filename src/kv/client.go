// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/client.go

package kv

import (
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	vlr "github.com/jeanfrancoisgratton/vaultlib/v2/kv"
)

// Client is a KV-engine-scoped Vault client. Construct it with NewClient,
// passing the KV engine (mount path); the token/address resolution and the
// underlying vaultlib client creation happen once. All methods operate
// against the engine baked into the client.
type Client struct {
	vc     *vlr.Client
	engine string
}

// NewClient resolves the Vault auth token and server address (flag → env →
// ~/.vault-token file), then builds the underlying vaultlib KV client
// scoped to kvengine.
func NewClient(kvengine string) (*Client, *ce.CustomError) {
	if err := shared.SetVaultToken(); err != nil {
		return nil, err
	}
	if err := shared.SetServerAddress(); err != nil {
		return nil, err
	}

	cfg := vlr.Config{Address: shared.VaultServerAddress, Token: shared.VaultAuthToken, MountPath: kvengine}
	vc, err := vlr.NewClient(cfg)
	if err != nil {
		return nil, &ce.CustomError{Title: "Error creating vault client", Message: err.Error()}
	}
	return &Client{vc: vc, engine: kvengine}, nil
}
