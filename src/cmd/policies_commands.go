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
	Long:    `Allowed commands are { read | write | list | delete | sample }`,
	Example: `  vclt policy read myapp-policy
  vclt policy write myapp-policy myapp-policy.hcl
  vclt policy list
  vclt policy rm myapp-policy`,
}

var policiesReadCmd = &cobra.Command{
	Use:     "read POLICY_NAME",
	Aliases: []string{"get"},
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
	Aliases: []string{"ls", "show"},
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

func init() {
	rootCmd.AddCommand(policiesCmd)
	policiesCmd.AddCommand(policiesReadCmd, policiesWriteCmd, policiesLsCmd, policiesRmCmd, policiesGenerateCmd)

}
