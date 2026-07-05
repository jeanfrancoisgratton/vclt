// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/policy_parser_test.go

package policies

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestParseTTLDuration is a "table-driven test": instead of one test
// function per input, we declare a slice of cases and loop over them.
// t.Run gives each case its own name, so a failure reports exactly which
// input broke (and `go test -run TestParseTTLDuration/plain_seconds`
// can re-run a single case).
func TestParseTTLDuration(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    time.Duration
		wantErr bool
	}{
		{name: "plain seconds", in: "90", want: 90 * time.Second},
		{name: "zero seconds", in: "0", want: 0},
		{name: "go duration", in: "1h30m", want: 90 * time.Minute},
		{name: "sub-second duration", in: "500ms", want: 500 * time.Millisecond},
		{name: "negative integer falls through to ParseDuration and fails", in: "-5", wantErr: true},
		{name: "garbage", in: "banana", wantErr: true},
		{name: "empty string", in: "", wantErr: true},
		{name: "seconds overflow", in: "99999999999999999999", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseTTLDuration(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("parseTTLDuration(%q) = %v, expected an error", tc.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseTTLDuration(%q) returned unexpected error: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("parseTTLDuration(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// TestParseJSON exercises the whole validation chain (capabilities,
// deny exclusivity, wildcard parameters, TTL ordering, mfa_methods)
// through the public entry point, using real policy documents as input.
func TestParseJSON(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantErr string // substring the error must contain; "" means expect success
	}{
		{
			name: "minimal valid policy",
			in:   `{"path": {"secret/data/app": {"capabilities": ["read"]}}}`,
		},
		{
			name: "full-featured valid policy",
			in: `{"path": {"secret/data/app/*": {
				"capabilities": ["create", "update"],
				"allowed_parameters": {"ttl": [], "*": []},
				"required_parameters": ["ttl"],
				"min_wrapping_ttl": "1s",
				"max_wrapping_ttl": "90m",
				"mfa_methods": ["duo"]
			}}}`,
		},
		{
			name:    "not JSON at all",
			in:      `path "secret/*" { capabilities = ["read"] }`,
			wantErr: "invalid policy JSON",
		},
		{
			name:    "no path entries",
			in:      `{"path": {}}`,
			wantErr: `no "path" entries`,
		},
		{
			name:    "empty capabilities",
			in:      `{"path": {"secret/foo": {"capabilities": []}}}`,
			wantErr: "capabilities cannot be empty",
		},
		{
			name:    "unknown capability",
			in:      `{"path": {"secret/foo": {"capabilities": ["fly"]}}}`,
			wantErr: `unknown capability "fly"`,
		},
		{
			name:    "deny mixed with read",
			in:      `{"path": {"secret/foo": {"capabilities": ["deny", "read"]}}}`,
			wantErr: "mutually exclusive",
		},
		{
			name:    "wildcard allowed_parameters with values",
			in:      `{"path": {"secret/foo": {"capabilities": ["read"], "allowed_parameters": {"*": ["x"]}}}}`,
			wantErr: `wildcard key "*" must have an empty value list`,
		},
		{
			name:    "min TTL not below max TTL",
			in:      `{"path": {"secret/foo": {"capabilities": ["read"], "min_wrapping_ttl": "1h", "max_wrapping_ttl": "30m"}}}`,
			wantErr: "must be less than",
		},
		{
			name:    "bad TTL string",
			in:      `{"path": {"secret/foo": {"capabilities": ["read"], "min_wrapping_ttl": "soon"}}}`,
			wantErr: "invalid TTL",
		},
		{
			name:    "explicit empty mfa_methods",
			in:      `{"path": {"secret/foo": {"capabilities": ["read"], "mfa_methods": []}}}`,
			wantErr: "mfa_methods must not be an empty list",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := ParseJSON([]byte(tc.in))
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got success:\n%s", tc.wantErr, out)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("error %q does not contain %q", err.Error(), tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if out == "" {
				t.Error("valid policy produced empty canonical output")
			}
		})
	}
}

// TestValidateErrorIsDeterministic pins down the documented behaviour that
// paths are validated in sorted order: with two broken paths, the error
// must always mention the alphabetically-first one.
func TestValidateErrorIsDeterministic(t *testing.T) {
	in := `{"path": {
		"zzz/broken": {"capabilities": []},
		"aaa/broken": {"capabilities": []}
	}}`
	for i := 0; i < 20; i++ {
		_, err := ParseJSON([]byte(in))
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), `"aaa/broken"`) {
			t.Fatalf("run %d: error mentions the wrong path: %v", i, err)
		}
	}
}

// TestParsePolicyFile checks extension-based dispatch using real temporary
// files. t.TempDir() hands us a directory that the test framework deletes
// automatically when the test ends — no cleanup code needed.
func TestParsePolicyFile(t *testing.T) {
	dir := t.TempDir()

	jsonPath := filepath.Join(dir, "policy.json")
	if err := os.WriteFile(jsonPath, []byte(`{"path": {"secret/foo": {"capabilities": ["read"]}}}`), 0600); err != nil {
		t.Fatal(err)
	}
	hclPath := filepath.Join(dir, "policy.HCL") // uppercase on purpose: dispatch must be case-insensitive
	if err := os.WriteFile(hclPath, []byte("path \"secret/foo\" {\n  capabilities = [\"read\"]\n}\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if _, err := ParsePolicyFile(jsonPath); err != nil {
		t.Errorf("JSON dispatch failed: %v", err)
	}
	if _, err := ParsePolicyFile(hclPath); err != nil {
		t.Errorf("HCL dispatch failed: %v", err)
	}
	if _, err := ParsePolicyFile(filepath.Join(dir, "missing.json")); err == nil {
		t.Error("expected an error for a nonexistent file")
	}
}
