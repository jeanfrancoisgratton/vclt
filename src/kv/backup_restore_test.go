// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/backup_restore_test.go

package kv

import (
	"os"
	"path/filepath"
	"testing"
)

// TestEncodeDecodeRoundTrip exercises encodefile (backup_engine.go) and
// decode2temp (restore_engine.go) back to back, both local-file-only
// operations with no Vault server involved: a plaintext backup file is
// encoded in place, then decoded into a fresh temp file, and the result
// must match the original plaintext.
func TestEncodeDecodeRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "backup.json")
	plaintext := `{"secret/foo":{"bar":"baz"}}`
	if err := os.WriteFile(path, []byte(plaintext), 0600); err != nil {
		t.Fatal(err)
	}

	if err := encodefile(path); err != nil {
		t.Fatalf("encodefile failed: %v", err)
	}

	encoded, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) == plaintext {
		t.Fatal("file content is unchanged after encodefile; expected it to be encoded")
	}
	// encodefile's own .enc scratch file must not survive.
	if _, statErr := os.Stat(path + ".enc"); !os.IsNotExist(statErr) {
		t.Errorf("expected %s.enc to be cleaned up, stat error = %v", path, statErr)
	}

	decodedPath, derr := decode2temp(path)
	if derr != nil {
		t.Fatalf("decode2temp failed: %v", derr)
	}
	t.Cleanup(func() { os.Remove(decodedPath) })

	decoded, err := os.ReadFile(decodedPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(decoded) != plaintext {
		t.Errorf("decoded content = %q, want original %q", decoded, plaintext)
	}
}

func TestEncodefileMissingSource(t *testing.T) {
	err := encodefile(filepath.Join(t.TempDir(), "does-not-exist.json"))
	if err == nil {
		t.Fatal("expected an error encoding a nonexistent file, got nil")
	}
}

// TestDecode2tempCleansUpOnFailure checks that decode2temp doesn't leak its
// scratch temp file when decoding fails (e.g. the input isn't validly
// encoded).
func TestDecode2tempCleansUpOnFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "not-encoded.json")
	if err := os.WriteFile(path, []byte("this is not encoded data"), 0600); err != nil {
		t.Fatal(err)
	}

	before, _ := os.ReadDir(dir)

	_, err := decode2temp(path)
	if err == nil {
		t.Fatal("expected an error decoding non-encoded content, got nil")
	}

	after, _ := os.ReadDir(dir)
	if len(after) != len(before) {
		t.Errorf("directory entry count changed from %d to %d; a scratch temp file may have leaked", len(before), len(after))
	}
}
