// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/parser_types.go
// Original timestamp: 2026/06/22 06:02:54

package policies

import (
	"encoding/json"
)

// Document mirrors Vault's JSON ACL policy schema: a single top-level
// "path" object keyed by path string, each holding that path's rule.
type Document struct {
	Path map[string]Rule `json:"path"`
}

// Rule is the full set of per-path fields Vault's ACL policy engine
// recognizes.  All fields except Capabilities are optional; omitempty
// ensures they are absent from the re-marshaled output when unset.
type Rule struct {
	// Capabilities is required and must be non-empty.  Valid values:
	// create, read, update, delete, list, patch, sudo, deny.
	// "deny" is mutually exclusive with all other values.
	Capabilities []string `json:"capabilities"`

	// AllowedParameters whitelists specific request parameter key/value
	// pairs.  Map key is the parameter name; the special key "*" matches
	// every parameter and must map to an empty slice.  An empty slice for
	// any key means "any value for this parameter is permitted".
	// Values may be any JSON-representable type (string, number, bool).
	AllowedParameters map[string][]json.RawMessage `json:"allowed_parameters,omitempty"`

	// DeniedParameters blacklists specific request parameter key/value
	// pairs.  The "*" wildcard key must map to an empty slice (denies all
	// parameters entirely).  An empty slice for any other key means
	// "reject the parameter regardless of its value".
	DeniedParameters map[string][]json.RawMessage `json:"denied_parameters,omitempty"`

	// RequiredParameters lists parameter names that must be present in the
	// request; Vault denies the request if any are absent.
	RequiredParameters []string `json:"required_parameters,omitempty"`

	// MinWrappingTTL is the minimum TTL clients may request when wrapping
	// a response for this path.  Setting it to "1s" effectively makes
	// response wrapping mandatory.  Accepts a Go duration string
	// ("5m", "1h30m") or a plain non-negative integer number of seconds.
	MinWrappingTTL string `json:"min_wrapping_ttl,omitempty"`

	// MaxWrappingTTL is the maximum TTL clients may request when wrapping
	// a response for this path.  Same format as MinWrappingTTL.
	// When both are set, MinWrappingTTL must be strictly less than
	// MaxWrappingTTL.
	MaxWrappingTTL string `json:"max_wrapping_ttl,omitempty"`

	// MFAMethods lists the MFA method configurations (defined in Vault's
	// MFA config) that must all be satisfied before access is granted.
	// When present the slice must not be empty.
	MFAMethods []string `json:"mfa_methods,omitempty"`
}

// validCapabilities is Vault's fixed capability vocabulary.
var validCapabilities = map[string]bool{
	"create": true, "read": true, "update": true, "delete": true,
	"list": true, "patch": true, "sudo": true, "deny": true,
}
