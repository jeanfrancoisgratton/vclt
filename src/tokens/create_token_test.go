// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/tokens/create_token_test.go

package tokens

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	tkn "github.com/jeanfrancoisgratton/vaultlib/v2/tokens"
)

// TestValidateTTL is a table-driven test covering both accepted forms (a
// bare integer number of seconds, or a Go duration string) and the
// rejected ones (negative, empty, garbage).
func TestValidateTTL(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{name: "bare seconds", in: "3600"},
		{name: "zero seconds", in: "0"},
		{name: "go duration hours", in: "24h"},
		{name: "go duration minutes", in: "30m"},
		{name: "whitespace padded", in: "  1h  "},
		{name: "negative integer", in: "-5", wantErr: true},
		{name: "negative duration", in: "-1h", wantErr: true},
		{name: "empty string", in: "", wantErr: true},
		{name: "garbage", in: "banana", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTTL(tc.in)
			if tc.wantErr && err == nil {
				t.Fatalf("validateTTL(%q) = nil, expected an error", tc.in)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("validateTTL(%q) returned unexpected error: %v", tc.in, err)
			}
		})
	}
}

// TestSplitPolicies checks the "default" fallback for an empty --policies
// flag, and comma-splitting otherwise.
func TestSplitPolicies(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{name: "empty falls back to default", in: "", want: []string{"default"}},
		{name: "single policy", in: "deploy-policy", want: []string{"deploy-policy"}},
		{name: "multiple policies", in: "deploy-policy,read-policy", want: []string{"deploy-policy", "read-policy"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := splitPolicies(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("splitPolicies(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// TestSaveTokenInfoToFile checks the .json suffix auto-append, the file's
// permissions (token material must stay owner-only), and that the written
// JSON round-trips to the original struct.
func TestSaveTokenInfoToFile(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "ci-deploy-token") // no .json extension on purpose

	original := &tkn.TokenAuth{
		ClientToken:   "hvs.notarealtoken",
		Accessor:      "accessor-id",
		Policies:      []string{"default", "deploy-policy"},
		TokenPolicies: []string{"default", "deploy-policy"},
		LeaseDuration: 3600,
		Renewable:     true,
		Orphan:        true,
	}

	if err := saveTokenInfoToFile(target, original); err != nil {
		t.Fatalf("saveTokenInfoToFile failed: %v", err)
	}

	savedPath := target + ".json"
	info, statErr := os.Stat(savedPath)
	if statErr != nil {
		t.Fatalf("expected %s to exist: %v", savedPath, statErr)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file permissions = %o, want 600", perm)
	}

	raw, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatal(err)
	}
	var loaded tkn.TokenAuth
	if err := json.Unmarshal(raw, &loaded); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}
	if !reflect.DeepEqual(*original, loaded) {
		t.Errorf("round trip mismatch:\nsaved:  %+v\nloaded: %+v", *original, loaded)
	}
}

// TestSaveTokenInfoToFileExtensionNotDoubled checks that a filename already
// ending in .json (any case) isn't given a second .json suffix.
func TestSaveTokenInfoToFileExtensionNotDoubled(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "ci-deploy-token.JSON")

	if err := saveTokenInfoToFile(target, &tkn.TokenAuth{ClientToken: "x"}); err != nil {
		t.Fatalf("saveTokenInfoToFile failed: %v", err)
	}
	if _, err := os.Stat(target); err != nil {
		t.Errorf("expected the original filename to be preserved as-is: %v", err)
	}
	if _, err := os.Stat(target + ".json"); err == nil {
		t.Error("extension was doubled: found target.JSON.json")
	}
}
