// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/secret_commands.go
// Original timestamp: 2026/06/14 12:54:55

package cmd

import (
	"fmt"
	"os"

	hftfx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
	"vclt/secrets"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "secret management subcommands",
	Long:  `Allowed commands are { read | write | list }`,
}

var secretsReadCmd = &cobra.Command{
	Use:     "vclt read KV_engine secret_path [-f field]",
	Aliases: []string{"get"},
	Short:   "Read the 'secret_path' secret from the 'KV_engine' secret engine",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if kvreadErr := secrets.ReadSecrets(args[0], args[1]); kvreadErr != nil {
			fmt.Println(hftfx.SkullBonesSign(kvreadErr.Error()))
			os.Exit(1)
		}
	},
}

var secretsWriteCmd = &cobra.Command{
	Use:     "vclt write KV_engine secret_path [-f field]",
	Aliases: []string{"put"},
	Short:   "Write the 'secret_path' secret to the 'KV_engine' secret engine",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(1)
	},
}

var secretsLsCmd = &cobra.Command{
	Use:     "vclt list KV_engine secret_path",
	Aliases: []string{"ls", "show"},
	Short:   "List the 'secret_path' secrets in the 'KV_engine' secret engine",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(secretsCmd)
	secretsCmd.AddCommand(secretsLsCmd, secretsReadCmd, secretsWriteCmd)

	secretsCmd.PersistentFlags().StringVarP(&secrets.SecretMountPath, "mount", "m", "", "KV v2 mount path (required)")
	secretsCmd.PersistentFlags().IntVarP(&secrets.SecretVersion, "version", "v", 0, "Secret version (0 = latest available)")
	secretsCmd.PersistentFlags().StringVarP(&secrets.SecretField, "field", "f", "", "Specific field to display")
}
