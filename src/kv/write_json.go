// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/write_json.go

package kv

import (
	"encoding/json"
	"fmt"
	"os"

	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// WriteJSON writes every field decoded from jsonFile's top-level JSON object
// to the secret at path in a single Vault API call, producing exactly one new
// KV v2 version. jsonFile is expected to hold a flat JSON object of fields --
// the same shape `kv read -o json` prints -- so a secret can round-trip out
// via `kv read -o json --out FILE` and back in via `kv write --json`. As with
// WriteBatch, fields already present in the secret are preserved: the JSON
// file's fields are overlaid on top of the current secret rather than
// replacing it outright.
func (c *Client) WriteJSON(path, jsonFile string) *ce.CustomError {
	fields, perr := parseJSONFile(jsonFile)
	if perr != nil {
		return perr
	}
	if len(fields) == 0 {
		return &ce.CustomError{Title: "Empty JSON secret file", Message: fmt.Sprintf("%s contains no fields", jsonFile), Code: shared.ErrReadFile}
	}

	if kvErr := c.mergeAndWrite(path, fields); kvErr != nil {
		return kvErr
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign(fmt.Sprintf("Wrote %d field(s) to %s", len(fields), path)))
	}
	return nil
}

// parseJSONFile reads path and decodes it as a flat JSON object of secret
// fields.
func parseJSONFile(path string) (map[string]interface{}, *ce.CustomError) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ce.CustomError{Title: shared.ErrorMessages[shared.ErrReadFile].Msg, Message: err.Error(), Code: shared.ErrReadFile}
	}

	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, &ce.CustomError{Title: "Invalid JSON secret file", Message: fmt.Sprintf("%s: %s", path, err.Error()), Code: shared.ErrReadFile}
	}

	return fields, nil
}
