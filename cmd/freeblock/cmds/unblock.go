package cmds

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// UnblockCmd is a command that unblocks domains in the hosts file.
var UnblockCmd = &cobra.Command{
	Use:   "unblock DOMAIN [DOMAIN...]",
	Short: "unblock domains",
	Long: `Unblock domains by commenting out 0.0.0.0 entries from the hosts file for each
domain.

For blocked hosts with a comment that has another IP address, the domain is
reverted back to to that IP address and the comment is deleted.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := Unblock(args, hostsFile); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	UnblockCmd.Flags().StringVar(
		&hostsFile, "hosts-file", defaultHostsFile, "Change the default hosts file.")
}

// Unblock blocks the domains in the hostsFile.
func Unblock(domains []string, hostsFile string) error {
	lines, err := readLines(hostsFile)
	if err != nil {
		return err
	}

	// Modify the lines in place.
	for i, line := range lines {
		hostnames := line.Hostnames()
		if len(hostnames) == 0 {
			continue
		}

		// See if this line refers to one or more of the domains we want to block.
		for _, hostname := range hostnames {
			if !contains(domains, hostname) {
				continue
			}

			oldIP := line.GetIP()
			if oldIP != blockedIP {
				// We don't want to comment this one out, because it's already unblocked.
				continue
			}

			// Check for a commented IP address.
			commentIdx := strings.Index(string(line), commentPrefix) //nolint:gocritic // false positive
			if commentIdx == -1 {
				// There's no IP address to revert to, so we'll just comment out the line.
				line.Comment()
			} else {
				// We'll revert to the commented IP address.
				commentFields := strings.Fields(string(line)[commentIdx+len(commentPrefix):])
				if len(commentFields) == 1 && net.ParseIP(commentFields[0]) != nil {
					line.SetIP(commentFields[0])
					line = line[:commentIdx]
				}
			}
			lines[i] = line

			break
		}
	}

	return writeLines(lines, hostsFile)
}
