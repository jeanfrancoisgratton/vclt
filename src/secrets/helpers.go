// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/secrets/helpers.go
// Original timestamp: 2026/06/14 13:33:05

package secrets

import (
	"encoding/json"
	"fmt"
	"os"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"vclt/shared"
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
			fmt.Printf("%s : %v\n", hftx.Green(k), v)
		}
	}
	return nil
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
