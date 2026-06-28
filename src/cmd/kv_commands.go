// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/kv_commands.go
// Original timestamp: 2026/06/14 12:54:55

package cmd

import (
	"fmt"
	"os"

	"vclt/kv"

	hftfx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
)

var kvCmd = &cobra.Command{
	Use:   "kv",
	Short: "kv secret management subcommands",
	Long:  `Allowed commands are { read | write | list | delete | destroy | backup | restore }`,
}

var kvReadCmd = &cobra.Command{
	Use:     "read KV_ENGINE SECRET_PATH",
	Aliases: []string{"get"},
	Short:   "Read the 'SECRET_PATH' secret from the 'KV_ENGINE' secret engine",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if kvreadErr := kv.ReadSecrets(args[0], args[1]); kvreadErr != nil {
			fmt.Println(hftfx.SkullBonesSign(kvreadErr.Error()))
			os.Exit(1)
		}
	},
}

var kvWriteCmd = &cobra.Command{
	Use:     "write KV_ENGINE SECRET_PATH KEY VALUE",
	Aliases: []string{"put"},
	Short:   "Write the KEY:VALUE pair SECRET in the 'SECRET_PATH' of the 'KV_ENGINE' secret engine",
	Args:    cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		if _, kvwriteErr := kv.WriteSecrets(args[0], args[1], args[2], args[3]); kvwriteErr != nil {
			fmt.Println(hftfx.SkullBonesSign(kvwriteErr.Error()))
			os.Exit(1)
		}
	},
}

var kvLsCmd = &cobra.Command{
	Use:     "list KV_ENGINE",
	Aliases: []string{"ls", "show"},
	Short:   "List the kv in the 'KV_ENGINE' secret engine",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := kv.ListSecrets(args[0], true); err != nil {
			fmt.Println(hftfx.SkullBonesSign(err.Error()))
		}
	},
}

var kvRmCmd = &cobra.Command{
	Use:     "rm KV_ENGINE SECRET_PATH",
	Aliases: []string{"delete"},
	Short:   "Delete a secret or a field in a secret",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := kv.DeleteSecret(args[0], args[1]); err != nil {
			fmt.Println(hftfx.SkullBonesSign(err.Error()))
		}
	},
}

var kvDestroyCmd = &cobra.Command{
	Use:   "destroy KV_ENGINE SECRET_PATH",
	Short: "Destroy a secret",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := kv.DestroySecret(args[0], args[1]); err != nil {
			fmt.Println(hftfx.SkullBonesSign(err.Error()))
		}
	},
}

var kvBackupCmd = &cobra.Command{
	Use:     "backup KV_ENGINE BACKUP_FILE[.json]",
	Aliases: []string{"dump"},
	Short:   "Backup a kv engine",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := kv.BackupEngine(args[0], args[1]); err != nil {
			fmt.Println(hftfx.SkullBonesSign(err.Error()))
		}
	},
}

var kvRestoreCmd = &cobra.Command{
	Use:     "restore KV_ENGINE BACKUP_FILE[.json]",
	Aliases: []string{"import"},
	Short:   "Restore a kv engine",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := kv.RestoreEngine(args[0], args[1]); err != nil {
			fmt.Println(hftfx.SkullBonesSign(err.Error()))
		}
	},
}

func init() {
	rootCmd.AddCommand(kvCmd)
	kvCmd.AddCommand(kvLsCmd, kvReadCmd, kvWriteCmd, kvRmCmd, kvDestroyCmd, kvBackupCmd, kvRestoreCmd)

	//secretsCmd.PersistentFlags().StringVarP(&kv.SecretMountPath, "mount", "m", "", "KV v2 mount path (required)")
	kvCmd.PersistentFlags().IntVarP(&kv.SecretVersion, "version", "v", 0, "Secret version (0 = latest available)")
	kvReadCmd.PersistentFlags().StringVarP(&kv.SecretField, "field", "f", "", "Specific field to manage")
	kvRmCmd.PersistentFlags().StringVarP(&kv.SecretField, "field", "f", "", "Specific field to manage")
	kvLsCmd.PersistentFlags().BoolVarP(&kv.ExtendedSecretsList, "extended", "x", false, "Show extended info")
}
