// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/hcl_parser_test.go

package policies

import (
	"strings"
	"testing"
)

func TestParseHCL(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantErr string
	}{
		{
			name: "valid policy",
			in: `path "secret/data/app/*" {
  capabilities = ["read", "list"]
}`,
		},
		{
			name: "sloppy formatting still accepted",
			in:   `path "secret/foo" {capabilities=["read"]}`,
		},
		{
			name:    "syntactically broken HCL",
			in:      `path "secret/foo" { capabilities = `,
			wantErr: "invalid policy HCL",
		},
		{
			name: "valid HCL but invalid policy semantics",
			in: `path "secret/foo" {
  capabilities = ["deny", "read"]
}`,
			wantErr: "mutually exclusive",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := ParseHCL([]byte(tc.in))
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
			// The function promises trimmed output.
			if out != strings.TrimSpace(out) {
				t.Errorf("output has leading/trailing whitespace: %q", out)
			}
		})
	}
}
