// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/admin/types.go
// Original timestamp: 2026/06/18 06:31:58

package admin

type VaultRootKeysStruct struct {
	MinimumRequired int      `json:"minimumRequired,omitempty"`
	Shards          []string `json:"shards"`
	InitialRootKey  string   `json:"initialRootKey,omitempty"`
	Comments        string   `json:"comments,omitempty"`
}
