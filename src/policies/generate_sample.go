// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/generate_sample.go
// Original timestamp: 2026/06/28 11:22:13

package policies

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"vclt/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

const samplePolicyJSON = `{
  "path": {
    "myengine1/*": {
      "capabilities": ["create", "read", "update", "delete", "list"]
    },
    "myengine2/data/*": {
      "capabilities": ["create", "read", "update", "delete", "list"]
    },
    "myengine2/metadata/*": {
      "capabilities": ["read", "list", "delete"]
    },
    "myengine2/delete/*": {
      "capabilities": ["delete"]
    },
    "sys/mounts/*": {
      "capabilities": ["read", "list", "sudo"]
    }
  }
}`

const samplePolicyHCL = `path "myengine1/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

path "myengine2/data/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

path "myengine2/metadata/*" {
  capabilities = ["read", "list", "delete"]
}

path "myengine2/delete/*" {
  capabilities = ["delete"]
}

path "sys/mounts/*" {
  capabilities = ["read", "list", "sudo"]
}`

// GenerateSamplePolicy writes a sample Vault ACL policy file.
// The format (JSON or HCL) is determined by the filename extension:
// .json produces a JSON file; anything else (including no extension)
// produces an HCL file, and the .hcl extension is appended if absent.
func GenerateSamplePolicy(samplePolicyName string) *ce.CustomError {
	switch strings.ToLower(filepath.Ext(samplePolicyName)) {
	case ".json":
		return generateSample(samplePolicyName, samplePolicyJSON)
	default:
		if !strings.EqualFold(filepath.Ext(samplePolicyName), ".hcl") {
			samplePolicyName += ".hcl"
		}
		return generateSample(samplePolicyName, samplePolicyHCL)
	}
}

func generateSample(path, content string) *ce.CustomError {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return &ce.CustomError{Title: "Error writing sample policy file", Message: err.Error()}
	}
	if !shared.QuietOutput {
		fmt.Println("Generated sample policy file:", hftx.Green(path))
	}
	return nil
}
