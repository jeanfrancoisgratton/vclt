// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/kv_commands.go
// Original timestamp: 2026/06/14 12:54:55

package cmd

import (
	"fmt"

	"vclt/shared"

	"vclt/kv"

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
		c, err := kv.NewClient(args[0])
		if err != nil {
			err.Die()
		}
		if kvErr := c.Read(args[1]); kvErr != nil {
			kvErr.Die()
		}
	},
}

var kvWriteCmd = &cobra.Command{
	Use:     "write KV_ENGINE SECRET_PATH KEY [VALUE]",
	Aliases: []string{"put"},
	Short:   "Write the KEY:VALUE pair SECRET in the 'SECRET_PATH' of the 'KV_ENGINE' secret engine",
	Long: `Write the KEY:VALUE pair SECRET in the 'SECRET_PATH' of the 'KV_ENGINE' secret engine.

VALUE is read from the --in FILE instead of the command line when --in is given,
in which case the VALUE argument must be omitted:
  vclt kv write KV_ENGINE SECRET_PATH KEY --in FILE

With --batch, KEY and VALUE are both omitted in favor of a single BATCH_FILE
argument holding one KEY/VALUE field per line, all written in a single Vault
API call (one new KV version, preserving any other fields already on the
secret):
  vclt kv write --batch KV_ENGINE SECRET_PATH BATCH_FILE`,
	// VALUE is positional unless --in or --batch is set, in which case it
	// must be omitted (the value(s) come from a file instead) to avoid
	// ambiguity about which one wins. --batch and --in cannot be combined:
	// they disagree about what the 3rd positional argument means (KEY vs.
	// BATCH_FILE).
	Args: func(cmd *cobra.Command, args []string) error {
		if kv.BatchMode && kv.SecretInputFile != "" {
			return fmt.Errorf("--batch and --in cannot be used together")
		}
		if kv.BatchMode || kv.SecretInputFile != "" {
			return cobra.ExactArgs(3)(cmd, args)
		}
		return cobra.ExactArgs(4)(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		c, err := kv.NewClient(args[0])
		if err != nil {
			err.Die()
		}
		if kv.BatchMode {
			if kvErr := c.WriteBatch(args[1], args[2]); kvErr != nil {
				kvErr.Die()
			}
			return
		}
		value := ""
		if kv.SecretInputFile == "" {
			value = args[3]
		}
		if _, kvErr := c.Write(args[1], args[2], value); kvErr != nil {
			kvErr.Die()
		}
	},
}

var kvLsCmd = &cobra.Command{
	Use:     "list KV_ENGINE",
	Aliases: []string{"ls", "show"},
	Short:   "List the kv in the 'KV_ENGINE' secret engine",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := kv.NewClient(args[0])
		if err != nil {
			err.Die()
		}
		if _, kvErr := c.List(true); kvErr != nil {
			kvErr.Die()
		}
	},
}

var kvRmCmd = &cobra.Command{
	Use:     "rm KV_ENGINE SECRET_PATH",
	Aliases: []string{"delete"},
	Short:   "Delete a secret or a field in a secret",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := kv.NewClient(args[0])
		if err != nil {
			err.Die()
		}
		if kvErr := c.Delete(args[1]); kvErr != nil {
			kvErr.Die()
		}
	},
}

var kvDestroyCmd = &cobra.Command{
	Use:   "destroy KV_ENGINE SECRET_PATH",
	Short: "Destroy a secret",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := kv.NewClient(args[0])
		if err != nil {
			err.Die()
		}
		if kvErr := c.Destroy(args[1]); kvErr != nil {
			kvErr.Die()
		}
	},
}

var kvBackupCmd = &cobra.Command{
	Use:     "backup KV_ENGINE BACKUP_FILE[.json]",
	Aliases: []string{"dump"},
	Short:   "Backup a kv engine",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := kv.NewClient(args[0])
		if err != nil {
			err.Die()
		}
		if kvErr := c.Backup(args[1]); kvErr != nil {
			kvErr.Die()
		}
	},
}

var kvRestoreCmd = &cobra.Command{
	Use:     "restore KV_ENGINE BACKUP_FILE[.json]",
	Aliases: []string{"import"},
	Short:   "Restore a kv engine",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := kv.NewClient(args[0])
		if err != nil {
			err.Die()
		}
		if kvErr := c.Restore(args[1]); kvErr != nil {
			kvErr.Die()
		}
	},
}

func init() {
	rootCmd.AddCommand(kvCmd)
	kvCmd.AddCommand(kvLsCmd, kvReadCmd, kvWriteCmd, kvRmCmd, kvDestroyCmd, kvBackupCmd, kvRestoreCmd)

	//secretsCmd.PersistentFlags().StringVarP(&kv.SecretMountPath, "mount", "m", "", "KV v2 mount path (required)")
	kvCmd.PersistentFlags().IntVarP(&kv.SecretVersion, "version", "v", 0, "Secret version (0 = latest available)")
	kvReadCmd.PersistentFlags().StringVarP(&kv.SecretField, "field", "f", "", "Specific field to manage")
	kvReadCmd.PersistentFlags().StringVarP(&shared.OutputFormat, "outputformat", "o", "text", "Output format: text|json")
	kvReadCmd.PersistentFlags().StringVar(&kv.SecretOutputFile, "out", "", "Write the secret to FILE (mode 0600) instead of stdout")
	kvWriteCmd.PersistentFlags().StringVar(&kv.SecretInputFile, "in", "", "Read the secret VALUE from FILE instead of the command line")
	kvWriteCmd.PersistentFlags().BoolVar(&kv.BatchMode, "batch", false, "Write multiple KEY/VALUE fields from a batch file in one call (usage: write --batch KV_ENGINE SECRET_PATH BATCH_FILE)")
	kvRmCmd.PersistentFlags().StringVarP(&kv.SecretField, "field", "f", "", "Specific field to manage")
	kvLsCmd.PersistentFlags().BoolVarP(&kv.ExtendedSecretsList, "extended", "x", false, "Show extended info")
	kvBackupCmd.PersistentFlags().BoolVarP(&kv.Cleartext, "cleartext", "c", false, "Backup cleartext (default: encrypted)")
	kvRestoreCmd.PersistentFlags().BoolVarP(&kv.Cleartext, "cleartext", "c", false, "Restore cleartext (default: encrypted)")
}
