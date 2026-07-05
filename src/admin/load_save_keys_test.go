// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/admin/load_save_keys_test.go

package admin

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestRootKeysRoundTrip saves a struct and loads it back, asserting we get
// the same data. Because save/load build their path from $HOME, we point
// HOME at a throwaway directory with t.Setenv — the framework restores the
// real value when the test ends, so nothing leaks into your actual
// ~/.config/JFG/vclt.
func TestRootKeysRoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, ".config", "JFG", "vclt"), 0700); err != nil {
		t.Fatal(err)
	}

	original := VaultRootKeysStruct{
		MinimumRequired: 3,
		Shards:          []string{"shard-one", "shard-two", "shard-three"},
		InitialRootKey:  "hvs.notarealtoken",
		Comments:        "unit test fixture",
	}

	if cerr := original.saveRootKeys("keys.json"); cerr != nil {
		t.Fatalf("saveRootKeys failed: %v", cerr)
	}

	// The file holds secrets; it must not be group/world readable.
	info, err := os.Stat(filepath.Join(home, ".config", "JFG", "vclt", "keys.json"))
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("root keys file has permissions %o, want 600", perm)
	}

	var loaded VaultRootKeysStruct
	if cerr := loaded.loadRootKeys("keys.json"); cerr != nil {
		t.Fatalf("loadRootKeys failed: %v", cerr)
	}
	if !reflect.DeepEqual(original, loaded) {
		t.Errorf("round trip mismatch:\nsaved:  %+v\nloaded: %+v", original, loaded)
	}
}

func TestLoadRootKeysMissingFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	var vk VaultRootKeysStruct
	if cerr := vk.loadRootKeys("does-not-exist.json"); cerr == nil {
		t.Error("expected an error for a missing file, got nil")
	}
}
