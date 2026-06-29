// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/hcl_parser.go

package policies

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
)

// ParseHCLFile reads the file at path, validates it as a Vault HCL ACL
// policy, and returns the canonical formatted HCL text ready to pass
// directly to policies.CreatePolicy as the rules argument.
func ParseHCLFile(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading policy file %s: %w", path, err)
	}
	return ParseHCL(raw)
}

// ParseHCL validates raw HCL bytes as a Vault ACL policy document and
// returns the canonical formatted HCL text.
//
// Validation reuses the same Document/Rule logic as the JSON path so
// capability, parameter-map, TTL, and MFA checks are identical regardless
// of input format.
func ParseHCL(raw []byte) (string, error) {
	var doc Document
	if err := hcl.Unmarshal(raw, &doc); err != nil {
		return "", fmt.Errorf("invalid policy HCL: %w", err)
	}
	if err := doc.validate(); err != nil {
		return "", err
	}

	// Re-format through the HCL printer so the output is canonical and
	// consistently indented regardless of how the input was laid out.
	// If the printer fails (unexpected after a clean Unmarshal), fall back
	// to the trimmed raw bytes — Vault is happy with either.
	formatted, err := printer.Format(raw)
	if err != nil {
		return strings.TrimSpace(string(raw)), nil
	}
	return strings.TrimSpace(string(formatted)), nil
}
