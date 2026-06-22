// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/policies_commands.go
// Original timestamp: 2026/06/20 20:16:25

package cmd

import (
	"fmt"
	"os"

	hftfx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
	"vclt/policies"
)

var policiesCmd = &cobra.Command{
	Use:   "policies",
	Short: "policies management subcommands",
	Long:  `Allowed commands are { read | write | list | delete }`,
}

var policiesReadCmd = &cobra.Command{
	Use:     "read POLICY_NAME",
	Aliases: []string{"get"},
	Short:   "Read the POLICY_NAME policies",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, polreadErr := policies.ReadPolicy(args[0], true); polreadErr != nil {
			fmt.Println(hftfx.SkullBonesSign(polreadErr.Error()))
			os.Exit(1)
		}
	},
}

var policiesWriteCmd = &cobra.Command{
	Use:     "write POLICY_NAME POLICY_FILE",
	Aliases: []string{"put"},
	Short:   "Write the POLICY_NAME policies from the POLICY_FILE file",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if polwriteErr := policies.WritePolicy(args[0], args[1]); polwriteErr != nil {
			fmt.Println(hftfx.SkullBonesSign(polwriteErr.Error()))
			os.Exit(1)
		}
	},
}

var policiesLsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "show"},
	Short:   "List the policies",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := policies.List(true); err != nil {
			fmt.Println(hftfx.SkullBonesSign(err.Error()))
		}
	},
}

var policiesRmCmd = &cobra.Command{
	Use:     "rm POLICY_NAME1...POLICY_NAME2...POLICY_NAMEx",
	Aliases: []string{"delete"},
	Short:   "Delete one or many policies",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := policies.DeletePolicies(args); err != nil {
			fmt.Println(hftfx.SkullBonesSign(err.Error()))
		}
	},
}

func init() {
	rootCmd.AddCommand(policiesCmd)
	policiesCmd.AddCommand(policiesReadCmd, policiesWriteCmd, policiesLsCmd, policiesRmCmd)

}
