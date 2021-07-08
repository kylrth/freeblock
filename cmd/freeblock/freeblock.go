package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kylrth/freeblock/cmd/freeblock/cmds"
)

// Cmd is the root freeblock command.
var Cmd = &cobra.Command{
	Use:   "freeblock",
	Short: "freeblock provides tools for blocking and unblocking websites using /etc/hosts",
}

func init() {
	Cmd.AddCommand(
		cmds.BlockCmd,
		cmds.OpenCmd,
		cmds.UnblockCmd,
	)
}

func main() {
	if err := Cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
