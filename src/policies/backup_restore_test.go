// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/backup_restore_test.go

package policies

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl"
)

// TestMarshalBundleHCLRoundTrip exercises marshalBundleHCL (the one place a
// policy backup's HCL form has to be hand-rendered, since hashicorp/hcl v1
// has no encoder) back to back with hcl.Unmarshal, covering rules text that
// is itself HCL, is itself JSON, and contains characters HCL's string
// syntax must escape correctly (backslashes, tabs, literal "${...}"
// sequences that must NOT be interpolated, non-ASCII text).
func TestMarshalBundleHCLRoundTrip(t *testing.T) {
	orig := policyBundle{Policy: map[string]policyEntry{
		"app-readonly": {Rules: "path \"secret/data/app/*\" {\n  capabilities = [\"read\", \"list\"]\n}"},
		"app-json":     {Rules: `{"path":{"secret/data/app2/*":{"capabilities":["read"]}}}`},
		"weird-chars":  {Rules: "line with \\ backslash and \t tab and ${not_an_interp} and unicode é"},
	}}

	hclText := marshalBundleHCL(orig)

	var decoded policyBundle
	if err := hcl.Unmarshal([]byte(hclText), &decoded); err != nil {
		t.Fatalf("hcl.Unmarshal failed on generated HCL: %v\n--- generated ---\n%s", err, hclText)
	}

	origJSON, _ := json.Marshal(orig)
	decodedJSON, _ := json.Marshal(decoded)
	if string(origJSON) != string(decodedJSON) {
		t.Errorf("round trip mismatch:\norig:    %s\ndecoded: %s", origJSON, decodedJSON)
	}
}

// TestPolicyBackupFileRoundTrip exercises writePolicyBackupFile and
// readPolicyBackupFile back to back, both local-file-only operations with
// no Vault server involved, under both Cleartext settings.
func TestPolicyBackupFileRoundTrip(t *testing.T) {
	orig := Cleartext
	t.Cleanup(func() { Cleartext = orig })

	for _, cleartext := range []bool{true, false} {
		t.Run(map[bool]string{true: "cleartext", false: "encoded"}[cleartext], func(t *testing.T) {
			Cleartext = cleartext

			dir := t.TempDir()
			path := filepath.Join(dir, "policies-backup.json")
			payload := []byte(`{"policy":{"app-readonly":{"rules":"path \"secret/*\" {}"}}}`)

			if err := writePolicyBackupFile(path, payload); err != nil {
				t.Fatalf("writePolicyBackupFile failed: %v", err)
			}

			onDisk, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if cleartext && string(onDisk) != string(payload) {
				t.Errorf("cleartext file content = %q, want unchanged %q", onDisk, payload)
			}
			if !cleartext && string(onDisk) == string(payload) {
				t.Error("file content is unchanged after writePolicyBackupFile with Cleartext=false; expected it to be encoded")
			}

			got, rerr := readPolicyBackupFile(path)
			if rerr != nil {
				t.Fatalf("readPolicyBackupFile failed: %v", rerr)
			}
			if string(got) != string(payload) {
				t.Errorf("readPolicyBackupFile() = %q, want %q", got, payload)
			}

			// No scratch temp files left behind either way.
			entries, _ := os.ReadDir(dir)
			if len(entries) != 1 {
				names := make([]string, len(entries))
				for i, e := range entries {
					names[i] = e.Name()
				}
				t.Errorf("expected exactly 1 file in %s, got %v", dir, names)
			}
		})
	}
}

func TestWritePolicyBackupFileBadDir(t *testing.T) {
	orig := Cleartext
	Cleartext = false
	t.Cleanup(func() { Cleartext = orig })

	err := writePolicyBackupFile(filepath.Join(t.TempDir(), "does-not-exist", "backup.json"), []byte("{}"))
	if err == nil {
		t.Fatal("expected an error writing into a nonexistent directory, got nil")
	}
}

func TestReadPolicyBackupFileMissing(t *testing.T) {
	orig := Cleartext
	t.Cleanup(func() { Cleartext = orig })

	for _, cleartext := range []bool{true, false} {
		t.Run(map[bool]string{true: "cleartext", false: "encoded"}[cleartext], func(t *testing.T) {
			Cleartext = cleartext
			_, err := readPolicyBackupFile(filepath.Join(t.TempDir(), "does-not-exist.json"))
			if err == nil {
				t.Fatal("expected an error reading a nonexistent backup file, got nil")
			}
		})
	}
}
