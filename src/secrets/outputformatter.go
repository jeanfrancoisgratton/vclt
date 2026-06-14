// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp : 2025/06/04 14:55
// Original filename : src/kv/outputformatter.go

package secrets

import (
	"encoding/json"
	"fmt"
	"os"

	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

func outputData(data map[string]interface{}, suppress bool) *ce.CustomError {
	if SecretField != "" {
		val, found := data[SecretField]
		if !found {
			title := "ReadSecret error"
			message := fmt.Sprintf("Field %s not found", SecretField)
			code := shared.ErrFieldNotFound
			if !shared.QuietOutput {
				fmt.Println(hftx.SkullBonesSign(title + " " + message))
			}
			return &ce.CustomError{Title: title, Message: message, Code: code}
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
			title := "JSON encoding error"
			message := err.Error()
			code := shared.ErrExtractData
			if !shared.QuietOutput {
				fmt.Println(hftx.SkullBonesSign(title + " " + message))
			}
			return &ce.CustomError{Title: title, Message: message, Code: code}
		}
	} else {
		for k, v := range data {
			fmt.Printf("%s: %v\n", k, v)
		}
	}
	return nil
}
