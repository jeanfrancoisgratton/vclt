// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/generate_sample_test.go

package policies

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerateSamplePolicyExtensionDispatch checks the filename-extension
// dispatch documented on GenerateSamplePolicy: .json produces JSON, any
// other extension (or none) produces HCL with .hcl appended if missing.
func TestGenerateSamplePolicyExtensionDispatch(t *testing.T) {
	cases := []struct {
		name       string
		filename   string
		wantSuffix string
		wantJSON   bool
	}{
		{name: "json extension", filename: "mypolicy.json", wantSuffix: "mypolicy.json", wantJSON: true},
		{name: "no extension gets .hcl appended", filename: "mypolicy", wantSuffix: "mypolicy.hcl"},
		{name: "explicit .hcl extension kept as-is", filename: "mypolicy.hcl", wantSuffix: "mypolicy.hcl"},
		{name: "uppercase .HCL kept as-is", filename: "mypolicy.HCL", wantSuffix: "mypolicy.HCL"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			target := filepath.Join(dir, tc.filename)

			if err := GenerateSamplePolicy(target); err != nil {
				t.Fatalf("GenerateSamplePolicy failed: %v", err)
			}

			wantPath := filepath.Join(dir, tc.wantSuffix)
			content, readErr := os.ReadFile(wantPath)
			if readErr != nil {
				t.Fatalf("expected file at %s: %v", wantPath, readErr)
			}

			if tc.wantJSON {
				if _, err := ParseJSON(content); err != nil {
					t.Errorf("generated JSON sample doesn't parse as a valid policy: %v", err)
				}
			} else if !strings.Contains(string(content), "capabilities") {
				t.Errorf("generated HCL sample doesn't look like a policy:\n%s", content)
			}
		})
	}
}
