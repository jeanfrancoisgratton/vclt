// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/policy_parser.go
// Original timestamp: 2026/06/22 06:04:26

package policies

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ParsePolicyFile dispatches to the correct parser based on the file
// extension: .hcl → HCL parser, .json (or anything else) → JSON parser.
// The returned string is ready to pass directly to policies.CreatePolicy.
func ParsePolicyFile(path string) (string, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".hcl":
		return ParseHCLFile(path)
	default:
		return ParseJSONFile(path)
	}
}

// ParseJSONFile reads the file at path, validates it as a Vault JSON ACL
// policy, and returns the canonical re-marshaled JSON text ready to pass
// directly to policies.CreatePolicy as the rules argument.
func ParseJSONFile(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading policy file %s: %w", path, err)
	}
	return ParseJSON(raw)
}

// ParseJSON validates raw JSON bytes as a Vault ACL policy document and
// returns the canonical re-marshaled JSON text.
func ParseJSON(raw []byte) (string, error) {
	var doc Document
	if err := json.Unmarshal(raw, &doc); err != nil {
		return "", fmt.Errorf("invalid policy JSON: %w", err)
	}
	if err := doc.validate(); err != nil {
		return "", err
	}
	canon, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("re-marshaling policy: %w", err)
	}
	return string(canon), nil
}

// validate checks the full Document for semantic correctness, returning
// the first error found.  Path names are evaluated in sorted order so
// error messages are deterministic regardless of Go's map iteration order.
func (d Document) validate() error {
	if len(d.Path) == 0 {
		return fmt.Errorf("invalid policy: no \"path\" entries defined")
	}

	paths := make([]string, 0, len(d.Path))
	for p := range d.Path {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	for _, p := range paths {
		if err := d.Path[p].validate(p); err != nil {
			return err
		}
	}
	return nil
}

// validate checks a single Rule for semantic correctness.
func (r Rule) validate(path string) error {
	if err := r.validateCapabilities(path); err != nil {
		return err
	}
	if err := r.validateParameterMap("allowed_parameters", path, r.AllowedParameters); err != nil {
		return err
	}
	if err := r.validateParameterMap("denied_parameters", path, r.DeniedParameters); err != nil {
		return err
	}
	if err := r.validateWrappingTTLs(path); err != nil {
		return err
	}
	if r.MFAMethods != nil && len(r.MFAMethods) == 0 {
		return fmt.Errorf("invalid policy: path %q: mfa_methods must not be an empty list when specified", path)
	}
	return nil
}

func (r Rule) validateCapabilities(path string) error {
	if len(r.Capabilities) == 0 {
		return fmt.Errorf("invalid policy: path %q: capabilities cannot be empty", path)
	}
	hasDeny := false
	for _, cap := range r.Capabilities {
		if !validCapabilities[cap] {
			return fmt.Errorf("invalid policy: path %q: unknown capability %q", path, cap)
		}
		if cap == "deny" {
			hasDeny = true
		}
	}
	if hasDeny && len(r.Capabilities) > 1 {
		return fmt.Errorf("invalid policy: path %q: \"deny\" is mutually exclusive with all other capabilities", path)
	}
	return nil
}

// validateParameterMap checks allowed_parameters or denied_parameters.
// The only constraint Vault enforces at policy-parse time beyond JSON
// validity is that the "*" wildcard key must carry an empty value list —
// any non-empty list for "*" would be silently meaningless, so we reject
// it early.
func (r Rule) validateParameterMap(field, path string, m map[string][]json.RawMessage) error {
	if vals, ok := m["*"]; ok && len(vals) != 0 {
		return fmt.Errorf("invalid policy: path %q: %s: wildcard key \"*\" must have an empty value list", path, field)
	}
	return nil
}

func (r Rule) validateWrappingTTLs(path string) error {
	if r.MinWrappingTTL != "" {
		if err := validateTTL(r.MinWrappingTTL); err != nil {
			return fmt.Errorf("invalid policy: path %q: min_wrapping_ttl: %w", path, err)
		}
	}
	if r.MaxWrappingTTL != "" {
		if err := validateTTL(r.MaxWrappingTTL); err != nil {
			return fmt.Errorf("invalid policy: path %q: max_wrapping_ttl: %w", path, err)
		}
	}
	if r.MinWrappingTTL != "" && r.MaxWrappingTTL != "" {
		min, _ := parseTTLDuration(r.MinWrappingTTL)
		max, _ := parseTTLDuration(r.MaxWrappingTTL)
		if min >= max {
			return fmt.Errorf(
				"invalid policy: path %q: min_wrapping_ttl (%s) must be less than max_wrapping_ttl (%s)",
				path, r.MinWrappingTTL, r.MaxWrappingTTL,
			)
		}
	}
	return nil
}

// validateTTL returns an error if s is not a valid TTL string.
func validateTTL(s string) error {
	if _, err := parseTTLDuration(s); err != nil {
		return fmt.Errorf(
			"invalid TTL %q: must be a Go duration string (e.g. \"5m\", \"1h30m\") or a non-negative integer number of seconds",
			s,
		)
	}
	return nil
}

// parseTTLDuration converts s to a time.Duration.  It accepts a plain
// non-negative integer (treating it as seconds) or any string accepted
// by time.ParseDuration.
func parseTTLDuration(s string) (time.Duration, error) {
	if n, err := strconv.ParseUint(s, 10, 64); err == nil {
		const maxDurationSeconds = uint64(math.MaxInt64) / uint64(time.Second)
		if n > maxDurationSeconds {
			return 0, fmt.Errorf("TTL %q exceeds maximum supported duration", s)
		}
		return time.Duration(n) * time.Second, nil
	}
	return time.ParseDuration(s)
}
