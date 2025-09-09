// vclt
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/kvCommands.go
// Original timestamp: 2024/06/28 14:20

package cmd

import (
	"fmt"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"vclt/kv"
	"vclt/sys"
)

var kvCmd = &cobra.Command{
	Use:     "kv",
	Example: "vclt kv { get | put | add | list }",
	Short:   "kv store sub-command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You need to specify one of the following subcommand: get | put | add | list")
		os.Exit(0)
	},
}

var kvGetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"read"},
	//Example: "vclt kv { get | put | add | list }",
	Short: "Read an entry in a secret",
	Run: func(cmd *cobra.Command, args []string) {
		nVer := 0 // this is the way to signify to the lib that we want the latest version of the secret
		if len(args) < 2 {
			fmt.Println("Usage: vclt kv get PATH FIELD VERSION")
		}
		
		if len(args) > 2 {
			nVer, _ = strconv.Atoi(args[2])
		}
		if res, ce := kv.Get(args[0], args[1], nVer); ce != nil {
			fmt.Println(ce.Error())
		} else {
			if sys.Verbose {
				fmt.Printf("%s: %s\n", args[1], hf.Green(fmt.Sprintf("%s", res)))
			} else {
				fmt.Printf("%s\n", res)
			}
		}

		os.Exit(0)
	},
}

//var kvLSfCmd = &cobra.Command{
//	Use:     "lsf",
//	Aliases: []string{"listfields"},
//	Short:   "Lists all fields in an entry",
//	Run: func(cmd *cobra.Command, args []string) {
//		nVer := 0
//		if len(args) < 1 {
//			fmt.Println("Usage: vclt lsf FIELD VERSION")
//		}
//		if len(args) > 1 {
//			nVer, _ = strconv.Atoi(args[1])
//		}
//		if res, ce := kv.ListFields(args[0], nVer); ce != nil {
//			ce.Error()
//		} else {
//			if !sys.Quiet {
//				fmt.Printf("Entry: %s\nFields:\n", hf.Green(args[0]))
//				for output, _ := range res {
//					fmt.Printf("\t%s\n", output)
//				}
//			} else {
//				for output, _ := range res {
//					fmt.Printf("%s\n", output)
//				}
//			}
//		}
//	},
//}
//
//var kvLSeCmd = &cobra.Command{
//	Use:     "lse",
//	Aliases: []string{"listentries"},
//	Short:   "Lists all entries in a store",
//	Run: func(cmd *cobra.Command, args []string) {
//		var res []string
//		var err *cerr.CustomError
//
//		if res, err = kv.ListEntries(); err != nil {
//			err.Error()
//		}
//
//		//if !sys.Quiet {
//		//	fmt.Println("Entries in current store")
//		//}
//		for _, output := range res {
//			fmt.Println(output)
//		}
//	},
//}

func init() {
	kvCmd.AddCommand(kvGetCmd /*, kvLSfCmd, kvLSeCmd*/)

	kvCmd.PersistentFlags().BoolVarP(&sys.Verbose, "verbose", "v", false, "Verbose output")
}
