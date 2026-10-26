// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/policies/types.go
// Original timestamp: 2026/06/20 20:14:11

package policies

// Cleartext controls whether Backup/Restore read and write their file in
// plain JSON/HCL (true) or in encoded form (false, the default) -- same
// flag and same default as kv.Cleartext.
var Cleartext = false
