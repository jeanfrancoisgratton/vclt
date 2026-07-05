// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/shared/errors_test.go

package shared

import (
	"testing"
)

// TestErrorMessagesComplete guards against adding a new error code to the
// const block and forgetting to add its entry in ErrorMessages: it walks
// every code from the first to the last and checks the map has it.
func TestErrorMessagesComplete(t *testing.T) {
	for code := ErrVaultAuthTokenMissing; code <= ErrVaultInvalidAuth; code++ {
		info, ok := ErrorMessages[code]
		if !ok {
			t.Errorf("error code %d has no entry in ErrorMessages", code)
			continue
		}
		if info.Int2StringCode == "" || info.Msg == "" {
			t.Errorf("error code %d has an empty name or description: %+v", code, info)
		}
	}
}

// TestErrorNamesUnique makes sure no two codes share the same string name,
// which would make log output ambiguous.
func TestErrorNamesUnique(t *testing.T) {
	seen := map[string]int{}
	for code, info := range ErrorMessages {
		if prev, dup := seen[info.Int2StringCode]; dup {
			t.Errorf("codes %d and %d share the same name %q", prev, code, info.Int2StringCode)
		}
		seen[info.Int2StringCode] = code
	}
}
