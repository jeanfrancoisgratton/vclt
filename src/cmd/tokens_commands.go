// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/tokens_commands.go
// Original timestamp: 2026/06/30 11:22:39

package cmd

import (
	"vclt/shared"
	"vclt/tokens"

	"github.com/spf13/cobra"
)

var tokensCmd = &cobra.Command{
	Use:     "token",
	Aliases: []string{"tokens"},
	Short:   "tokens management subcommands",
	Long:    `Allowed commands are { create | revoke | renew | lookup | self | accessors }`,
}

var tokenCreateCmd = &cobra.Command{
	Use:     "create TOKEN_NAME",
	Aliases: []string{"write"},
	Short:   "Create a token with the appropriate config policies with flags",
	Long: `By default, tokens are created as orphaned and renewable.
If no policies are specified the token will be bound to the default policy`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := tokens.NewClient()
		if err != nil {
			err.Die()
		}
		if tknErr := c.Create(args[0], true); tknErr != nil {
			tknErr.Die()
		}
	},
}

var tokenRevokeCmd = &cobra.Command{
	Use:     "revoke TOKEN_NAME",
	Aliases: []string{"remove", "delete"},
	Short:   "Permanently revoke a token and its children (if any)",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := tokens.NewClient()
		if err != nil {
			err.Die()
		}
		if tknErr := c.Revoke(args[0]); tknErr != nil {
			tknErr.Die()
		}
	},
}

var tokenRenewCmd = &cobra.Command{
	Use:   "renew TOKEN_NAME",
	Short: "Renew the token TOKEN_NAME. The -d flag sets the new lease duration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := tokens.NewClient()
		if err != nil {
			err.Die()
		}
		if tknErr := c.Renew(args[0]); tknErr != nil {
			tknErr.Die()
		}
	},
}

var tokenLookupCmd = &cobra.Command{
	Use:   "lookup TOKEN_NAME",
	Short: "Displays the info about the named token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := tokens.NewClient()
		if err != nil {
			err.Die()
		}
		if _, tknErr := c.LookupToken(args[0], true); tknErr != nil {
			tknErr.Die()
		}
	},
}

var tokenLookupSelfCmd = &cobra.Command{
	Use:   "self",
	Short: "Displays the info about the token currently in use",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := tokens.NewClient()
		if err != nil {
			err.Die()
		}
		if _, tknErr := c.LookupSelf(true); tknErr != nil {
			tknErr.Die()
		}
	},
}

var tokenListAccessorsCmd = &cobra.Command{
	Use:   "accessors",
	Short: "List all token accessors",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := tokens.NewClient()
		if err != nil {
			err.Die()
		}
		if tknErr := c.ListAccessors(); tknErr != nil {
			tknErr.Die()
		}
	},
}

func init() {
	rootCmd.AddCommand(tokensCmd)
	tokensCmd.AddCommand(tokenCreateCmd, tokenRevokeCmd, tokenRenewCmd, tokenLookupSelfCmd, tokenLookupCmd, tokenListAccessorsCmd)

	tokenCreateCmd.Flags().StringVarP(&tokens.TokenPolicies, "policies", "P", "", "comma-separated list of policies that should be used for the token")
	tokenCreateCmd.Flags().StringVarP(&tokens.TokenTTL, "ttl", "T", "1h", "token TTL (defaults to 1h)")
	tokenCreateCmd.Flags().BoolVarP(&tokens.TokenOrphaned, "orphaned", "o", true, "orphaned token")
	tokenCreateCmd.Flags().BoolVarP(&tokens.TokenRenewable, "renewable", "r", true, "renewable token")
	tokenCreateCmd.Flags().StringVarP(&tokens.TokenSavefile, "file", "f", "", "save token info to file")

	tokenRenewCmd.Flags().IntVarP(&tokens.TokenDuration, "duration", "d", 0, "token duration (defaults to 1h0m0s)")

	tokenLookupSelfCmd.Flags().StringVarP(&tokens.TokenSavefile, "file", "f", "", "save token info to file")
	tokenLookupSelfCmd.PersistentFlags().StringVarP(&shared.OutputFormat, "output", "o", "text", "Output format: text|json")

	tokenLookupCmd.Flags().StringVarP(&tokens.TokenSavefile, "file", "f", "", "save token info to file")
	tokenLookupCmd.PersistentFlags().StringVarP(&shared.OutputFormat, "output", "o", "text", "Output format: text|json")

}
