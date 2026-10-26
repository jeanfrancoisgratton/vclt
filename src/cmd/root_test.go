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

// TestKvWriteJSONFlagRegistered checks that `kv write` has a --json
// boolean flag with no shorthand.
func TestKvWriteJSONFlagRegistered(t *testing.T) {
	jsonFlag := kvWriteCmd.PersistentFlags().Lookup("json")
	if jsonFlag == nil {
		t.Fatal("expected --json to be registered on kv write")
	}
	if jsonFlag.Shorthand != "" {
		t.Errorf("--json has shorthand -%s, want no shorthand", jsonFlag.Shorthand)
	}
	if jsonFlag.Value.Type() != "bool" {
		t.Errorf("--json is type %s, want bool", jsonFlag.Value.Type())
	}
}

// TestKvWriteArgsValidation exercises kvWriteCmd's dynamic-arity Args
// function directly: VALUE is a 4th positional arg normally, dropped to 3
// when --in, --batch, or --json is set (the value(s) then come from a file
// instead), and any two of --in/--batch/--json together must be rejected
// outright since they disagree about what the 3rd positional argument means.
func TestKvWriteArgsValidation(t *testing.T) {
	reset := func() {
		kv.SecretInputFile = ""
		kv.BatchMode = false
		kv.JSONMode = false
	}
	t.Cleanup(reset)

	fourArgs := []string{"engine", "path", "key", "value"}
	threeArgs := []string{"engine", "path", "key-or-file"}

	cases := []struct {
		name    string
		in      bool
		batch   bool
		json    bool
		args    []string
		wantErr bool
	}{
		{name: "plain write needs 4 args", args: fourArgs},
		{name: "plain write rejects 3 args", args: threeArgs, wantErr: true},
		{name: "--in needs 3 args", in: true, args: threeArgs},
		{name: "--in rejects 4 args", in: true, args: fourArgs, wantErr: true},
		{name: "--batch needs 3 args", batch: true, args: threeArgs},
		{name: "--batch rejects 4 args", batch: true, args: fourArgs, wantErr: true},
		{name: "--json needs 3 args", json: true, args: threeArgs},
		{name: "--json rejects 4 args", json: true, args: fourArgs, wantErr: true},
		{name: "--in and --batch together are always rejected", in: true, batch: true, args: threeArgs, wantErr: true},
		{name: "--in and --json together are always rejected", in: true, json: true, args: threeArgs, wantErr: true},
		{name: "--batch and --json together are always rejected", batch: true, json: true, args: threeArgs, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset()
			if tc.in {
				kv.SecretInputFile = "some-file"
			}
			kv.BatchMode = tc.batch
			kv.JSONMode = tc.json

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

// TestGroupCommandRejectsUnknownSubcommand guards a Cobra pitfall: a command
// with no Run/RunE of its own is "not runnable", and Cobra's execute() skips
// Args validation entirely for those, returning flag.ErrHelp -- which the
// top-level Execute() treats as "just show help", exiting 0. That silently
// swallowed a typo'd subcommand (e.g. `vclt kv bogus`) as if it had
// succeeded, which is exactly the failure mode a calling shell script can't
// detect via $?. Each group command below was given a Run (so it's runnable)
// specifically so its NoArgs Args validator actually gets evaluated: a
// genuine bare invocation (no args) must still exit 0 (help), but a stray/
// unmatched token must now be a real, non-zero-exit error.
func TestGroupCommandRejectsUnknownSubcommand(t *testing.T) {
	// Cobra's Execute()/ExecuteC() always run starting from the true root
	// regardless of which command in the tree Execute() is called on, using
	// the root's own SetArgs -- so driving this through rootCmd (rather than
	// e.g. kvCmd directly) is required for SetArgs to actually take effect.
	groupNames := []string{"kv", "policy", "sys", "token", "admin", "completion"}

	run := func(t *testing.T, args []string) error {
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)
		t.Cleanup(func() {
			rootCmd.SetOut(nil)
			rootCmd.SetErr(nil)
			rootCmd.SetArgs(nil)
		})
		rootCmd.SetArgs(args)
		return rootCmd.Execute()
	}

	for _, name := range groupNames {
		t.Run(name, func(t *testing.T) {
			if err := run(t, []string{name, "bogus-subcommand"}); err == nil {
				t.Error("expected an error for an unmatched subcommand, got nil (would exit 0)")
			}

			if err := run(t, []string{name}); err != nil {
				t.Errorf("expected a bare invocation (no args) to still just show help, got error: %v", err)
			}
		})
	}
}
