// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/shared/setglobals_test.go

package shared

import (
	"os"
	"path/filepath"
	"testing"
)

// resetGlobals clears the package-level VaultAuthToken/VaultServerAddress
// before and after each test, since SetVaultToken/SetServerAddress mutate
// them as a side effect and every test in this file shares the same
// package-level state.
func resetGlobals(t *testing.T) {
	t.Helper()
	VaultAuthToken = ""
	VaultServerAddress = ""
	t.Cleanup(func() {
		VaultAuthToken = ""
		VaultServerAddress = ""
	})
}

func TestSetServerAddressFlagTakesPrecedence(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_ADDR", "https://from-env:8200")
	VaultServerAddress = "https://from-flag:8200"

	if err := SetServerAddress(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if VaultServerAddress != "https://from-flag:8200" {
		t.Errorf("got %q, want the flag value to win over $VAULT_ADDR", VaultServerAddress)
	}
}

func TestSetServerAddressFallsBackToEnv(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_ADDR", "https://from-env:8200")

	if err := SetServerAddress(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if VaultServerAddress != "https://from-env:8200" {
		t.Errorf("got %q, want $VAULT_ADDR value", VaultServerAddress)
	}
}

func TestSetServerAddressMissing(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_ADDR", "")

	err := SetServerAddress()
	if err == nil {
		t.Fatal("expected an error when neither -a nor $VAULT_ADDR is set, got nil")
	}
	if err.Code != ErrVaultServerAddressMissing {
		t.Errorf("error code = %d, want %d (ErrVaultServerAddressMissing)", err.Code, ErrVaultServerAddressMissing)
	}
}

func TestSetVaultTokenFlagTakesPrecedence(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_TOKEN", "token-from-env")
	VaultAuthToken = "token-from-flag"

	if err := SetVaultToken(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if VaultAuthToken != "token-from-flag" {
		t.Errorf("got %q, want the flag value to win over $VAULT_TOKEN", VaultAuthToken)
	}
}

func TestSetVaultTokenFallsBackToEnv(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_TOKEN", "token-from-env")

	if err := SetVaultToken(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if VaultAuthToken != "token-from-env" {
		t.Errorf("got %q, want $VAULT_TOKEN value", VaultAuthToken)
	}
}

// TestSetVaultTokenFallsBackToFile covers the third rung of the resolution
// chain: ~/.vault-token, whitespace-trimmed, used only when neither the
// flag nor $VAULT_TOKEN are set.
func TestSetVaultTokenFallsBackToFile(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_TOKEN", "")
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.WriteFile(filepath.Join(home, ".vault-token"), []byte("token-from-file\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := SetVaultToken(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if VaultAuthToken != "token-from-file" {
		t.Errorf("got %q, want trimmed contents of ~/.vault-token", VaultAuthToken)
	}
}

func TestSetVaultTokenMissing(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_TOKEN", "")
	t.Setenv("HOME", t.TempDir()) // empty dir: no .vault-token file

	err := SetVaultToken()
	if err == nil {
		t.Fatal("expected an error when no token source is available, got nil")
	}
	if err.Code != ErrVaultAuthTokenMissing {
		t.Errorf("error code = %d, want %d (ErrVaultAuthTokenMissing)", err.Code, ErrVaultAuthTokenMissing)
	}
}
