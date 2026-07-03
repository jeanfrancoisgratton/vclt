// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/policies_commands.go
// Original timestamp: 2026/06/20 20:16:25

package cmd

import (
	"fmt"
	"os"
	"vclt/tokens"

	hftfx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
)

var tokensCmd = &cobra.Command{
	Use:     "token",
	Aliases: []string{"tokens"},
	Short:   "tokens management subcommands",
	Long:    `Allowed commands are { read | write | list | revoke | lookup | lookupself }`,
}

var tokenCreateCmd = &cobra.Command{
	Use:     "create TOKEN_NAME",
	Aliases: []string{"write"},
	Short:   "Create a token with the appropriate config policies with flags",
	Long: `By default, tokens are created as orphaned and renewable.
If no policies are specified the token will be bound to the default policy`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if tknWriteErr := tokens.CreateToken(args[0], true); tknWriteErr != nil {
			fmt.Println(hftfx.SkullBonesSign(tknWriteErr.Error()))
			os.Exit(1)
		}
	},
}

var tokenRevokeCmd = &cobra.Command{
	Use:     "revoke TOKEN_NAME",
	Aliases: []string{"remove", "delete"},
	Short:   "Permanently revoke a token and its children (if any)",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if tknWriteErr := tokens.RevokeToken(args[0]); tknWriteErr != nil {
			fmt.Println(hftfx.SkullBonesSign(tknWriteErr.Error()))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tokensCmd)
	tokensCmd.AddCommand(tokenCreateCmd, tokenRevokeCmd)

	tokenCreateCmd.Flags().StringVarP(&tokens.TokenPolicies, "policies", "P", "", "comma-separated list of policies that should be used for the token")
	tokenCreateCmd.Flags().StringVarP(&tokens.TokenTTL, "ttl", "T", "1h", "token TTL (defaults to 1hour")
	tokenCreateCmd.Flags().BoolVarP(&tokens.TokenOrphaned, "orphaned", "o", true, "orphaned token")
	tokenCreateCmd.Flags().BoolVarP(&tokens.TokenRenewable, "renewable", "r", true, "renewable token")
	tokenCreateCmd.Flags().StringVarP(&tokens.SaveTokenToFile, "file", "f", "", "save token to file")
	//_ = tokenCreateCmd.MarkFlagRequired("policies")

}
