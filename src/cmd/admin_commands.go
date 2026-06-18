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
	Long: `It will create a root keys file in $HOME/.config/JFG/vclt.
If no filename is provided, $HOME/.config/JFG/vclt/rootkeys.json will be created`,
	Run: func(cmd *cobra.Command, args []string) {
		rkfile := "rootkeys.json"
		if len(args) > 0 {
			rkfile = args[0]
		}
		if admErr := admin.SetRootKeys(rkfile); admErr != nil {
			fmt.Println(hftfx.SkullBonesSign(admErr.Error()))
			os.Exit(1)
		}
	},
}

var adminSealCmd = &cobra.Command{
	Use:   "seal",
	Short: "Seal your Hashicorp Vault",
	Run: func(cmd *cobra.Command, args []string) {
		if admErr := admin.Seal(); admErr != nil {
			fmt.Println(hftfx.SkullBonesSign(admErr.Error()))
			os.Exit(1)
		}
	},
}

var adminUnsealCmd = &cobra.Command{
	Use:   "unseal [root key file]",
	Short: "Unseal your Hashicorp Vault",
	Long: `It will use the root keys file stored in $HOME/.config/JFG/vclt.
If none is provided, $HOME/.config/JFG/vclt/rootkeys.json will be used`,
	Run: func(cmd *cobra.Command, args []string) {
		rkfile := "rootkeys.json"
		if len(args) > 0 {
			rkfile = args[0]
		}
		if admErr := admin.Unseal(rkfile); admErr != nil {
			fmt.Println(hftfx.SkullBonesSign(admErr.Error()))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)
	adminCmd.AddCommand(adminSetKeysCmd, adminSealCmd, adminUnsealCmd)

}
