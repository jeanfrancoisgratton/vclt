// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/write_batch.go

package kv

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	vlr "github.com/jeanfrancoisgratton/vaultlib/v2/kv"
)

// WriteBatch writes every KEY/VALUE field parsed from batchFile to the
// secret at path in a single Vault API call, producing exactly one new KV
// v2 version (as opposed to one version per field, which looping Write
// would produce). Fields already present in the secret are preserved: the
// current secret is read first (a not-yet-existing secret is treated as
// empty, same as a single kv write against a brand-new path), the batch
// file's fields are overlaid on top, and the merged map is written back --
// the same upsert-one-field-leave-the-rest-alone behavior as a single kv
// write, just batched.
func (c *Client) WriteBatch(path, batchFile string) *ce.CustomError {
	fields, perr := parseBatchFile(batchFile)
	if perr != nil {
		return perr
	}
	if len(fields) == 0 {
		return &ce.CustomError{Title: "Empty batch file", Message: fmt.Sprintf("%s contains no KEY/VALUE fields", batchFile), Code: shared.ErrReadFile}
	}

	merged := map[string]interface{}{}
	if current, rerr := c.vc.ReadSecret(path, vlr.ReadOptions{FallbackToLatestAvailable: true}); rerr != nil {
		// A not-yet-existing secret is expected and not fatal: WriteBatch
		// must be able to create a brand-new secret from scratch, same as a
		// single `kv write` on a path that doesn't exist yet. Any other
		// failure (sealed, unauthorized, unavailable, ...) is real and must
		// not be silently papered over by proceeding as if the secret were
		// simply empty.
		if ceErr := classifyReadError(rerr); ceErr.Code != shared.ErrInvalidPath {
			return ceErr
		}
	} else if current != nil {
		for k, v := range current.Data {
			merged[k] = v
		}
	}

	for k, v := range fields {
		merged[k] = v
	}

	if _, kvErr := c.vc.WriteSecret(path, merged, vlr.WriteOptions{}); kvErr != nil {
		return &ce.CustomError{Title: "Error writing batch secret", Message: kvErr.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign(fmt.Sprintf("Wrote %d field(s) to %s", len(fields), path)))
	}
	return nil
}

// parseBatchFile reads batchFile and returns its KEY/VALUE fields. Blank
// lines and lines whose first non-space character is '#' are ignored.
// Each remaining line is split at the first '=' or whitespace character,
// whichever comes first: everything before it is the key, everything after
// it (leading/trailing whitespace trimmed) is the value -- since the split
// happens only once, a value may contain further spaces or '=' characters
// with no quoting needed. Wrapping a value in double quotes is only needed
// to preserve meaningful leading/trailing whitespace that would otherwise
// be trimmed.
func parseBatchFile(path string) (map[string]string, *ce.CustomError) {
	f, err := os.Open(path)
	if err != nil {
		return nil, &ce.CustomError{Title: shared.ErrorMessages[shared.ErrReadFile].Msg, Message: err.Error(), Code: shared.ErrReadFile}
	}
	defer f.Close()

	fields := map[string]string{}
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.IndexAny(line, " \t=")
		if idx <= 0 {
			return nil, &ce.CustomError{Title: "Invalid batch file line", Message: fmt.Sprintf("%s:%d: expected KEY followed by '=' or whitespace then VALUE, got %q", path, lineNum, line), Code: shared.ErrReadFile}
		}

		key := line[:idx]
		value := strings.TrimSpace(line[idx+1:])
		if len(value) >= 2 && strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			value = value[1 : len(value)-1]
		}

		fields[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, &ce.CustomError{Title: shared.ErrorMessages[shared.ErrReadFile].Msg, Message: err.Error(), Code: shared.ErrReadFile}
	}

	return fields, nil
}
