// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/policies_commands.go
// Original timestamp: 2026/06/20 20:16:25

package cmd

import (
	"vclt/policies"

	"github.com/spf13/cobra"
)

var policiesCmd = &cobra.Command{
	Use:     "policy",
	Aliases: []string{"policies"},
	Short:   "policies management subcommands",
	Long:    `Allowed commands are { read | write | list | delete | sample | backup | restore }`,
	Example: `  vclt policy read myapp-policy
  vclt policy write myapp-policy myapp-policy.hcl
  vclt policy list
  vclt policy rm myapp-policy
  vclt policy backup myapp-policy myapp-policy-backup.json
  vclt policy backup --all policies-backup.json`,
	// See the matching comment on kvCmd in kv_commands.go: Run makes this
	// command runnable so Args (NoArgs) actually gets evaluated instead of
	// Cobra silently treating a stray subcommand as "show help", exit 0.
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
	Args: cobra.NoArgs,
}

var policiesReadCmd = &cobra.Command{
	Use:     "read POLICY_NAME",
	Aliases: []string{"get", "show"},
	Short:   "Read the POLICY_NAME policies",
	Args:    cobra.ExactArgs(1),
	Example: `  vclt policy read myapp-policy
  vclt policy get myapp-policy`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := policies.NewClient()
		if err != nil {
			err.Die()
		}
		if _, polErr := c.Read(args[0], true); polErr != nil {
			polErr.Die()
		}
	},
}

var policiesWriteCmd = &cobra.Command{
	Use:     "write POLICY_NAME POLICY_FILE",
	Aliases: []string{"put"},
	Short:   "Write the POLICY_NAME policies from the POLICY_FILE file",
	Args:    cobra.ExactArgs(2),
	Example: `  vclt policy write myapp-policy myapp-policy.hcl
  vclt policy put myapp-policy myapp-policy.json`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := policies.NewClient()
		if err != nil {
			err.Die()
		}
		if polErr := c.Write(args[0], args[1]); polErr != nil {
			polErr.Die()
		}
	},
}

var policiesLsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the policies",
	Args:    cobra.ExactArgs(0),
	Example: `  vclt policy list
  vclt policy ls`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := policies.NewClient()
		if err != nil {
			err.Die()
		}
		if _, polErr := c.List(true); polErr != nil {
			polErr.Die()
		}
	},
}

var policiesRmCmd = &cobra.Command{
	Use:     "rm POLICY_NAME1...POLICY_NAME2...POLICY_NAMEx",
	Aliases: []string{"delete"},
	Short:   "Delete one or many policies",
	Args:    cobra.MinimumNArgs(1),
	Example: `  vclt policy rm myapp-policy
  vclt policy delete myapp-policy old-policy another-policy`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := policies.NewClient()
		if err != nil {
			err.Die()
		}
		if polErr := c.Delete(args); polErr != nil {
			polErr.Die()
		}
	},
}

var policiesGenerateCmd = &cobra.Command{
	Use:     "generate FILENAME",
	Aliases: []string{"gen", "sample"},
	Short:   "Generate a sample policy file",
	Long: `Generate a sample policy file in the FILENAME file that can be used as the basis for a new policy.
This is quite useful to understand how to write a policy in JSON or HCL format. It will also do a syntax check on the file before submitting it to Vault.`,
	Example: `  vclt policy generate sample-policy.hcl
  vclt policy gen sample-policy.json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if polErr := policies.GenerateSamplePolicy(args[0]); polErr != nil {
			polErr.Die()
		}
	},
}

var policiesBackupCmd = &cobra.Command{
	Use:     "backup POLICY_NAME BACKUP_FILE[.json|.hcl]",
	Aliases: []string{"dump"},
	Short:   "Backup one ACL policy, or every policy with --all",
	Long: `Dumps POLICY_NAME's rules to BACKUP_FILE, exactly as Vault returned them.

With --all, POLICY_NAME is omitted in favor of a single BACKUP_FILE argument
holding every ACL policy on the server in one call. The bundle's own format
is selected by BACKUP_FILE's extension: .hcl writes HCL, anything else
writes JSON (a single policy's own file has no such choice to make, since
its rules are dumped verbatim in whatever format Vault already stores them
in). Vault's built-in, virtual policies (e.g. "root") have no stored rules
to dump and are skipped rather than failing the whole backup.

By default the resulting file is encoded (via helperFunctions/v5 file
encoding) so policy rules are not left in plaintext on disk; pass
--cleartext to keep the raw JSON/HCL instead.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if policies.AllPolicies {
			return cobra.ExactArgs(1)(cmd, args)
		}
		return cobra.ExactArgs(2)(cmd, args)
	},
	Example: `  vclt policy backup myapp-policy myapp-policy-backup.json
  vclt policy backup --all policies-backup.json
  vclt policy backup --all policies-backup.hcl
  vclt policy dump --all policies-backup.json -c`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := policies.NewClient()
		if err != nil {
			err.Die()
		}
		if policies.AllPolicies {
			if polErr := c.BackupAll(args[0]); polErr != nil {
				polErr.Die()
			}
			return
		}
		if polErr := c.Backup(args[0], args[1]); polErr != nil {
			polErr.Die()
		}
	},
}

var policiesRestoreCmd = &cobra.Command{
	Use:     "restore POLICY_NAME BACKUP_FILE[.json|.hcl]",
	Aliases: []string{"import"},
	Short:   "Restore one ACL policy, or every policy with --all",
	Long: `Recreates POLICY_NAME from BACKUP_FILE, a file produced by 'policy backup'.
This is validated the same way 'policy write' validates its file before
submitting it to Vault.

With --all, POLICY_NAME is omitted in favor of a single BACKUP_FILE argument
holding every ACL policy from a bundle produced by 'policy backup --all'.
Unlike the single-policy case, a bundle's policies are restored exactly as
captured, not re-validated, since they already passed Vault's own
acceptance the first time they were written. The format is selected by
extension, same as backup.

By default the file is assumed to be encoded; pass --cleartext if the
backup was taken with -c.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if policies.AllPolicies {
			return cobra.ExactArgs(1)(cmd, args)
		}
		return cobra.ExactArgs(2)(cmd, args)
	},
	Example: `  vclt policy restore myapp-policy myapp-policy-backup.json
  vclt policy restore --all policies-backup.json
  vclt policy import --all policies-backup.hcl`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := policies.NewClient()
		if err != nil {
			err.Die()
		}
		if policies.AllPolicies {
			if polErr := c.RestoreAll(args[0]); polErr != nil {
				polErr.Die()
			}
			return
		}
		if polErr := c.Restore(args[0], args[1]); polErr != nil {
			polErr.Die()
		}
	},
}

func init() {
	rootCmd.AddCommand(policiesCmd)
	policiesCmd.AddCommand(policiesReadCmd, policiesWriteCmd, policiesLsCmd, policiesRmCmd, policiesGenerateCmd, policiesBackupCmd, policiesRestoreCmd)

	policiesBackupCmd.PersistentFlags().BoolVarP(&policies.Cleartext, "cleartext", "c", false, "Backup cleartext (default: encoded)")
	policiesRestoreCmd.PersistentFlags().BoolVarP(&policies.Cleartext, "cleartext", "c", false, "Restore cleartext (default: encoded)")
	// No shorthand: -a is already bound to the global --address flag, which
	// every command (including this one) inherits as a persistent flag --
	// reusing -a here would collide with it.
	policiesBackupCmd.PersistentFlags().BoolVar(&policies.AllPolicies, "all", false, "Backup every policy in one file instead of a single named policy")
	policiesRestoreCmd.PersistentFlags().BoolVar(&policies.AllPolicies, "all", false, "Restore every policy from one file instead of a single named policy")
}
