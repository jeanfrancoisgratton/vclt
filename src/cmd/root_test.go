// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/root_test.go

package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"vclt/kv"
)

// TestVersionCommandOutput checks that `vclt version` prints the
// ldflags-injected buildVersion/buildDate (rather than, e.g., silently
// falling back to the "dev"/"unknown" zero values) — this is the command's
// only behavior, and it needs no Vault server.
func TestVersionCommandOutput(t *testing.T) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	versionCmd.Run(versionCmd, nil)

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "vclt v"+buildVersion) {
		t.Errorf("output %q does not contain the build version %q", out, buildVersion)
	}
	if !strings.Contains(out, buildDate) {
		t.Errorf("output %q does not contain the build date %q", out, buildDate)
	}
}

// TestKvReadOutFlagHasNoShorthandCollision guards the fix applied when the
// file-output flag was added to `kv read` (originally --file, renamed to
// --out): it must NOT take -f, since -f is already bound to --field on the
// same command. Regressing this would make one of the two flags silently
// shadow the other.
func TestKvReadOutFlagHasNoShorthandCollision(t *testing.T) {
	fieldFlag := kvReadCmd.PersistentFlags().ShorthandLookup("f")
	if fieldFlag == nil {
		t.Fatal("expected -f to be registered as a shorthand on kv read")
	}
	if fieldFlag.Name != "field" {
		t.Errorf("-f is bound to %q, want it bound to \"field\"", fieldFlag.Name)
	}

	outFlag := kvReadCmd.PersistentFlags().Lookup("out")
	if outFlag == nil {
		t.Fatal("expected --out to be registered on kv read")
	}
	if outFlag.Shorthand != "" {
		t.Errorf("--out has shorthand -%s, want no shorthand (it would collide with -f/--field)", outFlag.Shorthand)
	}
}

// TestKvWriteInFlagRegistered checks that `kv write` has an --in flag
// (the file-input counterpart to `kv read`'s --out), also with no
// shorthand, for the same reason as --out above -- kvWriteCmd doesn't
// currently bind -i to anything else, but pinning "no shorthand" keeps the
// two file flags symmetric and leaves -i free for future use.
func TestKvWriteInFlagRegistered(t *testing.T) {
	inFlag := kvWriteCmd.PersistentFlags().Lookup("in")
	if inFlag == nil {
		t.Fatal("expected --in to be registered on kv write")
	}
	if inFlag.Shorthand != "" {
		t.Errorf("--in has shorthand -%s, want no shorthand", inFlag.Shorthand)
	}
}

// TestKvWriteBatchFlagRegistered checks that `kv write` has a --batch
// boolean flag with no shorthand.
func TestKvWriteBatchFlagRegistered(t *testing.T) {
	batchFlag := kvWriteCmd.PersistentFlags().Lookup("batch")
	if batchFlag == nil {
		t.Fatal("expected --batch to be registered on kv write")
	}
	if batchFlag.Shorthand != "" {
		t.Errorf("--batch has shorthand -%s, want no shorthand", batchFlag.Shorthand)
	}
	if batchFlag.Value.Type() != "bool" {
		t.Errorf("--batch is type %s, want bool", batchFlag.Value.Type())
	}
}

// TestKvWriteArgsValidation exercises kvWriteCmd's dynamic-arity Args
// function directly: VALUE is a 4th positional arg normally, dropped to 3
// when --in or --batch is set (the value(s) then come from a file instead),
// and --in + --batch together must be rejected outright since they disagree
// about what the 3rd positional argument means.
func TestKvWriteArgsValidation(t *testing.T) {
	reset := func() {
		kv.SecretInputFile = ""
		kv.BatchMode = false
	}
	t.Cleanup(reset)

	fourArgs := []string{"engine", "path", "key", "value"}
	threeArgs := []string{"engine", "path", "key-or-file"}

	cases := []struct {
		name    string
		in      bool
		batch   bool
		args    []string
		wantErr bool
	}{
		{name: "plain write needs 4 args", args: fourArgs},
		{name: "plain write rejects 3 args", args: threeArgs, wantErr: true},
		{name: "--in needs 3 args", in: true, args: threeArgs},
		{name: "--in rejects 4 args", in: true, args: fourArgs, wantErr: true},
		{name: "--batch needs 3 args", batch: true, args: threeArgs},
		{name: "--batch rejects 4 args", batch: true, args: fourArgs, wantErr: true},
		{name: "--in and --batch together are always rejected", in: true, batch: true, args: threeArgs, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset()
			if tc.in {
				kv.SecretInputFile = "some-file"
			}
			kv.BatchMode = tc.batch

			err := kvWriteCmd.Args(kvWriteCmd, tc.args)
			if tc.wantErr && err == nil {
				t.Error("expected an error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
