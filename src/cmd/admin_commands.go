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
	"vclt/admin"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "server admin subcommands",
	Long:  `Allowed commands are { setrootkeys | seal | unseal }`,
}

var adminSetKeysCmd = &cobra.Command{
	Use:   "setrootkeys",
	Short: "Store the necessary root keys needed to unseal a Vault in secure storage",
	Run: func(cmd *cobra.Command, args []string) {
		if admErr := admin.SetRootKeys(); admErr != nil {
			fmt.Println(hftfx.SkullBonesSign(admErr.Error()))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)
	adminCmd.AddCommand(adminSetKeysCmd)

}
