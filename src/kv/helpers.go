// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/helpers.go
// Original timestamp: 2026/06/14 13:33:05

package kv

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"vclt/shared"
)

// writeSecretFile writes a secret's rendered content to SecretOutputFile with
// owner-only permissions, since the file may contain sensitive material.
func writeSecretFile(content []byte) *ce.CustomError {
	if err := os.WriteFile(SecretOutputFile, content, 0600); err != nil {
		return &ce.CustomError{Title: shared.ErrorMessages[shared.ErrWriteFile].Msg, Message: err.Error(), Code: shared.ErrWriteFile}
	}
	if !shared.QuietOutput {
		fmt.Printf("Secret written to %s\n", SecretOutputFile)
	}
	return nil
}

// readSecretValueFromFile reads a secret's value from SecretInputFile, used
// by `kv write --in FILE` as an alternative to passing VALUE positionally.
// A single trailing newline (or CRLF) is stripped, mirroring the newline
// writeSecretFile appends when a single field is written to a file, so a
// value round-trips cleanly through `kv read --field ... --out` and back
// through `kv write --in`.
func readSecretValueFromFile() (string, *ce.CustomError) {
	data, err := os.ReadFile(SecretInputFile)
	if err != nil {
		return "", &ce.CustomError{Title: shared.ErrorMessages[shared.ErrReadFile].Msg, Message: err.Error(), Code: shared.ErrReadFile}
	}
	value := strings.TrimSuffix(string(data), "\n")
	value = strings.TrimSuffix(value, "\r")
	return value, nil
}

func outputData(data map[string]interface{}, suppress bool) *ce.CustomError {
	if SecretField != "" {
		val, found := data[SecretField]
		if !found {
			title := "ReadSecret error"
			message := fmt.Sprintf("Field %s not found", SecretField)
			return &ce.CustomError{Title: title, Message: message, Code: shared.ErrFieldNotFound}
		}
		if suppress {
			return nil
		}
		if shared.OutputFormat == "json" {
			out := map[string]interface{}{SecretField: val}
			json.NewEncoder(os.Stdout).Encode(out)
		} else {
			fmt.Printf("%v\n", val)
		}
		return nil
	}

	if suppress {
		return nil
	}

	if shared.OutputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(data); err != nil {
			return &ce.CustomError{Title: "JSON encoding error", Message: err.Error(), Code: shared.ErrExtractData}
		}
	} else {
		for k, v := range data {
			fmt.Printf("%s : %v\n", hftx.Green(k), v)
		}
	}
	return nil
}

// classifyReadError maps an error returned by vaultlib's KV read path onto a
// CustomError carrying the matching vclt error code, so that Die() can translate
// it into a meaningful, POSIX-safe exit status (e.g. a sealed vault exits with
// ErrVaultSealed rather than an indistinguishable generic 1).
//
// vaultlib normalises its failures into recognisable wrapped strings; we match
// on those. Order matters: the sealed case must be tested before the more
// general "unavailable" one, since vaultlib reports "sealed or unavailable".
func classifyReadError(err error) *ce.CustomError {
	msg := err.Error()
	var code int
	switch {
	case strings.Contains(msg, "sealed"):
		code = shared.ErrVaultSealed
	case strings.Contains(msg, "unavailable"), strings.Contains(msg, "connection refused"),
		strings.Contains(msg, "no such host"):
		code = shared.ErrVaultUnavailable
	case strings.Contains(msg, "unauthorized"), strings.Contains(msg, "invalid Vault token"),
		strings.Contains(msg, "permission denied"):
		code = shared.ErrVaultInvalidAuth
	case strings.Contains(msg, "does not exist"):
		code = shared.ErrInvalidPath
	default:
		code = shared.ErrReadSecret
	}
	return &ce.CustomError{Title: shared.ErrorMessages[code].Msg, Message: msg, Code: code}
}

// findLatestAvailableVersion :
// In a KV engine version 2+, multiple versions of a secret may exist.
// We may need to find the latest recorded version
//func findLatestAvailableVersion(client *api.Client, metaPath string) (int, *ce.CustomError) {
//	meta, err := client.Logical().Read(metaPath)
//	if err != nil || meta == nil {
//		return 0, &ce.CustomError{Title: "Unable to fetch metadata", Message: err.Error()}
//	}
//
//	rawVersions, ok := meta.Data["versions"].(map[string]interface{})
//	if !ok {
//		return 0, &ce.CustomError{Title: "Version metadata not found", Message: err.Error()}
//	}
//
//	var available []int
//	for verStr, vmetaAny := range rawVersions {
//		vmeta, ok := vmetaAny.(map[string]interface{})
//		if !ok {
//			continue
//		}
//		if destroyed, _ := vmeta["destroyed"].(bool); destroyed {
//			continue
//		}
//		if ver, err := strconv.Atoi(verStr); err == nil {
//			available = append(available, ver)
//		}
//	}
//
//	if len(available) == 0 {
//		return 0, nil
//	}
//
//	sort.Sort(sort.Reverse(sort.IntSlice(available)))
//	return available[0], nil
//}
