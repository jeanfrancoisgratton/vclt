// Hashicorp Vault Client (vclt)
// src/cmd/root.go

package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"vclt/shared"

	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vclt",
	Short: "Hashicorp Vault client",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(hftx.White("vclt 2.4.0 (2026.07.04), Go version = v" + strings.TrimPrefix(runtime.Version(), "go")))
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

	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().BoolVarP(&shared.QuietOutput, "quiet", "q", false, "Display output to stdout")
	rootCmd.PersistentFlags().BoolVarP(&shared.DebugMode, "debug", "d", false, "Debug mode")
	rootCmd.PersistentFlags().StringVarP(&shared.VaultAuthToken, "token", "t", "", "Vault token (or use VAULT_TOKEN)")
	rootCmd.PersistentFlags().StringVarP(&shared.VaultServerAddress, "address", "a", "", "Vault server address (or use VAULT_ADDR)")
	//rootCmd.PersistentFlags().StringVarP(&shared.OutputFormat, "output", "o", "text", "Output format: text|json")
}
