package cmds

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// OpenCmd is a command that temporarily unblocks domains in a hosts file.
var OpenCmd = &cobra.Command{
	Use:   "open DOMAIN [DOMAIN...]",
	Short: "open domains while the command is running",
	Long: `Temporarily unblock domains using the 'unblock' command, and then block them
again before exiting when a SIGINT is received.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

		if err := Open(args, hostsFile, osSignals); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	OpenCmd.Flags().StringVar(
		&hostsFile, "hosts-file", defaultHostsFile, "Change the default hosts file.")
}

// Open temporarily unblocks the domains in the hostsFile, and then closes them when it receives a
// signal on osSignals.
func Open(domains []string, hostsFile string, osSignals <-chan os.Signal) (err error) {
	// Back up the original lines in the hosts file, so that we can revert at the end.
	backupLines, err := readLines(hostsFile)
	if err != nil {
		return fmt.Errorf("backup hosts file: %w", err)
	}
	defer func() {
		var as *ErrBlockTiming
		if errors.As(err, &as) {
			// No need to restore because we were kept from unblocking.
			return
		}

		fmt.Fprintln(os.Stderr, "\nRestoring old hosts file...")
		e := writeLines(backupLines, hostsFile)
		if e != nil && err == nil {
			err = fmt.Errorf("restore backed up hosts file: %w", err)

			return
		}
		fmt.Fprintln(os.Stderr, "\tdone.")
	}()

	err = Unblock(domains, hostsFile, DefaultNower{})
	if err != nil {
		return fmt.Errorf("unblock domains: %w", err)
	}

	fmt.Fprintln(os.Stderr, "Domains temporarily unblocked:")
	for _, domain := range domains {
		fmt.Fprintf(os.Stderr, "- %s\n", domain)
	}

	// Wait for a SIGINT or SIGTERM before restoring the old file and exiting.
	<-osSignals

	return nil
}
