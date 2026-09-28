// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/write_json_test.go

package kv

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"vclt/shared"
)

// TestParseJSONFile covers the flat-object happy path plus the two failure
// modes: invalid JSON, and a top-level value that isn't an object.
func TestParseJSONFile(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "flat object of mixed value types",
			content: `{"username": "admin", "retries": 3, "enabled": true}`,
			want:    map[string]interface{}{"username": "admin", "retries": float64(3), "enabled": true},
		},
		{
			name:    "empty object",
			content: `{}`,
			want:    map[string]interface{}{},
		},
		{
			name:    "malformed JSON is an error",
			content: `{"username": `,
			wantErr: true,
		},
		{
			name:    "top-level array is not a flat object",
			content: `["username", "admin"]`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "secret.json")
			if err := os.WriteFile(path, []byte(tc.content), 0600); err != nil {
				t.Fatal(err)
			}

			got, err := parseJSONFile(path)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got fields: %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("parseJSONFile() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestParseJSONFileMissing(t *testing.T) {
	_, err := parseJSONFile(filepath.Join(t.TempDir(), "does-not-exist.json"))
	if err == nil {
		t.Fatal("expected an error reading a nonexistent JSON file, got nil")
	}
	if err.Code != shared.ErrReadFile {
		t.Errorf("error code = %d, want %d (ErrReadFile)", err.Code, shared.ErrReadFile)
	}
}
