// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/policies_commands_test.go

package cmd

import (
	"testing"

	"vclt/policies"

	"github.com/spf13/cobra"
)

// TestPoliciesBackupRestoreAllFlagRegistered checks that --all is
// registered on both commands with no shorthand: -a is already bound to
// the global --address flag, which every command inherits as a persistent
// flag, so reusing it here would collide.
func TestPoliciesBackupRestoreAllFlagRegistered(t *testing.T) {
	cases := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"policy backup", policiesBackupCmd},
		{"policy restore", policiesRestoreCmd},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := tc.cmd.PersistentFlags().Lookup("all")
			if f == nil {
				t.Fatalf("expected --all to be registered on %s", tc.name)
			}
			if f.Shorthand != "" {
				t.Errorf("--all has shorthand -%s, want no shorthand (it would collide with -a/--address)", f.Shorthand)
			}
			if f.Value.Type() != "bool" {
				t.Errorf("--all is type %s, want bool", f.Value.Type())
			}
		})
	}
}

// TestPoliciesBackupRestoreArgsValidation exercises policiesBackupCmd's and
// policiesRestoreCmd's dynamic-arity Args function: POLICY_NAME + FILE (2
// args) normally, dropped to a single FILE/BACKUP_FILE argument (1 arg)
// when --all is set, since --all's bundle has no single POLICY_NAME to
// take a positional slot.
func TestPoliciesBackupRestoreArgsValidation(t *testing.T) {
	t.Cleanup(func() { policies.AllPolicies = false })

	oneArg := []string{"backup-file.json"}
	twoArgs := []string{"my-policy", "backup-file.json"}

	cases := []struct {
		name    string
		cmd     *cobra.Command
		all     bool
		args    []string
		wantErr bool
	}{
		{name: "backup: plain needs 2 args", cmd: policiesBackupCmd, args: twoArgs},
		{name: "backup: plain rejects 1 arg", cmd: policiesBackupCmd, args: oneArg, wantErr: true},
		{name: "backup: --all needs 1 arg", cmd: policiesBackupCmd, all: true, args: oneArg},
		{name: "backup: --all rejects 2 args", cmd: policiesBackupCmd, all: true, args: twoArgs, wantErr: true},
		{name: "restore: plain needs 2 args", cmd: policiesRestoreCmd, args: twoArgs},
		{name: "restore: plain rejects 1 arg", cmd: policiesRestoreCmd, args: oneArg, wantErr: true},
		{name: "restore: --all needs 1 arg", cmd: policiesRestoreCmd, all: true, args: oneArg},
		{name: "restore: --all rejects 2 args", cmd: policiesRestoreCmd, all: true, args: twoArgs, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			policies.AllPolicies = tc.all

			err := tc.cmd.Args(tc.cmd, tc.args)
			if tc.wantErr && err == nil {
				t.Error("expected an error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
