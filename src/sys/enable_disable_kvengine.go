// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/sys/enable_disable_kvengine.go
// Original timestamp: 2026/06/25 20:47:29

package sys

import (
	"fmt"
	"os"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/jeanfrancoisgratton/vaultlib/v2/sys"
	"golang.org/x/term"

	"vclt/shared"
)

func (c *Client) EnableKVengine(kvEngine string) *ce.CustomError {
	err := c.vc.EnableKVEngine(kvEngine, sys.EnableKVOptions{
		Version: KVEngineVersion, Description: KVEngineDescription})
	if err != nil {
		return &ce.CustomError{Title: "Error enabling kv engine", Message: err.Error()}
	}

	if !shared.QuietOutput {
		fmt.Printf("KV engine %s %s\n", kvEngine, hftx.Green("ENABLED"))
	}
	return nil
}

func (c *Client) DisableKVengine(kvEngine string) *ce.CustomError {
	if !KVDisableConfirm {
		fmt.Print(hftx.WarningSign(" CAUTION. This operation is irreversible; are you sure you want to disable it (Y/N) ? "))
		if yesno, err := AskYesNo(); err != nil {
			return err
		} else {
			fmt.Println()
			if !yesno {
				return &ce.CustomError{Fatality: ce.Continuable, Title: " Engine not disabled", Message: "User cancelled"}
			}
		}
	}
	err := c.vc.DisableKVEngine(kvEngine)
	if err != nil {
		return &ce.CustomError{Title: "Error disabling kv engine", Message: err.Error()}
	}

	if !shared.QuietOutput {
		fmt.Printf("KV engine %s %s\n", kvEngine, hftx.Red("DISABLED"))
	}
	return nil
}

func AskYesNo() (bool, *ce.CustomError) {
	fd := int(os.Stdin.Fd())

	for {
		// Put terminal into raw mode to capture individual keypresses
		oldState, err := term.MakeRaw(fd)
		if err != nil {
			// Fallback/Safety: If not a terminal (e.g. piped input), panic or handle
			return false, &ce.CustomError{Title: "Error reading terminal", Message: err.Error()}
		}

		var b [1]byte
		_, err = os.Stdin.Read(b[:])

		// Immediately restore the terminal so it behaves normally
		term.Restore(fd, oldState)

		if err != nil {
			continue // Keep trying if a read error occurs
		}

		// Check the key pressed
		switch b[0] {
		case 'y', 'Y':
			return true, nil
		case 'n', 'N':
			return false, nil
		default:
			// Ignore any other keypress and loop back to wait again
			continue
		}
	}
}
