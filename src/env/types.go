// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/env/types.go
// Original timestamp: 2026/06/14 07:30:25

package env

var ConfigFile string
var EnvName, VAddress, VUserName, VPassword, KVstorePath, EnvComments string

type Config_s struct {
	EnvironmentName string `json:"EnvironmentName"`
	VaultAddress    string `json:"VaultAddress"`
	VaultToken      string `json:"VaultToken,omitempty"`
	VaultUsername   string `json:"VaultUsername"`
	VaultPassword   string `json:"VaultPassword,omitempty"`
	KVEnginePath    string `json:"KVEnginePath"`
	Comments        string `json:"Comments,omitempty"`
}
type rootKeysFile struct {
	ServerAddr string   `json:"ServerAddr"`
	Comment    string   `json:"Comment,omitempty"`
	Keypart    []string `json:"Keypart"`
}
