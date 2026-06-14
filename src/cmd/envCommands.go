// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/envCommands.go
// Original timestamp: 2023/09/15 13:40

package cmd

import (
	"fmt"
	"os"

	cerr "github.com/jeanfrancoisgratton/customError/v3"
	"github.com/spf13/cobra"
	"vclt/env"
	environment "vclt/env"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment sub-command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Valid subcommands are: { list | add | remove | createrootkeys | info }")
	},
}

var envListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Example: "cm env list [directory]",
	Short:   "Lists all env files",
	Run: func(cmd *cobra.Command, args []string) {
		argument := ""
		if len(args) > 0 {
			argument = args[0]
		}
		if err := environment.ListEnvironments(argument); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var envRmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"remove"},
	Example: "cm env remove { FILE[.json] | defaultEnv.json }",
	Short:   "Removes the env FILE",
	Run: func(cmd *cobra.Command, args []string) {
		fname := ""
		if len(args) == 0 {
			fname = "defaultEnv.json"
		} else {
			fname = args[0]
		}
		var err *cerr.CustomError
		if err = environment.RemoveEnvFile(fname); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var envAddCmd = &cobra.Command{
	Use: "add",
	//Aliases: []string{"create"},
	Example: "cm env add [FILE[.json]]",
	Short:   "Adds the env FILE",
	Long: `The extension (.json) is implied and will be added if missing. Moreover, not specifying a filename
Will create a defaultEnv.json file, which is the application's default file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fname := ""
		if len(args) == 0 {
			fname = "defaultEnv.json"
		} else {
			fname = args[0]
		}
		if err := environment.AddEnvFile(fname); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var envInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"explain"},
	Example: "cm env info FILE1[.json] FILE2[.json]... FILEn[.json]",
	Short:   "Prints the env FILE[12n] information",
	Long:    `You can list as many env files as you wish, here`,
	Run: func(cmd *cobra.Command, args []string) {
		envfiles := []string{"defaultEnv.json"}
		if len(args) != 0 {
			envfiles = args
		}
		var err *cerr.CustomError
		if err = environment.ExplainEnvFile(envfiles); err != nil {
			fmt.Println(err.Error())
		}
	},
}

// envAddCmd : prompts the user for the values, while this one uses the ones passed as flags
var envCreateCmd = &cobra.Command{
	Use:     "create",
	Example: "vclt env create FILE[.json] -n ENVIRONMENT_NAME -a VAULT_ADDRESS -u VAULT_USERNAME -p VAULT_PASSWORD -k KV_STORE_PATH -c COMMENTS",
	Short:   "Creates a new env FILE using command line parameters",
	Long: `This will create a new env FILE using command line parameters.
The following parameters are mandatory (unless written otherwise),
*after* you've provided a filename (no extension needed)' :
-n EnvironmentName : the environment name
-a VaultAddress : the Vault server address (ex: https://myvault.example.com:8200)
-t VaultToken : Auth token
-u VaultUserName : the Vault username; that is, who will be connect to the vault
-p VaultPassword : the Vault username password
-k KVstorePath : the secrets engine path (name)
(optional) -c Comments: says it all.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must specify a file to create")
			os.Exit(1)
		}
		var err *cerr.CustomError
		if err = environment.CreateEnvFile(args[0]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var envCreateRootKeysCmd = &cobra.Command{
	Use:     "createrootkeys",
	Example: "vclt env createrootkeys",
	Aliases: []string{"crk"},
	Short:   "Creates a file containing the root keys needed to unseal a Vault",
	Long:    `This file will be named rootkeys.json, in $HOME/.config/JFG/vclt. The entries will be encoded`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := environment.CreateRootKeys(3); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	envCmd.AddCommand(envListCmd, envRmCmd, envAddCmd, envInfoCmd, envCreateCmd, envCreateRootKeysCmd)

	//var EnvName, VaultAddress, VaultUserName, VaultPassword, KVstorePath, Comments string
	envCreateCmd.PersistentFlags().StringVarP(&env.EnvName, "name", "n", "", "The new environment name")
	envCreateCmd.PersistentFlags().StringVarP(&env.VAddress, "address", "a", "", "Your vault server address (ex: https://myvault.example.com:8200)")
	envCreateCmd.PersistentFlags().StringVarP(&env.VUserName, "token", "t", "", "User's Auth token")
	envCreateCmd.PersistentFlags().StringVarP(&env.VUserName, "user", "u", "", "The username")
	envCreateCmd.PersistentFlags().StringVarP(&env.VPassword, "password", "p", "", "That user's password")
	envCreateCmd.PersistentFlags().StringVarP(&env.KVstorePath, "kvstore", "k", "", "The KV mount path")
	envCreateCmd.PersistentFlags().StringVarP(&env.EnvComments, "comments", "c", "", "Optional comments")
}
