// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/backup_policies.go

package policies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// policyBundle is the on-disk shape of a full ACL-policy backup: every
// policy name mapped to its rules exactly as Vault returned them. Vault
// stores a policy's rules verbatim as whatever text it was last written
// with (HCL or JSON), so Backup/Restore don't reinterpret or reformat that
// text -- Restore hands it straight back to CreatePolicy unchanged. Unlike
// a single policy file (see Document/Rule in parser_types.go), the bundle
// itself can still be written as either JSON or HCL, selected the same way
// Write/GenerateSamplePolicy already do: by path's extension.
type policyBundle struct {
	Policy map[string]policyEntry `json:"policy"`
}

type policyEntry struct {
	Rules string `json:"rules"`
}

// Backup dumps every ACL policy on the server to path, in JSON or HCL
// depending on path's extension (.hcl -> HCL, anything else -> JSON). By
// default the file is written encoded (helperFunctions/v5), matching
// kv.Backup; pass --cleartext (Cleartext = true) to keep it as plain
// JSON/HCL instead.
func (c *Client) Backup(path string) *ce.CustomError {
	names, err := c.List(false)
	if err != nil {
		return err
	}

	bundle := policyBundle{Policy: make(map[string]policyEntry, len(names))}
	for _, name := range names {
		policy, rerr := c.vc.ReadPolicy(name)
		if rerr != nil {
			return &ce.CustomError{Title: "Error reading the policy", Message: rerr.Error()}
		}
		bundle.Policy[name] = policyEntry{Rules: policy.Rules}
	}

	var payload []byte
	if strings.EqualFold(filepath.Ext(path), ".hcl") {
		payload = []byte(marshalBundleHCL(bundle))
	} else {
		var merr error
		payload, merr = json.MarshalIndent(bundle, "", "  ")
		if merr != nil {
			return &ce.CustomError{Title: "Error serializing policy backup", Message: merr.Error()}
		}
	}

	if werr := writePolicyBackupFile(path, payload); werr != nil {
		return werr
	}

	if !shared.QuietOutput {
		suffix := "ies"
		if len(names) == 1 {
			suffix = "y"
		}
		fmt.Printf("%s %d polic%s to %s\n", hftx.EnabledSign("Successfully dumped"), len(names), suffix, hftx.Green(path))
	}
	return nil
}

// writePolicyBackupFile writes payload to path, either directly (cleartext)
// or encoded through a temp file (the default), mirroring kv.Backup's
// encode-after-write behavior without needing its rename dance, since here
// the cleartext content is already in memory rather than having been
// written to path by a separate bulk-dump call first.
func writePolicyBackupFile(path string, payload []byte) *ce.CustomError {
	if Cleartext {
		if err := os.WriteFile(path, payload, 0600); err != nil {
			return &ce.CustomError{Title: shared.ErrorMessages[shared.ErrWriteFile].Msg, Message: err.Error(), Code: shared.ErrWriteFile}
		}
		return nil
	}

	tmp, terr := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp-*")
	if terr != nil {
		return &ce.CustomError{Title: shared.ErrorMessages[shared.ErrWriteFile].Msg, Message: terr.Error(), Code: shared.ErrWriteFile}
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, werr := tmp.Write(payload); werr != nil {
		tmp.Close()
		return &ce.CustomError{Title: shared.ErrorMessages[shared.ErrWriteFile].Msg, Message: werr.Error(), Code: shared.ErrWriteFile}
	}
	if cerr := tmp.Close(); cerr != nil {
		return &ce.CustomError{Title: shared.ErrorMessages[shared.ErrWriteFile].Msg, Message: cerr.Error(), Code: shared.ErrWriteFile}
	}

	if eerr := hf.EncodeFile(tmpPath, path, ""); eerr != nil {
		return &ce.CustomError{Title: "Error encoding policy backup file", Message: eerr.Error()}
	}
	return nil
}

// marshalBundleHCL hand-renders bundle as HCL: the hashicorp/hcl v1 package
// this project uses for parsing has no encoder, so this is the one place a
// policy backup's HCL form has to be built by hand rather than reusing a
// library. Each policy's rules text is embedded via strconv.Quote, which
// produces the same escaping HCL's own string-literal syntax expects
// (verified to round-trip through hcl.Unmarshal, including backslashes,
// tabs, literal "${...}" sequences, and non-ASCII text). Policy names are
// sorted for a deterministic, diff-friendly file.
func marshalBundleHCL(bundle policyBundle) string {
	names := make([]string, 0, len(bundle.Policy))
	for name := range bundle.Policy {
		names = append(names, name)
	}
	sort.Strings(names)

	var b strings.Builder
	for i, name := range names {
		if i > 0 {
			b.WriteString("\n\n")
		}
		fmt.Fprintf(&b, "policy %s {\n  rules = %s\n}", strconv.Quote(name), strconv.Quote(bundle.Policy[name].Rules))
	}
	b.WriteString("\n")
	return b.String()
}
