// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/env/listInfoEnvs.go
// Original timestamp: 2023/09/13 16:01

package env

import (
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
	"path/filepath"
	"strings"
)

func ListEnvironments(envdir string) *cerr.CustomError {
	var err error
	var dirFH *os.File
	var finfo, fileInfos []os.FileInfo

	// list env files
	if envdir == "" {
		envdir = filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vclt")
	}
	if dirFH, err = os.Open(envdir); err != nil {
		ce := &cerr.CustomError{Title: "Unable to open config directory", Message: err.Error()}
		return ce
	}

	if fileInfos, err = dirFH.Readdir(0); err != nil {
		ce := &cerr.CustomError{Title: "Unable to read files in config directory", Message: err.Error()}
		return ce
	}

	for _, info := range fileInfos {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			finfo = append(finfo, info)
		}
	}

	fmt.Printf("Number of env files: %s\n", hf.Green(fmt.Sprintf("%d", len(finfo))))

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Environment file", "File size", "Modification time"})

	for _, fi := range finfo {
		t.AppendRow([]interface{}{hf.Green(fi.Name()), hf.Green(hf.SI(uint64(fi.Size()))),
			hf.Green(fmt.Sprintf("%v", fi.ModTime().Format("2006/01/02 15:04:05")))})
	}
	t.SortBy([]table.SortBy{
		{Name: "Environment file", Mode: table.Asc},
		{Name: "File size", Mode: table.Asc},
	})
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault
	t.Render()

	return nil
}

func ExplainEnvFile(envfiles []string) *cerr.CustomError {
	oldEnvFile := ConfigFile

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Environment file", "Env name", "Vault address", "Auth token", "Vault user", "Vault password",
		"KV Path", "Comments"})

	for _, envfile := range envfiles {
		if !strings.HasSuffix(envfile, ".json") {
			envfile += ".json"
		}
		ConfigFile = envfile

		if e, err := LoadEnvironmentFile(); err != nil {
			ConfigFile = oldEnvFile
			return err
		} else {
			vt := ""
			if e.VaultToken != "" {
				vt = "*ENCRYPTED*"
			}
			t.AppendRow([]interface{}{hf.Green(envfile), hf.Green(e.EnvironmentName), hf.Green(e.VaultAddress),
				hf.Yellow(vt), hf.Green(e.VaultUsername), hf.Yellow("*ENCRYPTED*"),
				hf.Green(e.KVEnginePath), hf.Green(e.Comments)})
		}

	}
	t.SortBy([]table.SortBy{
		{Name: "Environment file", Mode: table.Asc},
	})
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault
	t.Render()

	ConfigFile = oldEnvFile
	return nil
}
