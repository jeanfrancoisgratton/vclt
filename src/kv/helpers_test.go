// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/helpers_test.go

package kv

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"vclt/shared"
)

// TestClassifyReadError pins down the error-message-substring -> error-code
// mapping that translates vaultlib's wrapped errors into vclt's own exit
// codes. Order matters in the source (sealed must be checked before the
// more general "unavailable"), so this also guards against that ordering
// regressing.
func TestClassifyReadError(t *testing.T) {
	cases := []struct {
		name    string
		errText string
		want    int
	}{
		{"sealed", "vault is sealed or unavailable", shared.ErrVaultSealed},
		{"unavailable", "Vault is unavailable", shared.ErrVaultUnavailable},
		{"connection refused", "dial tcp: connection refused", shared.ErrVaultUnavailable},
		{"no such host", "dial tcp: no such host", shared.ErrVaultUnavailable},
		{"unauthorized", "unauthorized", shared.ErrVaultInvalidAuth},
		{"invalid token", "invalid Vault token", shared.ErrVaultInvalidAuth},
		{"permission denied", "permission denied", shared.ErrVaultInvalidAuth},
		{"missing path", "secret does not exist", shared.ErrInvalidPath},
		{"unrecognized", "something else entirely broke", shared.ErrReadSecret},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := classifyReadError(errors.New(tc.errText))
			if got.Code != tc.want {
				t.Errorf("classifyReadError(%q).Code = %d, want %d", tc.errText, got.Code, tc.want)
			}
		})
	}
}

// TestOutputDataFieldNotFound checks the one code path in outputData that
// returns an error regardless of the suppress flag: a requested field that
// isn't present in the secret's data.
func TestOutputDataFieldNotFound(t *testing.T) {
	old := SecretField
	SecretField = "missing-field"
	t.Cleanup(func() { SecretField = old })

	err := outputData(map[string]interface{}{"present-field": "value"}, true)
	if err == nil {
		t.Fatal("expected an error for a missing field, got nil")
	}
	if err.Code != shared.ErrFieldNotFound {
		t.Errorf("error code = %d, want %d (ErrFieldNotFound)", err.Code, shared.ErrFieldNotFound)
	}
}

// TestOutputDataSuppressed exercises the suppress=true path (used by kv rm
// after resolving a field) for both the whole-secret and single-field
// cases: neither should print anything or return an error.
func TestOutputDataSuppressed(t *testing.T) {
	old := SecretField
	t.Cleanup(func() { SecretField = old })

	SecretField = ""
	if err := outputData(map[string]interface{}{"a": 1}, true); err != nil {
		t.Errorf("whole-secret suppressed call returned an error: %v", err)
	}

	SecretField = "a"
	if err := outputData(map[string]interface{}{"a": 1}, true); err != nil {
		t.Errorf("single-field suppressed call returned an error: %v", err)
	}
}

// TestWriteSecretFile checks that writeSecretFile creates the file with
// owner-only permissions and the exact content it was given.
func TestWriteSecretFile(t *testing.T) {
	old := SecretOutputFile
	t.Cleanup(func() { SecretOutputFile = old })

	dir := t.TempDir()
	SecretOutputFile = filepath.Join(dir, "secret.txt")
	shared.QuietOutput = true
	t.Cleanup(func() { shared.QuietOutput = false })

	if err := writeSecretFile([]byte("s3cr3t\n")); err != nil {
		t.Fatalf("writeSecretFile failed: %v", err)
	}

	info, statErr := os.Stat(SecretOutputFile)
	if statErr != nil {
		t.Fatal(statErr)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file permissions = %o, want 600", perm)
	}

	got, readErr := os.ReadFile(SecretOutputFile)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(got) != "s3cr3t\n" {
		t.Errorf("file content = %q, want %q", got, "s3cr3t\n")
	}
}

// TestWriteSecretFileError checks that an unwritable destination surfaces
// as ErrWriteFile rather than panicking or being silently swallowed.
func TestWriteSecretFileError(t *testing.T) {
	old := SecretOutputFile
	t.Cleanup(func() { SecretOutputFile = old })

	// The parent directory doesn't exist, so the write must fail.
	SecretOutputFile = filepath.Join(t.TempDir(), "no-such-subdir", "secret.txt")

	err := writeSecretFile([]byte("s3cr3t"))
	if err == nil {
		t.Fatal("expected an error writing to a nonexistent directory, got nil")
	}
	if err.Code != shared.ErrWriteFile {
		t.Errorf("error code = %d, want %d (ErrWriteFile)", err.Code, shared.ErrWriteFile)
	}
}
