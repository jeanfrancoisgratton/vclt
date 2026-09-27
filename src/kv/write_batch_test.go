// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/write_batch_test.go

package kv

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"vclt/shared"
)

// TestParseBatchFile covers the delimiter rules (first '=' or whitespace
// wins, whichever comes first), comment/blank-line skipping, optional
// double-quote preservation of whitespace, and the "no delimiter found"
// error path.
func TestParseBatchFile(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "space delimiter",
			content: "username admin\n",
			want:    map[string]string{"username": "admin"},
		},
		{
			name:    "equals delimiter, value keeps internal spaces without quoting",
			content: "password=s3cr3t with spaces\n",
			want:    map[string]string{"password": "s3cr3t with spaces"},
		},
		{
			name:    "tab delimiter",
			content: "username\tadmin\n",
			want:    map[string]string{"username": "admin"},
		},
		{
			name:    "quoted value preserves leading/trailing whitespace",
			content: `quoted="  padded value  "` + "\n",
			want:    map[string]string{"quoted": "  padded value  "},
		},
		{
			name:    "value with its own '=' is kept whole (split happens once)",
			content: "conn=user=name\n",
			want:    map[string]string{"conn": "user=name"},
		},
		{
			name:    "blank lines and comments are skipped",
			content: "# a comment\n\nkey value\n\n# trailing comment\n",
			want:    map[string]string{"key": "value"},
		},
		{
			name:    "later duplicate key wins",
			content: "key first\nkey second\n",
			want:    map[string]string{"key": "second"},
		},
		{
			name:    "line with no delimiter is an error",
			content: "justatoken\n",
			wantErr: true,
		},
		{
			name:    "line starting with a delimiter (empty key) is an error",
			content: "=value\n",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "batch.txt")
			if err := os.WriteFile(path, []byte(tc.content), 0600); err != nil {
				t.Fatal(err)
			}

			got, err := parseBatchFile(path)
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
				t.Errorf("parseBatchFile() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestParseBatchFileMissing(t *testing.T) {
	_, err := parseBatchFile(filepath.Join(t.TempDir(), "does-not-exist.txt"))
	if err == nil {
		t.Fatal("expected an error reading a nonexistent batch file, got nil")
	}
	if err.Code != shared.ErrReadFile {
		t.Errorf("error code = %d, want %d (ErrReadFile)", err.Code, shared.ErrReadFile)
	}
}
