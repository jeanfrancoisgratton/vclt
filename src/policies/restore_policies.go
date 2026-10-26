// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/restore_policies.go

package policies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vclt/shared"

	"github.com/hashicorp/hcl"
	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// Restore recreates a single named ACL policy from a file produced by
// Backup. By default the file is assumed to be encoded; pass --cleartext
// (Cleartext = true) if the backup was taken with -c.
//
// Unlike RestoreAll, this file is just that one policy's rules text (not a
// bundle), the same shape 'policy write' expects as its POLICY_FILE -- so,
// per the decision to validate single-policy restores the same way
// 'policy write' does, this decodes the file (if needed) into a temp file
// and runs it through the exact same ParsePolicyFile validation Write uses,
// rather than handing Vault unvalidated content.
func (c *Client) Restore(policyName, path string) *ce.CustomError {
	payload, rerr := readPolicyBackupFile(path)
	if rerr != nil {
		return rerr
	}

	tmp, terr := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".restore-*"+filepath.Ext(path))
	if terr != nil {
		return &ce.CustomError{Title: shared.ErrorMessages[shared.ErrReadFile].Msg, Message: terr.Error(), Code: shared.ErrReadFile}
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, werr := tmp.Write(payload); werr != nil {
		tmp.Close()
		return &ce.CustomError{Title: shared.ErrorMessages[shared.ErrReadFile].Msg, Message: werr.Error(), Code: shared.ErrReadFile}
	}
	tmp.Close()

	rules, perr := ParsePolicyFile(tmpPath)
	if perr != nil {
		return &ce.CustomError{Title: "Invalid policy file", Message: perr.Error()}
	}

	if err := c.vc.CreatePolicy(policyName, rules); err != nil {
		return &ce.CustomError{Title: "Error restoring policy", Message: err.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign("Restored policy "+hftx.Green(policyName)) + " from " + path)
	}
	return nil
}

// RestoreAll recreates every ACL policy found in a bundle produced by
// BackupAll. Format (JSON or HCL) is selected the same way BackupAll chose
// it: by path's extension. By default the file is assumed to be encoded;
// pass --cleartext (Cleartext = true) if the backup was taken with -c.
//
// Each policy's rules text is handed straight back to CreatePolicy exactly
// as BackupAll captured it -- not re-validated through ParsePolicyFile --
// since it already passed Vault's own acceptance the first time it was
// written.
func (c *Client) RestoreAll(path string) *ce.CustomError {
	payload, rerr := readPolicyBackupFile(path)
	if rerr != nil {
		return rerr
	}

	var bundle policyBundle
	if strings.EqualFold(filepath.Ext(path), ".hcl") {
		if err := hcl.Unmarshal(payload, &bundle); err != nil {
			return &ce.CustomError{Title: "Invalid policy backup file", Message: fmt.Sprintf("%s: %s", path, err.Error())}
		}
	} else {
		if err := json.Unmarshal(payload, &bundle); err != nil {
			return &ce.CustomError{Title: "Invalid policy backup file", Message: fmt.Sprintf("%s: %s", path, err.Error())}
		}
	}

	if len(bundle.Policy) == 0 {
		return &ce.CustomError{Title: "Empty policy backup file", Message: fmt.Sprintf("%s contains no policies", path), Code: shared.ErrReadFile}
	}

	for name, entry := range bundle.Policy {
		if err := c.vc.CreatePolicy(name, entry.Rules); err != nil {
			return &ce.CustomError{Title: "Error restoring policy " + name, Message: err.Error()}
		}
	}

	if !shared.QuietOutput {
		suffix := "ies"
		if len(bundle.Policy) == 1 {
			suffix = "y"
		}
		fmt.Printf("%s %d polic%s from %s\n", hftx.EnabledSign("Successfully restored"), len(bundle.Policy), suffix, hftx.Green(path))
	}
	return nil
}

// readPolicyBackupFile reads path, either directly (cleartext) or decoded
// through a temp file (the default), mirroring kv.Restore's decode-before-
// read behavior.
func readPolicyBackupFile(path string) ([]byte, *ce.CustomError) {
	if Cleartext {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, &ce.CustomError{Title: shared.ErrorMessages[shared.ErrReadFile].Msg, Message: err.Error(), Code: shared.ErrReadFile}
		}
		return data, nil
	}

	tmp, terr := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".dec-*")
	if terr != nil {
		return nil, &ce.CustomError{Title: shared.ErrorMessages[shared.ErrReadFile].Msg, Message: terr.Error(), Code: shared.ErrReadFile}
	}
	tmpPath := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpPath)

	if derr := hf.DecodeFile(path, tmpPath, ""); derr != nil {
		return nil, &ce.CustomError{Title: "Error decoding policy backup file", Message: derr.Error()}
	}

	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, &ce.CustomError{Title: shared.ErrorMessages[shared.ErrReadFile].Msg, Message: err.Error(), Code: shared.ErrReadFile}
	}
	return data, nil
}
