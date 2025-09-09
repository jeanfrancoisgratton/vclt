package main

import (
	"fmt"
	"os"
	"path/filepath"
	"vclt/cmd"
)

func main() {
	var err error
	currentWorkingDir := ""
	// Whatever happens, we need to preserve the current pwd, and restore it on exit, however the software exits
	if currentWorkingDir, err = os.Getwd(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// First, we need to create a configuration directory. This is a per-user config dir
	if err = os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vclt"), os.ModePerm); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// We then launch the command loop
	cmd.Execute()

	// Software execution is complete, let's get the hell outta here
	_ = os.Chdir(currentWorkingDir)
}
