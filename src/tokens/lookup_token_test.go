// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/tokens/lookup_token_test.go

package tokens

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	tkn "github.com/jeanfrancoisgratton/vaultlib/v2/tokens"
)

// TestSaveTokenInfo2file mirrors TestSaveTokenInfoToFile but for the
// TokenInfo shape returned by LookupToken/LookupSelf (a different struct
// with its own JSON tags), saved by a separate function in lookup_token.go.
func TestSaveTokenInfo2file(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "prod-token-lookup")

	original := &tkn.TokenInfo{
		ID:          "hvs.notarealtoken",
		Accessor:    "accessor-id",
		DisplayName: "ci-deploy",
		TTL:         3600,
		Policies:    []string{"default"},
		Orphan:      true,
		Renewable:   true,
	}

	if err := saveTokenInfo2file(target, original); err != nil {
		t.Fatalf("saveTokenInfo2file failed: %v", err)
	}

	savedPath := target + ".json"
	info, statErr := os.Stat(savedPath)
	if statErr != nil {
		t.Fatalf("expected %s to exist: %v", savedPath, statErr)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file permissions = %o, want 600", perm)
	}

	raw, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatal(err)
	}
	var loaded tkn.TokenInfo
	if err := json.Unmarshal(raw, &loaded); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}
	if !reflect.DeepEqual(*original, loaded) {
		t.Errorf("round trip mismatch:\nsaved:  %+v\nloaded: %+v", *original, loaded)
	}
}
