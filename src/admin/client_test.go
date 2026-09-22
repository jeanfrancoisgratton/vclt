// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/admin/client_test.go

package admin

import (
	"testing"

	"vclt/shared"
)

// TestNewClientNoTokenAvailable checks that NewClient fails fast on the
// token-resolution step, before ever attempting to reach a Vault server,
// when none of -t / $VAULT_TOKEN / ~/.vault-token are available.
func TestNewClientNoTokenAvailable(t *testing.T) {
	shared.VaultAuthToken = ""
	shared.VaultServerAddress = ""
	t.Cleanup(func() {
		shared.VaultAuthToken = ""
		shared.VaultServerAddress = ""
	})

	t.Setenv("VAULT_TOKEN", "")
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("HOME", t.TempDir()) // empty dir: no .vault-token file

	c, err := NewClient()
	if err == nil {
		t.Fatal("expected an error when no token source is available, got nil")
	}
	if c != nil {
		t.Error("expected a nil client alongside the error")
	}
	if err.Code != shared.ErrVaultAuthTokenMissing {
		t.Errorf("error code = %d, want %d (ErrVaultAuthTokenMissing)", err.Code, shared.ErrVaultAuthTokenMissing)
	}
}

// TestNewUnsealClientNoTokenRequired checks the documented asymmetry:
// unsealing operates against a sealed, unauthenticated server, so
// NewUnsealClient must not require a token — only the server address.
// Without an address either, it should fail on that check instead
// (ErrVaultServerAddressMissing, not ErrVaultAuthTokenMissing).
func TestNewUnsealClientNoTokenRequired(t *testing.T) {
	shared.VaultAuthToken = ""
	shared.VaultServerAddress = ""
	t.Cleanup(func() {
		shared.VaultAuthToken = ""
		shared.VaultServerAddress = ""
	})

	t.Setenv("VAULT_TOKEN", "")
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("HOME", t.TempDir())

	c, err := NewUnsealClient()
	if err == nil {
		t.Fatal("expected an error when no server address is available, got nil")
	}
	if c != nil {
		t.Error("expected a nil client alongside the error")
	}
	if err.Code != shared.ErrVaultServerAddressMissing {
		t.Errorf("error code = %d, want %d (ErrVaultServerAddressMissing) -- NewUnsealClient must not require a token", err.Code, shared.ErrVaultServerAddressMissing)
	}
}
