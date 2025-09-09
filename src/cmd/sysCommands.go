// vclt
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename : sysCommands.go
// Original timestamp : 2024/06/30 18:44

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"vclt/sys"
)

var sysCmd = &cobra.Command{
	Use:     "sys",
	Example: "vclt sys { login }",
	Short:   "system commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You need to specify one of the following subcommand: login")
		os.Exit(0)
	},
}

var loginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"auth"},
	Example: "vclt sys login { -t | -u }",
	Short:   "Login to Vault, using the token or userpass method",
	Long: `The -t (auth token) method will take precedence if both or none are specified.
	If -u is used, the username & password will be fetched from the environment file`,
	Run: func(cmd *cobra.Command, args []string) {
		if sys.SysUsePassword {
			if err := sys.LoginUserPass(); err != nil {
				fmt.Println(err.Error())
			}
		} else {
			if err := sys.LoginToken(); err != nil {
				fmt.Println(err.Error())
			}
		}
	},
}

func init() {
	sysCmd.AddCommand(loginCmd)

	loginCmd.PersistentFlags().BoolVarP(&sys.Verbose, "verbose", "v", false, "Quiet (no) output")
	loginCmd.PersistentFlags().BoolVarP(&sys.SysStoreToken, "storetkn", "s", false, "Store token in environment file")
	loginCmd.PersistentFlags().BoolVarP(&sys.SysUsePassword, "userpass", "U", false, "Use the 'userpass' authentication method")
}
