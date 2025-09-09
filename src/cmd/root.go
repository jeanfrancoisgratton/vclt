// Hashicorp Vault Client (vclt)
// src/cmd/root.go

package cmd

import (
	"fmt"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"vclt/env"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "vclt",
	Short:   "Hashicorp Vault client",
	Version: hf.White(fmt.Sprintf("DEV2.00.00-0-%s (2025.03.17)", runtime.GOARCH)),
}

var clCmd = &cobra.Command{
	Use:     "changelog",
	Aliases: []string{"cl"},
	Short:   "Shows changelog",
	Run: func(cmd *cobra.Command, args []string) {
		changelog()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.DisableAutoGenTag = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(clCmd)
	//rootCmd.AddCommand(envCmd, kvCmd, loginCmd, sysCmd)
	rootCmd.AddCommand(envCmd, loginCmd, sysCmd, kvCmd)

	rootCmd.PersistentFlags().StringVarP(&env.ConfigFile, "env", "e", "defaultEnv.json", "Default env configuration file; this is a per-user setting.")
	//rootCmd.PersistentFlags().BoolVarP(&sys.SysUseToken, "authtoken", "T", false, "Use the 'auth token' authentication method")
}

func changelog() {
	//fmt.Printf("\x1b[2J")
	fmt.Printf("\x1bc")

	fmt.Println("CHANGELOG")
	fmt.Println()
	fmt.Println()

	fmt.Print(`
VERSION		DATE			COMMENT
-------		----			-------
DEV2.00.00	2025.03.17		Rewrite of the client, now using the hcpVaultLib package
1.02.00		2024.07.08		Added env create, to create new environment files from the CLI, updated to GO 1.22.5
1.01.00		2024.07.02		Completed kv get, kv lse, kv lsf
1.00.00		2024.06.27		Initial version
`)
}
