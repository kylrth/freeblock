// Package cmds provides Cobra commands for various bits of freeblock functionality.
package cmds

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kylrth/freeblock/pkg/hosts"
)

// BlockCmd is a command that blocks domains in the hosts file.
var BlockCmd = &cobra.Command{
	Use:   "block DOMAIN [DOMAIN...]",
	Short: "block domains",
	Long: `Block domains by adding a 0.0.0.0 entry to /etc/hosts for each domain.

For hosts already present in the file, the address is set to 0.0.0.0 and the old
address is kept as a comment at the end of the line. If there is a commented-out
line for a domain, that line is uncommented.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := Block(args, hostsFile); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

var hostsFile string

func init() {
	BlockCmd.Flags().StringVar(
		&hostsFile, "hosts-file", "/etc/hosts", "Change the default hosts file.")
}

const (
	blockedIP     = "0.0.0.0"
	commentPrefix = " # "
)

// Block blocks the domains in the hostsFile.
func Block(domains []string, hostsFile string) error {
	lines, err := readLines(hostsFile)
	if err != nil {
		return err
	}

	blocked := make(map[string]bool, len(domains))

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

			line.Uncomment()
			oldIP := line.GetIP()
			if oldIP != blockedIP {
				line.SetIP(blockedIP)
				line += hosts.Line(commentPrefix + oldIP)
			}
			lines[i] = line

			for _, h := range hostnames {
				blocked[h] = true
			}

			break
		}
	}

	// Add lines for sites that haven't been blocked yet.
	for _, domain := range domains {
		if blocked[domain] {
			continue
		}
		lines = append(lines, hosts.Line(blockedIP+" "+domain))
	}

	return writeLines(lines, hostsFile)
}

func readLines(hostsFile string) ([]hosts.Line, error) {
	f, err := os.Open(hostsFile)
	if err != nil {
		return nil, err
	}
	lines, err := hosts.ReadLines(f)
	if err != nil {
		f.Close()

		return lines, fmt.Errorf("read lines: %w", err)
	}

	return lines, f.Close()
}

func writeLines(lines []hosts.Line, hostsFile string) error {
	f, err := os.Create(hostsFile)
	if err != nil {
		return fmt.Errorf("write hosts file: %w", err)
	}
	err = hosts.WriteLines(f, lines)
	if err != nil {
		f.Close()

		return fmt.Errorf("write hosts file: %w", err)
	}

	return f.Close()
}

func contains(arr []string, v string) bool {
	for _, a := range arr {
		if a == v {
			return true
		}
	}

	return false
}
