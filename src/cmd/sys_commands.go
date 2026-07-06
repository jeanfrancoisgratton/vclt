// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/sys_commands.go
// Original timestamp: 2026/06/25 20:47:29

package cmd

import (
	"fmt"
	"os"

	"vclt/sys"

	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
)

var sysCmd = &cobra.Command{
	Use:   "sys",
	Short: "system management subcommands",
	Long:  `Allowed commands are { kvenable | kvdisable | listmounts }`,
}

var sysEnableKVCmd = &cobra.Command{
	Use:     "kvenable KVENGINE_NAME [-V version] [-D description]",
	Aliases: []string{"enablekv"},
	Short:   "Enable a KV secret engine",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := sys.NewClient()
		if err != nil {
			fmt.Println(hftx.SkullBonesSign(err.Error()))
			os.Exit(1)
		}
		if sysErr := c.EnableKVengine(args[0]); sysErr != nil {
			fmt.Println(hftx.SkullBonesSign(sysErr.Error()))
			os.Exit(1)
		}
	},
}

var sysDisableKVCmd = &cobra.Command{
	Use:     "kvdisable KVENGINE_NAME [-y]",
	Aliases: []string{"disablekv"},
	Short:   "Disable a KV secret engine",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := sys.NewClient()
		if err != nil {
			fmt.Println(hftx.SkullBonesSign(err.Error()))
			os.Exit(1)
		}
		if sysErr := c.DisableKVengine(args[0]); sysErr != nil {
			fmt.Println(hftx.SkullBonesSign(sysErr.Error()))
			os.Exit(1)
		}
	},
}

var listMountsCmd = &cobra.Command{
	Use:     "listmounts",
	Aliases: []string{"mounts"},
	Short:   "Lists all mounts (secret engines)",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := sys.NewClient()
		if err != nil {
			fmt.Println(hftx.SkullBonesSign(err.Error()))
			os.Exit(1)
		}
		if _, sysErr := c.ListMounts(true); sysErr != nil {
			fmt.Println(hftx.SkullBonesSign(sysErr.Error()))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(sysCmd)
	sysCmd.AddCommand(sysEnableKVCmd, sysDisableKVCmd, listMountsCmd)

	sysEnableKVCmd.Flags().IntVarP(&sys.KVEngineVersion, "version", "V", 2, "KV engine version")
	sysEnableKVCmd.Flags().StringVarP(&sys.KVEngineDescription, "desc", "D", "", "KV engine description")
	sysDisableKVCmd.Flags().BoolVarP(&sys.KVDisableConfirm, "yes", "y", false, "Disable confirmation (assume NO by default)")
}
