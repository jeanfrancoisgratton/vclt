// Hashicorp Vault Client (vclt)
// src/cmd/root.go

package cmd

import (
	"os"
	"runtime"
	"strings"

	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
	"vclt/shared"

	"vclt/env"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "vclt",
	Short:   "Hashicorp Vault client",
	Version: hftx.White("DEV2.00.00 (2026.06.10), Go version = v" + strings.TrimPrefix(runtime.Version(), "go")),
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
	//rootCmd.AddCommand(envCmd, kvCmd, loginCmd, sysCmd)
	rootCmd.AddCommand(envCmd, loginCmd, sysCmd, kvCmd)

	rootCmd.PersistentFlags().StringVarP(&env.ConfigFile, "env", "e", "defaultEnv.json", "Default env configuration file; this is a per-user setting")
	rootCmd.PersistentFlags().BoolVarP(&shared.QuietOutput, "quiet", "q", false, "Display output to stdout")
	rootCmd.PersistentFlags().BoolVarP(&shared.DebugMode, "debug", "d", false, "Debug mode")

	//rootCmd.PersistentFlags().BoolVarP(&sys.SysUseToken, "authtoken", "T", false, "Use the 'auth token' authentication method")
}
