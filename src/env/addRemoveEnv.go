// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/env/addRemoveEnv.go
// Original timestamp: 2023/09/15 08:23

package env

import (
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError"
	"os"
	"path/filepath"
	"strings"
)

func RemoveEnvFile(envfile string) *cerr.CustomError {
	if !strings.HasSuffix(envfile, ".json") {
		envfile += ".json"
	}
	if err := os.Remove(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vclt", envfile)); err != nil {
		return &cerr.CustomError{Title: "Error removing " + envfile, Message: err.Error()}
	}

	fmt.Printf("%s removed succesfully\n", envfile)
	return nil
}

func AddEnvFile(envfile string) *cerr.CustomError {
	if !strings.HasSuffix(envfile, ".json") {
		envfile += ".json"
	}

	if env, err := getEnvParams(envfile); err != nil {
		return err
	} else {
		es := &env
		return es.SaveEnvironmentFile(envfile)
	}
}
